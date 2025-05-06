CMDS := $(notdir $(wildcard cmd/*))
all: $(patsubst %,build-%, $(CMDS))

define BUILD_CMD
build-$(1):
  @$(GARBLE) build -o $(BIN_DIR)/$(1) ./cmd/$(1)
endef
$(foreach cmd,$(CMDS),$(eval $(call BUILD_CMD,$(cmd))))