package parser

import (
	"strings"
	"testing"
)

func TestParseRejectsUnresolvedPathPlaceholder(t *testing.T) {
	_, err := Parse([]byte(`openapi: 3.1.0
info: {title: Test, version: '1'}
paths:
  /m/{groupId}/notifications:
    get:
      x-terraform-list: Notification
      parameters: [{name: groupId, in: query, schema: {type: string}}]
      responses:
        '200':
          description: Notifications
          content:
            application/json:
              schema: {$ref: '#/components/schemas/Notifications'}
components:
  schemas:
    Notification:
      type: object
      properties:
        id: {type: string}
    Notifications:
      type: object
      properties:
        items: {type: array, items: {$ref: '#/components/schemas/Notification'}}
`))
	if err == nil || !strings.Contains(err.Error(), "unresolved path placeholder") {
		t.Fatalf("error = %v; want unresolved placeholder error", err)
	}
}
