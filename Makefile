.PHONY: *

OS=$(shell uname | tr "[:upper:]" "[:lower:]")
ARCH=$(shell uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
e2e-test:
	@cd tests/e2e; rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform local-plugins
	mkdir -p tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)
	go mod tidy
	go build -o tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)/terraform-provider-criblio_v999.99.9
	#the remote mirror won't have our custom version, so this will always fail, hence || true
	@cd tests/e2e; terraform providers mirror ./local-plugins || true
	@cd tests/e2e; ls -R local-plugins; terraform init -plugin-dir ./local-plugins; terraform apply -auto-approve; terraform destroy -auto-approve

unit-test:
	export TF_ACC=true; go test -v ./tests/unit

test-cleanup:
	@cd tests/e2e; rm -rf local-plugins .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

build-speakeasy: 
	speakeasy run --skip-versioning --output console --minimal

test-speakeasy: 
	speakeasy test
