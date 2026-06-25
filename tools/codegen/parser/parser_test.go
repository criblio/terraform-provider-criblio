package parser

import (
	"path/filepath"
	"testing"
)

func TestParseCertificateResource(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	certificate := resourceByName(t, resources, "certificate")
	if certificate.Create.Path != "/m/{groupId}/system/certificates" {
		t.Fatalf("create path = %q", certificate.Create.Path)
	}
	if certificate.Read.Path != "/m/{groupId}/system/certificates/{id}" {
		t.Fatalf("read path = %q", certificate.Read.Path)
	}
	if len(certificate.Create.Examples) != 1 {
		t.Fatalf("create example count = %d", len(certificate.Create.Examples))
	}
	if certificate.Create.Examples[0].Name != "CertificateCreateExample" {
		t.Fatalf("create example name = %q", certificate.Create.Examples[0].Name)
	}

	groupID := fieldByTFName(t, certificate.Fields, "group_id")
	if !groupID.Required || !groupID.ForceNew || !groupID.PathParam {
		t.Fatalf("group_id flags = required:%v forceNew:%v pathParam:%v", groupID.Required, groupID.ForceNew, groupID.PathParam)
	}

	id := fieldByTFName(t, certificate.Fields, "id")
	if !id.Required || !id.ForceNew || !id.SuppressDiff {
		t.Fatalf("id flags = required:%v forceNew:%v suppressDiff:%v", id.Required, id.ForceNew, id.SuppressDiff)
	}

	privKey := fieldByTFName(t, certificate.Fields, "priv_key")
	if !privKey.Sensitive || !privKey.PreferState || privKey.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("priv_key flags = sensitive:%v prefer:%v strategy:%q", privKey.Sensitive, privKey.PreferState, privKey.ApplyStrategy)
	}

	passphrase := fieldByTFName(t, certificate.Fields, "passphrase")
	if !passphrase.Sensitive || !passphrase.PreferState || passphrase.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("passphrase flags = sensitive:%v prefer:%v strategy:%q", passphrase.Sensitive, passphrase.PreferState, passphrase.ApplyStrategy)
	}

	expiry := fieldByTFName(t, certificate.Fields, "cert_expiry_date")
	if !expiry.Computed || expiry.Optional || expiry.Required {
		t.Fatalf("cert_expiry_date flags = computed:%v optional:%v required:%v", expiry.Computed, expiry.Optional, expiry.Required)
	}

	inUse := fieldByTFName(t, certificate.Fields, "in_use")
	if !inUse.Computed || inUse.Type != "array" {
		t.Fatalf("in_use computed/type = %v/%q", inUse.Computed, inUse.Type)
	}

	args := fieldByTFName(t, certificate.Fields, "args")
	if args.Type != "array" || args.ElementType != "object" || args.ApplyStrategy != "preferState" || !args.SuppressDiff {
		t.Fatalf("args type/element/strategy/suppressDiff = %q/%q/%q/%v", args.Type, args.ElementType, args.ApplyStrategy, args.SuppressDiff)
	}
	if args.NestedModelName != "CertificateArgsModel" || args.NestedAPIModelName != "CertificateArgsAPIModel" {
		t.Fatalf("args nested model names = %q/%q", args.NestedModelName, args.NestedAPIModelName)
	}
	if len(args.Fields) != 2 {
		t.Fatalf("args nested field count = %d", len(args.Fields))
	}
	argName := fieldByTFName(t, args.Fields, "name")
	if !argName.Required || argName.Type != "string" || !argName.SuppressDiff {
		t.Fatalf("args.name required/type/suppressDiff = %v/%q/%v", argName.Required, argName.Type, argName.SuppressDiff)
	}

	renamed := fieldByTFName(t, certificate.Fields, "display_name")
	if renamed.APIName != "displayName" {
		t.Fatalf("renamed APIName = %q", renamed.APIName)
	}

	custom := fieldByTFName(t, certificate.Fields, "conf")
	if custom.CustomType != "jsontypes.NormalizedType{}" {
		t.Fatalf("custom type = %q", custom.CustomType)
	}

	if hasField(certificate.Fields, "ignored") {
		t.Fatalf("ignored field should not be present")
	}
}

