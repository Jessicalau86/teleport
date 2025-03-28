# Make targets:
#
#  all    : builds all binaries in development mode, without web assets (default)
#  full   : builds all binaries for PRODUCTION use
#  release: prepares a release tarball
#  clean  : removes all build artifacts
#  test   : runs tests

# To update the Teleport version, update VERSION variable:
# Naming convention:
#   Stable releases:   "1.0.0"
#   Pre-releases:      "1.0.0-alpha.1", "1.0.0-beta.2", "1.0.0-rc.3"
#   Master/dev branch: "1.0.0-dev"
VERSION=15.0.0-dev

DOCKER_IMAGE ?= teleport

GOPATH ?= $(shell go env GOPATH)

# This directory will be the real path of the directory of the first Makefile in the list.
MAKE_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# If set to 1, webassets are not built.
WEBASSETS_SKIP_BUILD ?= 0

# These are standard autotools variables, don't change them please
ifneq ("$(wildcard /bin/bash)","")
SHELL := /bin/bash -o pipefail
endif
BUILDDIR ?= build
BINDIR ?= /usr/local/bin
DATADIR ?= /usr/local/share/teleport
ADDFLAGS ?=
PWD ?= `pwd`
GIT ?= git
TELEPORT_DEBUG ?= false
GITTAG=v$(VERSION)
CGOFLAG ?= CGO_ENABLED=1

# RELEASE_DIR is where the release artifacts (tarballs, pacakges, etc) are put. It
# should be an absolute directory as it is used by e/Makefile too, from the e/ directory.
RELEASE_DIR := $(CURDIR)/$(BUILDDIR)/artifacts

# When TELEPORT_DEBUG is true, set flags to produce
# debugger-friendly builds.
ifeq ("$(TELEPORT_DEBUG)","true")
BUILDFLAGS ?= $(ADDFLAGS) -gcflags=all="-N -l"
else
BUILDFLAGS ?= $(ADDFLAGS) -ldflags '-w -s' -trimpath
endif

GO_ENV_OS := $(shell go env GOOS)
OS ?= $(GO_ENV_OS)

GO_ENV_ARCH := $(shell go env GOARCH)
ARCH ?= $(GO_ENV_ARCH)

FIPS ?=
RELEASE = teleport-$(GITTAG)-$(OS)-$(ARCH)-bin

# Include common makefile shared between OSS and Ent.
include common.mk

# FIPS support must be requested at build time.
FIPS_MESSAGE := without-FIPS-support
ifneq ("$(FIPS)","")
FIPS_TAG := fips
FIPS_MESSAGE := with-FIPS-support
RELEASE = teleport-$(GITTAG)-$(OS)-$(ARCH)-fips-bin
endif

# PAM support will only be built into Teleport if headers exist at build time.
PAM_MESSAGE := without-PAM-support
ifneq ("$(wildcard /usr/include/security/pam_appl.h)","")
PAM_TAG := pam
PAM_MESSAGE := with-PAM-support
else
# PAM headers for Darwin live under /usr/local/include/security instead, as SIP
# prevents us from modifying/creating /usr/include/security on newer versions of MacOS
ifneq ("$(wildcard /usr/local/include/security/pam_appl.h)","")
PAM_TAG := pam
PAM_MESSAGE := with-PAM-support
endif
endif

# darwin universal (Intel + Apple Silicon combined) binary support
RELEASE_darwin_arm64 = $(RELEASE_DIR)/teleport-$(GITTAG)-darwin-arm64-bin.tar.gz
RELEASE_darwin_amd64 = $(RELEASE_DIR)/teleport-$(GITTAG)-darwin-amd64-bin.tar.gz
BUILDDIR_arm64 = $(BUILDDIR)/arm64
BUILDDIR_amd64 = $(BUILDDIR)/amd64
# TARBINS is the path of the binaries in the release tarballs
TARBINS = $(addprefix teleport/,$(BINS))

# Check if rust and cargo are installed before compiling
CHECK_CARGO := $(shell cargo --version 2>/dev/null)
CHECK_RUST := $(shell rustc --version 2>/dev/null)

RUST_TARGET_ARCH ?= $(CARGO_TARGET_$(OS)_$(ARCH))

CARGO_TARGET_darwin_amd64 := x86_64-apple-darwin
CARGO_TARGET_darwin_arm64 := aarch64-apple-darwin
CARGO_TARGET_linux_arm64 := aarch64-unknown-linux-gnu
CARGO_TARGET_linux_amd64 := x86_64-unknown-linux-gnu

CARGO_TARGET := --target=${CARGO_TARGET_${OS}_${ARCH}}

# If set to 1, Windows RDP client is not built.
RDPCLIENT_SKIP_BUILD ?= 0

# Enable Windows RDP client build?
with_rdpclient := no
RDPCLIENT_MESSAGE := without-Windows-RDP-client

ifeq ($(RDPCLIENT_SKIP_BUILD),0)
ifneq ($(CHECK_RUST),)
ifneq ($(CHECK_CARGO),)

# Do not build RDP client on ARM or 386.
ifneq ("$(ARCH)","arm")
ifneq ("$(ARCH)","386")
with_rdpclient := yes
RDPCLIENT_MESSAGE := with-Windows-RDP-client
RDPCLIENT_TAG := desktop_access_rdp
endif
endif

endif
endif
endif

# Set C_ARCH for building libfido2 and dependencies. ARCH is the Go
# architecture which uses different names for architectures than C
# uses. Export it for the build.assets/build-fido2-macos.sh script.
C_ARCH_amd64 = x86_64
C_ARCH = $(or $(C_ARCH_$(ARCH)),$(ARCH))
export C_ARCH

# Enable libfido2 for testing?
# Eagerly enable if we detect the package, we want to test as much as possible.
ifeq ("$(shell pkg-config libfido2 2>/dev/null; echo $$?)", "0")
LIBFIDO2_TEST_TAG := libfido2
endif

# Build tsh against libfido2?
# FIDO2=yes and FIDO2=static enable static libfido2 builds.
# FIDO2=dynamic enables dynamic libfido2 builds.
LIBFIDO2_MESSAGE := without-libfido2
ifneq (, $(filter $(FIDO2), yes static))
LIBFIDO2_MESSAGE := with-libfido2
LIBFIDO2_BUILD_TAG := libfido2 libfido2static
else ifeq ("$(FIDO2)", "dynamic")
LIBFIDO2_MESSAGE := with-libfido2
LIBFIDO2_BUILD_TAG := libfido2
endif

# Enable Touch ID builds?
# Only build if TOUCHID=yes to avoid issues when cross-compiling to 'darwin'
# from other systems.
TOUCHID_MESSAGE := without-Touch-ID
ifeq ("$(TOUCHID)", "yes")
TOUCHID_MESSAGE := with-Touch-ID
TOUCHID_TAG := touchid
endif

# Enable PIV for testing?
# Eagerly enable if we detect the dynamic libpcsclite library, we want to test as much as possible.
ifeq ("$(shell pkg-config libpcsclite 2>/dev/null; echo $$?)", "0")
# This test tag should not be used for builds/releases, only tests.
PIV_TEST_TAG := piv
endif

# Build teleport/api with PIV? This requires the libpcsclite library for linux.
#
# PIV=yes and PIV=static enable static piv builds. This is used by the build
# process to link a static library of libpcsclite for piv-go to connect to.
#
# PIV=dynamic enables dynamic piv builds. This can be used for local
# builds and runs utilizing a dynamic libpcsclite library - `apt get install libpcsclite-dev`
PIV_MESSAGE := without-PIV-support
ifneq (, $(filter $(PIV), yes static dynamic))
PIV_MESSAGE := with-PIV-support
PIV_BUILD_TAG := piv
ifneq ("$(PIV)", "dynamic")
# Link static pcsc libary. By default, piv-go will look for the dynamic library.
# https://github.com/go-piv/piv-go/blob/master/piv/pcsc_unix.go#L23
STATIC_LIBS += -lpcsclite
STATIC_LIBS_TSH += -lpcsclite
endif
endif

