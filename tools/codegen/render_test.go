package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestRendererHonorsCodegenIgnore(t *testing.T) {
	resources := parseFixture(t)
	ignored, err := readIgnoreFile(filepath.Join("testdata", ".codegen-ignore"))
	if err != nil {
		t.Fatalf("readIgnoreFile returned error: %v", err)
	}
	dir := t.TempDir()
	files, err := newRenderer(dir, ignored).render([]parser.ResourceDef{resourceByName(t, resources, "certificate")})
	if err != nil {
		t.Fatalf("render returned error: %v", err)
	}

	if !hasSkipped(files, filepath.Join(dir, "internal/provider/certificate_resource.go")) {
		t.Fatalf("resource output was not skipped")
	}
	if _, err := os.Stat(filepath.Join(dir, "internal/provider/certificate_resource.go")); !os.IsNotExist(err) {
		t.Fatalf("ignored resource file exists or stat returned unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "internal/provider/certificate_types.go")); err != nil {
		t.Fatalf("types output missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tests/acceptance/certificate_test.go")); err != nil {
		t.Fatalf("acceptance test output missing: %v", err)
	}
	if hasSkipped(files, filepath.Join(dir, "tests/acceptance/certificate_test.go")) {
		t.Fatalf("acceptance test output should not be skipped by the ignore fixture")
	}
}

func TestRendererSkipsCustomAcceptanceTests(t *testing.T) {
	resources := parseFixture(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "tests/acceptance/certificate_test.go")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create acceptance test directory: %v", err)
	}
	content := []byte("package tests\n\nfunc TestCustomCertificate(t *testing.T) {}\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write custom acceptance test: %v", err)
	}

	files, err := newRenderer(dir, ignoreSet{}).render([]parser.ResourceDef{resourceByName(t, resources, "certificate")})
	if err != nil {
		t.Fatalf("render returned error: %v", err)
	}

	if !hasSkipped(files, path) {
		t.Fatalf("custom acceptance test output was not skipped")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read custom acceptance test: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("custom acceptance test was overwritten")
	}
}

func TestRendererOverwritesGeneratedAcceptanceTests(t *testing.T) {
	resources := parseFixture(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "tests/acceptance/certificate_test.go")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create acceptance test directory: %v", err)
	}
	content := []byte("// " + generatedHeader + "\n\npackage tests\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write generated acceptance test: %v", err)
	}

	files, err := newRenderer(dir, ignoreSet{}).render([]parser.ResourceDef{resourceByName(t, resources, "certificate")})
	if err != nil {
		t.Fatalf("render returned error: %v", err)
	}

	if hasSkipped(files, path) {
		t.Fatalf("generated acceptance test output should not be skipped")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated acceptance test: %v", err)
	}
	if !strings.Contains(string(got), generatedHeader) {
		t.Fatalf("generated acceptance test header missing")
	}
	if string(got) == string(content) {
		t.Fatalf("generated acceptance test was not overwritten")
	}
}

