#!/bin/bash

checkWorkspaceStatus() {
   curl -s -X GET 'https://api.cribl-playground.cloud/v2/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider2' \
       --header "Authorization: Bearer $BEARER_TOKEN"  | jq '.state'
}

#get the bearer token
BEARER_TOKEN=`curl -s -L -X POST 'https://login.cribl-playground.cloud/oauth/token' -H 'Content-Type: application/json' -d "{
    \"grant_type\": \"client_credentials\",
    \"client_id\": \"$CRIBL_CLIENT_ID\",
    \"client_secret\": \"$CRIBL_CLIENT_SECRET\",
    \"audience\": \"https://api.cribl-playground.cloud\"
}" | jq '.access_token' | tr -d '"'`

if [[ $BEARER_TOKEN == ""  ]]; then
    echo "Bearer token get failed!"
    exit 1
fi

#perform delete
curl -s -X DELETE 'https://api.cribl-playground.cloud/v2/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider2' \
     --header "Authorization: Bearer $BEARER_TOKEN"

if [[ $? -ne 0  ]]; then
    echo "Workspace delete failed!"
    exit 1
fi

echo

END_TIME=$(( $(date +%s) + 1800 ))
success=false
while [[ $(date +%s) -lt $END_TIME ]]; do
    #check status in loop until deleted
    if [[ `checkWorkspaceStatus` == "null" ]]; then
	success=true
        break
    else 
        echo "waiting for workspace to be deleted, sleeping for 60 seconds..."
	sleep 60
    fi
done

if [[ $success == "false" ]]; then
    echo "Workspace delete timed out!"
    exit 1
fi

#create it here
curl -s -X POST "https://api.cribl-playground.cloud/v2/organizations/beautiful-nguyen-y8y4azd/workspaces" \
     -H 'Content-Type: application/json' \
     --header "Authorization: Bearer $BEARER_TOKEN" \
     -d '{"workspaceId": "tfprovider2", "alias": "e2e-tests", "region": "us-west-2", "description": "", "tags": []}'

if [[ $? -ne 0  ]]; then
    echo "Workspace create failed!"
    exit 1
fi

echo

END_TIME=$(( $(date +%s) + 1800 ))
success=false
while [[ $(date +%s) -lt $END_TIME ]]; do
    #check status in loop until created
    if [[ `checkWorkspaceStatus` == '"Workspace-Active"' ]]; then
	success=true
        break
    else 
        echo "waiting for workspace to be created, sleeping for 60 seconds..."
	sleep 60
    fi
done

if [[ $success == "false" ]]; then
    echo "Workspace create timed out!"
    checkWorkspaceStatus
    exit 1
fi


#create the needed single pack
curl -s -X POST "https://tfprovider2-beautiful-nguyen-y8y4azd.cribl-playground.cloud/api/v1/m/default/packs" \
     -H 'Content-Type: application/json' \
     --header "Authorization: Bearer $BEARER_TOKEN" \
     -d '{"version": "0.0.1", "tags": {"streamtags": []}, "exports": [], "displayName": "HelloPacks", "id": "HelloPacks"}'

if [[ $? -ne 0  ]]; then
    echo "Workspace create failed!"
    exit 1
fi
