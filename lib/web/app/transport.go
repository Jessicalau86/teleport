/*
Copyright 2020 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"

	"github.com/gravitational/teleport/api/constants"
	"github.com/gravitational/teleport/api/types"
	apiutils "github.com/gravitational/teleport/api/utils"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/reversetunnelclient"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/tlsca"
	"github.com/gravitational/teleport/lib/utils"
)

// transportConfig is configuration for a rewriting transport.
type transportConfig struct {
	proxyClient  reversetunnelclient.Tunnel
	accessPoint  auth.ReadProxyAccessPoint
	cipherSuites []uint16
	identity     *tlsca.Identity
	servers      []types.AppServer
	ws           types.WebSession
	clusterName  string
	log          logrus.FieldLogger
}

// Check validates configuration.
func (c *transportConfig) Check() error {
	if c.proxyClient == nil {
		return trace.BadParameter("proxy client missing")
	}
	if c.accessPoint == nil {
		return trace.BadParameter("access point missing")
	}
	if len(c.cipherSuites) == 0 {
		return trace.BadParameter("cipe suites missing")
	}
	if c.identity == nil {
		return trace.BadParameter("identity missing")
	}
	if len(c.servers) == 0 {
		return trace.BadParameter("servers missing")
	}
	if c.ws == nil {
		return trace.BadParameter("web session missing")
	}
	if c.clusterName == "" {
		return trace.BadParameter("cluster name missing")
	}

	return nil
}

// isRemoteApp returns true if the route to app points at remote cluster.
func (c *transportConfig) isRemoteApp() bool {
	return c.clusterName != c.identity.RouteToApp.ClusterName
}

// transport is a rewriting http.RoundTripper that can forward requests to
// an application service.
type transport struct {
	c *transportConfig

	// tr is used for forwarding http connections.
	tr http.RoundTripper

	// clientTLSConfig is the TLS config used for mutual authentication.
	clientTLSConfig *tls.Config

	// mu protects access to servers in the transportConfig
	mu sync.Mutex
}

// newTransport creates a new transport.
func newTransport(c *transportConfig) (*transport, error) {
	var err error
	if err := c.Check(); err != nil {
		return nil, trace.Wrap(err)
	}

	t := &transport{c: c}

	t.clientTLSConfig, err = configureTLS(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// Clone and configure the transport.
	tr, err := defaults.Transport()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	tr.DialContext = t.DialContext
	tr.TLSClientConfig = t.clientTLSConfig

	t.tr = tr
	return t, nil
}

// RoundTrip will rewrite the request, forward the request to the target
// application, emit an event to the audit log, then rewrite the response.
func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Clone the request so we can modify it without affecting the original.
	// This is necessary because the request cookies are deleted when the web
	// handler forward the request to the app proxy. When this happens, the
	// cookies are lost and the error handler will not be able to find the
	// session based on cookies.
	r = r.Clone(r.Context())
	// Perform any request rewriting needed before forwarding the request.
	if err := t.rewriteRequest(r); err != nil {
		return nil, trace.Wrap(err)
	}

	// Forward the request to the target application.
	resp, err := t.tr.RoundTrip(r)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// When proxying app in leaf cluster, the app will have PublicAddr
	// possibly unreachable to the clients connecting to the root cluster.
	// We want to rewrite any redirects to that PublicAddr to equivalent redirect for root cluster.
	if err = t.rewriteRedirect(resp); err != nil {
		return nil, trace.Wrap(err)
	}

	return resp, nil
}

// rewriteRedirect rewrites redirect to a public address of downstream app.
func (t *transport) rewriteRedirect(resp *http.Response) error {
	if !t.c.isRemoteApp() {
		// Not a downstream app, short circuit.
		return nil
	}

	if !utils.IsRedirect(resp.StatusCode) {
		return nil
	}

	location := resp.Header.Get("Location")
	// nothing to do without Location.
	if location == "" {
		return nil
	}

	// Parse the "Location" header.
	u, err := url.Parse(location)
	if err != nil {
		return trace.Wrap(err)
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	// If this is a redirect to a public addr of a leaf cluster, rewrite it.
	// We want the rewrite to happen using our own public address.
	if host == t.c.identity.RouteToApp.PublicAddr {
		// drop scheme and host, leaving only the relative path.
		u.Host = ""
		u.Scheme = ""

		// since the path can be an empty string, canonicalize it to "/".
		if u.Path == "" {
			u.Path = "/"
		}

		resp.Header.Set("Location", u.String())
	}
	return nil
}

// rewriteRequest applies any rewriting rules to the request before it's forwarded.
func (t *transport) rewriteRequest(r *http.Request) error {
	// Set dummy values for the request forwarder. Dialing through the tunnel is
	// actually performed using the transport created for this session but these
	// are needed for the forwarder.
	r.URL.Scheme = "https"
	r.URL.Host = constants.APIDomain

	// Remove the application session cookie from the header. This is done by
	// first wiping out the "Cookie" header then adding back all cookies
	// except the application session cookies. This appears to be the safest way
	// to serialize cookies.
	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, cookie := range cookies {
		switch cookie.Name {
		case CookieName, SubjectCookieName:
			continue
		default:
			r.AddCookie(cookie)
		}
	}

	return nil
}

// DialContext dials and connect to the application service over the reverse
// tunnel subsystem.
func (t *transport) DialContext(ctx context.Context, _, _ string) (conn net.Conn, err error) {
	t.mu.Lock()
	servers := make([]types.AppServer, len(t.c.servers))
	copy(servers, t.c.servers)
	t.mu.Unlock()

	var i int
	for ; i < len(servers); i++ {
		appServer := servers[i]
		conn, err = dialAppServer(ctx, t.c.proxyClient, t.c.identity.RouteToApp.ClusterName, appServer)
		if err != nil && isReverseTunnelDownError(err) {
			t.c.log.Warnf("Failed to connect to application server %q: %v.", appServer, err)
			// Continue to the next server if there is an issue
			// establishing a connection because the tunnel is not
			// healthy. Reset the error to avoid returning it if
			// this is the last server.
			err = nil
			continue
		}

		break
	}

	// eliminate any servers from the head of the list that were unreachable
	t.mu.Lock()
	t.c.servers = t.c.servers[i:]
	t.mu.Unlock()

	if conn != nil || err != nil {
		return conn, trace.Wrap(err)
	}

	return nil, trace.ConnectionProblem(nil, "no application servers remaining to connect")
}

// DialWebsocket dials a websocket connection over the transport's reverse
// tunnel.
func (t *transport) DialWebsocket(network, address string) (net.Conn, error) {
	conn, err := t.DialContext(context.Background(), network, address)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	// App access connections over reverse tunnel use mutual TLS.
	return tls.Client(conn, t.clientTLSConfig), nil
}

// dialAppServer dial and connect to the application service over the reverse
// tunnel subsystem.
func dialAppServer(ctx context.Context, proxyClient reversetunnelclient.Tunnel, clusterName string, server types.AppServer) (net.Conn, error) {
	clusterClient, err := proxyClient.GetSite(clusterName)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var from net.Addr
	from = &utils.NetAddr{AddrNetwork: "tcp", Addr: "@web-proxy"}
	clientSrcAddr, originalDst := utils.ClientAddrFromContext(ctx)
	if clientSrcAddr != nil {
		from = clientSrcAddr
	}

	conn, err := clusterClient.Dial(reversetunnelclient.DialParams{
		From:                  from,
		To:                    &utils.NetAddr{AddrNetwork: "tcp", Addr: reversetunnelclient.LocalNode},
		OriginalClientDstAddr: originalDst,
		ServerID:              fmt.Sprintf("%v.%v", server.GetHostID(), clusterName),
		ConnType:              server.GetTunnelType(),
		ProxyIDs:              server.GetProxyIDs(),
	})
	return conn, trace.Wrap(err)
}

// configureTLS creates and configures a *tls.Config that will be used for
// mutual authentication.
func configureTLS(c *transportConfig) (*tls.Config, error) {
	tlsConfig := utils.TLSConfig(c.cipherSuites)

	// Configure the pool of certificates that will be used to verify the
	// identity of the server. This allows the client to verify the identity of
	// the server it is connecting to.
	ca, err := c.accessPoint.GetCertAuthority(context.TODO(), types.CertAuthID{
		Type:       types.HostCA,
		DomainName: c.identity.RouteToApp.ClusterName,
	}, false)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	certPool, err := services.CertPool(ca)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	tlsConfig.RootCAs = certPool

	// Configure the identity that will be used to connect to the server. This
	// allows the server to verify the identity of the caller.
	certificate, err := tls.X509KeyPair(c.ws.GetTLSCert(), c.ws.GetPriv())
	if err != nil {
		return nil, trace.Wrap(err, "failed to parse certificate or key")
	}
	tlsConfig.Certificates = []tls.Certificate{certificate}

	// Use SNI to tell the other side which cluster signed the CA so it doesn't
	// have to fetch all CAs when verifying the cert.
	tlsConfig.ServerName = apiutils.EncodeClusterName(c.clusterName)

	return tlsConfig, nil
}

// isReverseTunnelDownError returns true if the provided error indicates that
// the reverse tunnel connection is down e.g. because the agent is down.
func isReverseTunnelDownError(err error) bool {
	return trace.IsConnectionProblem(err) ||
		strings.Contains(err.Error(), reversetunnelclient.NoApplicationTunnel)
}
