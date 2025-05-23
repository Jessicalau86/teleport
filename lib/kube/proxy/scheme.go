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

package proxy

import (
	"errors"
	"strings"

	"github.com/gravitational/trace"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/gravitational/teleport/lib/utils"
)

const (
	// listSuffix is the suffix added to the name of the type to create the name
	// of the list type.
	// For example: "Role" -> "RoleList"
	listSuffix = "List"
)

var (
	// globalKubeScheme is the runtime Scheme that holds information about supported
	// message types.
	globalKubeScheme = runtime.NewScheme()
	// globalKubeCodecs creates a serializer/deserizalier for the different codecs
	// supported by the Kubernetes API.
	globalKubeCodecs = serializer.NewCodecFactory(globalKubeScheme)
)

// Register all groups in the schema's registry.
// It manually registers support for `metav1.Table` because go-client does not
// support it but `kubectl` calls require support for it.
func init() {
	// Register external types for Scheme
	utilruntime.Must(registerDefaultKubeTypes(globalKubeScheme))
}

// registerDefaultKubeTypes registers the default types for the Kubernetes API into
// the given scheme.
func registerDefaultKubeTypes(s *runtime.Scheme) error {
	// Register external types for Scheme
	metav1.AddToGroupVersion(s, schema.GroupVersion{Version: "v1"})
	if err := metav1.AddMetaToScheme(s); err != nil {
		return trace.Wrap(err)
	}
	if err := metav1beta1.AddMetaToScheme(s); err != nil {
		return trace.Wrap(err)
	}
	if err := scheme.AddToScheme(s); err != nil {
		return trace.Wrap(err)
	}
	err := s.SetVersionPriority(corev1.SchemeGroupVersion)
	return trace.Wrap(err)
}

// newClientNegotiator creates a negotiator that based on `Content-Type` header
// from the Kubernetes API response is able to create a different encoder/decoder.
// Supported content types:
// - "application/json"
// - "application/yaml"
// - "application/vnd.kubernetes.protobuf"
func newClientNegotiator(codecFactory *serializer.CodecFactory) runtime.ClientNegotiator {
	return runtime.NewClientNegotiator(
		codecFactory.WithoutConversion(),
		schema.GroupVersion{
			// create a serializer for Kube API v1
			Version: "v1",
		},
	)
}

// newClusterSchemaBuilder creates a new schema builder for the given cluster.
// This schema includes all well-known Kubernetes types and all namespaced
// custom resources.
// It also returns a map of resources that we support RBAC restrictions for.
func newClusterSchemaBuilder(client kubernetes.Interface) (serializer.CodecFactory, rbacSupportedResources, error) {
	kubeScheme := runtime.NewScheme()
	kubeCodecs := serializer.NewCodecFactory(kubeScheme)
	supportedResources := maps.Clone(defaultRBACResources)

	if err := registerDefaultKubeTypes(kubeScheme); err != nil {
		return serializer.CodecFactory{}, nil, trace.Wrap(err)
	}
	// discoveryErr is returned when the discovery of one or more API groups fails.
	var discoveryErr *discovery.ErrGroupDiscoveryFailed
	// register all namespaced custom resources
	_, apiGroups, err := client.Discovery().ServerGroupsAndResources()
	switch {
	case errors.As(err, &discoveryErr) && len(discoveryErr.Groups) == 1:
		// If the discovery error is of type `ErrGroupDiscoveryFailed` and it
		// contains only one group, it it's possible that the group is the metrics
		// group. If that's the case, we can ignore the error and continue.
		// This is a workaround for the metrics group not being registered because
		// the metrics pod is not running. It's common for Kubernetes clusters without
		// nodes to not have the metrics pod running.
		const metricsAPIGroup = "metrics.k8s.io"
		for k := range discoveryErr.Groups {
			if k.Group != metricsAPIGroup {
				return serializer.CodecFactory{}, nil, trace.Wrap(err)
			}
		}
	case err != nil:
		return serializer.CodecFactory{}, nil, trace.Wrap(err)
	}

	for _, apiGroup := range apiGroups {
		group, version := getKubeAPIGroupAndVersion(apiGroup.GroupVersion)
		// Skip well-known Kubernetes API groups because they are already registered
		// in the scheme.
		if _, ok := knownKubernetesGroups[group]; ok {
			continue
		}

		groupVersion := schema.GroupVersion{Group: group, Version: version}
		for _, apiResource := range apiGroup.APIResources {
			// Skip cluster-scoped resources because we don't support RBAC restrictions
			// for them.
			if !apiResource.Namespaced {
				continue
			}
			// build the resource key to be able to look it up later and check if
			// if we support RBAC restrictions for it.
			resourceKey := allowedResourcesKey{
				apiGroup:     group,
				resourceKind: apiResource.Name,
			}
			// Namespaced custom resources are allowed if the user has access to
			// the namespace where the resource is located.
			// This means that we need to map the resource to the namespace kind.
			supportedResources[resourceKey] = utils.KubeCustomResource
			// create the group version kind for the resource
			gvk := groupVersion.WithKind(apiResource.Kind)
			// check if the resource is already registered in the scheme
			// if it is, we don't need to register it again.
			if _, err := kubeScheme.New(gvk); err == nil {
				continue
			}
			// register the resource with the scheme to be able to decode it
			// into an unstructured object
			kubeScheme.AddKnownTypeWithName(
				gvk,
				&unstructured.Unstructured{},
			)
			// register the resource list with the scheme to be able to decode it
			// into an unstructured object.
			// Resource lists follow the naming convention: <resource-kind>List
			kubeScheme.AddKnownTypeWithName(
				groupVersion.WithKind(apiResource.Kind+listSuffix),
				&unstructured.Unstructured{},
			)
		}
	}

	return kubeCodecs, supportedResources, nil
}

// getKubeAPIGroupAndVersion returns the API group and version from the given
// groupVersion string.
// The groupVersion string can be in the following formats:
// - "v1" -> group: "", version: "v1"
// - "<group>/<version>" -> group: "<group>", version: "<version>"
func getKubeAPIGroupAndVersion(groupVersion string) (group string, version string) {
	splits := strings.Split(groupVersion, "/")
	switch {
	case len(splits) == 1:
		return "", splits[0]
	case len(splits) >= 2:
		return splits[0], splits[1]
	default:
		return "", ""
	}
}

// knownKubernetesGroups is a map of well-known Kubernetes API groups that
// are already registered in the scheme and we don't need to register them
// again.
var knownKubernetesGroups = map[string]struct{}{
	// core group
	"":                             {},
	"apiregistration.k8s.io":       {},
	"apps":                         {},
	"events.k8s.io":                {},
	"authentication.k8s.io":        {},
	"authorization.k8s.io":         {},
	"autoscaling":                  {},
	"batch":                        {},
	"certificates.k8s.io":          {},
	"networking.k8s.io":            {},
	"policy":                       {},
	"rbac.authorization.k8s.io":    {},
	"storage.k8s.io":               {},
	"admissionregistration.k8s.io": {},
	"apiextensions.k8s.io":         {},
	"scheduling.k8s.io":            {},
	"coordination.k8s.io":          {},
	"node.k8s.io":                  {},
	"discovery.k8s.io":             {},
	"flowcontrol.apiserver.k8s.io": {},
	"metrics.k8s.io":               {},
}
