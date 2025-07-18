resource "criblio_schema" "my_schema" {
  description = "test schema"
  group_id    = "default"
  id          = "my_schema"
  schema      = <<-EOT
{
  "$id": "https://example.com/person.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Person",
  "type": "object",
  "required": ["firstName", "lastName", "age"],
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be greater than zero, less than 42.",
      "type": "integer",
      "minimum": 0,
      "maximum": 42
    }
  }
}
EOT
}

output "schema" {
  value = criblio_schema.my_schema
}
