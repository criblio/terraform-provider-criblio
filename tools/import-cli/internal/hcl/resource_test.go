package hcl

import (
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBlock_simple_attributes(t *testing.T) {
	attrs := map[string]Value{
		"id":       {Kind: KindString, String: "input-1"},
		"group_id": {Kind: KindString, String: "default"},
	}
	block, err := ResourceBlock("criblio_source", "hec", attrs, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, block)
	assert.Equal(t, "resource", block.Type())
	labels := block.Labels()
	require.Len(t, labels, 2)
	assert.Equal(t, "criblio_source", labels[0])
	assert.Equal(t, "hec", labels[1])
}

func TestResourceBlockBytes_generates_valid_hcl(t *testing.T) {
	attrs := map[string]Value{
		"id":       {Kind: KindString, String: "input-1"},
		"group_id": {Kind: KindString, String: "default"},
	}
	bytes, err := ResourceBlockBytes("criblio_source", "hec", attrs, nil)
	require.NoError(t, err)
	require.NotEmpty(t, bytes)
	src := string(bytes)
	assert.Contains(t, src, `resource "criblio_source" "hec"`)
	assert.Contains(t, src, `id`)
	assert.Contains(t, src, `group_id`)

	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err, "generated HCL must parse successfully")
}

func TestResourceBlock_nested_attributes(t *testing.T) {
	attrs := map[string]Value{
		"id":   {Kind: KindString, String: "pipe-1"},
		"conf": {
			Kind: KindMap,
			Map: map[string]Value{
				"description": {Kind: KindString, String: "my pipeline"},
				"async_func_timeout": {Kind: KindNumber, Number: 5000},
			},
		},
	}
	block, err := ResourceBlock("criblio_pipeline", "main", attrs, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, block)
	bytes, _ := ResourceBlockBytes("criblio_pipeline", "main", attrs, nil)
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err, "nested attributes must render and parse")
	assert.Contains(t, string(bytes), "conf")
}

func TestResourceBlock_list_and_null(t *testing.T) {
	attrs := map[string]Value{
		"id":   {Kind: KindString, String: "x"},
		"tags": {Kind: KindList, List: []Value{
			{Kind: KindString, String: "a"},
			{Kind: KindNull},
			{Kind: KindString, String: "b"},
		}},
		"optional": {Kind: KindNull},
	}
	bytes, err := ResourceBlockBytes("criblio_source", "example", attrs, nil)
	require.NoError(t, err)
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), "tags")
}

func TestResourceBlock_skip_null_option(t *testing.T) {
	attrs := map[string]Value{
		"id":       {Kind: KindString, String: "x"},
		"optional": {Kind: KindNull},
	}
	opts := &ResourceBlockOptions{SkipNullAttributes: true}
	bytes, err := ResourceBlockBytes("criblio_source", "example", attrs, opts)
	require.NoError(t, err)
	src := string(bytes)
	assert.Contains(t, src, "id")
	assert.NotContains(t, src, "optional")
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err)
}

func TestResourceBlock_requires_type_and_name(t *testing.T) {
	_, err := ResourceBlock("", "name", map[string]Value{}, nil, nil)
	require.Error(t, err)
	_, err = ResourceBlock("type", "", map[string]Value{}, nil, nil)
	require.Error(t, err)
}

func TestFileWithResources_multiple_blocks(t *testing.T) {
	resources := []ResourceInput{
		{TypeName: "criblio_source", Name: "one", Attrs: map[string]Value{"id": {Kind: KindString, String: "1"}, "group_id": {Kind: KindString, String: "default"}}},
		{TypeName: "criblio_source", Name: "two", Attrs: map[string]Value{"id": {Kind: KindString, String: "2"}, "group_id": {Kind: KindString, String: "default"}}},
	}
	f, err := FileWithResources(resources, nil)
	require.NoError(t, err)
	require.NotNil(t, f)
	bytes := f.Bytes()
	err = ParseHCL(bytes, "main.tf")
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), `"one"`)
	assert.Contains(t, string(bytes), `"two"`)
}

