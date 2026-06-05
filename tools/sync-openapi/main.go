package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

const (
	defaultDocsURL = "https://docs.cribl.io/cribl-as-code/api-reference/control-plane/cribl-core/"
	defaultOutput  = "upstream-openapi.yml"

	maxDocsSize = 10 << 20
	maxSpecSize = 100 << 20
)

type config struct {
	source  string
	docsURL string
	version string
	url     string
	file    string
	output  string
	timeout time.Duration
	client  *http.Client
}

func main() {
	cfg := parseConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	resolved, size, err := syncOpenAPI(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync openapi: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s from %s (%d bytes)\n", cfg.output, resolved, size)
}

func parseConfig() config {
	var cfg config
	flag.StringVar(&cfg.source, "source", envOrDefault("CRIBL_OPENAPI_SOURCE", "docs"), "source type: docs, url, or file")
	flag.StringVar(&cfg.docsURL, "docs-url", envOrDefault("CRIBL_OPENAPI_DOCS_URL", defaultDocsURL), "Cribl API docs page containing the OpenAPI versions map")
	flag.StringVar(&cfg.version, "version", envOrDefault("CRIBL_OPENAPI_VERSION", "latest"), "docs OpenAPI version to download, for example latest or 4.18.0")
	flag.StringVar(&cfg.url, "url", os.Getenv("CRIBL_OPENAPI_URL"), "direct OpenAPI YAML URL when -source=url")
	flag.StringVar(&cfg.file, "file", os.Getenv("CRIBL_OPENAPI_FILE"), "local OpenAPI YAML path when -source=file")
	flag.StringVar(&cfg.output, "out", envOrDefault("CRIBL_OPENAPI_OUT", defaultOutput), "output OpenAPI YAML path")
	flag.DurationVar(&cfg.timeout, "timeout", 60*time.Second, "overall sync timeout")
	flag.Parse()
	cfg.client = http.DefaultClient
	return cfg
}

func syncOpenAPI(ctx context.Context, cfg config) (string, int, error) {
	if cfg.output == "" {
		return "", 0, fmt.Errorf("output path is required")
	}

	spec, resolved, err := loadSpec(ctx, cfg)
	if err != nil {
		return "", 0, err
	}
	if err := validateOpenAPI(spec); err != nil {
		return "", 0, err
	}
	if err := writeFileAtomic(cfg.output, spec); err != nil {
		return "", 0, err
	}
	return resolved, len(spec), nil
}

func loadSpec(ctx context.Context, cfg config) ([]byte, string, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.source)) {
	case "docs", "":
		specURL, err := resolveDocsSpecURL(ctx, cfg.client, cfg.docsURL, cfg.version)
		if err != nil {
			return nil, "", err
		}
		spec, err := download(ctx, cfg.client, specURL, maxSpecSize)
		if err != nil {
			return nil, "", err
		}
		return spec, specURL, nil
	case "url":
		if cfg.url == "" {
			return nil, "", fmt.Errorf("-url or CRIBL_OPENAPI_URL is required when source=url")
		}
		spec, err := download(ctx, cfg.client, cfg.url, maxSpecSize)
		if err != nil {
			return nil, "", err
		}
		return spec, cfg.url, nil
	case "file":
		if cfg.file == "" {
			return nil, "", fmt.Errorf("-file or CRIBL_OPENAPI_FILE is required when source=file")
		}
		spec, err := os.ReadFile(cfg.file)
		if err != nil {
			return nil, "", fmt.Errorf("read %s: %w", cfg.file, err)
		}
		return spec, cfg.file, nil
	default:
		return nil, "", fmt.Errorf("unsupported source %q; use docs, url, or file", cfg.source)
	}
}

func resolveDocsSpecURL(ctx context.Context, client *http.Client, docsURL, version string) (string, error) {
	if docsURL == "" {
		return "", fmt.Errorf("docs URL is required")
	}
	body, err := download(ctx, client, docsURL, maxDocsSize)
	if err != nil {
		return "", err
	}

	versions, err := parseVersionsMap(string(body))
	if err != nil {
		return "", fmt.Errorf("parse docs versions: %w", err)
	}

	key, specURL, err := selectVersion(versions, version)
	if err != nil {
		return "", err
	}
	if _, err := url.ParseRequestURI(specURL); err != nil {
		return "", fmt.Errorf("docs version %s has invalid OpenAPI URL %q: %w", key, specURL, err)
	}
	return specURL, nil
}

