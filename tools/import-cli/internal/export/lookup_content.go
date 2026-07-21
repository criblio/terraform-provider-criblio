package export

import (
	"context"
	"encoding/json"
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

func lookupContentAsset(ctx context.Context, client *importclient.Client, typeName string, attrs map[string]hcl.Value, resourceName string, idMap map[string]string) ([]generator.ResourceFile, bool, error) {
	if !isLookupFileType(typeName) {
		return nil, true, nil
	}

	content, ok := attrs["content"]
	if !ok || content.Kind == hcl.KindNull {
		downloaded, found, err := fetchLookupContent(ctx, client, typeName, idMap)
		if err != nil || !found {
			return nil, false, err
		}
		content = hcl.Value{Kind: hcl.KindString, String: downloaded}
	}
	if content.Kind != hcl.KindString {
		return nil, false, nil
	}

	assetPath := path.Join("files", resourceName, lookupContentFilename(attrs))
	attrs["content"] = hcl.Value{
		Kind: hcl.KindExpression,
		Expr: "file(" + strconv.Quote("${path.module}/"+assetPath) + ")",
	}
	return []generator.ResourceFile{{
		Path:    assetPath,
		Content: []byte(content.String),
	}}, true, nil
}

func lookupContentFilename(attrs map[string]hcl.Value) string {
	id, ok := attrs["id"]
	if !ok || id.Kind != hcl.KindString || id.String == "" {
		return "lookup.csv"
	}
	name := safeLookupContentFilename(id.String)
	if name == "" || name == "." {
		return "lookup.csv"
	}
	return name
}

func safeLookupContentFilename(name string) string {
	name = strings.TrimSpace(name)
	var b strings.Builder
	lastUnderscore := false
	for _, r := range name {
		allowed := (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' ||
			r == '-' ||
			r == '_'
		if allowed {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(b.String(), "._")
}

func fetchLookupContent(ctx context.Context, client *importclient.Client, typeName string, idMap map[string]string) (string, bool, error) {
	if client == nil || client.REST == nil {
		return "", false, nil
	}
	id := idMap["id"]
	groupID := idMap["group_id"]
	if id == "" || groupID == "" {
		return "", false, nil
	}

	filename := lookupDownloadFilename(id)
	var requestPath string
	switch typeName {
	case "criblio_lookup_file":
		requestPath = fmt.Sprintf("/m/%s/system/lookups/%s/content?raw=true", url.PathEscape(groupID), url.PathEscape(filename))
	case "criblio_pack_lookups":
		pack := idMap["pack"]
		if pack == "" {
			return "", false, nil
		}
		requestPath = fmt.Sprintf("/m/%s/p/%s/system/lookups/%s/content?raw=true", url.PathEscape(groupID), url.PathEscape(pack), url.PathEscape(filename))
	default:
		return "", false, nil
	}

	body, err := restclient.GetRaw(ctx, client.REST, requestPath)
	if err != nil {
		if restclient.IsNotFound(err) {
			return "", false, nil
		}
		return "", false, err
	}
	if content, ok, isEnvelope := lookupContentFromJSONResponse(body); isEnvelope {
		if !ok {
			return "", false, nil
		}
		return content, true, nil
	}
	if len(body) == 0 {
		return "", false, nil
	}
	return string(body), true, nil
}

func lookupContentFromJSONResponse(body []byte) (string, bool, bool) {
	var envelope struct {
		Items []struct {
			Content *string `json:"content"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "", false, false
	}
	if envelope.Items == nil {
		return "", false, false
	}
	if len(envelope.Items) == 0 || envelope.Items[0].Content == nil {
		return "", false, true
	}
	return *envelope.Items[0].Content, true, true
}

func lookupDownloadFilename(id string) string {
	if lookupIDHasKnownExtension(id) {
		return id
	}
	return id + ".csv"
}

func lookupIDHasKnownExtension(id string) bool {
	lower := strings.ToLower(id)
	return strings.HasSuffix(lower, ".csv") ||
		strings.HasSuffix(lower, ".gz") ||
		strings.HasSuffix(lower, ".csv.gz") ||
		strings.HasSuffix(lower, ".mmdb")
}
