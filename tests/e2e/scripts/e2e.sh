#!/bin/bash

if [[ ! `echo "$PWD"| rev | cut -d '/' -f 1 | rev` == terraform-provider-criblio ]]; then
    echo "Script must be run from project root directory"
    return 1
fi

cd tests/e2e

#remove stale terraform files
echo -n "Removing stale Terraform files..."
rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform 
echo " Done!" 

#the remote mirror won't have our custom version, so this will always fail, hence || true
terraform providers mirror ./local-plugins || true
terraform init -plugin-dir ./local-plugins

./scripts/import.sh init

terraform apply -auto-approve 
tfApply=$? 

#remove our state files to test if we can import everything into a fresh state
rm -rf .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup .terraform 
terraform providers mirror ./local-plugins || true
terraform init -plugin-dir ./local-plugins

#because imports leverage the data sources to read in state files, 
#this will exersize both import and data functionality
./scripts/import.sh refresh_state

#make sure we didn't break something
terraform refresh
tfRefresh=$?

terraform destroy -auto-approve 
tfDestroy=$? 
if [[ $tfApply -ne 0 ]] || [[ $tfRefresh -ne 0 ]] || [[ $tfDestroy -ne 0  ]]; then 
    echo echo "***FAILURE IN TERRAFORM OPS***" 
    echo 
    echo "Exit Codes -> Apply: $tfApply, Refresh: $tfRefresh, Destroy: $tfDestroy"
    echo 
    exit 1 
fi

exit 0
