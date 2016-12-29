BUILD_TYPE?=build

bui:
	go $(BUILD_TYPE) .


ARTIFACTS := artifacts/bui-{{.OS}}-{{.Arch}}
LDFLAGS := -X main.Version=$(VERSION)
release:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"

	@echo "Checking that TARGETS was defined in the calling environment"
	@test -n "$(TARGETS)"
	@echo "OK.  TARGETS='$(TARGETS)'"
	rm -rf artifacts
	gox -osarch="$(TARGETS)" -ldflags="$(LDFLAGS)" --output="$(ARTIFACTS)/bui"      .

	cd artifacts && for x in bui-*; do cp -a ../ui/ $$x/ui; tar -czvf $$x.tar.gz $$x; rm -r $$x;  done