/**
 * Copyright 2021 Gravitational, Inc.
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
import * as Alerts from 'design/Alert';
import { ButtonIcon, Text, ButtonSecondary, Image, Flex, Box } from 'design';
import DialogConfirmation, {
  DialogContent,
  DialogHeader,
  DialogFooter,
} from 'design/DialogConfirmation';
import { Attempt } from 'shared/hooks/useAsync';
import * as Icons from 'design/Icon';

import LinearProgress from 'teleterm/ui/components/LinearProgress';
import svgHardwareKey from 'teleterm/ui/ClusterConnect/ClusterLogin/FormLogin/PromptWebauthn/hardware.svg';

import type * as tsh from 'teleterm/services/tshd/types';

export type HeadlessPromptProps = {
  cluster: tsh.Cluster;
  clientIp: string;
  skipConfirm: boolean;
  onApprove(): Promise<void>;
  abortApproval(): void;
  /**
   * onReject updates the state of the request by rejecting it.
   */
  onReject(): Promise<void>;
  headlessAuthenticationId: string;
  updateHeadlessStateAttempt: Attempt<void>;
  /**
   * onCancel simply closes the modal and ignores the request. The user is still able to confirm or
   * reject the request from the Web UI.
   */
  onCancel(): void;
};

export function HeadlessPrompt({
  cluster,
  clientIp,
  skipConfirm,
  onApprove,
  abortApproval,
  onReject,
  headlessAuthenticationId,
  updateHeadlessStateAttempt,
  onCancel,
}: HeadlessPromptProps) {
  // skipConfirm automatically attempts to approve a headless auth attempt,
  // so let's show waitForMfa from the very beginning in that case.
  const [waitForMfa, setWaitForMfa] = useState(skipConfirm);

  return (
    <DialogConfirmation
      dialogCss={() => ({
        maxWidth: '480px',
        width: '100%',
      })}
      disableEscapeKeyDown={false}
      open={true}
    >
      <DialogHeader justifyContent="space-between" mb={0} alignItems="baseline">
        <Text typography="h4">
          Headless command on <b>{cluster.name}</b>
        </Text>
        <ButtonIcon
          type="button"
          color="text.slightlyMuted"
          onClick={() => {
            abortApproval();
            onCancel();
          }}
        >
          <Icons.Cross size="medium" />
        </ButtonIcon>
      </DialogHeader>
      {updateHeadlessStateAttempt.status === 'error' && (
        <Alerts.Danger mb={0}>
          {updateHeadlessStateAttempt.statusText}
        </Alerts.Danger>
      )}
      <DialogContent>
        <Text color="text.slightlyMuted">
          Someone initiated a headless command from <b>{clientIp}</b>.
          <br />
          If it was not you, click Reject and contact your administrator.
        </Text>
        <Text color="text.muted" mt={1} fontSize="12px">
          Request ID: {headlessAuthenticationId}
        </Text>
      </DialogContent>
      {waitForMfa && (
        <DialogContent mb={2}>
          <Text color="text.slightlyMuted">
            Complete MFA verification to approve the Headless Login.
          </Text>

          <Image mt={4} mb={4} width="200px" src={svgHardwareKey} mx="auto" />
          <Box textAlign="center" style={{ position: 'relative' }}>
            <Text bold>Insert your security key and tap it</Text>
            <LinearProgress />
          </Box>

          <Flex justifyContent="flex-end" mt={4} gap={3}>
            {/*
              The Reject button is there so that if skipping confirmation is enabled (see
              HeadlessAuthenticationService) then the user still has the ability to reject the
              request from the screen that prompts for key touch.
            */}
            <ButtonSecondary
              type="button"
              onClick={() => {
                abortApproval();
                onReject();
              }}
            >
              Reject
            </ButtonSecondary>
            <ButtonSecondary
              type="button"
              onClick={() => {
                abortApproval();
                onCancel();
              }}
            >
              Cancel
            </ButtonSecondary>
          </Flex>
        </DialogContent>
      )}
      {!waitForMfa && (
        <DialogFooter>
          <ButtonSecondary
            autoFocus
            mr={3}
            type="submit"
            onClick={e => {
              e.preventDefault();
              setWaitForMfa(true);
              onApprove();
            }}
          >
            Approve
          </ButtonSecondary>
          <ButtonSecondary
            type="button"
            onClick={e => {
              e.preventDefault();
              onReject();
            }}
          >
            Reject
          </ButtonSecondary>
        </DialogFooter>
      )}
    </DialogConfirmation>
  );
}
