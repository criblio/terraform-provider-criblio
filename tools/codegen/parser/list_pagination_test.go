package parser

import (
	"fmt"
	"testing"
)

func TestParseListDisablePagination(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("disabled=%t", disabled), func(t *testing.T) {
			definitions, err := Parse([]byte(fmt.Sprintf(`openapi: 3.1.0
info: {title: Test, version: '1'}
paths:
  /system/keys:
    get:
      x-terraform-list: Key
      x-terraform-disable-pagination: %t
      parameters:
        - {name: offset, in: query, schema: {type: integer}}
        - {name: limit, in: query, schema: {type: integer}}
      responses:
        '200':
          description: Keys
          content:
            application/json:
              schema: {$ref: '#/components/schemas/Keys'}
components:
  schemas:
    Key:
      type: object
      properties:
        keyId: {type: string}
    Keys:
      type: object
      properties:
        totalCount: {type: integer}
        items: {type: array, items: {$ref: '#/components/schemas/Key'}}
`, disabled)))
			if err != nil {
				t.Fatal(err)
			}
			if len(definitions) != 1 {
				t.Fatalf("got %d definitions, want 1", len(definitions))
			}
			if got := definitions[0].List.Paginated; got == disabled {
				t.Fatalf("Paginated = %t, want %t", got, !disabled)
			}
		})
	}
}
