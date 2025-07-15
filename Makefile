.PHONY: *

OS=$(shell uname | tr "[:upper:]" "[:lower:]")
ARCH=$(shell uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
e2e-test:
	export CRIBL_SERVER_URL="https://app.cribl-playground.cloud" && \
	export CRIBL_ORGANIZATION_ID="beautiful-nguyen-y8y4azd" &&  \
	export CRIBL_WORKSPACE_ID="tfprovider2" 
	@cd tests/e2e; rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform local-plugins
	mkdir -p tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)
	go mod tidy
	go build -o tests/e2e/local-plugins/registry.terraform.io/criblio/criblio/999.99.9/$(OS)_$(ARCH)/terraform-provider-criblio_v999.99.9
	#the remote mirror won't have our custom version, so this will always fail, hence || true
	@cd tests/e2e; terraform providers mirror ./local-plugins || true
	#echo lines required in order to setup imports
	@echo 'import { ' >> tests/e2e/edge_fleet.tf; echo 'to = criblio_group.my_edge_fleet ' >> tests/e2e/edge_fleet.tf; echo 'id = "my-edge-fleet" ' >> tests/e2e/edge_fleet.tf; echo '}' >> tests/e2e/edge_fleet.tf
	@echo 'import { ' >> tests/e2e/stream_syslog_to_lake.tf; echo 'to = criblio_group.syslog_worker_group ' >> tests/e2e/stream_syslog_to_lake.tf; echo 'id = "syslog-workers" ' >> tests/e2e/stream_syslog_to_lake.tf; echo '}' >> tests/e2e/stream_syslog_to_lake.tf
	@cd tests/e2e; ls -R local-plugins; terraform init -plugin-dir ./local-plugins; flag1=$$?; terraform apply -auto-approve; flag2=$$?; terraform destroy -auto-approve; flag3=$$?; if [[ $$flag1 -ne 0 ]] || [[ $$flag2 -ne 0 ]] || [[ $$flag3 -ne 0 ]]; then echo; echo "***FAILURE IN TERRAFORM OPS***"; echo; exit 1; fi

acceptance-test:
	export CRIBL_SERVER_URL="https://app.cribl-playground.cloud" && \
	export CRIBL_ORGANIZATION_ID="beautiful-nguyen-y8y4azd" &&  \
	export CRIBL_WORKSPACE_ID="tfprovider" && \
	export TF_ACC=true && \
	go test -v ./tests/acceptance

test-cleanup:
	@cd tests/e2e; rm -rf local-plugins .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

build-speakeasy: 
	speakeasy run --skip-versioning --output console --minimal

unit-test: 
	go test -v ./internal/sdk/internal/hooks

test-speakeasy: 
	speakeasy test
