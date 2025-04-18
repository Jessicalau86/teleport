/*
Copyright 2019 Gravitational, Inc.

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

import React from 'react';
import { useTheme } from 'styled-components';
import Image from 'design/Image';

import LogoHeroLight from './LogoHeroLight.svg';
import LogoHeroDark from './LogoHeroDark.svg';

const LogoHero = ({ ...rest }) => {
  const theme = useTheme();
  const src = theme.type === 'light' ? LogoHeroLight : LogoHeroDark;
  return <Image {...rest} src={src} />;
};

LogoHero.defaultProps = {
  src: LogoHeroDark,
  maxHeight: '120px',
  maxWidth: '200px',
  my: '48px',
  mx: 'auto',
};

export default LogoHero;