func TestRenderedSnippets(t *testing.T) {
	resources := parseFixture(t)
	certificate := resourceByName(t, resources, "certificate")
	resourceContent := renderTemplate(t, "resource", certificate)
	assertContains(t, resourceContent, "applyCertificateAPIToState(apiModel, &model, true, false)")
	assertContains(t, resourceContent, "applyCertificateAPIToState(apiModel, &model, true, isCertificateImportState(&model))")
	assertContains(t, resourceContent, "apiModel, err := r.api.Read(ctx, model)")
	assertContains(t, resourceContent, "applyCertificateAPIToState(apiModel, &model, false, false)")
	assertContains(t, resourceContent, "if !preserveInputs || (fillMissingInputs && (state.Cert.IsNull() || state.Cert.IsUnknown()))")
	assertContains(t, resourceContent, "api.DisplayName.IsNull()")
	assertContains(t, resourceContent, "if !api.InUse.IsNull() && !api.InUse.IsUnknown()")
	assertContains(t, resourceContent, "stringFromAPIOrPrior(api.Passphrase.ValueString(), state.Passphrase)")
	assertContains(t, resourceContent, "stringplanmodifier.RequiresReplaceIfConfigured()")
	assertContains(t, resourceContent, "custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "custom_listplanmodifier.SuppressDiff(custom_listplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "custom_objectplanmodifier.SuppressDiff(custom_objectplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "state.InUse = types.ListValueMust(types.StringType, nil)")
	assertContains(t, resourceContent, "if elementType := state.InUse.ElementType(context.Background()); elementType == nil")
	assertContains(t, resourceContent, "state.Args = types.ListNull(types.ObjectType{AttrTypes: CertificateArgsAttrTypes()})")
	assertContains(t, resourceContent, "state.Args = types.ListValueMust(types.ObjectType{AttrTypes: CertificateArgsAttrTypes()}, nil)")
	assertContains(t, resourceContent, "len(state.Args.Elements()) == 0")
	assertContains(t, resourceContent, "clients, ok := req.ProviderData.(*ProviderClients)")
	assertContains(t, resourceContent, "r.client = clients.RC")
	assertContains(t, resourceContent, `json:"group_id"`)
	assertContains(t, resourceContent, `path.Root("group_id")`)
	assertNotContains(t, resourceContent, "speakeasy_")
	assertNotContains(t, resourceContent, "internal/sdk")

	clientContent := renderTemplate(t, "client", certificate)
	assertContains(t, clientContent, "return deleteCertificateSoft(ctx, a.client, model)")

	typesContent := renderTemplate(t, "types", certificate)
	assertContains(t, typesContent, "Conf jsontypes.Normalized")
	assertContains(t, typesContent, "InUse types.List")
	assertContains(t, typesContent, "InUse []types.String")
	assertContains(t, typesContent, "func (m CertificateModel) MarshalJSON()")
	assertContains(t, typesContent, "func CertificateTerraformNameToAPIName(name string) string")
	assertContains(t, typesContent, "output[CertificateTerraformNameToAPIName(key)] = value")
	assertContains(t, typesContent, "func CertificateAPIValueToTerraformValue(value any, typ attr.Type) (attr.Value, error)")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.StringType, input.InUse)")

	dataSourceContent := renderTemplate(t, "data_source", certificate)
	assertContains(t, dataSourceContent, "applyCertificateAPIToState(apiModel, &model, false, false)")

	docContent := renderTemplate(t, "doc", certificate)
	assertContains(t, docContent, `resource "criblio_certificate" "my_certificate"`)
	assertContains(t, docContent, `display_name = "Upstream Certificate"`)
	assertContains(t, docContent, `group_id = "default"`)
	assertContains(t, docContent, `name = "precision"`)

	destination := resourceByName(t, resources, "destination")
	destinationTypes := renderTemplate(t, "types", destination)
	assertContains(t, destinationTypes, "OutputAzureBlob *OutputAzureBlobModel")
	assertContains(t, destinationTypes, "OutputElastic *OutputElasticModel")
	assertContains(t, destinationTypes, "OutputS3 *OutputS3Model")

	destinationResource := renderTemplate(t, "resource", destination)
	assertContains(t, destinationResource, "if api.OutputAzureBlob != nil")
	assertContains(t, destinationResource, "state.OutputAzureBlob = &OutputAzureBlobModel{}")
	assertContains(t, destinationResource, "stringFromAPIOrPrior(api.OutputAzureBlob.AccountKey.ValueString(), state.OutputAzureBlob.AccountKey)")

	mappingRuleset := resourceByName(t, resources, "mapping_ruleset")
	mappingRulesetResource := renderTemplate(t, "resource", mappingRuleset)
	assertContains(t, mappingRulesetResource, `"conf": schema.SingleNestedAttribute{`)
	assertContains(t, mappingRulesetResource, `"add": schema.ListNestedAttribute{`)
	assertContains(t, mappingRulesetResource, `"name": schema.StringAttribute{`)
	assertContains(t, mappingRulesetResource, `"value": schema.StringAttribute{`)
	assertContains(t, mappingRulesetResource, `"id": schema.StringAttribute{`)
	assertContains(t, mappingRulesetResource, "Optional: true,")

	mappingRulesetDataSource := renderTemplate(t, "data_source", mappingRuleset)
	assertContains(t, mappingRulesetDataSource, `"conf": schema.SingleNestedAttribute{`)
	assertContains(t, mappingRulesetDataSource, `"labels": schema.ListAttribute{`)
	assertContains(t, mappingRulesetDataSource, "ElementType: types.StringType,")
	assertContains(t, mappingRulesetDataSource, `"add": schema.ListNestedAttribute{`)
	assertContains(t, mappingRulesetDataSource, `"name": schema.StringAttribute{`)
	assertContains(t, mappingRulesetDataSource, `"value": schema.StringAttribute{`)

	mappingRulesetTypes := renderTemplate(t, "types", mappingRuleset)
	assertContains(t, mappingRulesetTypes, "func mappingRulesetID(model MappingRulesetModel) string")
	assertContains(t, mappingRulesetTypes, `return "default"`)
	assertContains(t, mappingRulesetTypes, `function["id"] = "eval"`)
	assertContains(t, mappingRulesetTypes, `function["final"] = true`)
	assertContains(t, mappingRulesetTypes, `add["name"] = "groupId"`)
	mappingRulesetDoc := renderTemplate(t, "doc", mappingRuleset)
	assertContains(t, mappingRulesetDoc, "- `final` (Boolean)")
	assertContains(t, mappingRulesetDoc, "- `name` (String)")

	key := resourceByName(t, resources, "Key")
	keyClient := renderTemplate(t, "client", key)
	assertContains(t, keyClient, `"net/url"`)
	assertContains(t, keyClient, `id := model.ID.ValueString()`)
	assertContains(t, keyClient, `fmt.Sprintf("/m/%s/system/keys?id=%s", model.GroupID.ValueString(), url.QueryEscape(id))`)
	assertContains(t, keyClient, `return normalizeKeyAPIModel(apiModel, id), err`)
	assertContains(t, keyClient, `restclient.Get[[]KeyModel](ctx, a.client, fmt.Sprintf("/m/%s/system/keys", model.GroupID.ValueString()))`)
	assertContains(t, keyClient, `id := keyAPIID(model)`)
	assertContains(t, keyClient, `if item.ID.ValueString() == id`)
	assertContains(t, keyClient, `model.ID = types.StringValue(apiID)`)
	assertContains(t, keyClient, `The keys API does not support deleting key metadata`)
	assertContains(t, keyClient, `func keyAPIID(model KeyModel) string`)
	assertContains(t, keyClient, `func normalizeKeyAPIModel(model *KeyModel, terraformID string) *KeyModel`)
	keyResource := renderTemplate(t, "resource", key)
	assertContains(t, keyResource, `"algorithm": schema.StringAttribute{`)
	assertContains(t, keyResource, "Optional: true,")
	assertContains(t, keyResource, "Computed: true,")
	assertContains(t, keyResource, "state.Algorithm = types.StringNull()")
	keyTypes := renderTemplate(t, "types", key)
	assertContains(t, keyTypes, `output["algorithm"] = value`)

	jsonString := parser.ResourceDef{
		StructName: "Schema",
		Fields: []parser.FieldDef{
			{
				APIName:       "schema",
				TerraformName: "schema",
				GoName:        "Schema",
				Type:          "string",
				CustomType:    "jsontypes.NormalizedType{}",
				RequestField:  true,
			},
		},
	}
	jsonStringTypes := renderTemplate(t, "types", jsonString)
	assertContains(t, jsonStringTypes, `value := m.Schema.ValueString()`)
	assertContains(t, jsonStringTypes, `output["schema"] = value`)
	assertNotContains(t, jsonStringTypes, `convert schema to API value`)

	noRead := parser.ResourceDef{
		StructName: "LakehouseDatasetConnection",
		NoRead:     true,
		Create: parser.OperationDef{
			Method: "POST",
			Path:   "/products/lake/lakes/default/lakehouses/{lakehouse_id}/datasets/{lake_dataset_id}/connections",
			PathParams: []parser.FieldDef{
				{APIName: "lakehouse_id", GoName: "LakehouseID"},
				{APIName: "lake_dataset_id", GoName: "LakeDatasetID"},
			},
		},
		Update: parser.OperationDef{
			Method: "PATCH",
			Path:   "/products/lake/lakes/default/lakehouses/{lakehouse_id}/datasets/{lake_dataset_id}/connections",
			PathParams: []parser.FieldDef{
				{APIName: "lakehouse_id", GoName: "LakehouseID"},
				{APIName: "lake_dataset_id", GoName: "LakeDatasetID"},
			},
		},
	}
	noReadClient := renderTemplate(t, "client", noRead)
	assertContains(t, noReadClient, "restclient.PostNoResponse(ctx, a.client")
	assertContains(t, noReadClient, "restclient.PatchNoResponse(ctx, a.client")
	assertNotContains(t, noReadClient, "response envelope")

	mappingRulesetResource = renderTemplate(t, "resource", mappingRuleset)
	assertContains(t, mappingRulesetResource, `state.Conf = types.ObjectNull(MappingRulesetConfAttrTypes())`)

	fixedIdentity := parser.ResourceDef{
		StructName: "Routes",
		TypeName:   "criblio_routes",
		Create: parser.OperationDef{
			Method: "PATCH",
			Path:   "/m/{groupId}/routes/{id}",
		},
		Read: parser.OperationDef{
			Method: "GET",
			Path:   "/m/{groupId}/routes/{id}",
		},
		Update: parser.OperationDef{
			Method: "PATCH",
			Path:   "/m/{groupId}/routes/{id}",
		},
		Delete: parser.OperationDef{
			Method: "PATCH",
			Path:   "/m/{groupId}/routes/{id}",
			ResetBody: map[string]any{
				"id": "default",
				"routes": []any{
					map[string]any{
						"name":     "Default",
						"filter":   "true",
						"pipeline": "main",
						"final":    true,
						"disabled": false,
					},
				},
			},
		},
		Fields: []parser.FieldDef{
			{
				APIName:       "id",
				TerraformName: "id",
				GoName:        "ID",
				Type:          "string",
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				PathParam:     true,
				FixedValue:    "default",
			},
		},
	}
	fixedResource := renderTemplate(t, "resource", fixedIdentity)
	assertContains(t, fixedResource, `"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"`)
	assertContains(t, fixedResource, `"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"`)
	assertContains(t, fixedResource, `Default: stringdefault.StaticString("default")`)
	assertContains(t, fixedResource, `stringvalidator.OneOf("default")`)
	assertContains(t, fixedResource, `state.ID = types.StringValue("default")`)
	fixedClient := renderTemplate(t, "client", fixedIdentity)
	assertContains(t, fixedClient, `body := map[string]any{"id": "default", "routes": []any{map[string]any{"disabled": false, "filter": "true", "final": true, "name": "Default", "pipeline": "main"}}}`)
	assertContains(t, fixedClient, `restclient.Patch[map[string]any, RoutesModel]`)
	fixedDataSource := renderTemplate(t, "data_source", fixedIdentity)
	assertContains(t, fixedDataSource, `Optional: true`)
	assertContains(t, fixedDataSource, `Computed: true`)
	assertContains(t, fixedDataSource, `model.ID = types.StringValue("default")`)
	assertNotContains(t, fixedDataSource, `Required: true`)

	packPipeline := parser.ResourceDef{
		StructName: "PackPipeline",
		TypeName:   "criblio_pack_pipeline",
		Create: parser.OperationDef{
			Method: "POST",
			Path:   "/m/{groupId}/p/{pack}/pipelines",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
			},
		},
		Read: parser.OperationDef{
			Method: "GET",
			Path:   "/m/{groupId}/p/{pack}/pipelines/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
				{APIName: "id", TerraformName: "id", GoName: "ID"},
			},
		},
		Update: parser.OperationDef{
			Method:         "PATCH",
			Path:           "/m/{groupId}/p/{pack}/pipelines/{id}",
			ReadAfterWrite: true,
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
				{APIName: "id", TerraformName: "id", GoName: "ID"},
			},
		},
	}
	packPipelineClient := renderTemplate(t, "client", packPipeline)
	assertContains(t, packPipelineClient, `"errors"`)
	assertContains(t, packPipelineClient, `if packPipelineAlreadyExists(err) {`)
	assertContains(t, packPipelineClient, `return a.Update(ctx, model)`)
	assertContains(t, packPipelineClient, `func packPipelineAlreadyExists(err error) bool`)

	packLookups := parser.ResourceDef{
		StructName: "PackLookups",
		TypeName:   "criblio_pack_lookups",
		Create: parser.OperationDef{
			Method: "POST",
			Path:   "/m/{groupId}/p/{pack}/system/lookups",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
			},
		},
		Read: parser.OperationDef{
			Method: "GET",
			Path:   "/m/{groupId}/p/{pack}/system/lookups/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
				{APIName: "id", TerraformName: "id", GoName: "ID"},
			},
		},
		Update: parser.OperationDef{
			Method: "PATCH",
			Path:   "/m/{groupId}/p/{pack}/system/lookups/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
				{APIName: "id", TerraformName: "id", GoName: "ID"},
			},
		},
		Delete: parser.OperationDef{
			Method: "DELETE",
			Path:   "/m/{groupId}/p/{pack}/system/lookups/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
				{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
				{APIName: "id", TerraformName: "id", GoName: "ID"},
			},
		},
	}
	packLookupsClient := renderTemplate(t, "client", packLookups)
	assertContains(t, packLookupsClient, `"net/url"`)
	assertContains(t, packLookupsClient, `if err := uploadPackLookupFileContent(ctx, a.client, model, id); err != nil`)
	assertContains(t, packLookupsClient, `for _, id := range lookupFileAPIIDs(configuredID)`)
	assertContains(t, packLookupsClient, `return normalizePackLookupsAPIModel(apiModel, configuredID), nil`)
	assertContains(t, packLookupsClient, `func uploadPackLookupFileContent(ctx context.Context, client *restclient.Client, model PackLookupsModel, id string) error`)
	assertContains(t, packLookupsClient, `url.QueryEscape(filename)`)

	fixedAPIField := parser.ResourceDef{
		Name:     "searchengine",
		FileStem: "search_engine",
		TypeName: "criblio_search_engine",
		Fields: []parser.FieldDef{
			{
				APIName:       "engineType",
				TerraformName: "engine_type",
				Type:          "string",
				Optional:      true,
				Computed:      true,
				FixedValue:    "local",
			},
			{APIName: "id", TerraformName: "id", Type: "string", Required: true},
			{APIName: "tierSize", TerraformName: "tier_size", Type: "string", Required: true},
		},
		Create: parser.OperationDef{
			Examples: []parser.ExampleDef{
				{
					Name: "minimal",
					Value: map[string]any{
						"engineType": "local",
						"id":         "engine-01",
						"tierSize":   "small",
					},
				},
			},
		},
	}
	assertNotContains(t, docSchema(fixedAPIField), "engine_type")
	fixedAPIFieldExample, ok := upstreamExampleUsage(fixedAPIField)
	if !ok {
		t.Fatalf("upstreamExampleUsage returned no example")
	}
	assertNotContains(t, fixedAPIFieldExample, "engine_type")
	assertContains(t, fixedAPIFieldExample, `id = "engine-01"`)
}