# Reproducible builds are only available on select targets, and only when OS=linux.
REPRODUCIBLE ?=
ifneq ("$(OS)","linux")
REPRODUCIBLE = no
endif

# On Windows only build tsh. On all other platforms build teleport, tctl,
# and tsh.
BINS_default = teleport tctl tsh tbot
BINS_windows = tsh
BINS = $(or $(BINS_$(OS)),$(BINS_default))
BINARIES = $(addprefix $(BUILDDIR)/,$(BINS))

# Joins elements of the list in arg 2 with the given separator.
#   1. Element separator.
#   2. The list.
EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
join-with = $(subst $(SPACE),$1,$(strip $2))

# Separate TAG messages into comma-separated WITH and WITHOUT lists for readability.

COMMA := ,
MESSAGES := $(PAM_MESSAGE) $(FIPS_MESSAGE) $(BPF_MESSAGE) $(RDPCLIENT_MESSAGE) $(LIBFIDO2_MESSAGE) $(TOUCHID_MESSAGE) $(PIV_MESSAGE)
WITH := $(subst -," ",$(call join-with,$(COMMA) ,$(subst with-,,$(filter with-%,$(MESSAGES)))))
WITHOUT := $(subst -," ",$(call join-with,$(COMMA) ,$(subst without-,,$(filter without-%,$(MESSAGES)))))
RELEASE_MESSAGE := "Building with GOOS=$(OS) GOARCH=$(ARCH) REPRODUCIBLE=$(REPRODUCIBLE) and with $(WITH) and without $(WITHOUT)."

# On platforms that support reproducible builds, ensure the archive is created in a reproducible manner.
TAR_FLAGS ?=
ifeq ("$(REPRODUCIBLE)","yes")
TAR_FLAGS = --sort=name --owner=root:0 --group=root:0 --mtime='UTC 2015-03-02' --format=gnu
endif

VERSRC = version.go gitref.go api/version.go

KUBECONFIG ?=
TEST_KUBE ?=
export

TEST_LOG_DIR = ${abspath ./test-logs}

# Set CGOFLAG and BUILDFLAGS as needed for the OS/ARCH.
ifeq ("$(OS)","linux")
# True if $ARCH == amd64 || $ARCH == arm64
ifeq ("$(ARCH)","arm64")
	ifeq ($(IS_NATIVE_BUILD),"no")
		CGOFLAG += CC=aarch64-linux-gnu-gcc
	endif
else ifeq ("$(ARCH)","arm")
CGOFLAG = CGO_ENABLED=1

# ARM builds need to specify the correct C compiler
ifeq ($(IS_NATIVE_BUILD),"no")
CC=arm-linux-gnueabihf-gcc
endif

# Add -debugtramp=2 to work around 24 bit CALL/JMP instruction offset.
BUILDFLAGS = $(ADDFLAGS) -ldflags '-w -s -debugtramp=2' -trimpath
endif
endif # OS == linux

# Windows requires extra parameters to cross-compile with CGO.
ifeq ("$(OS)","windows")
ARCH ?= amd64
ifneq ("$(ARCH)","amd64")
$(error "Building for windows requires ARCH=amd64")
endif
CGOFLAG = CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++
BUILDFLAGS = $(ADDFLAGS) -ldflags '-w -s' -trimpath -buildmode=exe
endif

CGOFLAG_TSH ?= $(CGOFLAG)

# Map ARCH into the architecture flag for electron-builder if they
# are different to the Go $(ARCH) we use as an input.
ELECTRON_BUILDER_ARCH_amd64 = x64
ELECTRON_BUILDER_ARCH = $(or $(ELECTRON_BUILDER_ARCH_$(ARCH)),$(ARCH))

#
# 'make all' builds all 4 executables and places them in the current directory.
#
# NOTE: Works the same as `make`. Left for legacy reasons.
.PHONY: all
all: version
	@echo "---> Building OSS binaries."
	$(MAKE) $(BINARIES)

#
# make binaries builds all binaries defined in the BINARIES environment variable
#
.PHONY: binaries
binaries:
	$(MAKE) $(BINARIES)

# By making these 3 targets below (tsh, tctl and teleport) PHONY we are solving
# several problems:
# * Build will rely on go build internal caching https://golang.org/doc/go1.10 at all times
# * Manual change detection was broken on a large dependency tree
# If you are considering changing this behavior, please consult with dev team first
.PHONY: $(BUILDDIR)/tctl
$(BUILDDIR)/tctl:
	GOOS=$(OS) GOARCH=$(ARCH) $(CGOFLAG) go build -tags "$(PAM_TAG) $(FIPS_TAG) $(PIV_BUILD_TAG)" -o $(BUILDDIR)/tctl $(BUILDFLAGS) ./tool/tctl

.PHONY: $(BUILDDIR)/teleport
$(BUILDDIR)/teleport: ensure-webassets bpf-bytecode rdpclient
	GOOS=$(OS) GOARCH=$(ARCH) $(CGOFLAG) go build -tags "webassets_embed $(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(WEBASSETS_TAG) $(RDPCLIENT_TAG) $(PIV_BUILD_TAG)" -o $(BUILDDIR)/teleport $(BUILDFLAGS) ./tool/teleport

# NOTE: Any changes to the `tsh` build here must be copied to `windows.go` in Dronegen until
# 		we can use this Makefile for native Windows builds.
.PHONY: $(BUILDDIR)/tsh
$(BUILDDIR)/tsh:
	GOOS=$(OS) GOARCH=$(ARCH) $(CGOFLAG_TSH) go build -tags "$(FIPS_TAG) $(LIBFIDO2_BUILD_TAG) $(TOUCHID_TAG) $(PIV_BUILD_TAG)" -o $(BUILDDIR)/tsh $(BUILDFLAGS) ./tool/tsh

.PHONY: $(BUILDDIR)/tbot
$(BUILDDIR)/tbot:
	GOOS=$(OS) GOARCH=$(ARCH) $(CGOFLAG) go build -tags "$(FIPS_TAG)" -o $(BUILDDIR)/tbot $(BUILDFLAGS) ./tool/tbot

#
# BPF support (IF ENABLED)
# Requires a recent version of clang and libbpf installed.
#
ifeq ("$(with_bpf)","yes")
$(ER_BPF_BUILDDIR):
	mkdir -p $(ER_BPF_BUILDDIR)

$(RS_BPF_BUILDDIR):
	mkdir -p $(RS_BPF_BUILDDIR)

