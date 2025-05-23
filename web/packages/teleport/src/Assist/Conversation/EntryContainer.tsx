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

import React, { PropsWithChildren } from 'react';
import styled, { keyframes } from 'styled-components';

import { Author } from 'teleport/Assist/types';

interface EntryContainerProps {
  author: Author;
  index: number;
  length: number;
  streaming: boolean;
  lastMessage: boolean;
  hideOverflow?: boolean;
}

const Container = styled.div<EntryContainerProps>`
  display: flex;
  flex-direction: column;
  align-items: ${p =>
    p.author === Author.Teleport ? 'flex-start' : 'flex-end'};
  justify-content: ${p => (p.author === Author.Teleport ? '' : 'flex-end')};
  position: relative;
  font-size: 14px;
  margin-bottom: 5px;

  --content-overflow: ${p => (p.hideOverflow ? 'hidden' : 'visible')};
  --content-background: ${p =>
    p.author === Author.Teleport
      ? p.theme.colors.levels.popout
      : p.theme.colors.buttons.primary.default};
  --content-color: ${p =>
    p.author === Author.Teleport
      ? p.theme.colors.text.main
      : p.theme.colors.buttons.primary.text};
  --content-border-radius: ${p =>
    getBorderRadius(p.author === Author.Teleport, p.index, p.length)};
`;

const blink = keyframes`
  to {
    visibility: hidden;
  }
`;

const Content = styled.div`
  background: var(--content-background);
  color: var(--content-color);
  border-radius: var(--content-border-radius);
  box-shadow: 0 6px 12px -2px rgba(50, 50, 93, 0.05),
    0 3px 7px -3px rgba(0, 0, 0, 0.1);
  max-width: 90%;
  border: 1px solid ${p => p.theme.colors.spotBackground[1]};
  overflow: var(--content-overflow);

  &.streaming {
    > div > :not(ol):not(ul):not(pre):last-child:not(button):after,
    > div > ol:last-child li:last-child:after,
    > div > pre:last-child code:after,
    > div > ul:last-child li:last-child:after {
      animation: ${blink} 1s steps(5, start) infinite;
      content: '▋';
      margin-left: 0.25rem;
      vertical-align: baseline;
      opacity: 0.8;
    }
  }
`;

export function EntryContainer(props: PropsWithChildren<EntryContainerProps>) {
  const authorIsTeleport = props.author === Author.Teleport;
  const streaming = props.streaming && props.lastMessage && authorIsTeleport;

  return (
    <Container {...props}>
      <Content className={streaming ? 'streaming' : null}>
        {props.children}
      </Content>
    </Container>
  );
}

function getBorderRadius(isTeleport: boolean, index: number, length: number) {
  const isLast = index === length - 1;
  const isFirst = index === 0;

  if (isTeleport) {
    return `${isFirst ? '14px' : '5px'} 14px 14px ${isLast ? '14px' : '5px'}`;
  }

  return `14px ${isFirst ? '14px' : '5px'} ${isLast ? '14px' : '5px'} 14px`;
}
