package export

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type packDestinationTransport struct{ groupID string }

func (tr packDestinationTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body := `{"items":[]}`
	switch req.URL.Path {
	case "/api/v1/products/stream/groups":
		body = `{"items":[{"id":"` + tr.groupID + `"}]}`
	case "/api/v1/m/" + tr.groupID + "/packs":
		body = `{"items":[{"id":"custom-pack"}]}`
	case "/api/v1/m/" + tr.groupID + "/p/custom-pack/system/outputs":
		body = `{"items":[{"id":"default","type":"default"},{"id":"devnull","type":"devnull"},{"id":"custom-output","type":"webhook"}]}`
	case "/api/v1/m/" + tr.groupID + "/p/custom-pack/system/outputs/custom-output":
		body = `{"items":[{"id":"custom-output","type":"webhook","url":"https://example.com/events"}]}`
	}
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
}

func TestPackDestinationExportAcrossGroups(t *testing.T) {
	reg := buildTestRegistry(t)
	for _, groupID := range []string{"default", "custom-group"} {
		for _, excludeDefaults := range []bool{false, true} {
			name := groupID
			if excludeDefaults {
				name += "/exclude-defaults"
			}
			t.Run(name, func(t *testing.T) {
				ctx := context.Background()
				client := &importclient.Client{REST: restclient.New(restclient.Config{
					BaseURL: "https://example.invalid", BearerToken: "test",
					HTTPClient: &http.Client{Transport: packDestinationTransport{groupID: groupID}},
				})}
				groups := []string{groupID}
				results, err := discovery.Discover(ctx, client, reg, []string{"criblio_pack_destination"}, nil, groups, false)
				require.NoError(t, err)
				require.Len(t, results, 1)
				require.NoError(t, results[0].Err)
				require.Equal(t, 1, results[0].Count)
				result, err := ToResourceItems(ctx, client, reg, results, groups, groups, 1, excludeDefaults, IncludeOverride{}, nil)
				require.NoError(t, err)
				require.Empty(t, result.ConvertSkipped)
				require.Len(t, result.Items, 1)
				item := result.Items[0]
				assert.Equal(t, groupID, item.GroupID)
				assert.JSONEq(t, `{"group_id":"`+groupID+`","id":"custom-output","pack":"custom-pack"}`, item.ImportID)
				assert.Equal(t, "custom-pack", item.Attrs["pack"].String)
				assert.Equal(t, "https://example.com/events", item.Attrs["output_webhook"].Map["url"].String)
				file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
				require.NoError(t, err)
				assert.Contains(t, string(file.Bytes()), `resource "criblio_pack_destination"`)
				assert.Contains(t, string(file.Bytes()), "https://example.com/events")
			})
		}
	}
}
