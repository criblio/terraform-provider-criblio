package parser

import "testing"

func TestWorkspaceSettingsPartialObjectUpdates(t *testing.T) {
	defs, err := Parse([]byte(`openapi: 3.1.0
info: {title: test, version: '1'}
paths:
  /system/settings/conf:
    patch:
      x-terraform-resource: true
      x-terraform-resource-name: WorkspaceSystemSettings
      x-terraform-partial-object-updates: true
      x-terraform-resource-factory-hook: newWorkspaceSystemSettingsManagedResource
      requestBody:
        content:
          application/json:
            schema: {$ref: '#/components/schemas/Settings'}
components:
  schemas:
    Settings:
      type: object
      properties:
        id:
          type: string
          const: conf
          x-terraform-fixed-value: conf
          x-terraform-local-only: true
        tls:
          type: object
          required: [minVersion, rejectUnauthorized]
          properties:
            minVersion: {type: string}
            rejectUnauthorized: {type: boolean}
        upgradeSettings:
          type: object
          properties:
            packageUrls:
              type: array
              items:
                type: object
                required: [packageUrl]
                properties:
                  packageUrl: {type: string}
`))
	if err != nil {
		t.Fatal(err)
	}
	settings := resourceByName(t, defs, "WorkspaceSystemSettings")
	if settings.ResourceFactoryHook != "newWorkspaceSystemSettingsManagedResource" || len(settings.Create.PathParams) != 0 {
		t.Fatal("workspace settings lost their factory or gained a path selector")
	}
	id := fieldByTFName(t, settings.Fields, "id")
	if id.RequestField || id.FixedValue != "conf" || !id.Computed || !id.Optional {
		t.Fatalf("incorrect local singleton identity: %+v", id)
	}
	tls := fieldByTFName(t, settings.Fields, "tls")
	for _, field := range tls.Fields {
		if field.Required || !field.Optional {
			t.Fatalf("partial settings require the full section: %+v", field)
		}
	}
	upgrades := fieldByTFName(t, settings.Fields, "upgrade_settings")
	packages := fieldByTFName(t, upgrades.Fields, "package_urls")
	if !fieldByTFName(t, packages.Fields, "package_url").Required {
		t.Fatal("replacement array elements must retain required fields")
	}
}
