package export

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// AttachPackAssets exports every pack as a portable .crbl archive and updates
// its Terraform configuration to upload that local archive to the target.
func AttachPackAssets(ctx context.Context, client *importclient.Client, items []generator.ResourceItem) (int, error) {
	if client == nil || client.REST == nil {
		return 0, fmt.Errorf("REST client is nil")
	}
	attached := 0
	for i := range items {
		item := &items[i]
		if item.TypeName != "criblio_pack" {
			continue
		}
		packID, ok := packStringAttr(item.Attrs, "id")
		if !ok || strings.TrimSpace(packID) == "" || strings.TrimSpace(item.GroupID) == "" {
			return attached, fmt.Errorf("pack %q is missing group_id or id", item.Name)
		}
		filename := portablePackFilename(packID)
		query := url.Values{}
		query.Set("mode", "merge")
		query.Set("filename", filename)
		requestPath := fmt.Sprintf("/m/%s/packs/%s/export?%s", url.PathEscape(item.GroupID), url.PathEscape(packID), query.Encode())
		content, err := restclient.GetRaw(ctx, client.REST, requestPath)
		if err != nil {
			return attached, fmt.Errorf("export pack %q: %w", packID, err)
		}
		if len(content) == 0 {
			return attached, fmt.Errorf("export pack %q: API returned an empty archive", packID)
		}

		assetPath := path.Join("files", item.Name, filename)
		item.Files = append(item.Files, generator.ResourceFile{Path: assetPath, Content: content})
		item.Attrs["filename"] = hcl.Value{
			Kind: hcl.KindExpression,
			Expr: strconv.Quote("${path.module}/" + assetPath),
		}
		delete(item.Attrs, "source")
		delete(item.Attrs, "description")
		delete(item.Attrs, "display_name")
		attached++
	}
	return attached, nil
}

// PreparePackResources configures packs as empty containers. Their supported
// contents are recreated by the exported criblio_pack_* resources.
func PreparePackResources(items []generator.ResourceItem) int {
	prepared := 0
	for i := range items {
		if items[i].TypeName != "criblio_pack" {
			continue
		}
		delete(items[i].Attrs, "source")
		delete(items[i].Attrs, "filename")
		items[i].Files = nil
		prepared++
	}
	return prepared
}

// RemovePackChildResources removes resources already contained in exported
// pack archives, preventing duplicate creation after the archive is installed.
func RemovePackChildResources(items []generator.ResourceItem) ([]generator.ResourceItem, int) {
	filtered := make([]generator.ResourceItem, 0, len(items))
	removed := 0
	for _, item := range items {
		if strings.HasPrefix(item.TypeName, "criblio_pack_") {
			removed++
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered, removed
}

func portablePackFilename(packID string) string {
	name := safeLookupContentFilename(packID)
	if name == "" {
		name = "pack"
	}
	if !strings.HasSuffix(strings.ToLower(name), ".crbl") {
		name += ".crbl"
	}
	return name
}

func packStringAttr(attrs map[string]hcl.Value, name string) (string, bool) {
	value, ok := attrs[name]
	if !ok || value.Kind != hcl.KindString {
		return "", false
	}
	return value.String, true
}