func TestParseOneOfVariants(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	destination := resourceByName(t, resources, "destination")
	if len(destination.OneOfVariants) != 3 {
		t.Fatalf("variant count = %d", len(destination.OneOfVariants))
	}
	azure := variantByName(t, destination.OneOfVariants, "OutputAzureBlob")
	accountKey := fieldByTFName(t, azure.Fields, "account_key")
	if accountKey.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("account_key strategy = %q", accountKey.ApplyStrategy)
	}
}

func TestParseMappingRulesetBackwardCompatibleDefaults(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	mappingRuleset := resourceByName(t, resources, "mapping_ruleset")
	id := fieldByTFName(t, mappingRuleset.Fields, "id")
	if id.Required || !id.Optional || !id.Computed {
		t.Fatalf("mapping ruleset id flags = required:%v optional:%v computed:%v", id.Required, id.Optional, id.Computed)
	}

	conf := fieldByTFName(t, mappingRuleset.Fields, "conf")
	functions := fieldByTFName(t, conf.Fields, "functions")
	functionID := fieldByTFName(t, functions.Fields, "id")
	if functionID.Required || !functionID.Optional || functionID.Computed {
		t.Fatalf("mapping function id flags = required:%v optional:%v computed:%v", functionID.Required, functionID.Optional, functionID.Computed)
	}
	final := fieldByTFName(t, functions.Fields, "final")
	if final.Required || !final.Optional || final.Computed {
		t.Fatalf("mapping function final flags = required:%v optional:%v computed:%v", final.Required, final.Optional, final.Computed)
	}
	if final.FixedValue != "" {
		t.Fatalf("mapping function final fixed value = %q", final.FixedValue)
	}
	functionConf := fieldByTFName(t, functions.Fields, "conf")
	add := fieldByTFName(t, functionConf.Fields, "add")
	name := fieldByTFName(t, add.Fields, "name")
	if name.Required || !name.Optional || name.Computed {
		t.Fatalf("mapping add name flags = required:%v optional:%v computed:%v", name.Required, name.Optional, name.Computed)
	}
	if name.FixedValue != "" {
		t.Fatalf("mapping add name fixed value = %q", name.FixedValue)
	}
}

func TestSearchResourcesHideGroupIDForBackwardCompatibility(t *testing.T) {
	resource := &ResourceDef{
		TypeName: "criblio_search_saved_query",
		Fields: []FieldDef{
			{
				APIName:       "groupId",
				TerraformName: "group_id",
				GoName:        "GroupID",
				Type:          "string",
				Required:      true,
				ForceNew:      true,
				PathParam:     true,
			},
		},
	}

	applyResourceCompatibility(resource)

	if hasField(resource.Fields, "group_id") {
		t.Fatalf("search group_id should be hidden from Terraform schema")
	}
}

func TestNotificationHidesGroupIDForBackwardCompatibility(t *testing.T) {
	resource := &ResourceDef{
		StructName: "Notification",
		TypeName:   "criblio_notification",
		Fields: []FieldDef{
			{
				APIName:       "groupId",
				TerraformName: "group_id",
				GoName:        "GroupID",
				Type:          "string",
				Required:      true,
				ForceNew:      true,
				PathParam:     true,
			},
		},
	}

	applyResourceCompatibility(resource)

	if hasField(resource.Fields, "group_id") {
		t.Fatalf("notification group_id should be hidden from Terraform schema")
	}
}

