config.plist: Sample.plist ssdt kexts tools drivers
	bash generate-plist.sh Sample.plist \
		--ssdt ssdt/*.aml \
    	--kext kext/Lilu.kext kext/*.kext \
    	--tool tool/*.efi \
    	--driver driver/*.efi > config.plist
.PHONY: clean
clean:
	rm -f ./config.plist
	rm -f ./ssdt/*.aml

ssdt: $(patsubst %.dsl, %.aml, $(wildcard ssdt/*.dsl))
ssdt/%.aml: ssdt/%.dsl
	iasl $^
kexts: $(wildcard kext/*.kext)
tools: $(wildcard tool/*.tool)
drivers: $(wildcard driver/*.efi)
