package main

var serverData = `- url: https://app.cribl.cloud
  x-speakeasy-server-id: cloud
  variables:
    workspaceName:
      default: main
      description: The Workspace name
    organizationId:
      default: ian description: the Organization ID
    cloudDomain:
      default: cribl.cloud
      description: Cribl Cloud domain name
- url: https://{workspaceName}-{organizationId}.{cloudDomain}/api/v1/m/{groupName}
  x-speakeasy-server-id: cloud-group
  variables:
    workspaceName:
      default: main
      description: The Workspace name
    organizationId:
      default: ian
      description: The Organization ID
    cloudDomain:
      default: cribl.cloud
      description: Cribl Cloud domain name
    groupName:
      default: default
      description: The name of the Worker Group or Fleet
- url: https://{hostname}:{port}/api/v1
  x-speakeasy-server-id: managed
  variables:
    hostname:
      default: localhost
      description: The hostname of the managed API server
    port:
      default: '9000'
      description: The port of the managed API server
- url: https://{hostname}:{port}/api/v1/m/{groupName}
  x-speakeasy-server-id: managed-group
  variables:
    hostname:
      default: localhost
      description: The hostname of the managed API server
    port:
      default: '9000'
      description: The port of the managed API server
    groupName:
      default: default
      description: The name of the Worker Group or Fleet`

var schemaPackData = `x-speakeasy-entity: Pack
type: object
properties:
  displayName:
    type: string
  id:
    type: string
  description:
    type: string
  version:
    type: string
  source:
    type: string
    x-speakeasy-xor-with:
      - filename
  disabled:
    type: boolean
required:
  - id`

//this map is intended to add a single key and value to existing yaml structures
//in '"foo.bar.biz" = "yaml"'; foo.bar must exist!
var pathSpeakeasyOperation = map[string]string{
	"paths./products/{product}/groups.post.x-speakeasy-entity-operation": "Group#create",
	"paths./m/{groupId}/lib/appscope-configs/{id}.patch.x-speakeasy-entity-operation": "Group#create",
}

