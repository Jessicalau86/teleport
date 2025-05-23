/*
Copyright 2017-2018 Gravitational, Inc.

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

package native

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/lib/utils"
)

func TestMain(m *testing.M) {
	utils.InitLoggerForTests()
	os.Exit(m.Run())
}

// TestPrecomputeMode verifies that package enters precompute mode when
// PrecomputeKeys is called.
func TestPrecomputeMode(t *testing.T) {
	t.Parallel()

	PrecomputeKeys()

	select {
	case <-precomputedKeys:
	case <-time.After(time.Second * 10):
		t.Fatal("Key precompute routine failed to start.")
	}
}

// TestGenerateRSAPKSC1Keypair tests that GeneratePrivateKey generates
// a valid PKCS1 rsa key.
func TestGeneratePKSC1RSAKey(t *testing.T) {
	t.Parallel()

	priv, err := GeneratePrivateKey()
	require.NoError(t, err)

	block, rest := pem.Decode(priv.PrivateKeyPEM())
	require.NoError(t, err)
	require.Empty(t, rest)

	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	require.NoError(t, err)
}

func TestGenerateEICEKey_when_boringbinary(t *testing.T) {
	if !IsBoringBinary() {
		t.Skip()
	}

	publicKey, privateKey, err := GenerateEICEKey()
	require.NoError(t, err)

	// We expect an RSA Key because boringcrypto doesn't yet support generating ED25519 keys.
	require.IsType(t, rsa.PublicKey{}, publicKey)
	require.IsType(t, rsa.PrivateKey{}, privateKey)
}

func TestGenerateEICEKey(t *testing.T) {
	if IsBoringBinary() {
		t.Skip()
	}

	publicKey, privateKey, err := GenerateEICEKey()
	require.NoError(t, err)

	// We expect an ED25519 key
	require.IsType(t, ed25519.PublicKey{}, publicKey)
	require.IsType(t, ed25519.PrivateKey{}, privateKey)
}
