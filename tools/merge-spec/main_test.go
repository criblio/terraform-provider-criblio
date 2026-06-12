package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestRunPrefixesUnwrapsAndAppliesOverlay(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "upstream-openapi.yml")
	overlay := filepath.Join(dir, "terraform-overlay.yml")
	output := filepath.Join(dir, "merged-spec.yml")
	cloudOnlyOutput := filepath.Join(dir, "cloud_only_paths.go")

	writeFile(t, input, `openapi: 3.1.0
paths:
  /search/datasets:
    get:
      x-cribl-availability: cloud
      responses:
        "200":
          description: ok
  /search/jobs:
    get:
      x-cribl-availability: both
      responses:
        "200":
          description: ok
    post:
      x-cribl-availability: cloud
      responses:
        "200":
          description: ok
  /system/certificates:
    post:
      x-cribl-availability: both
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CountedCertificate"
    get:
      x-cribl-availability: both
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CountedCertificate"
  /system/certificates/{id}:
    get:
      x-cribl-availability: both
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CountedCertificate"
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
  /system/strings:
    get:
      x-cribl-availability: both
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CountedString"
  /master/groups:
    get:
      x-cribl-availability: both
      responses:
        "200":
          description: ok
components:
  schemas:
    Certificate:
      type: object
      properties:
        id:
          type: string
    CountedCertificate:
      type: object
      properties:
        count:
          type: integer
        items:
          type: array
          items:
            $ref: "#/components/schemas/Certificate"
    CountedString:
      type: object
      properties:
        count:
          type: integer
        items:
          type: array
          items:
            type: string
`)

	writeFile(t, overlay, `actions:
  - target: "$.paths./m/{groupId}/system/certificates.post"
    update:
      x-terraform-resource: true
      x-terraform-resource-name: certificate
  - target: "$.paths./m/{groupId}/system/certificates/{id}.get"
    update:
      x-terraform-read: certificate
  - target: "$.components.schemas.Certificate.properties.id"
    update:
      x-terraform-force-new: true
`)

	if err := run(input, overlay, output, cloudOnlyOutput); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	merged := readFile(t, output)
	assertContains(t, merged, "/m/{groupId}/system/certificates:")
	assertContains(t, merged, "/m/{groupId}/system/strings:")
	assertContains(t, merged, "/master/groups:")
	assertNotContains(t, merged, "  /system/certificates:")
	assertContains(t, merged, "name: groupId")
	assertContains(t, merged, `$ref: "#/components/schemas/Certificate"`)
	assertContains(t, merged, "type: string")
	assertContains(t, merged, "x-terraform-resource-name: certificate")
	assertContains(t, merged, "x-terraform-force-new: true")
	assertNotContains(t, merged, `$ref: "#/components/schemas/CountedCertificate"`)
	assertNotContains(t, merged, `$ref: "#/components/schemas/CountedString"`)

	outputNode, err := readYAML(output)
	if err != nil {
		t.Fatalf("failed to parse merged output: %v", err)
	}
	assertSchemaRef(t, outputNode, "$.paths./m/{groupId}/system/certificates.post.responses.200.content.application/json.schema", "#/components/schemas/Certificate")
	assertSchemaRef(t, outputNode, "$.paths./m/{groupId}/system/certificates/{id}.get.responses.200.content.application/json.schema", "#/components/schemas/Certificate")
	assertArrayItemsRef(t, outputNode, "$.paths./m/{groupId}/system/certificates.get.responses.200.content.application/json.schema", "#/components/schemas/Certificate")
	assertArrayItemsType(t, outputNode, "$.paths./m/{groupId}/system/strings.get.responses.200.content.application/json.schema", "string")

	cloudOnlyPaths := readFile(t, cloudOnlyOutput)
	assertContains(t, cloudOnlyPaths, `"/search/datasets": true`)
	assertNotContains(t, cloudOnlyPaths, `"/search/jobs": true`)
	assertNotContains(t, cloudOnlyPaths, `"/system/certificates": true`)
}

func writeFile(t *testing.T, filename, content string) {
	t.Helper()

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", filename, err)
	}
}

func readFile(t *testing.T, filename string) string {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %s: %v", filename, err)
	}
	return string(content)
}

func assertContains(t *testing.T, s, want string) {
	t.Helper()

	if !strings.Contains(s, want) {
		t.Fatalf("expected output to contain %q", want)
	}
}

func assertNotContains(t *testing.T, s, want string) {
	t.Helper()

	if strings.Contains(s, want) {
		t.Fatalf("expected output not to contain %q", want)
	}
}

func assertSchemaRef(t *testing.T, root *yaml.Node, target, want string) {
	t.Helper()

	schema, err := lookupTarget(documentMapping(root), target)
	if err != nil {
		t.Fatalf("failed to look up %s: %v", target, err)
	}
	ref, ok := mappingValue(schema, "$ref")
	if !ok {
		t.Fatalf("%s missing $ref, expected %q", target, want)
	}
	if ref.Value != want {
		t.Fatalf("%s $ref = %q, expected %q", target, ref.Value, want)
	}
}

func assertArrayItemsRef(t *testing.T, root *yaml.Node, target, want string) {
	t.Helper()

	items := assertArraySchema(t, root, target)
	ref, ok := mappingValue(items, "$ref")
	if !ok {
		t.Fatalf("%s items missing $ref, expected %q", target, want)
	}
	if ref.Value != want {
		t.Fatalf("%s items $ref = %q, expected %q", target, ref.Value, want)
	}
}

func assertArrayItemsType(t *testing.T, root *yaml.Node, target, want string) {
	t.Helper()

	items := assertArraySchema(t, root, target)
	typ, ok := mappingValue(items, "type")
	if !ok {
		t.Fatalf("%s items missing type, expected %q", target, want)
	}
	if typ.Value != want {
		t.Fatalf("%s items type = %q, expected %q", target, typ.Value, want)
	}
}

func assertArraySchema(t *testing.T, root *yaml.Node, target string) *yaml.Node {
	t.Helper()

	schema, err := lookupTarget(documentMapping(root), target)
	if err != nil {
		t.Fatalf("failed to look up %s: %v", target, err)
	}
	typ, ok := mappingValue(schema, "type")
	if !ok {
		t.Fatalf("%s missing type, expected array", target)
	}
	if typ.Value != "array" {
		t.Fatalf("%s type = %q, expected array", target, typ.Value)
	}
	items, ok := mappingValue(schema, "items")
	if !ok {
		t.Fatalf("%s missing items schema", target)
	}
	return items
}
