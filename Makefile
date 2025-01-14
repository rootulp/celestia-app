# GIT_TAG is an environment variable that is set to the latest git tag on the
# current commit with the following example priority: v2.2.0, v2.2.0-mocha,
# v2.2.0-arabica, v2.2.0-rc0, v2.2.0-beta, v2.2.0-alpha. If no tag points to the
# current commit, git describe is used. The priority in this command is
# necessary because `git tag --sort=-creatordate` only works for annotated tags
# with metadata. Git tags created via GitHub releases are not annotated and do
# not have metadata like creatordate. Therefore, this command is a hacky attempt
# to get the most recent tag on the current commit according to Celestia's
# testnet versioning scheme + SemVer.
GIT_TAG := $(shell git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' \
    || git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-mocha$$' \
    || git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-arabica$$' \
    || git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]*$$' \
    || git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-(beta)$$' \
    || git tag --points-at HEAD --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-(alpha)$$' \
    || git describe --tags)
VERSION := $(shell echo $(GIT_TAG) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
IMAGE := ghcr.io/tendermint/docker-build-proto:latest
DOCKER_PROTO_BUILDER := docker run -v $(shell pwd):/workspace --workdir /workspace $(IMAGE)
PROJECTNAME=$(shell basename "$(PWD)")
HTTPS_GIT := https://github.com/celestiaorg/celestia-app.git
PACKAGE_NAME          := github.com/celestiaorg/celestia-app/v3
# Before upgrading the GOLANG_CROSS_VERSION, please verify that a Docker image exists with the new tag.
# See https://github.com/goreleaser/goreleaser-cross/pkgs/container/goreleaser-cross
GOLANG_CROSS_VERSION  ?= v1.23.1
# Set this to override the max square size of the binary
OVERRIDE_MAX_SQUARE_SIZE ?=
# Set this to override the upgrade height delay of the binary
OVERRIDE_UPGRADE_HEIGHT_DELAY ?=

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=celestia-app \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=celestia-appd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/celestiaorg/celestia-app/v3/pkg/appconsts.OverrideSquareSizeUpperBoundStr=$(OVERRIDE_MAX_SQUARE_SIZE) \
		  -X github.com/celestiaorg/celestia-app/v3/pkg/appconsts.OverrideUpgradeHeightDelayStr=$(OVERRIDE_UPGRADE_HEIGHT_DELAY)

BUILD_FLAGS := -tags "ledger" -ldflags '$(ldflags)'

## help: Get more info on make commands.
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help

## build: Build the celestia-appd binary into the ./build directory.
build: mod
	@cd ./cmd/celestia-appd
	@mkdir -p build/
	@echo "--> Building build/celestia-appd"
	@go build $(BUILD_FLAGS) -o build/ ./cmd/celestia-appd
.PHONY: build

## install: Build and install the celestia-appd binary into the $GOPATH/bin directory.
install: check-bbr
	@echo "--> Installing celestia-appd"
	@go install $(BUILD_FLAGS) ./cmd/celestia-appd
.PHONY: install

## mod: Update all go.mod files.
mod:
	@echo "--> Updating go.mod"
	@go mod tidy
	@echo "--> Updating go.mod in ./test/interchain"
	@(cd ./test/interchain && go mod tidy)
.PHONY: mod

## mod-verify: Verify dependencies have expected content.
mod-verify: mod
	@echo "--> Verifying dependencies have expected content"
	GO111MODULE=on go mod verify
.PHONY: mod-verify

## proto-gen: Generate protobuf files. Requires docker.
proto-gen:
	@echo "--> Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen:v0.7 sh ./scripts/protocgen.sh
.PHONY: proto-gen

## proto-lint: Lint protobuf files. Requires docker.
proto-lint:
	@echo "--> Linting Protobuf files"
	@$(DOCKER_BUF) lint --error-format=json
.PHONY: proto-lint

## proto-check-breaking: Check if there are any breaking change to protobuf definitions.
proto-check-breaking:
	@echo "--> Checking if Protobuf definitions have any breaking changes"
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main
.PHONY: proto-check-breaking

## proto-format: Format protobuf files. Requires Docker.
proto-format:
	@echo "--> Formatting Protobuf files"
	@$(DOCKER_PROTO_BUILDER) find . -name '*.proto' -path "./proto/*" -exec clang-format -i {} \;
.PHONY: proto-format

## build-docker: Build the celestia-appd docker image from the current branch. Requires docker.
build-docker:
	@echo "--> Building Docker image"
	$(DOCKER) build -t celestiaorg/celestia-app -f docker/Dockerfile .
.PHONY: build-docker

## build-ghcr-docker: Build the celestia-appd docker image from the last commit. Requires docker.
build-ghcr-docker:
	@echo "--> Building Docker image"
	$(DOCKER) build -t ghcr.io/celestiaorg/celestia-app:$(COMMIT) -f docker/Dockerfile .
.PHONY: build-ghcr-docker

## publish-ghcr-docker: Publish the celestia-appd docker image. Requires docker.
publish-ghcr-docker:
# Make sure you are logged in and authenticated to the ghcr.io registry.
	@echo "--> Publishing Docker image"
	$(DOCKER) push ghcr.io/celestiaorg/celestia-app:$(COMMIT)
.PHONY: publish-ghcr-docker

## lint: Run all linters; golangci-lint, markdownlint, hadolint, yamllint.
lint:
	@echo "--> Running golangci-lint"
	@golangci-lint run
	@echo "--> Running markdownlint"
	@markdownlint --config .markdownlint.yaml '**/*.md'
	@echo "--> Running hadolint"
	@hadolint docker/Dockerfile
	@hadolint docker/txsim/Dockerfile
	@echo "--> Running yamllint"
	@yamllint --no-warnings . -c .yamllint.yml
.PHONY: lint

## markdown-link-check: Check all markdown links.
markdown-link-check:
	@echo "--> Running markdown-link-check"
	@find . -name \*.md -print0 | xargs -0 -n1 markdown-link-check
.PHONY: markdown-link-check


## fmt: Format files per linters golangci-lint and markdownlint.
fmt:
	@echo "--> Running golangci-lint --fix"
	@golangci-lint run --fix
	@echo "--> Running markdownlint --fix"
	@markdownlint --fix --quiet --config .markdownlint.yaml .
.PHONY: fmt

## test: Run tests.
test:
	@echo "--> Running tests"
	@go test -timeout 30m ./...
.PHONY: test

## test-short: Run tests in short mode.
test-short:
	@echo "--> Running tests in short mode"
	@go test ./... -short -timeout 1m
.PHONY: test-short

## test-e2e: Run end to end tests via knuu. This command requires a kube/config file to configure kubernetes.
test-e2e:
	@echo "--> Running end to end tests"
	go run ./test/e2e $(filter-out $@,$(MAKECMDGOALS))
.PHONY: test-e2e

## test-race: Run tests in race mode.
test-race:
# TODO: Remove the -skip flag once the following tests no longer contain data races.
# https://github.com/celestiaorg/celestia-app/issues/1369
	@echo "--> Running tests in race mode"
	@go test -timeout 15m ./... -v -race -skip "TestPrepareProposalConsistency|TestIntegrationTestSuite|TestBlobstreamRPCQueries|TestSquareSizeIntegrationTest|TestStandardSDKIntegrationTestSuite|TestTxsimCommandFlags|TestTxsimCommandEnvVar|TestMintIntegrationTestSuite|TestBlobstreamCLI|TestUpgrade|TestMaliciousTestNode|TestBigBlobSuite|TestQGBIntegrationSuite|TestSignerTestSuite|TestPriorityTestSuite|TestTimeInPrepareProposalContext|TestBlobstream|TestCLITestSuite|TestLegacyUpgrade|TestSignerTwins|TestConcurrentTxSubmission|TestTxClientTestSuite|Test_testnode|TestEvictions"
.PHONY: test-race

## test-bench: Run unit tests in bench mode.
test-bench:
	@echo "--> Running tests in bench mode"
	@go test -bench=. ./...
.PHONY: test-bench

## test-coverage: Generate test coverage.txt
test-coverage:
	@echo "--> Generating coverage.txt"
	@export VERSION=$(VERSION); bash -x scripts/test_cover.sh
.PHONY: test-coverage

## test-fuzz: Run all fuzz tests.
test-fuzz:
	bash -x scripts/test_fuzz.sh
.PHONY: test-fuzz

## txsim-install: Install the tx simulator.
txsim-install:
	@echo "--> Installing tx simulator"
	@go install $(BUILD_FLAGS) ./test/cmd/txsim
.PHONY: txsim-install

## txsim-build: Build the tx simulator binary into the ./build directory.
txsim-build:
	@echo "--> Building tx simulator"
	@cd ./test/cmd/txsim
	@mkdir -p build/
	@go build $(BUILD_FLAGS) -o build/ ./test/cmd/txsim
	@go mod tidy
.PHONY: txsim-build

## txsim-build-docker: Build the tx simulator Docker image. Requires Docker.
txsim-build-docker:
	docker build -t ghcr.io/celestiaorg/txsim -f docker/txsim/Dockerfile  .
.PHONY: txsim-build-docker

## adr-gen: Download the ADR template from the celestiaorg/.github repo.
adr-gen:
	@echo "--> Downloading ADR template"
	@curl -sSL https://raw.githubusercontent.com/celestiaorg/.github/main/adr-template.md > docs/architecture/adr-template.md
.PHONY: adr-gen

## goreleaser-check: Check the .goreleaser.yaml config file.
goreleaser-check:
	@if [ ! -f ".release-env" ]; then \
		echo "A .release-env file was not found but is required to create prebuilt binaries. This command is expected to be run in CI where a .release-env file exists. If you need to run this command locally to attach binaries to a release, you need to create a .release-env file with a Github token (classic) that has repo:public_repo scope."; \
		exit 1;\
	fi
	docker run \
		--rm \
		--env CGO_ENABLED=1 \
		--env GORELEASER_CURRENT_TAG=${GIT_TAG} \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		check
.PHONY: goreleaser-check

## prebuilt-binary: Create prebuilt binaries and attach them to GitHub release. Requires Docker.
prebuilt-binary:
	@if [ ! -f ".release-env" ]; then \
		echo "A .release-env file was not found but is required to create prebuilt binaries. This command is expected to be run in CI where a .release-env file exists. If you need to run this command locally to attach binaries to a release, you need to create a .release-env file with a Github token (classic) that has repo:public_repo scope."; \
		exit 1;\
	fi
	docker run \
		--rm \
		--env CGO_ENABLED=1 \
		--env GORELEASER_CURRENT_TAG=${GIT_TAG} \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean
.PHONY: prebuilt-binary

## check-bbr: Check if your system uses BBR congestion control algorithm. Only works on Linux.
check-bbr:
	@echo "Checking if BBR is enabled..."
	@if [ "$$(sysctl net.ipv4.tcp_congestion_control | awk '{print $$3}')" != "bbr" ]; then \
		echo "WARNING: BBR is not enabled. Please enable BBR for optimal performance. Call make enable-bbr or see Usage section in the README."; \
	else \
		echo "BBR is enabled."; \
	fi
.PHONY: check-bbr

## enable-bbr: Enable BBR congestion control algorithm. Only works on Linux.
enable-bbr:
	@echo "Configuring system to use BBR..."
	@if [ "$(sysctl net.ipv4.tcp_congestion_control | awk '{print $3}')" != "bbr" ]; then \
	    echo "BBR is not enabled. Configuring BBR..."; \
	    sudo modprobe tcp_bbr; \
            echo tcp_bbr | sudo tee -a /etc/modules; \
	    echo "net.core.default_qdisc=fq" | sudo tee -a /etc/sysctl.conf; \
	    echo "net.ipv4.tcp_congestion_control=bbr" | sudo tee -a /etc/sysctl.conf; \
	    sudo sysctl -p; \
	    echo "BBR has been enabled."; \
	else \
	    echo "BBR is already enabled."; \
	fi
.PHONY: enable-bbr

## disable-bbr: Disable BBR congestion control algorithm and revert to default.
disable-bbr:
	@echo "Disabling BBR and reverting to default congestion control algorithm..."
	@if [ "$$(sysctl net.ipv4.tcp_congestion_control | awk '{print $$3}')" = "bbr" ]; then \
	    echo "BBR is currently enabled. Reverting to Cubic..."; \
	    sudo sed -i '/^net.core.default_qdisc=fq/d' /etc/sysctl.conf; \
	    sudo sed -i '/^net.ipv4.tcp_congestion_control=bbr/d' /etc/sysctl.conf; \
	    sudo modprobe -r tcp_bbr; \
	    echo "net.ipv4.tcp_congestion_control=cubic" | sudo tee -a /etc/sysctl.conf; \
	    sudo sysctl -p; \
	    echo "BBR has been disabled, and Cubic is now the default congestion control algorithm."; \
	else \
	    echo "BBR is not enabled. No changes made."; \
	fi
.PHONY: disable-bbr

## enable-mptcp: Enable mptcp over multiple ports (not interfaces). Only works on Linux Kernel 5.6 and above.
enable-mptcp:
	@echo "Configuring system to use mptcp..."
	@sudo sysctl -w net.mptcp.enabled=1
	@sudo sysctl -w net.mptcp.mptcp_path_manager=ndiffports
	@sudo sysctl -w net.mptcp.mptcp_ndiffports=16
	@echo "Making MPTCP settings persistent across reboots..."
	@echo "net.mptcp.enabled=1" | sudo tee -a /etc/sysctl.conf
	@echo "net.mptcp.mptcp_path_manager=ndiffports" | sudo tee -a /etc/sysctl.conf
	@echo "net.mptcp.mptcp_ndiffports=16" | sudo tee -a /etc/sysctl.conf
	@echo "MPTCP configuration complete and persistent!"

.PHONY: enable-mptcp

## disable-mptcp: Disables mptcp over multiple ports. Only works on Linux Kernel 5.6 and above.
disable-mptcp:
	@echo "Disabling MPTCP..."
	@sudo sysctl -w net.mptcp.enabled=0
	@sudo sysctl -w net.mptcp.mptcp_path_manager=default
	@echo "Removing MPTCP settings from /etc/sysctl.conf..."
	@sudo sed -i '/net.mptcp.enabled=1/d' /etc/sysctl.conf
	@sudo sed -i '/net.mptcp.mptcp_path_manager=ndiffports/d' /etc/sysctl.conf
	@sudo sed -i '/net.mptcp.mptcp_ndiffports=16/d' /etc/sysctl.conf
	@echo "MPTCP configuration reverted!"

.PHONY: disable-mptcp

CONFIG_FILE ?= ${HOME}/.celestia-app/config/config.toml
SEND_RECV_RATE ?= 10485760  # 10 MiB

## configure-v3: Modifies config file in-place to conform to v3.x recommendations.
configure-v3:
	@echo "Using config file at: $(CONFIG_FILE)"
	@if [ "$$(uname)" = "Darwin" ]; then \
		sed -i '' "s/^recv_rate = .*/recv_rate = $(SEND_RECV_RATE)/" $(CONFIG_FILE); \
		sed -i '' "s/^send_rate = .*/send_rate = $(SEND_RECV_RATE)/" $(CONFIG_FILE); \
		sed -i '' "s/ttl-num-blocks = .*/ttl-num-blocks = 12/" $(CONFIG_FILE); \
	else \
		sed -i "s/^recv_rate = .*/recv_rate = $(SEND_RECV_RATE)/" $(CONFIG_FILE); \
		sed -i "s/^send_rate = .*/send_rate = $(SEND_RECV_RATE)/" $(CONFIG_FILE); \
		sed -i "s/ttl-num-blocks = .*/ttl-num-blocks = 12/" $(CONFIG_FILE); \
	fi
.PHONY: configure-v3


## debug-version: Print the git tag and version.
debug-version:
	@echo "GIT_TAG: $(GIT_TAG)"
	@echo "VERSION: $(VERSION)"
.PHONY: debug-version
