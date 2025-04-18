/*
Copyright 2019-2021 Gravitational, Inc.

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

import React, { lazy } from 'react';

import {
  AddCircle,
  Application,
  CirclePlay,
  ClipboardUser,
  Cluster,
  Database,
  Desktop,
  EqualizersVertical,
  Integrations as IntegrationsIcon,
  Kubernetes,
  Laptop,
  ListThin,
  Lock,
  Question,
  Server,
  ShieldCheck,
  SlidersVertical,
  Terminal,
  UserCircleGear,
  Users as UsersIcon,
} from 'design/Icon';

import cfg from 'teleport/config';

import localStorage from 'teleport/services/localStorage';

import {
  ManagementSection,
  NavigationCategory,
} from 'teleport/Navigation/categories';

import { NavTitle } from './types';

import type { FeatureFlags, TeleportFeature } from './types';

const Audit = lazy(() => import('./Audit'));
const Nodes = lazy(() => import('./Nodes'));
const Sessions = lazy(() => import('./Sessions'));
const UnifiedResources = lazy(() => import('./UnifiedResources'));
const Account = lazy(() => import('./Account'));
const Applications = lazy(() => import('./Apps'));
const Kubes = lazy(() => import('./Kubes'));
const Support = lazy(() => import('./Support'));
const Clusters = lazy(() => import('./Clusters'));
const Trust = lazy(() => import('./TrustedClusters'));
const Users = lazy(() => import('./Users'));
const Roles = lazy(() => import('./Roles'));
const DeviceTrust = lazy(() => import('./DeviceTrust'));
const Recordings = lazy(() => import('./Recordings'));
const AuthConnectors = lazy(() => import('./AuthConnectors'));
const Locks = lazy(() => import('./LocksV2/Locks'));
const NewLock = lazy(() => import('./LocksV2/NewLock'));
const Databases = lazy(() => import('./Databases'));
const Desktops = lazy(() => import('./Desktops'));
const Discover = lazy(() => import('./Discover'));
const LockedAccessRequests = lazy(() => import('./AccessRequests'));
const Integrations = lazy(() => import('./Integrations'));
const IntegrationEnroll = lazy(
  () => import('@gravitational/teleport/src/Integrations/Enroll')
);

// ****************************
// Resource Features
// ****************************

class AccessRequests implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Access Requests',
    path: cfg.routes.accessRequest,
    exact: true,
    component: LockedAccessRequests,
  };

  hasAccess() {
    return true;
  }

  navigationItem = {
    title: NavTitle.AccessRequests,
    icon: <EqualizersVertical />,
    exact: true,
    getLink() {
      return cfg.routes.accessRequest;
    },
  };
}

export class FeatureNodes implements TeleportFeature {
  route = {
    title: 'Servers',
    path: cfg.routes.nodes,
    exact: true,
    component: Nodes,
  };

  navigationItem = {
    title: NavTitle.Servers,
    icon: <Server />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getNodesRoute(clusterId);
    },
  };

  hideFromNavigation = localStorage.areUnifiedResourcesEnabled();

  category = NavigationCategory.Resources;

  hasAccess(flags: FeatureFlags) {
    return flags.nodes;
  }
}

export class FeatureUnifiedResources implements TeleportFeature {
  route = {
    title: 'Resources',
    path: cfg.routes.unifiedResources,
    exact: true,
    component: UnifiedResources,
  };

  navigationItem = {
    title: NavTitle.Resources,
    icon: <Server />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getUnifiedResourcesRoute(clusterId);
    },
  };

  hideFromNavigation = !localStorage.areUnifiedResourcesEnabled();

  category = NavigationCategory.Resources;

  hasAccess() {
    return true;
  }
}

export class FeatureApps implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Applications',
    path: cfg.routes.apps,
    exact: true,
    component: Applications,
  };

  hideFromNavigation = localStorage.areUnifiedResourcesEnabled();

  hasAccess(flags: FeatureFlags) {
    return flags.applications;
  }

  navigationItem = {
    title: NavTitle.Applications,
    icon: <Application />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getAppsRoute(clusterId);
    },
  };
}

export class FeatureKubes implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Kubernetes',
    path: cfg.routes.kubernetes,
    exact: true,
    component: Kubes,
  };

  hideFromNavigation = localStorage.areUnifiedResourcesEnabled();

  hasAccess(flags: FeatureFlags) {
    return flags.kubernetes;
  }

  navigationItem = {
    title: NavTitle.Kubernetes,
    icon: <Kubernetes />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getKubernetesRoute(clusterId);
    },
  };
}

export class FeatureDatabases implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Databases',
    path: cfg.routes.databases,
    exact: true,
    component: Databases,
  };

  hideFromNavigation = localStorage.areUnifiedResourcesEnabled();

  hasAccess(flags: FeatureFlags) {
    return flags.databases;
  }

  navigationItem = {
    title: NavTitle.Databases,
    icon: <Database />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getDatabasesRoute(clusterId);
    },
  };
}

export class FeatureDesktops implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Desktops',
    path: cfg.routes.desktops,
    exact: true,
    component: Desktops,
  };

  hideFromNavigation = localStorage.areUnifiedResourcesEnabled();

  hasAccess(flags: FeatureFlags) {
    return flags.desktops;
  }

  navigationItem = {
    title: NavTitle.Desktops,
    icon: <Desktop />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getDesktopsRoute(clusterId);
    },
  };
}

export class FeatureSessions implements TeleportFeature {
  category = NavigationCategory.Resources;

  route = {
    title: 'Active Sessions',
    path: cfg.routes.sessions,
    exact: true,
    component: Sessions,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.activeSessions;
  }

  navigationItem = {
    title: NavTitle.ActiveSessions,
    icon: <Terminal />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getSessionsRoute(clusterId);
    },
  };
}

// ****************************
// Management Features
// ****************************

// - Access

export class FeatureUsers implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  route = {
    title: 'Manage Users',
    path: cfg.routes.users,
    exact: true,
    component: () => <Users />,
  };

  hasAccess(flags: FeatureFlags): boolean {
    return flags.users;
  }

  navigationItem = {
    title: NavTitle.Users,
    icon: <UsersIcon />,
    exact: true,
    getLink() {
      return cfg.getUsersRoute();
    },
  };

  getRoute() {
    return this.route;
  }
}

export class FeatureRoles implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  route = {
    title: 'Manage User Roles',
    path: cfg.routes.roles,
    exact: true,
    component: Roles,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.roles;
  }

  navigationItem = {
    title: NavTitle.Roles,
    icon: <ClipboardUser />,
    exact: true,
    getLink() {
      return cfg.routes.roles;
    },
  };
}

export class FeatureAuthConnectors implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  route = {
    title: 'Manage Auth Connectors',
    path: cfg.routes.sso,
    exact: false,
    component: AuthConnectors,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.authConnector;
  }

  navigationItem = {
    title: NavTitle.AuthConnectors,
    icon: <ShieldCheck />,
    exact: false,
    getLink() {
      return cfg.routes.sso;
    },
  };
}

export class FeatureLocks implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Identity;

  route = {
    title: 'Manage Session & Identity Locks',
    path: cfg.routes.locks,
    exact: true,
    component: Locks,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.locks;
  }

  navigationItem = {
    title: NavTitle.SessionAndIdentityLocks,
    icon: <Lock />,
    exact: false,
    getLink() {
      return cfg.getLocksRoute();
    },
  };
}

export class FeatureNewLock implements TeleportFeature {
  route = {
    title: 'Create New Lock',
    path: cfg.routes.newLock,
    exact: true,
    component: NewLock,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.newLocks;
  }

  // getRoute allows child class extending this
  // parent class to refer to this parent's route.
  getRoute() {
    return this.route;
  }
}

export class FeatureDiscover implements TeleportFeature {
  route = {
    title: 'Enroll New Resource',
    path: cfg.routes.discover,
    exact: true,
    component: Discover,
  };

  navigationItem = {
    title: NavTitle.EnrollNewResource,
    icon: <AddCircle />,
    exact: true,
    getLink() {
      return cfg.routes.discover;
    },
  };

  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  hasAccess(flags: FeatureFlags) {
    return flags.discover;
  }

  getRoute() {
    return this.route;
  }
}

export class FeatureIntegrations implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  hasAccess(flags: FeatureFlags) {
    return flags.integrations;
  }

  route = {
    title: 'Manage Integrations',
    path: cfg.routes.integrations,
    exact: true,
    component: () => <Integrations />,
  };

  navigationItem = {
    title: NavTitle.Integrations,
    icon: <IntegrationsIcon />,
    exact: true,
    getLink() {
      return cfg.routes.integrations;
    },
  };

  getRoute() {
    return this.route;
  }
}

export class FeatureIntegrationEnroll implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Access;

  route = {
    title: 'Enroll New Integration',
    path: cfg.routes.integrationEnroll,
    exact: false,
    component: () => <IntegrationEnroll />,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.enrollIntegrations;
  }

  navigationItem = {
    title: NavTitle.EnrollNewIntegration,
    icon: <AddCircle />,
    getLink() {
      return cfg.getIntegrationEnrollRoute(null);
    },
  };

  // getRoute allows child class extending this
  // parent class to refer to this parent's route.
  getRoute() {
    return this.route;
  }
}

// - Activity

export class FeatureRecordings implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Activity;

  route = {
    title: 'Session Recordings',
    path: cfg.routes.recordings,
    exact: true,
    component: Recordings,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.recordings;
  }

  navigationItem = {
    title: NavTitle.SessionRecordings,
    icon: <CirclePlay />,
    exact: true,
    getLink(clusterId: string) {
      return cfg.getRecordingsRoute(clusterId);
    },
  };
}

export class FeatureAudit implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Activity;

  route = {
    title: 'Audit Log',
    path: cfg.routes.audit,
    component: Audit,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.audit;
  }

  navigationItem = {
    title: NavTitle.AuditLog,
    icon: <ListThin />,
    getLink(clusterId: string) {
      return cfg.getAuditRoute(clusterId);
    },
  };
}

// - Clusters

export class FeatureClusters implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Clusters;

  route = {
    title: 'Clusters',
    path: cfg.routes.clusters,
    exact: false,
    component: Clusters,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.trustedClusters;
  }

  navigationItem = {
    title: NavTitle.ManageClusters,
    icon: <SlidersVertical />,
    exact: false,
    getLink() {
      return cfg.routes.clusters;
    },
  };
}

export class FeatureTrust implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Clusters;

  route = {
    title: 'Trusted Clusters',
    path: cfg.routes.trustedClusters,
    component: Trust,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.trustedClusters;
  }

  navigationItem = {
    title: NavTitle.TrustedClusters,
    icon: <Cluster />,
    getLink() {
      return cfg.routes.trustedClusters;
    },
  };
}

class FeatureDeviceTrust implements TeleportFeature {
  category = NavigationCategory.Management;
  section = ManagementSection.Identity;
  route = {
    title: 'Manage Trusted Devices',
    path: cfg.routes.deviceTrust,
    exact: true,
    component: DeviceTrust,
  };

  hasAccess(flags: FeatureFlags) {
    return flags.deviceTrust;
  }

  navigationItem = {
    title: NavTitle.TrustedDevices,
    icon: <Laptop />,
    exact: true,
    getLink() {
      return cfg.routes.deviceTrust;
    },
  };
}

// ****************************
// Other Features
// ****************************

export class FeatureAccount implements TeleportFeature {
  route = {
    title: 'Account Settings',
    path: cfg.routes.account,
    component: Account,
  };

  hasAccess() {
    return true;
  }

  topMenuItem = {
    title: NavTitle.AccountSettings,
    icon: <UserCircleGear />,
    getLink() {
      return cfg.routes.account;
    },
  };
}

export class FeatureHelpAndSupport implements TeleportFeature {
  route = {
    title: 'Help & Support',
    path: cfg.routes.support,
    exact: true,
    component: Support,
  };

  hasAccess() {
    return true;
  }

  topMenuItem = {
    title: NavTitle.HelpAndSupport,
    icon: <Question />,
    exact: true,
    getLink() {
      return cfg.routes.support;
    },
  };
}

export function getOSSFeatures(): TeleportFeature[] {
  return [
    // Resources
    new FeatureUnifiedResources(),
    new FeatureNodes(),
    new FeatureApps(),
    new FeatureKubes(),
    new FeatureDatabases(),
    new FeatureDesktops(),
    new AccessRequests(),
    new FeatureSessions(),

    // Management

    // - Access
    new FeatureUsers(),
    new FeatureRoles(),
    new FeatureAuthConnectors(),
    new FeatureIntegrations(),
    new FeatureDiscover(),
    new FeatureIntegrationEnroll(),

    // - Identity
    new FeatureLocks(),
    new FeatureNewLock(),
    new FeatureDeviceTrust(),

    // - Activity
    new FeatureRecordings(),
    new FeatureAudit(),

    // - Clusters
    new FeatureClusters(),
    new FeatureTrust(),

    // Other
    new FeatureAccount(),
    new FeatureHelpAndSupport(),
  ];
}
