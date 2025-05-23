/**
 * Copyright 2022 Gravitational, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ClusterResource } from 'teleport/services/userPreferences/types';

import type { JoinRole } from 'teleport/services/joinToken';

export enum ResourceKind {
  Application,
  Database,
  Desktop,
  Kubernetes,
  Server,
  SamlApplication,
  Discovery,
}

export function resourceKindToJoinRole(kind: ResourceKind): JoinRole {
  switch (kind) {
    case ResourceKind.Application:
      return 'App';
    case ResourceKind.Database:
      return 'Db';
    case ResourceKind.Desktop:
      return 'WindowsDesktop';
    case ResourceKind.Kubernetes:
      return 'Kube';
    case ResourceKind.Server:
      return 'Node';
    case ResourceKind.Discovery:
      return 'Discovery';
  }
}

export function resourceKindToPreferredResource(
  kind: ResourceKind
): ClusterResource {
  switch (kind) {
    case ResourceKind.Application:
      return ClusterResource.RESOURCE_WEB_APPLICATIONS;
    case ResourceKind.Database:
      return ClusterResource.RESOURCE_DATABASES;
    case ResourceKind.Desktop:
      return ClusterResource.RESOURCE_WINDOWS_DESKTOPS;
    case ResourceKind.Kubernetes:
      return ClusterResource.RESOURCE_KUBERNETES;
    case ResourceKind.Server:
      return ClusterResource.RESOURCE_SERVER_SSH;
    case ResourceKind.Discovery:
      return ClusterResource.RESOURCE_UNSPECIFIED;
  }
}