func TestParseFixedSingletonIdentity(t *testing.T) {
	resources, err := Parse([]byte(`
openapi: 3.1.0
paths:
  /m/{groupId}/routes/{id}:
    patch:
      x-terraform-resource: true
      x-terraform-resource-name: Routes
      parameters:
        - name: groupId
          in: path
          required: true
          schema:
            type: string
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Routes"
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Routes"
components:
  schemas:
    Routes:
      type: object
      required:
        - id
        - routes
      properties:
        id:
          type: string
          const: default
          x-terraform-fixed-value: default
        routes:
          type: array
          items:
            type: string
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	routes := resourceByName(t, resources, "Routes")
	id := fieldByTFName(t, routes.Fields, "id")
	if id.FixedValue != "default" {
		t.Fatalf("id fixed value = %q", id.FixedValue)
	}
	if id.Required || !id.Optional || !id.Computed || !id.PathParam || !id.ForceNew {
		t.Fatalf("id flags = required:%v optional:%v computed:%v path:%v force:%v", id.Required, id.Optional, id.Computed, id.PathParam, id.ForceNew)
	}
}

func TestParseKeyQueryID(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	key := resourceByName(t, resources, "Key")
	if key.TypeName != "criblio_key" {
		t.Fatalf("key TypeName = %q", key.TypeName)
	}
	if len(key.Create.QueryParams) != 1 || key.Create.QueryParams[0].TerraformName != "id" {
		t.Fatalf("key query params = %#v", key.Create.QueryParams)
	}
	id := fieldByTFName(t, key.Fields, "id")
	if id.APIName != "keyId" || !id.Required || !id.ForceNew || !id.PathParam || !id.QueryParam || !id.SuppressDiff || !id.PreferState {
		t.Fatalf("key id = api:%q required:%v forceNew:%v path:%v query:%v suppress:%v prefer:%v", id.APIName, id.Required, id.ForceNew, id.PathParam, id.QueryParam, id.SuppressDiff, id.PreferState)
	}
	keyID := fieldByTFName(t, key.Fields, "key_id")
	if keyID.APIName != "terraformKeyId" || !keyID.Computed {
		t.Fatalf("key_id = api:%q computed:%v", keyID.APIName, keyID.Computed)
	}
	algorithm := fieldByTFName(t, key.Fields, "algorithm")
	if !algorithm.Optional || !algorithm.Computed || algorithm.Required {
		t.Fatalf("algorithm flags = required:%v optional:%v computed:%v", algorithm.Required, algorithm.Optional, algorithm.Computed)
	}
	if hasField(key.Fields, "plain_key") || hasField(key.Fields, "cipher_key") {
		t.Fatalf("key material fields should be ignored")
	}
}

func TestParseManagementWorkspace(t *testing.T) {
	resources, err := Parse([]byte(`
openapi: 3.1.0
paths:
  /v1/organizations/{organizationId}/workspaces:
    post:
      x-terraform-resource: true
      x-terraform-resource-name: Workspace
      parameters:
        - name: organizationId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/WorkspaceCreateRequestDTO"
      responses:
        "201":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/WorkspaceSchema"
    get:
      x-terraform-list: Workspace
      x-terraform-list-name: Workspaces
      parameters:
        - name: organizationId
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/WorkspacesListResponseDTO"
  /v1/organizations/{organizationId}/workspaces/{workspaceId}:
    get:
      x-terraform-read: Workspace
      parameters:
        - name: organizationId
          in: path
          required: true
          schema:
            type: string
        - name: workspaceId
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/WorkspaceSchema"
    patch:
      x-terraform-update: Workspace
      x-terraform-read-after-write: true
      parameters:
        - name: organizationId
          in: path
          required: true
          schema:
            type: string
        - name: workspaceId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/WorkspacePatchRequestDTO"
      responses:
        "204":
          description: no body
components:
  schemas:
    WorkspaceCreateRequestDTO:
      type: object
      required: [workspaceId]
      properties:
        workspaceId:
          type: string
          x-terraform-name: workspace_id
        alias:
          type: string
        tags:
          type: array
          items:
            type: string
    WorkspacePatchRequestDTO:
      type: object
      properties:
        alias:
          type: string
        tags:
          type: array
          items:
            type: string
    WorkspaceSchema:
      type: object
      required: [workspaceId, region, state]
      properties:
        workspaceId:
          type: string
          x-terraform-name: workspace_id
        alias:
          type: string
        tags:
          type: array
          items:
            type: string
        region:
          type: string
        state:
          type: string
    WorkspacesListResponseDTO:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: "#/components/schemas/WorkspaceSchema"
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	workspace := resourceByName(t, resources, "Workspace")
	if workspace.Create.ResponseSchema != "WorkspaceSchema" {
		t.Fatalf("create response schema = %q", workspace.Create.ResponseSchema)
	}
	if !workspace.Update.ReadAfterWrite {
		t.Fatalf("workspace update should read after write")
	}
	if workspace.ListName != "Workspaces" || workspace.ListFileStem != "workspaces" || workspace.ListStructName != "Workspaces" {
		t.Fatalf("list metadata = %q/%q/%q", workspace.ListName, workspace.ListFileStem, workspace.ListStructName)
	}
	if len(workspace.List.PathParams) != 1 || workspace.List.PathParams[0].TerraformName != "organization_id" {
		t.Fatalf("list path params = %#v", workspace.List.PathParams)
	}

	workspaceID := fieldByTFName(t, workspace.Fields, "workspace_id")
	if !workspaceID.Required || !workspaceID.PathParam {
		t.Fatalf("workspace_id required/path = %v/%v", workspaceID.Required, workspaceID.PathParam)
	}
	alias := fieldByTFName(t, workspace.Fields, "alias")
	if !alias.RequestField || !alias.UpdateField {
		t.Fatalf("alias request/update = %v/%v", alias.RequestField, alias.UpdateField)
	}
	region := fieldByTFName(t, workspace.Fields, "region")
	if !region.Computed || region.RequestField || region.UpdateField {
		t.Fatalf("region computed/request/update = %v/%v/%v", region.Computed, region.RequestField, region.UpdateField)
	}
}

func TestParseNestedRefsAndObjectOneOf(t *testing.T) {
	resources, err := Parse([]byte(`
openapi: 3.0.0
info:
  title: test
  version: test
paths:
  /m/{groupId}/settings:
    patch:
      x-terraform-resource: true
      x-terraform-resource-name: GroupSystemSettings
      parameters:
        - name: groupId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Settings"
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Settings"
components:
  schemas:
    Settings:
      type: object
      required: [api, backups]
      properties:
        api:
          $ref: "#/components/schemas/API"
        backups:
          $ref: "#/components/schemas/Backups"
        packages:
          type: array
          items:
            $ref: "#/components/schemas/Package"
    API:
      type: object
      properties:
        host:
          type: string
    Backups:
      oneOf:
        - type: object
          properties:
            directory:
              type: string
        - type: object
          properties: {}
    Package:
      allOf:
        - type: object
          properties:
            url:
              type: string
`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	settings := resourceByName(t, resources, "GroupSystemSettings")
	api := fieldByTFName(t, settings.Fields, "api")
	if api.Type != "object" || api.NestedModelName != "GroupSystemSettingsAPIFieldModel" {
		t.Fatalf("api type/nested = %q/%q", api.Type, api.NestedModelName)
	}
	if fieldByTFName(t, api.Fields, "host").Type != "string" {
		t.Fatalf("api.host should be string")
	}

	backups := fieldByTFName(t, settings.Fields, "backups")
	if backups.Type != "object" || fieldByTFName(t, backups.Fields, "directory").Type != "string" {
		t.Fatalf("backups should resolve oneOf object fields")
	}

	packages := fieldByTFName(t, settings.Fields, "packages")
	if packages.Type != "array" || packages.ElementType != "object" || fieldByTFName(t, packages.Fields, "url").Type != "string" {
		t.Fatalf("packages should resolve array item allOf object fields, got type=%q elem=%q fields=%#v", packages.Type, packages.ElementType, packages.Fields)
	}
}

func resourceByName(t *testing.T, resources []ResourceDef, name string) ResourceDef {
	t.Helper()
	for _, resource := range resources {
		if resource.Name == name {
			return resource
		}
	}
	t.Fatalf("resource %q not found", name)
	return ResourceDef{}
}

func fieldByTFName(t *testing.T, fields []FieldDef, name string) FieldDef {
	t.Helper()
	for _, field := range fields {
		if field.TerraformName == name {
			return field
		}
	}
	t.Fatalf("field %q not found", name)
	return FieldDef{}
}

func variantByName(t *testing.T, variants []OneOfVariantDef, name string) OneOfVariantDef {
	t.Helper()
	for _, variant := range variants {
		if variant.APIName == name {
			return variant
		}
	}
	t.Fatalf("variant %q not found", name)
	return OneOfVariantDef{}
}

func hasField(fields []FieldDef, name string) bool {
	for _, field := range fields {
		if field.TerraformName == name {
			return true
		}
	}
	return false
}