# Build BPF code
$(ER_BPF_BUILDDIR)/%.bpf.o: bpf/enhancedrecording/%.bpf.c $(wildcard bpf/*.h) | $(ER_BPF_BUILDDIR)
	$(CLANG) -g -O2 -target bpf -D__TARGET_ARCH_$(KERNEL_ARCH) -I/usr/libbpf-${LIBBPF_VER}/include $(INCLUDES) $(CLANG_BPF_SYS_INCLUDES) -c $(filter %.c,$^) -o $@
	$(LLVM_STRIP) -g $@ # strip useless DWARF info

# Build BPF code
$(RS_BPF_BUILDDIR)/%.bpf.o: bpf/restrictedsession/%.bpf.c $(wildcard bpf/*.h) | $(RS_BPF_BUILDDIR)
	$(CLANG) -g -O2 -target bpf -D__TARGET_ARCH_$(KERNEL_ARCH) -I/usr/libbpf-${LIBBPF_VER}/include $(INCLUDES) $(CLANG_BPF_SYS_INCLUDES) -c $(filter %.c,$^) -o $@
	$(LLVM_STRIP) -g $@ # strip useless DWARF info

.PHONY: bpf-rs-bytecode
bpf-rs-bytecode: $(RS_BPF_BUILDDIR)/restricted.bpf.o

.PHONY: bpf-er-bytecode
bpf-er-bytecode: $(ER_BPF_BUILDDIR)/command.bpf.o $(ER_BPF_BUILDDIR)/disk.bpf.o $(ER_BPF_BUILDDIR)/network.bpf.o $(ER_BPF_BUILDDIR)/counter_test.bpf.o

.PHONY: bpf-bytecode
bpf-bytecode: bpf-er-bytecode bpf-rs-bytecode

# Generate vmlinux.h based on the installed kernel
.PHONY: update-vmlinux-h
update-vmlinux-h:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c >bpf/vmlinux.h

else
.PHONY: bpf-bytecode
bpf-bytecode:
endif

ifeq ("$(with_rdpclient)", "yes")
.PHONY: rdpclient
rdpclient:
ifneq ("$(FIPS)","")
	cargo build -p rdp-client --features=fips --release $(CARGO_TARGET)
else
	cargo build -p rdp-client --release $(CARGO_TARGET)
endif
else
.PHONY: rdpclient
rdpclient:
endif

# Build libfido2 and dependencies for MacOS. Uses exported C_ARCH variable defined earlier.
.PHONY: build-fido2
build-fido2:
	./build.assets/build-fido2-macos.sh build

.PHONY: print-fido2-pkg-path
print-fido2-pkg-path:
	@./build.assets/build-fido2-macos.sh pkg_config_path

#
# make full - Builds Teleport binaries with the built-in web assets and
# places them into $(BUILDDIR). On Windows, this target is skipped because
# only tsh is built.
#
.PHONY:full
full: WEBASSETS_SKIP_BUILD = 0
full: ensure-webassets
ifneq ("$(OS)", "windows")
	export WEBASSETS_SKIP_BUILD=0
	$(MAKE) all
endif

#
# make full-ent - Builds Teleport enterprise binaries
#
.PHONY:full-ent
full-ent: ensure-webassets-e
ifneq ("$(OS)", "windows")
	@if [ -f e/Makefile ]; then $(MAKE) -C e full; fi
endif

#
# make clean - Removes all build artifacts.
#
.PHONY: clean
clean: clean-ui clean-build

.PHONY: clean-build
clean-build:
	@echo "---> Cleaning up OSS build artifacts."
	rm -rf $(BUILDDIR)
# Check if the variable is set to prevent calling remove on the root directory.
ifneq ($(ER_BPF_BUILDDIR),)
	rm -f $(ER_BPF_BUILDDIR)/*.o
endif
ifneq ($(RS_BPF_BUILDDIR),)
	rm -f $(RS_BPF_BUILDDIR)/*.o
endif
	-cargo clean
	-go clean -cache
	rm -f *.gz
	rm -f *.zip
	rm -f gitref.go
	rm -rf build.assets/tooling/bin

.PHONY: clean-ui
clean-ui:
	rm -rf webassets/*
	find . -type d -name node_modules -prune -exec rm -rf {} \;

#
# make release - Produces a binary release tarball.
#

# RELEASE_DIR is where release artifact files are put, such as tarballs, packages, etc.
$(RELEASE_DIR):
	mkdir $@

.PHONY:
export
release:
	@echo "---> OSS $(RELEASE_MESSAGE)"
ifeq ("$(OS)", "windows")
	$(MAKE) --no-print-directory release-windows
else ifeq ("$(OS)", "darwin")
	$(MAKE) --no-print-directory release-darwin
else
	$(MAKE) --no-print-directory release-unix
endif

# These are aliases used to make build commands uniform.
.PHONY: release-amd64
release-amd64:
	$(MAKE) release ARCH=amd64

.PHONY: release-386
release-386:
	$(MAKE) release ARCH=386

.PHONY: release-arm
release-arm:
	$(MAKE) release ARCH=arm

.PHONY: release-arm64
release-arm64:
	$(MAKE) release ARCH=arm64

#
# make build-archive - Packages the results of a build into a release tarball
#
.PHONY: build-archive
build-archive: | $(RELEASE_DIR)
	@echo "---> Creating OSS release archive."
	mkdir teleport
	cp -rf $(BINARIES) \
		examples \
		build.assets/install\
		README.md \
		CHANGELOG.md \
		teleport/
	echo $(GITTAG) > teleport/VERSION
	tar $(TAR_FLAGS) -c teleport | gzip -n > $(RELEASE).tar.gz
	cp $(RELEASE).tar.gz $(RELEASE_DIR)
	rm -rf teleport
	@echo "---> Created $(RELEASE).tar.gz."

#
# make release-unix - Produces binary release tarballs for both OSS and
# Enterprise editions, containing teleport, tctl, tbot and tsh.
#
.PHONY: release-unix
release-unix: clean full build-archive
	@if [ -f e/Makefile ]; then $(MAKE) -C e release; fi

# release-unix-preserving-webassets cleans just the build and not the UI
# allowing webassets to be built in a prior step before building the release.
.PHONY: release-unix-preserving-webassets
release-unix-preserving-webassets: clean-build full build-archive
	@if [ -f e/Makefile ]; then $(MAKE) -C e release; fi

include darwin-signing.mk

.PHONY: release-darwin-unsigned
release-darwin-unsigned: RELEASE:=$(RELEASE)-unsigned
release-darwin-unsigned: full build-archive

.PHONY: release-darwin
ifneq ($(ARCH),universal)
release-darwin: ABSOLUTE_BINARY_PATHS:=$(addprefix $(CURDIR)/,$(BINARIES))
release-darwin: release-darwin-unsigned
	$(NOTARIZE_BINARIES)
	$(MAKE) build-archive
	@if [ -f e/Makefile ]; then $(MAKE) -C e release; fi
else

# release-darwin for ARCH == universal does not build binaries, but instead
# combines previously-built binaries. For this, it depends on the ARM64 and
# AMD64 signed tarballs being built into $(RELEASE_DIR). The dependencies
# expressed here will not make that happen as this is typically done on CI
# where these two tarballs are built in separate pipelines, and copied in for
# the universal build.
#
# For local manual runs, create these tarballs with:
#   make ARCH=arm64 release-darwin
#   make ARCH=amd64 release-darwin
# Ensure you have the rust toolchains for these installed by running
#   make ARCH=arm64 rustup-install-target-toolchain
#   make ARCH=amd64 rustup-install-target-toolchain
release-darwin: $(RELEASE_darwin_arm64) $(RELEASE_darwin_amd64)
	mkdir -p $(BUILDDIR_arm64) $(BUILDDIR_amd64)
	tar -C $(BUILDDIR_arm64) -xzf $(RELEASE_darwin_arm64) --strip-components=1 $(TARBINS)
	tar -C $(BUILDDIR_amd64) -xzf $(RELEASE_darwin_amd64) --strip-components=1 $(TARBINS)
	lipo -create -output $(BUILDDIR)/teleport $(BUILDDIR_arm64)/teleport $(BUILDDIR_amd64)/teleport
	lipo -create -output $(BUILDDIR)/tctl $(BUILDDIR_arm64)/tctl $(BUILDDIR_amd64)/tctl
	lipo -create -output $(BUILDDIR)/tsh $(BUILDDIR_arm64)/tsh $(BUILDDIR_amd64)/tsh
	lipo -create -output $(BUILDDIR)/tbot $(BUILDDIR_arm64)/tbot $(BUILDDIR_amd64)/tbot
	$(MAKE) ARCH=universal build-archive
	@if [ -f e/Makefile ]; then $(MAKE) -C e release; fi
endif

#
# make release-windows-unsigned - Produces a binary release archive containing only tsh.
#
.PHONY: release-windows-unsigned
release-windows-unsigned: clean all
	@echo "---> Creating OSS release archive."
	mkdir teleport
	cp -rf $(BUILDDIR)/* \
		README.md \
		CHANGELOG.md \
		teleport/
	mv teleport/tsh teleport/tsh-unsigned.exe
	echo $(GITTAG) > teleport/VERSION
	zip -9 -y -r -q $(RELEASE)-unsigned.zip teleport/
	rm -rf teleport/
	@echo "---> Created $(RELEASE)-unsigned.zip."

#
# make release-windows - Produces an archive containing a signed release of
# tsh.exe
#
.PHONY: release-windows
release-windows: release-windows-unsigned
	@if [ ! -f "windows-signing-cert.pfx" ]; then \
		echo "windows-signing-cert.pfx is missing or invalid, cannot create signed archive."; \
		exit 1; \
	fi

	rm -rf teleport
	@echo "---> Extracting $(RELEASE)-unsigned.zip"
	unzip $(RELEASE)-unsigned.zip

	@echo "---> Signing Windows binary."
	@osslsigncode sign \
		-pkcs12 "windows-signing-cert.pfx" \
		-n "Teleport" \
		-i https://goteleport.com \
		-t http://timestamp.digicert.com \
		-h sha2 \
		-in teleport/tsh-unsigned.exe \
		-out teleport/tsh.exe; \
	success=$$?; \
	rm -f teleport/tsh-unsigned.exe; \
	if [ "$${success}" -ne 0 ]; then \
		echo "Failed to sign tsh.exe, aborting."; \
		exit 1; \
	fi

	zip -9 -y -r -q $(RELEASE).zip teleport/
	rm -rf teleport/
	@echo "---> Created $(RELEASE).zip."

#
# make release-connect produces a release package of Teleport Connect.
# It is used only for MacOS releases. Windows releases do not use this
# Makefile. Linux uses the `teleterm` target in build.assets/Makefile.
#
# Only export the CSC_NAME (developer key ID) when the recipe is run, so
# that we do not shell out and run the `security` command if not necessary.
#
# Either CONNECT_TSH_BIN_PATH or CONNECT_TSH_APP_PATH environment variable
# should be defined for the `yarn package-term` command to succeed. CI sets
# this appropriately depending on whether a push build is running, or a
# proper release (a proper release needs the APP_PATH as that points to
# the complete signed package). See web/packages/teleterm/README.md for
# details.

.PHONY: release-connect
release-connect: | $(RELEASE_DIR)
	$(eval export CSC_NAME)
	yarn install --frozen-lockfile
	yarn build-term
	yarn package-term -c.extraMetadata.version=$(VERSION) --$(ELECTRON_BUILDER_ARCH)
	# Only copy proper builds with tsh.app to $(RELEASE_DIR)
	# Drop -universal "arch" from dmg name when copying to $(RELEASE_DIR)
	if [ -n "$$CONNECT_TSH_APP_PATH" ]; then \
		TARGET_NAME="Teleport Connect-$(VERSION)-$(ARCH).dmg"; \
		if [ "$(ARCH)" = 'universal' ]; then \
			TARGET_NAME="$${TARGET_NAME/-universal/}"; \
		fi; \
		cp web/packages/teleterm/build/release/"Teleport Connect-$(VERSION)-$(ELECTRON_BUILDER_ARCH).dmg" "$(RELEASE_DIR)/$${TARGET_NAME}"; \
	fi

#
# Remove trailing whitespace in all markdown files under docs/.
#
# Note: this runs in a busybox container to avoid incompatibilities between
# linux and macos CLI tools.
#
.PHONY:docs-fix-whitespace
docs-fix-whitespace:
	docker run --rm -v $(PWD):/teleport busybox \
		find /teleport/docs/ -type f -name '*.md' -exec sed -E -i 's/\s+$$//g' '{}' \;

#
# Test docs for trailing whitespace and broken links
#
.PHONY:docs-test
docs-test: docs-test-whitespace

#
# Check for trailing whitespace in all markdown files under docs/
#
.PHONY:docs-test-whitespace
docs-test-whitespace:
	if find docs/ -type f -name '*.md' | xargs grep -E '\s+$$'; then \
		echo "trailing whitespace found in docs/ (see above)"; \
		echo "run 'make docs-fix-whitespace' to fix it"; \
		exit 1; \
	fi

#
# Builds some tooling for filtering and displaying test progress/output/etc
#
# Deprecated: Use gotestsum instead.
TOOLINGDIR := ${abspath ./build.assets/tooling}
RENDER_TESTS := $(TOOLINGDIR)/bin/render-tests
$(RENDER_TESTS): $(wildcard $(TOOLINGDIR)/cmd/render-tests/*.go)
	cd $(TOOLINGDIR) && go build -o "$@" ./cmd/render-tests

#
# Install gotestsum to parse test output.
#
.PHONY: ensure-gotestsum
ensure-gotestsum:
# Install gotestsum if it's not already installed
 ifeq (, $(shell command -v gotestsum))
	go install gotest.tools/gotestsum@latest
endif

DIFF_TEST := $(TOOLINGDIR)/bin/difftest
$(DIFF_TEST): $(wildcard $(TOOLINGDIR)/cmd/difftest/*.go)
	cd $(TOOLINGDIR) && go build -o "$@" ./cmd/difftest

RERUN := $(TOOLINGDIR)/bin/rerun
$(RERUN): $(wildcard $(TOOLINGDIR)/cmd/rerun/*.go)
	cd $(TOOLINGDIR) && go build -o "$@" ./cmd/rerun

.PHONY: tooling
tooling: ensure-gotestsum $(DIFF_TEST)

#
# Runs all Go/shell tests, called by CI/CD.
#
.PHONY: test
test: test-helm test-sh test-api test-go test-rust test-operator

$(TEST_LOG_DIR):
	mkdir $(TEST_LOG_DIR)

.PHONY: helmunit/installed
helmunit/installed:
	@if ! helm unittest -h >/dev/null; then \
		echo 'Helm unittest plugin is required to test Helm charts. Run `helm plugin install https://github.com/quintush/helm-unittest --version 0.2.11` to install it'; \
		exit 1; \
	fi

# The CI environment is responsible for setting HELM_PLUGINS to a directory where
# quintish/helm-unittest is installed.
#
# Github Actions build uses /workspace as homedir and Helm can't pick up plugins by default there,
# so override the plugin location via environemnt variable when running in CI. Github Actions provide CI=true
# environment variable.
.PHONY: test-helm
test-helm: helmunit/installed
	helm unittest -3 examples/chart/teleport-cluster
	helm unittest -3 examples/chart/teleport-kube-agent

.PHONY: test-helm-update-snapshots
test-helm-update-snapshots: helmunit/installed
	helm unittest -3 -u examples/chart/teleport-cluster
	helm unittest -3 -u examples/chart/teleport-kube-agent

#
# Runs all Go tests except integration, called by CI/CD.
#
.PHONY: test-go
test-go: test-go-prepare test-go-unit test-go-touch-id test-go-tsh test-go-chaos

#
# Runs a test to ensure no environment variable leak into build binaries.
# This is typically done as part of the bloat test in CI, but this
# target exists for local testing.
#
.PHONY: test-env-leakage
test-env-leakage:
	$(eval export BUILD_SECRET=FAKE_SECRET)
	$(MAKE) full
	failed=0; \
	for binary in $(BINARIES); do \
		if strings $$binary | grep -q 'FAKE_SECRET'; then \
			echo "Error: $$binary contains FAKE_SECRET"; \
			failed=1; \
		fi; \
	done; \
	if [ $$failed -eq 1 ]; then \
		echo "Environment leak failure"; \
		exit 1; \
	else \
		echo "No environment leak, PASS"; \
	fi

# Runs test prepare steps
.PHONY: test-go-prepare
test-go-prepare: ensure-webassets bpf-bytecode rdpclient $(TEST_LOG_DIR) ensure-gotestsum $(VERSRC)

# Runs base unit tests
.PHONY: test-go-unit
test-go-unit: FLAGS ?= -race -shuffle on
test-go-unit: SUBJECT ?= $(shell go list ./... | grep -vE 'teleport/(e2e|integration|tool/tsh|integrations/operator|integrations/access|integrations/lib)')
test-go-unit:
	$(CGOFLAG) go test -cover -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(RDPCLIENT_TAG) $(LIBFIDO2_TEST_TAG) $(TOUCHID_TAG) $(PIV_TEST_TAG)" $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/unit.json \
		| gotestsum --raw-command -- cat

# Runs tbot unit tests
.PHONY: test-go-unit-tbot
test-go-unit-tbot: FLAGS ?= -race -shuffle on
test-go-unit-tbot:
	$(CGOFLAG) go test -cover -json $(FLAGS) $(ADDFLAGS) ./tool/tbot/... ./lib/tbot/... \
		| tee $(TEST_LOG_DIR)/unit.json \
		| gotestsum --raw-command -- cat

# Make sure untagged touchid code build/tests.
.PHONY: test-go-touch-id
test-go-touch-id: FLAGS ?= -race -shuffle on
test-go-touch-id: SUBJECT ?= ./lib/auth/touchid/...
test-go-touch-id:
ifneq ("$(TOUCHID_TAG)", "")
	$(CGOFLAG) go test -cover -json $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/unit.json \
		| gotestsum --raw-command -- cat
endif

# Runs ci tsh tests
.PHONY: test-go-tsh
test-go-tsh: FLAGS ?= -race -shuffle on
test-go-tsh: SUBJECT ?= github.com/gravitational/teleport/tool/tsh/...
test-go-tsh:
	$(CGOFLAG_TSH) go test -cover -json -tags "$(PAM_TAG) $(FIPS_TAG) $(LIBFIDO2_TEST_TAG) $(TOUCHID_TAG) $(PIV_TEST_TAG)" $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/unit.json \
		| gotestsum --raw-command -- cat

# Chaos tests have high concurrency, run without race detector and have TestChaos prefix.
.PHONY: test-go-chaos
test-go-chaos: CHAOS_FOLDERS = $(shell find . -type f -name '*chaos*.go' | xargs dirname | uniq)
test-go-chaos:
	$(CGOFLAG) go test -cover -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(RDPCLIENT_TAG)" -test.run=TestChaos $(CHAOS_FOLDERS) \
		| tee $(TEST_LOG_DIR)/chaos.json \
		| gotestsum --raw-command -- cat

#
# Runs all Go tests except integration, end-to-end, and chaos, called by CI/CD.
#
UNIT_ROOT_REGEX := ^TestRoot
.PHONY: test-go-root
test-go-root: ensure-webassets bpf-bytecode rdpclient $(TEST_LOG_DIR) ensure-gotestsum
test-go-root: FLAGS ?= -race -shuffle on
test-go-root: PACKAGES = $(shell go list $(ADDFLAGS) ./... | grep -v -e e2e -e integration -e integrations/operator)
test-go-root: $(VERSRC)
	$(CGOFLAG) go test -json -run "$(UNIT_ROOT_REGEX)" -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(RDPCLIENT_TAG)" $(PACKAGES) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/unit-root.json \
		| gotestsum --raw-command -- cat

#
# Runs Go tests on the api module. These have to be run separately as the package name is different.
#
.PHONY: test-api
test-api: $(VERSRC) $(TEST_LOG_DIR) ensure-gotestsum
test-api: FLAGS ?= -race -shuffle on
test-api: SUBJECT ?= $(shell cd api && go list ./...)
test-api:
	cd api && $(CGOFLAG) go test -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG)" $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/api.json \
		| gotestsum --raw-command -- cat

#
# Runs Teleport Operator tests.
# We have to run them using the makefile to ensure the installation of the k8s test tools (envtest)
#
.PHONY: test-operator
test-operator:
	make -C integrations/operator test
#
# Runs Go tests on the integrations/kube-agent-updater module. These have to be run separately as the package name is different.
#
.PHONY: test-kube-agent-updater
test-kube-agent-updater: $(VERSRC) $(TEST_LOG_DIR) ensure-gotestsum
test-kube-agent-updater: FLAGS ?= -race -shuffle on
test-kube-agent-updater: SUBJECT ?= $(shell cd integrations/kube-agent-updater && go list ./...)
test-kube-agent-updater:
	cd integrations/kube-agent-updater && $(CGOFLAG) go test -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG)" $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/kube-agent-updater.json \
		| gotestsum --raw-command --format=testname -- cat

.PHONY: test-access-integrations
test-access-integrations:
	make -C integrations test-access

.PHONY: test-integrations-lib
test-integrations-lib:
	make -C integrations test-lib

#
# Runs Go tests on the examples/teleport-usage module. These have to be run separately as the package name is different.
#
.PHONY: test-teleport-usage
test-teleport-usage: $(VERSRC) $(TEST_LOG_DIR) ensure-gotestsum
test-teleport-usage: FLAGS ?= -race -shuffle on
test-teleport-usage: SUBJECT ?= $(shell cd examples/teleport-usage && go list ./...)
test-teleport-usage:
	cd examples/teleport-usage && $(CGOFLAG) go test -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG)" $(PACKAGES) $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| tee $(TEST_LOG_DIR)/teleport-usage.json \
		| gotestsum --raw-command -- cat

#
# Flaky test detection. Usually run from CI nightly, overriding these default parameters
# This runs the same tests as test-go-unit but repeatedly to try to detect flaky tests.
#
# TODO(jakule): Migrate to gotestsum
.PHONY: test-go-flaky
FLAKY_RUNS ?= 3
FLAKY_TIMEOUT ?= 1h
FLAKY_TOP_N ?= 20
FLAKY_SUMMARY_FILE ?= /tmp/flaky-report.txt
test-go-flaky: FLAGS ?= -race -shuffle on
test-go-flaky: SUBJECT ?= $(shell go list ./... | grep -v -e e2e -e integration -e tool/tsh -e integrations/operator -e integrations/access -e integrations/lib )
test-go-flaky: GO_BUILD_TAGS ?= $(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(RDPCLIENT_TAG) $(TOUCHID_TAG) $(PIV_TEST_TAG)
test-go-flaky: RENDER_FLAGS ?= -report-by flakiness -summary-file $(FLAKY_SUMMARY_FILE) -top $(FLAKY_TOP_N)
test-go-flaky: test-go-prepare $(RENDER_TESTS) $(RERUN)
	$(CGOFLAG) $(RERUN) -n $(FLAKY_RUNS) -t $(FLAKY_TIMEOUT) \
		go test -count=1 -cover -json -tags "$(GO_BUILD_TAGS)" $(SUBJECT) $(FLAGS) $(ADDFLAGS) \
		| $(RENDER_TESTS) $(RENDER_FLAGS)

#
# Runs cargo test on our Rust modules.
# (a no-op if cargo and rustc are not installed)
#
ifneq ($(CHECK_RUST),)
ifneq ($(CHECK_CARGO),)
.PHONY: test-rust
test-rust:
	cargo test
else
.PHONY: test-rust
test-rust:
endif
endif

# Find and run all shell script unit tests (using https://github.com/bats-core/bats-core)
.PHONY: test-sh
test-sh:
	@if ! type bats 2>&1 >/dev/null; then \
		echo "Not running 'test-sh' target as 'bats' is not installed."; \
		if [ "$${DRONE}" = "true" ]; then echo "This is a failure when running in CI." && exit 1; fi; \
		exit 0; \
	fi; \
	find . -iname "*.bats" -exec dirname {} \; | uniq | xargs -t -L1 bats $(BATSFLAGS)


.PHONY: test-e2e
test-e2e:
	make -C e2e test

.PHONY: run-etcd
run-etcd:
	docker build -f .github/services/Dockerfile.etcd -t etcdbox --build-arg=ETCD_VERSION=3.5.9 .
	docker run -it --rm -p'2379:2379' etcdbox

#
# Integration tests. Need a TTY to work.
# Any tests which need to run as root must be skipped during regular integration testing.
#
.PHONY: integration
integration: FLAGS ?= -v -race
integration: PACKAGES = $(shell go list ./... | grep 'integration\([^s]\|$$\)' | grep -v integrations/lib/testing/integration )
integration:  $(TEST_LOG_DIR) ensure-gotestsum
	@echo KUBECONFIG is: $(KUBECONFIG), TEST_KUBE: $(TEST_KUBE)
	$(CGOFLAG) go test -timeout 30m -json -tags "$(PAM_TAG) $(FIPS_TAG) $(BPF_TAG) $(RDPCLIENT_TAG)" $(PACKAGES) $(FLAGS) \
		| tee $(TEST_LOG_DIR)/integration.json \
		| gotestsum --raw-command --format=testname -- cat

#
# Integration tests that run Kubernetes tests in order to complete successfully
# are run separately to all other integration tests.
#
INTEGRATION_KUBE_REGEX := TestKube.*
.PHONY: integration-kube
integration-kube: FLAGS ?= -v -race
integration-kube: PACKAGES = $(shell go list ./... | grep 'integration\([^s]\|$$\)')
integration-kube: $(TEST_LOG_DIR) ensure-gotestsum
	@echo KUBECONFIG is: $(KUBECONFIG), TEST_KUBE: $(TEST_KUBE)
	$(CGOFLAG) go test -json -run "$(INTEGRATION_KUBE_REGEX)" $(PACKAGES) $(FLAGS) \
		| tee $(TEST_LOG_DIR)/integration-kube.json \
		| gotestsum --raw-command --format=testname -- cat

#
# Integration tests which need to be run as root in order to complete successfully
# are run separately to all other integration tests. Need a TTY to work.
#
INTEGRATION_ROOT_REGEX := ^TestRoot
.PHONY: integration-root
integration-root: FLAGS ?= -v -race
integration-root: PACKAGES = $(shell go list ./... | grep 'integration\([^s]\|$$\)')
integration-root: $(TEST_LOG_DIR) ensure-gotestsum
	$(CGOFLAG) go test -json -run "$(INTEGRATION_ROOT_REGEX)" $(PACKAGES) $(FLAGS) \
		| tee $(TEST_LOG_DIR)/integration-root.json \
		| gotestsum --raw-command --format=testname -- cat


.PHONY: e2e-aws
e2e-aws: FLAGS ?= -v -race
e2e-aws: PACKAGES = $(shell go list ./... | grep 'e2e/aws')
e2e-aws: $(TEST_LOG_DIR) ensure-gotestsum
	@echo TEST_KUBE: $(TEST_KUBE) TEST_AWS_DB: $(TEST_AWS_DB)
	$(CGOFLAG) go test -json $(PACKAGES) $(FLAGS) $(ADDFLAGS)\
		| tee $(TEST_LOG_DIR)/e2e-aws.json \
		| gotestsum --raw-command --format=testname -- cat

#
# Lint the source code.
# By default lint scans the entire repo. Pass GO_LINT_FLAGS='--new' to only scan local
# changes (or last commit).
#
.PHONY: lint
lint: lint-api lint-go lint-kube-agent-updater lint-tools lint-protos lint-no-actions

#
# Lints everything but Go sources.
# Similar to lint.
#
.PHONY: lint-no-actions
lint-no-actions: lint-sh lint-helm lint-license lint-rust

.PHONY: lint-tools
lint-tools: lint-build-tooling lint-backport

#
# Runs the clippy linter on our rust modules
# (a no-op if cargo and rustc are not installed)
#
ifneq ($(CHECK_RUST),)
ifneq ($(CHECK_CARGO),)
.PHONY: lint-rust
lint-rust:
	cargo clippy --locked --all-targets -- -D warnings \
		&& cargo fmt -- --check
else
.PHONY: lint-rust
lint-rust:
endif
endif

.PHONY: lint-go
lint-go: GO_LINT_FLAGS ?=
lint-go:
	golangci-lint run -c .golangci.yml --build-tags='$(LIBFIDO2_TEST_TAG) $(TOUCHID_TAG) $(PIV_TEST_TAG)' $(GO_LINT_FLAGS)

.PHONY: fix-imports
fix-imports:
ifndef TELEPORT_DEVBOX
	$(MAKE) -C build.assets/ fix-imports
else
	$(MAKE) fix-imports/host
endif

.PHONY: fix-imports/host
fix-imports/host:
	@if ! type gci >/dev/null 2>&1; then\
		echo 'gci is not installed or is missing from PATH, consider installing it ("go install github.com/daixiang0/gci@latest") or use "make -C build.assets/ fix-imports"';\
		exit 1;\
	fi
	gci write -s standard -s default -s 'prefix(github.com/gravitational/teleport)' --skip-generated .

.PHONY: lint-build-tooling
lint-build-tooling: GO_LINT_FLAGS ?=
lint-build-tooling:
	cd build.assets/tooling && golangci-lint run -c ../../.golangci.yml $(GO_LINT_FLAGS)

.PHONY: lint-backport
lint-backport: GO_LINT_FLAGS ?=
lint-backport:
	cd assets/backport && golangci-lint run -c ../../.golangci.yml $(GO_LINT_FLAGS)

# api is no longer part of the teleport package, so golangci-lint skips it by default
.PHONY: lint-api
lint-api: GO_LINT_API_FLAGS ?=
lint-api:
	cd api && golangci-lint run -c ../.golangci.yml $(GO_LINT_API_FLAGS)

.PHONY: lint-kube-agent-updater
lint-kube-agent-updater: GO_LINT_API_FLAGS ?=
lint-kube-agent-updater:
	cd integrations/kube-agent-updater && golangci-lint run -c ../../.golangci.yml $(GO_LINT_API_FLAGS)

# TODO(awly): remove the `--exclude` flag after cleaning up existing scripts
.PHONY: lint-sh
lint-sh: SH_LINT_FLAGS ?=
lint-sh:
	find . -type f \( -name '*.sh' -or -name '*.sh.tmpl' \) -not -path "*/node_modules/*" | xargs \
		shellcheck \
		--exclude=SC2086 \
		--exclude=SC1091 \
		$(SH_LINT_FLAGS)

	# lint AWS AMI scripts
	# SC1091 prints errors when "source" directives are not followed
	find assets/aws/files/bin -type f | xargs \
		shellcheck \
		--exclude=SC2086 \
		--exclude=SC1091 \
		--exclude=SC2129 \
		$(SH_LINT_FLAGS)

# Lints all the Helm charts found in directories under examples/chart and exits on failure
# If there is a .lint directory inside, the chart gets linted once for each .yaml file in that directory
# We inherit yamllint's 'relaxed' configuration as it's more compatible with Helm output and will only error on
# show-stopping issues. Kubernetes' YAML parser is not particularly fussy.
# If errors are found, the file is printed with line numbers to aid in debugging.
.PHONY: lint-helm
lint-helm:
	@if ! type yamllint 2>&1 >/dev/null; then \
		echo "Not running 'lint-helm' target as 'yamllint' is not installed."; \
		if [ "$${DRONE}" = "true" ]; then echo "This is a failure when running in CI." && exit 1; fi; \
		exit 0; \
	fi; \
	for CHART in $$(find examples/chart -mindepth 1 -maxdepth 1 -type d); do \
		if [ -d $${CHART}/.lint ]; then \
			for VALUES in $${CHART}/.lint/*.yaml; do \
				export HELM_TEMP=$$(mktemp); \
				echo -n "Using values from '$${VALUES}': "; \
				yamllint -c examples/chart/.lint-config.yaml $${VALUES} || { cat -en $${VALUES}; exit 1; }; \
				helm lint --quiet --strict $${CHART} -f $${VALUES} || exit 1; \
				helm template test $${CHART} -f $${VALUES} 1>$${HELM_TEMP} || exit 1; \
				yamllint -c examples/chart/.lint-config.yaml $${HELM_TEMP} || { cat -en $${HELM_TEMP}; exit 1; }; \
			done \
		else \
			export HELM_TEMP=$$(mktemp); \
			helm lint --quiet --strict $${CHART} || exit 1; \
			helm template test $${CHART} 1>$${HELM_TEMP} || exit 1; \
			yamllint -c examples/chart/.lint-config.yaml $${HELM_TEMP} || { cat -en $${HELM_TEMP}; exit 1; }; \
		fi; \
	done

ADDLICENSE := $(GOPATH)/bin/addlicense
ADDLICENSE_ARGS := -c 'Gravitational, Inc' -l apache \
		-ignore '**/*.c' \
		-ignore '**/*.h' \
		-ignore '**/*.html' \
		-ignore '**/*.js' \
		-ignore '**/*.py' \
		-ignore '**/*.sh' \
		-ignore '**/*.tf' \
		-ignore '**/*.yaml' \
		-ignore '**/*.yml' \
		-ignore '**/*.sql' \
		-ignore '**/Dockerfile' \
		-ignore 'api/version.go' \
		-ignore 'docs/pages/includes/**/*.go' \
		-ignore 'e/**' \
		-ignore 'gen/**' \
		-ignore 'gitref.go' \
		-ignore 'lib/srv/desktop/rdp/rdpclient/target/**' \
		-ignore 'lib/web/build/**' \
		-ignore 'version.go' \
		-ignore 'webassets/**' \
		-ignore '**/node_modules/**' \
		-ignore 'web/packages/design/src/assets/icomoon/style.css' \
		-ignore 'ignoreme'

.PHONY: lint-license
lint-license: $(ADDLICENSE)
	$(ADDLICENSE) $(ADDLICENSE_ARGS) -check * 2>/dev/null

.PHONY: fix-license
fix-license: $(ADDLICENSE)
	$(ADDLICENSE) $(ADDLICENSE_ARGS) * 2>/dev/null

$(ADDLICENSE):
	cd && go install github.com/google/addlicense@v1.0.0

# This rule updates version files and Helm snapshots based on the Makefile
# VERSION variable.
#
# Used prior to a release by bumping VERSION in this Makefile and then
# running "make update-version".
.PHONY: update-version
update-version: version test-helm-update-snapshots

# This rule triggers re-generation of version files if Makefile changes.
.PHONY: version
version: $(VERSRC)

# This rule triggers re-generation of version files specified if Makefile changes.
$(VERSRC): Makefile
	VERSION=$(VERSION) $(MAKE) -f version.mk setver

# make tag - prints a tag to use with git for the current version
# 	To put a new release on Github:
# 		- bump VERSION variable
# 		- run make setver
# 		- commit changes to git
# 		- build binaries with 'make release'
# 		- run `make tag` and use its output to 'git tag' and 'git push --tags'
.PHONY: update-tag
update-tag: TAG_REMOTE ?= origin
update-tag:
	@test $(VERSION)
	cd build.assets/tooling && CGO_ENABLED=0 go run ./cmd/check -check valid -tag $(GITTAG)
	git tag $(GITTAG)
	git tag api/$(GITTAG)
	(cd e && git tag $(GITTAG) && git push origin $(GITTAG))
	git push $(TAG_REMOTE) $(GITTAG) && git push $(TAG_REMOTE) api/$(GITTAG)

.PHONY: test-package
test-package: remove-temp-files
	go test -v ./$(p)

.PHONY: test-grep-package
test-grep-package: remove-temp-files
	go test -v ./$(p) -check.f=$(e)

.PHONY: cover-package
cover-package: remove-temp-files
	go test -v ./$(p)  -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

.PHONY: profile
profile:
	go tool pprof http://localhost:6060/debug/pprof/profile

.PHONY: sloccount
sloccount:
	find . -o -name "*.go" -print0 | xargs -0 wc -l

.PHONY: remove-temp-files
remove-temp-files:
	find . -name flymake_* -delete

#
# print-go-version outputs Go version as a semver without "go" prefix
#
.PHONY: print-go-version
print-go-version:
	@$(MAKE) -C build.assets print-go-version | sed "s/go//"

# Dockerized build: useful for making Linux releases on OSX
.PHONY:docker
docker:
	make -C build.assets build

# Dockerized build: useful for making Linux binaries on macOS
.PHONY:docker-binaries
docker-binaries: clean
	make -C build.assets build-binaries PIV=$(PIV)

# Interactively enters a Docker container (which you can build and run Teleport inside of)
.PHONY:enter
enter:
	make -C build.assets enter

# Interactively enters a Docker container, as root (which you can build and run Teleport inside of)
.PHONY:enter-root
enter-root:
	make -C build.assets enter-root

# Interactively enters a Docker container (which you can build and run Teleport inside of).
# Similar to `enter`, but uses the centos7 container.
.PHONY:enter/centos7
enter/centos7:
	make -C build.assets enter/centos7

.PHONY:enter/grpcbox
enter/grpcbox:
	make -C build.assets enter/grpcbox

BUF := buf

# protos/all runs build, lint and format on all protos.
# Use `make grpc` to regenerate protos inside buildbox.
.PHONY: protos/all
protos/all: protos/build protos/lint protos/format

.PHONY: protos/build
protos/build: buf/installed
	$(BUF) build

.PHONY: protos/format
protos/format: buf/installed
	$(BUF) format -w

.PHONY: protos/lint
protos/lint: buf/installed
	$(BUF) lint
	$(BUF) lint --config=api/proto/buf-legacy.yaml api/proto

.PHONY: protos/breaking
protos/breaking: BASE=origin/master
protos/breaking: buf/installed
	@echo Checking compatibility against BASE=$(BASE)
	buf breaking . --against '.git#branch=$(BASE)'

.PHONY: lint-protos
lint-protos: protos/lint

.PHONY: lint-breaking
lint-breaking: protos/breaking

.PHONY: buf/installed
buf/installed:
	@if ! type -p $(BUF) >/dev/null; then \
		echo 'Buf is required to build/format/lint protos. Follow https://docs.buf.build/installation.'; \
		exit 1; \
	fi

# grpc generates gRPC stubs from service definitions.
# This target runs in the buildbox container.
.PHONY: grpc
grpc:
ifndef TELEPORT_DEVBOX
	$(MAKE) -C build.assets grpc
else
	$(MAKE) grpc/host
endif

# grpc/host generates gRPC stubs.
# Unlike grpc, this target runs locally.
.PHONY: grpc/host
grpc/host: protos/all
	@build.assets/genproto.sh

# protos-up-to-date checks if the generated gRPC stubs are up to date.
# This target runs in the buildbox container.
.PHONY: protos-up-to-date
protos-up-to-date:
ifndef TELEPORT_DEVBOX
	$(MAKE) -C build.assets protos-up-to-date
else
	$(MAKE) protos-up-to-date/host
endif

# protos-up-to-date/host checks if the generated gRPC stubs are up to date.
# Unlike protos-up-to-date, this target runs locally.
.PHONY: protos-up-to-date/host
protos-up-to-date/host: must-start-clean/host grpc/host
	@if ! $(GIT) diff --quiet; then \
		echo 'Please run make grpc.'; \
		exit 1; \
	fi

.PHONY: must-start-clean/host
must-start-clean/host:
	@if ! $(GIT) diff --quiet; then \
		echo 'This must be run from a repo with no unstaged commits.'; \
		exit 1; \
	fi

# crds-up-to-date checks if the generated CRDs from the protobuf stubs are up to date.
.PHONY: crds-up-to-date
crds-up-to-date: must-start-clean/host
	$(MAKE) -C integrations/operator manifests
	@if ! $(GIT) diff --quiet; then \
		echo 'Please run make -C integrations/operator manifests.'; \
		exit 1; \
	fi

print/env:
	env

.PHONY: goinstall
goinstall:
	go install $(BUILDFLAGS) \
		github.com/gravitational/teleport/tool/tsh \
		github.com/gravitational/teleport/tool/teleport \
		github.com/gravitational/teleport/tool/tctl \
		github.com/gravitational/teleport/tool/tbot

# make install will installs system-wide teleport
.PHONY: install
install: build
	@echo "\n** Make sure to run 'make install' as root! **\n"
	cp -f $(BUILDDIR)/tctl      $(BINDIR)/
	cp -f $(BUILDDIR)/tsh       $(BINDIR)/
	cp -f $(BUILDDIR)/tbot      $(BINDIR)/
	cp -f $(BUILDDIR)/teleport  $(BINDIR)/
	mkdir -p $(DATADIR)

# Docker image build. Always build the binaries themselves within docker (see
# the "docker" rule) to avoid dependencies on the host libc version.
.PHONY: image
image: OS=linux
image: TARBALL_PATH_SECTION:=-s "$(shell pwd)"
image: clean docker-binaries build-archive oss-deb
	cp ./build.assets/charts/Dockerfile $(BUILDDIR)/
	cd $(BUILDDIR) && docker build --no-cache . -t $(DOCKER_IMAGE):$(VERSION)-$(ARCH) --target teleport \
		--build-arg DEB_PATH="./teleport_$(VERSION)_$(ARCH).deb"
	if [ -f e/Makefile ]; then $(MAKE) -C e image PIV=$(PIV); fi

.PHONY: print-version
print-version:
	@echo $(VERSION)

.PHONY: chart-ent
chart-ent:
	$(MAKE) -C e chart

RUNTIME_SECTION ?=
TARBALL_PATH_SECTION ?=

ifneq ("$(RUNTIME)", "")
	RUNTIME_SECTION := -r $(RUNTIME)
endif
ifneq ("$(OSS_TARBALL_PATH)", "")
	TARBALL_PATH_SECTION := -s $(OSS_TARBALL_PATH)
endif

# build .pkg
.PHONY: pkg
pkg: | $(RELEASE_DIR)
	$(eval export DEVELOPER_ID_APPLICATION DEVELOPER_ID_INSTALLER)
	mkdir -p $(BUILDDIR)/
	cp ./build.assets/build-package.sh ./build.assets/build-common.sh $(BUILDDIR)/
	chmod +x $(BUILDDIR)/build-package.sh
	# runtime is currently ignored on OS X
	# we pass it through for consistency - it will be dropped by the build script
	cd $(BUILDDIR) && ./build-package.sh -t oss -v $(VERSION) -p pkg -b $(TELEPORT_BUNDLEID) -a $(ARCH) $(RUNTIME_SECTION) $(TARBALL_PATH_SECTION)
	cp $(BUILDDIR)/teleport-*.pkg $(RELEASE_DIR)
	if [ -f e/Makefile ]; then $(MAKE) -C e pkg; fi

# build tsh client-only .pkg
.PHONY: pkg-tsh
pkg-tsh: | $(RELEASE_DIR)
	$(eval export DEVELOPER_ID_APPLICATION DEVELOPER_ID_INSTALLER)
	./build.assets/build-pkg-tsh.sh -t oss -v $(VERSION) -b $(TSH_BUNDLEID) -a $(ARCH) $(TARBALL_PATH_SECTION)
	mkdir -p $(BUILDDIR)/
	mv tsh*.pkg* $(BUILDDIR)/
	cp $(BUILDDIR)/tsh-*.pkg $(RELEASE_DIR)

# build .rpm
.PHONY: rpm
rpm:
	mkdir -p $(BUILDDIR)/
	cp ./build.assets/build-package.sh ./build.assets/build-common.sh $(BUILDDIR)/
	chmod +x $(BUILDDIR)/build-package.sh
	cp -a ./build.assets/rpm $(BUILDDIR)/
	cp -a ./build.assets/rpm-sign $(BUILDDIR)/
	cd $(BUILDDIR) && ./build-package.sh -t oss -v $(VERSION) -p rpm -a $(ARCH) $(RUNTIME_SECTION) $(TARBALL_PATH_SECTION)
	if [ -f e/Makefile ]; then $(MAKE) -C e rpm; fi

# build unsigned .rpm (for testing)
.PHONY: rpm-unsigned
rpm-unsigned:
	$(MAKE) UNSIGNED_RPM=true rpm

# build open source .deb only
.PHONY: oss-deb
oss-deb:
	mkdir -p $(BUILDDIR)/
	cp ./build.assets/build-package.sh ./build.assets/build-common.sh $(BUILDDIR)/
	chmod +x $(BUILDDIR)/build-package.sh
	cd $(BUILDDIR) && ./build-package.sh -t oss -v $(VERSION) -p deb -a $(ARCH) $(RUNTIME_SECTION) $(TARBALL_PATH_SECTION)

# build .deb
.PHONY: deb
deb: oss-deb
	if [ -f e/Makefile ]; then $(MAKE) -C e deb; fi

# check binary compatibility with different OSes
.PHONY: test-compat
test-compat:
	./build.assets/build-test-compat.sh

.PHONY: ensure-webassets
ensure-webassets:
	@if [[ "${WEBASSETS_SKIP_BUILD}" -eq 1 ]]; then mkdir -p webassets/teleport && mkdir -p webassets/teleport/app && cp web/packages/build/index.ejs webassets/teleport/index.html; \
	else MAKE="$(MAKE)" "$(MAKE_DIR)/build.assets/build-webassets-if-changed.sh" OSS webassets/oss-sha build-ui web; fi

.PHONY: ensure-webassets-e
ensure-webassets-e:
	@if [[ "${WEBASSETS_SKIP_BUILD}" -eq 1 ]]; then mkdir -p webassets/teleport && mkdir -p webassets/e/teleport/app && cp web/packages/build/index.ejs webassets/e/teleport/index.html; \
	else MAKE="$(MAKE)" "$(MAKE_DIR)/build.assets/build-webassets-if-changed.sh" Enterprise webassets/e/e-sha build-ui-e web e/web; fi

.PHONY: init-submodules-e
init-submodules-e:
	git submodule init e
	git submodule update

# dronegen generates .drone.yml config
#
#    Usage:
#    - tsh login --proxy=platform.teleport.sh
#    - tsh apps login drone
#    - set $DRONE_TOKEN and $DRONE_SERVER (http://localhost:8080)
#    - tsh proxy app --port=8080 drone
#    - make dronegen
.PHONY: dronegen
dronegen:
	go run ./dronegen

# backport will automatically create backports for a given PR as long as you have the "gh" tool
# installed locally. To backport, type "make backport PR=1234 TO=branch/1,branch/2".
.PHONY: backport
backport:
	(cd ./assets/backport && go run main.go -pr=$(PR) -to=$(TO))

.PHONY: ensure-js-deps
ensure-js-deps:
	@if [[ "${WEBASSETS_SKIP_BUILD}" -eq 1 ]]; then mkdir -p webassets/teleport && touch webassets/teleport/index.html; \
	else yarn install --ignore-scripts; fi

.PHONY: build-ui
build-ui: ensure-js-deps
	@[ "${WEBASSETS_SKIP_BUILD}" -eq 1 ] || yarn build-ui-oss

.PHONY: build-ui-e
build-ui-e: ensure-js-deps
	@[ "${WEBASSETS_SKIP_BUILD}" -eq 1 ] || yarn build-ui-e

.PHONY: docker-ui
docker-ui:
	$(MAKE) -C build.assets ui

# rustup-install-target-toolchain ensures the required rust compiler is
# installed to build for $(ARCH)/$(OS) for the version of rust we use, as
# defined in build.assets/Makefile. It assumes that `rustup` is already
# installed for managing the rust toolchain.
.PHONY: rustup-install-target-toolchain
rustup-install-target-toolchain: RUST_VERSION := $(shell $(MAKE) --no-print-directory -C build.assets print-rust-version)
rustup-install-target-toolchain:
	rustup override set $(RUST_VERSION)
	rustup target add $(RUST_TARGET_ARCH)

# changelog generates PR changelog between the provided base tag and the tip of
# the specified branch.
#
# usage: BASE_BRANCH=branch/v13 BASE_TAG=13.2.0 make changelog
.PHONY: changelog
changelog:
	@./build.assets/changelog.sh BASE_BRANCH=$(BASE_BRANCH) BASE_TAG=$(BASE_TAG)
