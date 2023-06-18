VERSION := $(shell echo $(shell git describe --tags || git branch --show-current) | sed 's/^v//')

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/EscanBE/house-keeper/constants.VERSION=$(VERSION) \
            -X github.com/EscanBE/house-keeper/constants.BUILD_FROM_SOURCE=yes

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "building hkd binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go build -mod=readonly $(BUILD_FLAGS) -o build/hkd ./cmd/hkd
	@echo "Builded successfully"
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing hkd binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/hkd
	@echo "Installed successfully"
.PHONY: install