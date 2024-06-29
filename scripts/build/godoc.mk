vangen_version=v1.4.0

vangen-install:
	@echo "--> Installing vangen $(vangen_version)"
	@go install 4d63.com/vangen@$(vangen_version)

update-godocs:
	@echo "--> Running update-godocs"
	$(MAKE) vangen-install
	@vangen -config vangen.json

.PHONY: update-godocs