DIST   := ./dist
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

BINARIES := command-server query-server gateway-server publisher-server

# ── Platforms ────────────────────────────────────────────────────────────────
# Edit the list below to add or remove build targets.
PLATFORMS := \
	linux/amd64   \
	linux/arm64   \
	linux/386     \
	darwin/amd64  \
	darwin/arm64  \
	windows/amd64 \
	windows/arm64 \
	windows/386

# ── Helpers ──────────────────────────────────────────────────────────────────
bin_name = $(1)-$(2)-$(3)$(if $(filter windows,$(2)),.exe)

# ── Targets ──────────────────────────────────────────────────────────────────

.PHONY: build build-all clean \
        $(addprefix build-,$(BINARIES)) \
        $(addprefix build-all-,$(BINARIES))

## build: compile every binary for the host platform
build: $(addprefix build-,$(BINARIES))

## build-<binary>: compile a single binary for the host platform
$(addprefix build-,$(BINARIES)):
	@mkdir -p $(DIST)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "-s -w" \
		-o $(DIST)/$(call bin_name,$(patsubst build-%,%,$@),$(GOOS),$(GOARCH)) \
		./cmd/$(patsubst build-%,%,$@)

## build-all: cross-compile every binary for every PLATFORM
build-all: $(addprefix build-all-,$(BINARIES))

## build-all-<binary>: cross-compile a single binary for every PLATFORM
$(addprefix build-all-,$(BINARIES)):
	@mkdir -p $(DIST)
	@set -e; $(foreach p,$(PLATFORMS), \
	  $(eval os   = $(word 1,$(subst /, ,$(p)))) \
	  $(eval arch = $(word 2,$(subst /, ,$(p)))) \
	  GOOS=$(os) GOARCH=$(arch) go build -trimpath -ldflags "-s -w" \
	      -o $(DIST)/$(call bin_name,$(patsubst build-all-%,%,$@),$(os),$(arch)) \
	      ./cmd/$(patsubst build-all-%,%,$@) ; \
	)

## clean: remove build artifacts
clean:
	rm -rf $(DIST)
