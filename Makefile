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
	go test -v -timeout 30m ./tests/acceptance -parallel 1

test-cleanup:
	@cd tests/e2e; rm -rf local-plugins .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

unit-test: 
	go test -v ./internal/auth/... ./internal/restclient/... ./internal/sdk/credentials ./internal/sdk/internal/hooks ./tools/sync-openapi ./tools/merge-spec/...

sync-openapi:
	go run ./tools/sync-openapi

merge:
	go run ./tools/merge-spec

generate: merge
	go run ./tools/codegen --spec merged-spec.yml
	gofmt -w internal/provider/*_client.go internal/provider/*_types.go internal/provider/*_resource.go internal/provider/*_data_source.go tests/acceptance/*_test.go
	go run ./tools/bump-doc-version

migrate:
	@scripts/migrate.sh migrate

migrate-batch:
	@scripts/migrate.sh batch

unit-test-import-cli:
	go test -v ./tools/import-cli/...

integration-test-import-cli: build-import-cli
	go test -v -tags=integration ./tools/import-cli/integration/...

build-import-cli:
	go build -o goatify ./tools/import-cli

test-speakeasy: 
	speakeasy test && speakeasy lint openapi --non-interactive -s openapi.yml

# Fast/minimal Speakeasy run (CI-friendly).
build-speakeasy:
	speakeasy run --skip-versioning --output console --minimal

# Full Speakeasy run
build-speakeasy-full:
	GOTOOLCHAIN=go1.25.0 speakeasy run

install-sbom-tools:
	go install github.com/anchore/syft/cmd/syft@latest

sbom: install-sbom-tools
	@echo "Generating SBOM files..."
	syft scan . -o spdx-json=sbom-spdx.json
	syft scan . -o cyclonedx-json=sbom-cyclonedx.json
	@echo "SBOM files generated:"
	@ls -la sbom-*.json 
