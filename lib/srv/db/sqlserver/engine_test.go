/*
Copyright 2022 Gravitational, Inc.

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

package sqlserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/constants"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/types/events"
	libevents "github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/events/eventstest"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/srv/db/common"
	"github.com/gravitational/teleport/lib/srv/db/sqlserver/protocol"
	"github.com/gravitational/teleport/lib/srv/db/sqlserver/protocol/fixtures"
)

// TestHandleConnectionAuditEvents checks audit events emitted during HandleConnection execution.
func TestHandleConnectionAuditEvents(t *testing.T) {
	type check func(t *testing.T, err error, ee []events.AuditEvent)
	hasNoErr := func() check {
		return func(t *testing.T, err error, ee []events.AuditEvent) {
			require.NoError(t, err)
		}
	}
	hasAuditEventCode := func(want string) check {
		return func(t *testing.T, err error, ee []events.AuditEvent) {
			for _, v := range ee {
				if v.GetCode() == want {
					return
				}
			}
			require.Failf(t, "event not found", "event code: %s", want)
		}
	}
	hasAuditEvent := func(i int, want events.AuditEvent) check {
		return func(t *testing.T, err error, ee []events.AuditEvent) {
			diff := cmp.Diff(want, ee[i])
			require.Empty(t, diff)
		}
	}

	tests := []struct {
		name    string
		packets [][]byte
		checks  []check
	}{
		{
			name:    "rpc request procedure",
			packets: [][]byte{fixtures.GenerateCustomRPCCallPacket("foo3")},
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Procname: "foo3",
				}),
			},
		},
		{
			name:    "rpc request param",
			packets: [][]byte{fixtures.GenerateExecuteSQLRPCPacket("select @@version")},
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select @@version"},
					Procname:   "Sp_ExecuteSql",
				}),
			},
		},
		{
			name:    "sql batch",
			packets: [][]byte{fixtures.GenerateBatchQueryPacket("\nselect 'foo' as 'bar'\n        ")},
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.DatabaseSessionQuery{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionQueryEvent,
						Code: libevents.DatabaseSessionQueryCode,
					},
					DatabaseQuery: "\nselect 'foo' as 'bar'\n        ",
					Status: events.Status{
						Success: true,
					},
				}),
			},
		},
		{
			name:    "malformed packet",
			packets: [][]byte{fixtures.MalformedPacketTest},
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEventCode(libevents.DatabaseSessionMalformedPacketCode),
			},
		},
		{
			name:    "sql batch chunked packets",
			packets: fixtures.GenerateBatchQueryChunkedPacket(5, "select 'foo' as 'bar'"),
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.DatabaseSessionQuery{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionQueryEvent,
						Code: libevents.DatabaseSessionQueryCode,
					},
					DatabaseQuery: "select 'foo' as 'bar'",
					Status: events.Status{
						Success: true,
					},
				}),
			},
		},
		{
			name:    "rpc request param chunked",
			packets: fixtures.GenerateExecuteSQLRPCChunkedPacket(5, "select @@version"),
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select @@version"},
					Procname:   "Sp_ExecuteSql",
				}),
			},
		},
		{
			name: "intercalated chunked messages",
			packets: intercalateChunkedPacketMessages(
				fixtures.GenerateExecuteSQLRPCChunkedPacket(5, "select @@version"),
				fixtures.GenerateExecuteSQLRPCPacket("select 1"),
				2,
			),
			checks: []check{
				hasNoErr(),
				hasAuditEventCode(libevents.DatabaseSessionStartCode),
				hasAuditEventCode(libevents.DatabaseSessionEndCode),
				hasAuditEvent(1, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select @@version"},
					Procname:   "Sp_ExecuteSql",
				}),
				hasAuditEvent(2, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select 1"},
					Procname:   "Sp_ExecuteSql",
				}),
				hasAuditEvent(3, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select @@version"},
					Procname:   "Sp_ExecuteSql",
				}),
				hasAuditEvent(4, &events.SQLServerRPCRequest{
					DatabaseMetadata: events.DatabaseMetadata{
						DatabaseUser:     "sa",
						DatabaseType:     "self-hosted",
						DatabaseService:  "dummy",
						DatabaseURI:      "uri",
						DatabaseProtocol: "test",
					},
					Metadata: events.Metadata{
						Type: libevents.DatabaseSessionSQLServerRPCRequestEvent,
						Code: libevents.SQLServerRPCRequestCode,
					},
					Parameters: []string{"select 1"},
					Procname:   "Sp_ExecuteSql",
				}),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			_, err := b.Write(fixtures.Login7)
			require.NoError(t, err)

			db, err := types.NewDatabaseV3(types.Metadata{
				Name:   "dummy",
				Labels: map[string]string{"env": "prod"},
			}, types.DatabaseSpecV3{
				Protocol: "test",
				URI:      "uri",
			})
			require.NoError(t, err)

			for _, packet := range tc.packets {
				_, err = b.Write(packet)
				require.NoError(t, err)
			}

			emitterMock := &eventstest.MockRecorderEmitter{}
			audit, err := common.NewAudit(common.AuditConfig{
				Emitter:  emitterMock,
				Recorder: libevents.WithNoOpPreparer(libevents.NewDiscardRecorder()),
				Database: db,
			})
			require.NoError(t, err)

			e := Engine{
				EngineConfig: common.EngineConfig{
					Audit:   audit,
					Log:     logrus.New(),
					Auth:    &mockDBAuth{},
					Context: context.Background(),
				},
				Connector: &mockConnector{
					conn: &mockConn{
						buff: bytes.Buffer{},
					},
				},
				clientConn: &mockConn{
					buff: b,
				},
			}

			err = e.HandleConnection(context.Background(), &common.Session{
				Checker:  &mockChecker{},
				Database: db,
			})
			for _, ch := range tc.checks {
				ch(t, err, emitterMock.Events())
			}
		})
	}
}

// intercalateChunkedPacketMessages intercalates a chunked packet with a regular packet a specified number of times.
func intercalateChunkedPacketMessages(chunkedPacket [][]byte, regularPacket []byte, repeat int) [][]byte {
	var result [][]byte
	for i := 0; i < repeat; i++ {
		result = append(result, chunkedPacket...)
		result = append(result, regularPacket)
	}

	return result
}

type mockConn struct {
	net.Conn
	buff bytes.Buffer
}

func (o *mockConn) Read(p []byte) (n int, err error) {
	return o.buff.Read(p)
}

func (o *mockConn) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (o *mockConn) Close() error {
	return nil
}

type mockDBAuth struct {
	common.Auth
	// GetAzureIdentityResourceID mocks.
	azureIdentityResourceID    string
	azureIdentityResourceIDErr error
}

func (m *mockDBAuth) GetAuthPreference(ctx context.Context) (types.AuthPreference, error) {
	return types.NewAuthPreference(types.AuthPreferenceSpecV2{
		Type:         constants.Local,
		SecondFactor: constants.SecondFactorWebauthn,
		Webauthn: &types.Webauthn{
			RPID: "localhost",
		},
		RequireMFAType: types.RequireMFAType_SESSION,
	})
}

func (m *mockDBAuth) GetTLSConfig(_ context.Context, _ *common.Session) (*tls.Config, error) {
	return &tls.Config{}, nil
}

func (m *mockDBAuth) GetAzureIdentityResourceID(_ context.Context, _ string) (string, error) {
	return m.azureIdentityResourceID, m.azureIdentityResourceIDErr
}

type mockChecker struct {
	services.AccessChecker
}

func (m *mockChecker) CheckAccess(r services.AccessCheckable, state services.AccessState, matchers ...services.RoleMatcher) error {
	return nil
}

func (m *mockChecker) GetAccessState(authPref types.AuthPreference) services.AccessState {
	if authPref.GetRequireMFAType().IsSessionMFARequired() {
		return services.AccessState{
			MFARequired: services.MFARequiredAlways,
		}
	}
	return services.AccessState{
		MFARequired: services.MFARequiredNever,
	}
}

type mockConnector struct {
	conn io.ReadWriteCloser
}

func (m *mockConnector) Connect(context.Context, *common.Session, *protocol.Login7Packet) (io.ReadWriteCloser, []mssql.Token, error) {
	return m.conn, []mssql.Token{}, nil
}
