
# Go parameters
GO=go
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 
BUILDDIR = build

# All targets
all: test

# Rules for building
.PHONY: test
test:
	@PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) test -tags "$(TAGS)" ./...

.PHONY: mkdir
mkdir:
	install -d $(BUILDDIR)

.PHONY: clean
clean: 
	rm -fr $(BUILDDIR)
	$(GO) clean
	$(GO) mod tidy