func TestUpstreamExampleUsagePrefersRichestExample(t *testing.T) {
	resource := parser.ResourceDef{
		Name:     "searchmacro",
		FileStem: "search_macro",
		TypeName: "criblio_search_macro",
		Create: parser.OperationDef{
			Examples: []parser.ExampleDef{
				{
					Name: "minimal",
					Value: map[string]any{
						"id":          "all_events",
						"replacement": "true",
					},
				},
				{
					Name: "full",
					Value: map[string]any{
						"id":          "error_filter",
						"description": "Filters to high-severity events.",
						"replacement": `severity >= "Error"`,
						"tags":        "errors,prod",
					},
				},
			},
		},
		Fields: []parser.FieldDef{
			{APIName: "description", TerraformName: "description", Type: "string", Optional: true},
			{APIName: "groupId", TerraformName: "group_id", Type: "string", Required: true, PathParam: true},
			{APIName: "id", TerraformName: "id", Type: "string", Required: true},
			{APIName: "replacement", TerraformName: "replacement", Type: "string", Required: true},
			{APIName: "tags", TerraformName: "tags", Type: "string", Optional: true},
		},
	}

	got, ok := upstreamExampleUsage(resource)
	if !ok {
		t.Fatalf("upstreamExampleUsage returned no example")
	}
	assertContains(t, got, `id = "error_filter"`)
	assertContains(t, got, `description = "Filters to high-severity events."`)
	assertContains(t, got, `replacement = "severity >= \"Error\""`)
	assertContains(t, got, `tags = "errors,prod"`)
	assertContains(t, got, `group_id = "default_search"`)
	assertNotContains(t, got, `id = "all_events"`)
}

