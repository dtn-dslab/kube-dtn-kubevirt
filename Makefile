# Image URL to use all building/pushing image targets
IMG ?= harbor.sail.se.sjtu.edu.cn/kubedtn
SIDECAR_IMG = $(IMG)/kubedtn-sidecar
HACK_IMG = $(IMG)/kubedtn-hack
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25.0

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git describe --always)
SUBVERSION := $(shell echo $$DTN_KV_SUB)
TAG := $(BRANCH)-$(COMMIT)-$(SUBVERSION)

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

KUSTOMIZE ?= $(LOCALBIN)/kustomize
KUSTOMIZE_VERSION ?= v5.3.0
KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"

.PHONY: kustomize 
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: sidecar-docker
sidecar-docker:
	docker build . -t $(SIDECAR_IMG):$(TAG)

.PHONY: sidecar-push
sidecar-push:
	docker push $(SIDECAR_IMG):$(TAG)

.PHONY: sidecar-all
sidecar-all: sidecar-docker sidecar-push
