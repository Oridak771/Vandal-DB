# Makefile for Phantom-DB

.PHONY: all build test clean run

all: build

# =================================================================================================
# Main Targets
# =================================================================================================

# Build manager binary
build: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# =================================================================================================
# Code Generation and Formatting
# =================================================================================================

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# =================================================================================================
# Tooling
# =================================================================================================

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e; \
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d); \
	cd $$CONTROLLER_GEN_TMP_DIR; \
	go mod init tmp; \
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1; \
	rm -rf $$CONTROLLER_GEN_TMP_DIR; \
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