func TestUpstreamExampleUsageEscapesTerraformInterpolation(t *testing.T) {
	resource := parser.ResourceDef{
		Name:     "destination",
		FileStem: "destination",
		TypeName: "criblio_destination",
		Create: parser.OperationDef{
			Examples: []parser.ExampleDef{
				{
					Name: "s3",
					Value: map[string]any{
						"id":       "out-s3-main",
						"destPath": "`logs/${C.Time.strftime(_time, '%Y/%m/%d')}`",
					},
				},
			},
		},
		Fields: []parser.FieldDef{
			{APIName: "groupId", TerraformName: "group_id", Type: "string", Required: true, PathParam: true},
			{APIName: "id", TerraformName: "id", Type: "string", Required: true},
			{APIName: "destPath", TerraformName: "dest_path", Type: "string", Optional: true},
		},
	}

	got, ok := upstreamExampleUsage(resource)
	if !ok {
		t.Fatalf("upstreamExampleUsage returned no example")
	}
	assertContains(t, got, "dest_path = \"`logs/$${C.Time.strftime(_time, '%Y/%m/%d')}`\"")
	assertNotContains(t, got, "dest_path = \"`logs/${C.Time.strftime")
}

func TestPathExprOrdersParamsByPathTemplate(t *testing.T) {
	resource := parser.ResourceDef{
		StructName: "PackBreakers",
		TypeName:   "criblio_pack_breakers",
	}
	op := parser.OperationDef{
		Path: "/m/{groupId}/p/{pack}/lib/breakers/{id}",
		PathParams: []parser.FieldDef{
			{APIName: "pack", TerraformName: "pack", GoName: "Pack"},
			{APIName: "id", TerraformName: "id", GoName: "ID"},
			{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID"},
		},
	}

	got := pathExpr(resource, op)
	assertContains(t, got, `fmt.Sprintf("/m/%s/p/%s/lib/breakers/%s", model.GroupID.ValueString(), resolvePackIDForRestAPI(ctx, a.client, model.GroupID.ValueString(), model.Pack.ValueString()), model.ID.ValueString())`)
}

func TestExampleUsageUsesCuratedSourceExample(t *testing.T) {
	resource := parser.ResourceDef{
		Name:     "source",
		FileStem: "source",
		TypeName: "criblio_source",
		OneOfVariants: []parser.OneOfVariantDef{
			{GoName: "InputHttp"},
		},
	}

	got := exampleUsage(resource)
	assertContains(t, got, `resource "criblio_source" "my_http_source"`)
	assertContains(t, got, `resource "criblio_source" "my_source"`)
	assertNotContains(t, got, `input_appscope = {`)
}

func TestExampleUsagePrefersCuratedPackExample(t *testing.T) {
	resource := parser.ResourceDef{
		Name:     "packroutes",
		FileStem: "pack_routes",
		TypeName: "criblio_pack_routes",
		Create: parser.OperationDef{
			Examples: []parser.ExampleDef{
				{
					Name: "upstream",
					Value: map[string]any{
						"groupId": "Cribl",
						"pack":    "observability-pack",
					},
				},
			},
		},
	}

	got := exampleUsage(resource)
	assertContains(t, got, `resource "criblio_pack" "routes_pack"`)
	assertContains(t, got, `group_id = "default"`)
	assertNotContains(t, got, `group_id = "Cribl"`)
}

func TestGeneratedImportUsesPathParams(t *testing.T) {
	resource := parser.ResourceDef{
		FileStem: "lakehouse_dataset_connection",
		TypeName: "criblio_lakehouse_dataset_connection",
		Fields: []parser.FieldDef{
			{TerraformName: "lakehouse_id", PathParam: true},
			{TerraformName: "lake_dataset_id", PathParam: true},
		},
	}

	block := generatedImportBlock(resource)
	assertContains(t, block, `lakehouse_id = "lakehouse-01"`)
	assertContains(t, block, `lake_dataset_id = "web-logs"`)
	assertNotContains(t, block, `group_id`)
	assertNotContains(t, block, `cert-001`)

	command := generatedImportCommand(resource)
	assertContains(t, command, `"lakehouse_id": "lakehouse-01"`)
	assertContains(t, command, `"lake_dataset_id": "web-logs"`)
	assertNotContains(t, command, `cert-001`)
}

func TestSearchResourcePathUsesInternalDefaultSearchGroup(t *testing.T) {
	resource := parser.ResourceDef{
		Name:       "search_saved_query",
		FileStem:   "search_saved_query",
		TypeName:   "criblio_search_saved_query",
		StructName: "SearchSavedQuery",
		Create: parser.OperationDef{
			Path: "/m/{groupId}/search/saved",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true},
			},
		},
		Read: parser.OperationDef{
			Path: "/m/{groupId}/search/saved/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true},
				{APIName: "id", TerraformName: "id", GoName: "ID", Type: "string", PathParam: true},
			},
		},
		Fields: []parser.FieldDef{
			{APIName: "id", TerraformName: "id", GoName: "ID", Type: "string", Required: true},
		},
	}

	content, err := executeTemplate("client", resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	got := string(content)

	assertNotContains(t, got, `"group_id": schema.StringAttribute{`)
	assertNotContains(t, got, `model.GroupID.ValueString()`)
	assertContains(t, got, `fmt.Sprintf("/m/%s/search/saved", "default_search")`)
	assertContains(t, got, `fmt.Sprintf("/m/%s/search/saved/%s", "default_search", model.ID.ValueString())`)
}

