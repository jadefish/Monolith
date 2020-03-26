RUBY ?= $(shell command -v ruby)
TEMPLATES_DIR ?= ./configuration/templates
OUT_DIR ?= ./configuration/out
SERIALS_FILE ?= ./configuration/serials.yml

TEMPLATES := $(wildcard $(TEMPLATES_DIR)/*.plist)

default: $(RUBY) $(TEMPLATES) $(OUT_DIR)
	$(RUBY) ./configuration/make_config.rb \
			--in="$(TEMPLATES_DIR)" \
			--out="$(OUT_DIR)" \
			--serials="$(SERIALS_FILE)"

$(OUT_DIR):
	test -d $(BIN_DIR) || $(error "directory not found: $(OUT_DIR)")
