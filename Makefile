VERSION         := $(shell cat VERSION)

.PHONY: release

release:
	git tag -a $(VERSION) -m "Release" && git push origin $(VERSION)
