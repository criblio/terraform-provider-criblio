.PHONY: *

OS=$(shell uname | tr "[:upper:]" "[:lower:]")
ARCH=$(shell uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
e2e-test:
	mkdir -p tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)
	go mod tidy
	go build -o tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)/terraform-provider-criblio_v999.99.9
	@#run our wrapper script targeting our new plugin
	./tests/e2e/scripts/e2e.sh 

acceptance-test:
	export CRIBL_ORGANIZATION_ID="beautiful-nguyen-y8y4azd" &&  \
	export CRIBL_WORKSPACE_ID="tfprovider" && \
	export CRIBL_CLOUD_DOMAIN="cribl-playground.cloud" && \
	export TF_ACC=true && \
	go test -v ./tests/acceptance

test-cleanup:
	@cd tests/e2e; rm -rf local-plugins .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

unit-test: 
	go test -v ./internal/sdk/internal/hooks

test-speakeasy: 
	speakeasy test && speakeasy lint openapi --non-interactive -s openapi.yml

build-speakeasy: 
	speakeasy run --skip-versioning --output console --minimal

install-sbom-tools:
	go install github.com/anchore/syft/cmd/syft@latest

sbom: install-sbom-tools
	@echo "Generating SBOM files..."
	syft scan . -o spdx-json=sbom-spdx.json
	syft scan . -o cyclonedx-json=sbom-cyclonedx.json
	@echo "SBOM files generated:"
	@ls -la sbom-*.json 

