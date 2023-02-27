include Makefile.Common

.PHONY: install-tools
install-tools: $(TOOLS_BIN_NAMES)

$(TOOLS_BIN_DIR):
	mkdir -p $@

$(TOOLS_BIN_NAMES): $(TOOLS_BIN_DIR) $(TOOLS_MOD_DIR)/go.mod
	cd $(TOOLS_MOD_DIR) && $(GOCMD) build -o $@ -trimpath $(filter %/$(notdir $@),$(TOOLS_PKG_NAMES))

# Build the Collector executable.
.PHONY: stdinotel
stdinotel:
	cd ./cmd/stdinotel && GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -trimpath -o ../../bin/stdinotel_$(GOOS)_$(GOARCH)$(EXTENSION) \
		$(BUILD_INFO) -tags $(GO_BUILD_TAGS) .


.PHONY: lint
lint: $(LINT) checklicense misspell
	$(LINT) run --allow-parallel-runners --build-tags integration

.PHONY: tidy
tidy:
	rm -fr go.sum
	$(GOCMD) mod tidy -compat=1.19
