.PHONY: *

OS=$(shell uname | tr "[:upper:]" "[:lower:]")
ARCH=$(shell uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
e2e-test:
	@cd tests/e2e; rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform local-plugins
	mkdir -p tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)
	go mod tidy
	go build -o tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)/terraform-provider-criblio_v999.99.9
	#the remote mirror won't have our custom version, so this will always fail, hence || true
	@cd tests/e2e; terraform providers mirror ./local-plugins || true; ls -R local-plugins; terraform init -plugin-dir ./local-plugins; 
	#imports can be flakey, pass them anyway
	@cd tests/e2e; terraform import criblio_group.syslog_worker_group "syslog-workers"; terraform import criblio_group.my_edge_fleet "my-edge-fleet"; terraform import criblio_appscope_config.my_appscopeconfig "default"; terraform import criblio_grok.my_grok "default"; terraform import criblio_global_var.my_globalvar "default"; terraform import criblio_subscription.my_subscription "default"; terraform import criblio_regex.my_regex "default"; terraform import criblio_subscription.my_subscription_with_enabled "default" || true
	@cd tests/e2e; terraform apply -auto-approve; flag2=$$?; terraform destroy -auto-approve; flag3=$$?; if [ $$flag2 -ne 0 ] || [ $$flag3 -ne 0 ]; then echo; echo "***FAILURE IN TERRAFORM OPS***"; echo; exit 1; fi

acceptance-test:
	export CRIBL_SERVER_URL="https://app.cribl-playground.cloud" && \
	export CRIBL_ORGANIZATION_ID="beautiful-nguyen-y8y4azd" &&  \
	export CRIBL_WORKSPACE_ID="tfprovider" && \
	export TF_ACC=true && \
	go test -v ./tests/acceptance

test-cleanup:
	@cd tests/e2e; rm -rf local-plugins .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

unit-test: 
	go test -v ./internal/sdk/internal/hooks

test-speakeasy: 
	speakeasy test && speakeasy lint openapi --non-interactive -s openapi.yml

e2e-test-speakeasy: 
	speakeasy run --skip-versioning --output console --minimal --skip-upload-spec --skip-versioning --skip-compile

