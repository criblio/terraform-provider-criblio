package provider

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTagsAPIMapUsesEmptyArrays(t *testing.T) {
	tags, err := tagsAPIMap(context.Background(), packRequestBodyTagsObjectFromAPI(&packTagsAPI{}))
	if err != nil {
		t.Fatalf("build tags API map: %v", err)
	}

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
	tags, err := tagsAPIMap(context.Background(), packRequestBodyTagsObjectFromAPI(&packTagsAPI{
		DataType:   []string{"logs"},
		Domain:     []string{"security"},
		Streamtags: []string{"prod"},
		Technology: []string{"cribl"},
	}))
	if err != nil {
		t.Fatalf("build tags API map: %v", err)
	}

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
