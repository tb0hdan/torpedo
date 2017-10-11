.PHONY: tmp

PWD=$(shell pwd)
SOURCE=$(PWD)/tmp/torpedo

all: $(SOURCE) build

$(SOURCE):
	@mkdir -p $(PWD)/tmp
	@git clone https://github.com/tb0hdan/torpedo $(PWD)/tmp/torpedo

build:
	@rm -f $(PWD)/tmp/torpedo/src/torpedobot/main_plugins.go
	@cp -r $(PWD)/torpedo_text_plugin $(PWD)/tmp/torpedo/src
	@cp $(PWD)/main_plugins_text.go $(PWD)/tmp/torpedo/src/torpedobot/
	@$(MAKE) -C $(PWD)/tmp/torpedo
	@cd $(PWD)
	@mv $(PWD)/tmp/torpedo/bin/torpedobot $(PWD)/

clean:
	@rm -rf $(PWD)/tmp