func parseVersionsMap(html string) (map[string]string, error) {
	idx := strings.Index(html, "versions=")
	if idx < 0 {
		return nil, fmt.Errorf("versions map not found")
	}
	brace := strings.IndexByte(html[idx:], '{')
	if brace < 0 {
		return nil, fmt.Errorf("versions map opening brace not found")
	}
	start := idx + brace
	end, err := matchingBraceEnd(html, start)
	if err != nil {
		return nil, err
	}

	var versions map[string]string
	if err := json.Unmarshal([]byte(html[start:end]), &versions); err != nil {
		return nil, fmt.Errorf("decode versions map: %w", err)
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("versions map is empty")
	}
	return versions, nil
}

func matchingBraceEnd(s string, start int) (int, error) {
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case ch == '\\':
				escaped = true
			case ch == '"':
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, nil
			}
		}
	}
	return 0, fmt.Errorf("versions map closing brace not found")
}

func selectVersion(versions map[string]string, requested string) (string, string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" || strings.EqualFold(requested, "latest") {
		key, ok := latestVersionKey(versions)
		if !ok {
			return "", "", fmt.Errorf("no parseable docs versions found")
		}
		return key, versions[key], nil
	}

	request, err := parseVersion(requested)
	if err != nil {
		return "", "", err
	}
	if len(request.parts) >= 3 {
		if url, ok := versions[request.key]; ok {
			return request.key, url, nil
		}
	}

	var matches []versionKey
	for key := range versions {
		candidate, err := parseVersion(key)
		if err != nil {
			continue
		}
		if candidate.matches(request) {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 0 {
		return "", "", fmt.Errorf("version %q not found in docs", requested)
	}
	sort.Sort(sort.Reverse(byVersion(matches)))
	key := matches[0].key
	return key, versions[key], nil
}

func latestVersionKey(versions map[string]string) (string, bool) {
	keys := make([]versionKey, 0, len(versions))
	for key := range versions {
		parsed, err := parseVersion(key)
		if err != nil {
			continue
		}
		keys = append(keys, parsed)
	}
	if len(keys) == 0 {
		return "", false
	}
	sort.Sort(sort.Reverse(byVersion(keys)))
	return keys[0].key, true
}

type versionKey struct {
	key   string
	parts []int
}

type byVersion []versionKey

func (v byVersion) Len() int {
	return len(v)
}

func (v byVersion) Less(i, j int) bool {
	a := v[i].parts
	b := v[j].parts
	for idx := 0; idx < len(a) && idx < len(b); idx++ {
		if a[idx] != b[idx] {
			return a[idx] < b[idx]
		}
	}
	return len(a) < len(b)
}

func (v byVersion) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func parseVersion(raw string) (versionKey, error) {
	key := normalizeVersionKey(raw)
	trimmed := strings.TrimPrefix(key, "v")
	if trimmed == "" {
		return versionKey{}, fmt.Errorf("invalid version %q", raw)
	}

	segments := strings.Split(trimmed, ".")
	parts := make([]int, 0, len(segments))
	for _, segment := range segments {
		if segment == "" {
			return versionKey{}, fmt.Errorf("invalid version %q", raw)
		}
		part, err := strconv.Atoi(segment)
		if err != nil {
			return versionKey{}, fmt.Errorf("invalid version %q: %w", raw, err)
		}
		parts = append(parts, part)
	}
	return versionKey{key: key, parts: parts}, nil
}

func normalizeVersionKey(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return version
	}
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func (v versionKey) matches(requested versionKey) bool {
	if len(requested.parts) > len(v.parts) {
		return false
	}
	for i := range requested.parts {
		if v.parts[i] != requested.parts[i] {
			return false
		}
	}
	return true
}

func download(ctx context.Context, client *http.Client, rawURL string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request %s: %w", rawURL, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: status %s", rawURL, resp.Status)
	}

	reader := &io.LimitedReader{R: resp.Body, N: limit + 1}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", rawURL, err)
	}
	if int64(len(body)) > limit {
		return nil, fmt.Errorf("download %s exceeded %d byte limit", rawURL, limit)
	}
	return body, nil
}

func validateOpenAPI(spec []byte) error {
	var root map[string]any
	if err := yaml.Unmarshal(spec, &root); err != nil {
		return fmt.Errorf("parse OpenAPI YAML: %w", err)
	}
	openapi, ok := root["openapi"].(string)
	if !ok || strings.TrimSpace(openapi) == "" {
		return fmt.Errorf("OpenAPI YAML missing root openapi field")
	}
	if _, ok := root["paths"].(map[string]any); !ok {
		return fmt.Errorf("OpenAPI YAML missing root paths map")
	}
	return nil
}

func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp output: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp output: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp output: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