func TestNotificationResourcePathUsesCompatibleGroupFallback(t *testing.T) {
	resource := parser.ResourceDef{
		Name:       "notification",
		FileStem:   "notification",
		TypeName:   "criblio_notification",
		StructName: "Notification",
		Create: parser.OperationDef{
			Path: "/m/{groupId}/notifications",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true},
			},
		},
		Read: parser.OperationDef{
			Path: "/m/{groupId}/notifications/{id}",
			PathParams: []parser.FieldDef{
				{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true},
				{APIName: "id", TerraformName: "id", GoName: "ID", Type: "string", PathParam: true},
			},
		},
		Fields: []parser.FieldDef{
			{APIName: "group", TerraformName: "group", GoName: "Group", Type: "string", Optional: true, Computed: true},
			{APIName: "id", TerraformName: "id", GoName: "ID", Type: "string", Required: true, PathParam: true},
		},
	}

	content, err := executeTemplate("client", resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	got := string(content)

	assertContains(t, got, `fmt.Sprintf("/m/%s/notifications", notificationGroupID(model))`)
	assertContains(t, got, `fmt.Sprintf("/m/%s/notifications/%s", notificationGroupID(model), model.ID.ValueString())`)
	assertNotContains(t, got, `model.GroupID.ValueString()`)
	assertContains(t, got, `func notificationGroupID(model NotificationModel) string`)
	assertContains(t, got, `return "default"`)

	resourceContent, err := executeTemplate("resource", resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	got = string(resourceContent)

	assertContains(t, got, `Group string `+"`json:\"group\"`")
	assertContains(t, got, `GroupID string `+"`json:\"group_id\"`")
	assertContains(t, got, `model.Group = types.StringValue(data.Group)`)
	assertNotContains(t, got, `resource.ImportStatePassthroughID`)
}

func TestResourceSchemaUsesPlanModifierHook(t *testing.T) {
	resource := parser.ResourceDef{
		Name:       "collector",
		FileStem:   "collector",
		TypeName:   "criblio_collector",
		StructName: "Collector",
		Fields: []parser.FieldDef{
			{
				APIName:          "environment",
				TerraformName:    "environment",
				GoName:           "Environment",
				Type:             "string",
				Optional:         true,
				Computed:         true,
				PlanModifierHook: "collectorEnvironmentPlanModifiers",
			},
		},
	}

	content, err := executeTemplate("resource", resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	got := string(content)

	assertContains(t, got, `PlanModifiers: collectorEnvironmentPlanModifiers(),`)
	assertNotContains(t, got, `PlanModifiers: []planmodifier.String{`)
}

func TestPrimitiveArrayFieldsPreserveElementTypes(t *testing.T) {
	resource := parser.ResourceDef{
		Name:       "array_resource",
		FileStem:   "array_resource",
		TypeName:   "criblio_array_resource",
		StructName: "ArrayResource",
		Fields: []parser.FieldDef{
			{APIName: "codes", TerraformName: "codes", GoName: "Codes", Type: "array", ElementType: "integer", Optional: true},
			{APIName: "ratios", TerraformName: "ratios", GoName: "Ratios", Type: "array", ElementType: "number", Optional: true},
			{APIName: "flags", TerraformName: "flags", GoName: "Flags", Type: "array", ElementType: "boolean", Optional: true},
			{APIName: "names", TerraformName: "names", GoName: "Names", Type: "array", ElementType: "string", Optional: true},
		},
	}

	resourceContent := renderTemplate(t, "resource", resource)
	assertContains(t, resourceContent, "ElementType: types.Int64Type,")
	assertContains(t, resourceContent, "ElementType: types.Float64Type,")
	assertContains(t, resourceContent, "ElementType: types.BoolType,")
	assertContains(t, resourceContent, "ElementType: types.StringType,")

	typesContent := renderTemplate(t, "types", resource)
	assertContains(t, typesContent, "Codes []int64")
	assertContains(t, typesContent, "Ratios []float64")
	assertContains(t, typesContent, "Flags []bool")
	assertContains(t, typesContent, "Names []string")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.Int64Type, input.Codes)")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.Float64Type, input.Ratios)")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.BoolType, input.Flags)")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.StringType, input.Names)")
}

