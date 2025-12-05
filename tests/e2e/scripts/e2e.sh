#!/bin/bash

if [[ ! $(echo "$PWD" | rev | cut -d '/' -f 1 | rev) == terraform-provider-criblio ]]; then
	echo "Script must be run from project root directory"
	exit 1
fi

flags=""
if [[ -z $CRIBL_CLOUD_DOMAIN ]]; then
	flags='-var=onprem=true -parallelism=1'
fi

cd tests/e2e

exitCode=0

#remove stale terraform files
echo -n "Removing stale Terraform files..."
rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform
echo " Done!"

#the remote mirror won't have our custom version, so this will always fail, hence || true
terraform providers mirror ./local-plugins || true
terraform init -plugin-dir ./local-plugins

terraform apply -auto-approve $flags
tfApply=$?

#remove our state files to test if we can import everything into a fresh state
rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform
terraform providers mirror ./local-plugins || true
terraform init -plugin-dir ./local-plugins

#because imports leverage the data sources to read in state files,
#this will exersize both import and data functionality
./scripts/import.sh
tfImport=$?

#make sure we didn't break something
terraform refresh $flags
tfRefresh=$?

terraform destroy -auto-approve $flags
tfDestroy=$?
if [[ $tfApply -ne 0 ]] || [[ $tfImport -ne 0 ]] || [[ $tfRefresh -ne 0 ]] || [[ $tfDestroy -ne 0 ]]; then
	echo echo "***FAILURE IN TERRAFORM OPS***"
	echo
	echo "Exit Codes -> Apply: $tfApply, Import: $tfImport, Refresh: $tfRefresh, Destroy: $tfDestroy"
	echo
	exitCode=1
else
	exitCode=0
fi

exit $exitCode
