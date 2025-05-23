//go:build !windows
// +build !windows

/*
Copyright 2018 Gravitational, Inc.

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

package utils

import (
	"os"
	"syscall"
)

// On non-windows we just lock the target file itself.
func getPlatformLockFilePath(path string) string {
	return path
}

func getHardLinkCount(fi os.FileInfo) (uint64, bool) {
	if statT, ok := fi.Sys().(*syscall.Stat_t); ok {
		// we must do a cast here because this will be uint16 on OSX
		//nolint:unconvert // the cast is only necessary for macOS
		return uint64(statT.Nlink), true
	} else {
		return 0, false
	}
}