func TestObjectAsJSONAndMapNestedFields(t *testing.T) {
	resource := parser.ResourceDef{
		Name:       "pipeline",
		FileStem:   "pipeline",
		TypeName:   "criblio_pipeline",
		StructName: "Pipeline",
		Fields: []parser.FieldDef{
			{
				APIName:       "functions",
				TerraformName: "functions",
				GoName:        "Functions",
				Type:          "array",
				ElementType:   "object",
				Fields: []parser.FieldDef{
					{
						APIName:          "conf",
						TerraformName:    "conf",
						GoName:           "Conf",
						Type:             "object",
						CustomType:       "jsontypes.NormalizedType{}",
						ObjectAsJSON:     true,
						PlanModifierHook: "pipelineConfPlanModifiers",
						NotNull:          true,
						ValidJSON:        true,
						Description:      "Function configuration as JSON.",
					},
					{
						APIName:            "id",
						TerraformName:      "id",
						GoName:             "ID",
						Type:               "string",
						NotNull:            true,
						PipelineFunctionID: true,
						Description:        "Function ID.",
					},
				},
				NestedAttrTypes: "PipelineFunctionsAttrTypes",
			},
			{
				APIName:       "groups",
				TerraformName: "groups",
				GoName:        "Groups",
				Type:          "object",
				ElementType:   "object",
				Fields: []parser.FieldDef{
					{APIName: "name", TerraformName: "name", GoName: "Name", Type: "string", NotNull: true, Description: "Group name"},
				},
				NestedAttrTypes: "PipelineGroupsAttrTypes",
			},
		},
	}

	resourceContent := renderTemplate(t, "resource", resource)
	assertContains(t, resourceContent, `"conf": schema.StringAttribute{`)
	assertContains(t, resourceContent, `CustomType: jsontypes.NormalizedType{},`)
	assertContains(t, resourceContent, `PlanModifiers: pipelineConfPlanModifiers(),`)
	assertContains(t, resourceContent, `Validators: []validator.String{`)
	assertContains(t, resourceContent, `custom_stringvalidators.NotNull(),`)
	assertContains(t, resourceContent, `custom_validators.IsValidJSON(),`)
	assertContains(t, resourceContent, `custom_stringvalidators.IsCriblPipelineFunctionIDWithRestClient(&r.client),`)
	assertContains(t, resourceContent, `"groups": schema.MapNestedAttribute{`)
	assertContains(t, resourceContent, `state.Groups = types.MapNull(types.ObjectType{AttrTypes: PipelineGroupsAttrTypes()})`)
	assertNotContains(t, resourceContent, `state.Groups = types.MapNull(types.StringType)`)
	assertContains(t, resourceContent, `PipelineValueWithKnownNulls(state.Conf, types.ObjectType{AttrTypes: PipelineConfAttrTypes()})`)
	assertNotContains(t, resourceContent, `ElementType:   types.StringType,`)

	docsContent := renderTemplate(t, "doc", resource)
	assertContains(t, docsContent, "`conf` (String)")
	assertContains(t, docsContent, "`id` (String) Function ID. Not Null")
	assertContains(t, docsContent, "`groups` (Attributes Map)")
	assertContains(t, docsContent, "`name` (String) Group name Not Null")

	typesContent := renderTemplate(t, "types", resource)
	assertContains(t, typesContent, `Conf jsontypes.Normalized`)
	assertContains(t, typesContent, `Conf any `+"`json:\"conf,omitempty\"`")
	assertContains(t, typesContent, `types.MapType{ElemType: types.ObjectType{AttrTypes: PipelineGroupsAttrTypes()}}`)
	assertContains(t, typesContent, `value, err = PipelineObjectJSONFromTerraformValue(attribute)`)
	assertContains(t, typesContent, `if typ.Equal(jsontypes.NormalizedType{}) {`)
	assertContains(t, typesContent, `return jsontypes.NewNormalizedValue(string(raw)), nil`)
	assertContains(t, typesContent, `func PipelineValueWithKnownNulls(value attr.Value, typ attr.Type) (attr.Value, error)`)
	assertContains(t, typesContent, `types.MapType{ElemType: types.ObjectType{AttrTypes: PipelineGroupsAttrTypes()}}`)
}

