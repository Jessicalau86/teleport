/**
 * Copyright 2023 Gravitational, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import process from 'node:process';
import childProcess from 'node:child_process';
import { setTimeout } from 'node:timers/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const logsDir = process.argv[2];
// sendPidsImmediately controls whether this process is going to report children PIDs to its parent
// immediately or only after children report being ready.
const sendPidsImmediately = process.argv[3] === 'sendPidsImmediately';
// ignoreSigterm controls whether the agent process ignores SIGTERM or not.
const ignoreSigterm = process.argv[4] === 'ignoreSigterm';

if (!logsDir) {
  throw new Error(
    'Logs directory must be passed over argv as the first argument'
  );
}

// Workaround for the lack of __dirname in ESM modules.
const __dirname = path.dirname(fileURLToPath(import.meta.url));

const agent = childProcess.fork(
  path.join(__dirname, 'agentTestProcess.mjs'),
  ignoreSigterm ? ['ignoreSigterm'] : [],
  { stdio: 'inherit' }
);
const agentCleanupDaemon = childProcess.fork(
  path.join(__dirname, 'agentCleanupDaemon.js'),
  // Use a shorter timeout in tests. Each test needs to wait for the cleanup daemon to terminate,
  // so we don't want to spend full 5s on that.
  [agent.pid, process.pid, '/clusters/foo', logsDir, 50 /* timeToSigkill */],
  { stdio: 'inherit' }
);

const onceMessage = process =>
  new Promise(resolve => {
    process.once('message', resolve);
  });

const timeout = async ms => {
  await setTimeout(ms);
  throw new Error('timeout');
};

// parentTestProcess.mjs can be run directly from a terminal for debugging purposes, make sure to
// check process.send before using it.
if (process.send) {
  if (sendPidsImmediately) {
    process.send({
      agent: agent.pid,
      agentCleanupDaemon: agentCleanupDaemon.pid,
    });
  } else {
    // Wait to get messages from both children, signalling that they're set up and ready.
    Promise.race([
      Promise.all([onceMessage(agent), onceMessage(agentCleanupDaemon)]),
      timeout(2000),
    ]).then(
      () => {
        process.send({
          agent: agent.pid,
          agentCleanupDaemon: agentCleanupDaemon.pid,
        });
      },
      () => {
        process.send('timeout waiting for children to send a message');
      }
    );
  }
}

// For debugging purposes when running standalone from a terminal.
console.log(
  `parent: ${process.pid}, agent: ${agent.pid}, agentCleanupDaemon: ${agentCleanupDaemon.pid}`
);

// Needs to be bigger than the Jest test timeout (5s by default).
await setTimeout(10000);
