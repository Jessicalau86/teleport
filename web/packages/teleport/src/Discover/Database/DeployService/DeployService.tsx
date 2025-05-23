/**
 * Copyright 2023 Gravitational, Inc.
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

import React, { useState } from 'react';

import { ServiceDeployMethod } from '../common';

import { ManualDeploy } from './ManualDeploy';
import { AutoDeploy } from './AutoDeploy';

export function DeployService() {
  const [deployMethod, setDeployMethod] = useState<ServiceDeployMethod>('auto');

  function toggleDeployMethod() {
    setDeployMethod(deployMethod === 'auto' ? 'manual' : 'auto');
  }

  if (deployMethod === 'auto') {
    return <AutoDeploy toggleDeployMethod={toggleDeployMethod} />;
  }

  return <ManualDeploy toggleDeployMethod={toggleDeployMethod} />;
}

export type DeployServiceProp = {
  toggleDeployMethod?(): void;
};