func TestRestWriteCall(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{method: "POST", want: "Post"},
		{method: "PATCH", want: "Patch"},
		{method: "PUT", want: "Put"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			got := restWriteCall(parser.OperationDef{Method: tt.method})
			if got != tt.want {
				t.Fatalf("restWriteCall(%q) = %q, want %q", tt.method, got, tt.want)
			}
		})
	}
}

func parseFixture(t *testing.T) []parser.ResourceDef {
	t.Helper()
	resources, err := parser.ParseFile(filepath.Join("testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}
	return resources
}

func resourceByName(t *testing.T, resources []parser.ResourceDef, name string) parser.ResourceDef {
	t.Helper()
	for _, resource := range resources {
		if resource.Name == name {
			return resource
		}
	}
	t.Fatalf("resource %q not found", name)
	return parser.ResourceDef{}
}

func renderTemplate(t *testing.T, kind string, resource parser.ResourceDef) string {
	t.Helper()
	content, err := executeTemplate(kind, resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	return string(content)
}

func hasSkipped(files []renderedFile, path string) bool {
	for _, file := range files {
		if file.Path == path && file.Skipped {
			return true
		}
	}
	return false
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()
	if !strings.Contains(content, want) {
		t.Fatalf("expected content to contain %q", want)
	}
}

func assertNotContains(t *testing.T, content, want string) {
	t.Helper()
	if strings.Contains(content, want) {
		t.Fatalf("expected content not to contain %q", want)
	}
}