func TestResourceBlock_from_converted_model(t *testing.T) {
	model := &provider.SourceResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("input-hec-1"),
	}
	attrs, err := ModelToValue(model, nil)
	require.NoError(t, err)
	bytes, err := ResourceBlockBytes("criblio_source", "imported", attrs, nil)
	require.NoError(t, err)
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err)
	src := string(bytes)
	assert.True(t, strings.Contains(src, "criblio_source") && strings.Contains(src, "imported"))
	assert.Contains(t, src, "default")
	assert.Contains(t, src, "input-hec-1")
}

func TestFileWithResources_deterministic_output(t *testing.T) {
	attrs := map[string]Value{"id": {Kind: KindString, String: "x"}, "group_id": {Kind: KindString, String: "default"}}
	// Input in different order
	res1 := []ResourceInput{
		{TypeName: "criblio_source", Name: "second", Attrs: attrs},
		{TypeName: "criblio_source", Name: "first", Attrs: attrs},
	}
	res2 := []ResourceInput{
		{TypeName: "criblio_source", Name: "first", Attrs: attrs},
		{TypeName: "criblio_source", Name: "second", Attrs: attrs},
	}
	f1, err := FileWithResources(res1, nil)
	require.NoError(t, err)
	f2, err := FileWithResources(res2, nil)
	require.NoError(t, err)
	assert.Equal(t, string(f1.Bytes()), string(f2.Bytes()), "same resources in different order must produce same output")
	// First resource in output must be "first" (sorted by name)
	assert.Contains(t, string(f1.Bytes()), `"criblio_source" "first"`)
	assert.Contains(t, string(f1.Bytes()), `"criblio_source" "second"`)
}

func TestResourceBlock_lifecycle_ignore_changes(t *testing.T) {
	attrs := map[string]Value{
		"group_id": {Kind: KindString, String: "default"},
		"api":      {Kind: KindMap, Map: map[string]Value{"host": {Kind: KindString, String: ""}}},
	}
	block, err := ResourceBlock("criblio_group_system_settings", "default", attrs, nil, []string{"api"})
	require.NoError(t, err)
	require.NotNil(t, block)
	f := hclwrite.NewEmptyFile()
	f.Body().AppendBlock(block)
	src := string(f.Bytes())
	assert.Contains(t, src, "lifecycle")
	assert.Contains(t, src, "ignore_changes")
	assert.Contains(t, src, "api")
	err = ParseHCL(f.Bytes(), "test.tf")
	assert.NoError(t, err, "lifecycle block must parse")
}

func TestResourceBlock_from_pipeline_model_with_nested_conf(t *testing.T) {
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pipe-1"),
	}
	attrs, err := ModelToValue(model, nil)
	require.NoError(t, err)
	bytes, err := ResourceBlockBytes("criblio_pipeline", "main", attrs, nil)
	require.NoError(t, err)
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err, "pipeline resource with nested conf must parse")
	assert.Contains(t, string(bytes), `resource "criblio_pipeline" "main"`)
	assert.Contains(t, string(bytes), "conf")
}

func TestResourceBlock_pack_tags_preserves_domain_with_skip_empty_lists(t *testing.T) {
	// criblio_pack tags require "domain" when tags is set; PruneEmptyLists must not remove domain=[]
	attrs := map[string]Value{
		"id":       {Kind: KindString, String: "billing_pipeline"},
		"group_id": {Kind: KindString, String: "default"},
		"tags": {
			Kind: KindMap,
			Map: map[string]Value{
				"data_type":  {Kind: KindList, List: []Value{{Kind: KindString, String: "logs"}}},
				"domain":     {Kind: KindList, List: []Value{}}, // required but empty; must not be pruned
				"streamtags": {Kind: KindList, List: []Value{{Kind: KindString, String: "PaloAlto"}}},
				"technology": {Kind: KindList, List: []Value{{Kind: KindString, String: "paloalto"}}},
			},
		},
	}
	opts := DefaultResourceBlockOptions()
	bytes, err := ResourceBlockBytes("criblio_pack", "billing", attrs, opts)
	require.NoError(t, err)
	src := string(bytes)
	assert.Contains(t, src, "domain", "criblio_pack tags must include domain when SkipEmptyListAttributes is true")
	assert.Contains(t, src, "data_type")
	assert.Contains(t, src, "streamtags")
	err = ParseHCL(bytes, "test.tf")
	assert.NoError(t, err, "pack resource with tags must parse")
}
