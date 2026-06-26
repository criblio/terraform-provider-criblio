package provider

import (
	"encoding/json"
	"testing"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTagsAPIMapUsesEmptyArrays(t *testing.T) {
	tags := tagsAPIMap(&tfTypes.PackRequestBodyTags{})

	body, err := json.Marshal(map[string]any{
		"package": map[string]any{
			"tags": tags,
		},
	})
	if err != nil {
		t.Fatalf("marshal tags body: %v", err)
	}

	want := `{"package":{"tags":{"dataType":[],"domain":[],"streamtags":[],"technology":[]}}}`
	if string(body) != want {
		t.Fatalf("unexpected tags body\nwant: %s\n got: %s", want, string(body))
	}
}

func TestTagsAPIMapPreservesConfiguredTags(t *testing.T) {
	tags := tagsAPIMap(&tfTypes.PackRequestBodyTags{
		DataType:   []types.String{types.StringValue("logs")},
		Domain:     []types.String{types.StringValue("security")},
		Streamtags: []types.String{types.StringValue("prod")},
		Technology: []types.String{types.StringValue("cribl")},
	})

	body, err := json.Marshal(map[string]any{
		"package": map[string]any{
			"tags": tags,
		},
	})
	if err != nil {
		t.Fatalf("marshal tags body: %v", err)
	}

	want := `{"package":{"tags":{"dataType":["logs"],"domain":["security"],"streamtags":["prod"],"technology":["cribl"]}}}`
	if string(body) != want {
		t.Fatalf("unexpected tags body\nwant: %s\n got: %s", want, string(body))
	}
}
