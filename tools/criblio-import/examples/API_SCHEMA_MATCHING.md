# API Schema Matching in Import Tool

## OpenAPI Schema Structure

From `openapi.yml`:
```yaml
responses:
  "200":
    description: a list of Global Variable objects
    content:
      application/json:
        schema:
          type: object
          properties:
            count:
              type: integer
            items:
              type: array
              items:
                type: object
                additionalProperties: true
```

This means:
- Response is an object with `items` array
- Each item in `items` is an object with `additionalProperties: true` (dynamic properties)

## How Our Tool Matches This

### ✅ Our Tool Uses `GetGlobalVariable` (List Endpoint)

```go
request := operations.GetGlobalVariableRequest{
    GroupID: groupID,
}
res, err := client.GlobalVariables.GetGlobalVariable(ctx, request)
```

**Response Type:** `GetGlobalVariableResponseBody`
- `Items []shared.GlobalVar` - Already parsed to typed struct! ✅

The SDK automatically maps the `items` array to `[]shared.GlobalVar`, which has:
- `ID string`
- `Description *string`
- `Lib *string`
- `Tags *string`
- `Type *GlobalVarType`
- `Value *string`

### ✅ Our Tool Already Matches the Schema

We're correctly:
1. Using the list endpoint that returns `items` array
2. Getting typed `shared.GlobalVar` objects (SDK handles parsing)
3. Extracting individual fields from each item
4. Generating Terraform config with all fields

## Why Provider's Read Function Doesn't Match

The provider uses `GetGlobalVariableByID` (single item endpoint):

```go
res, err := r.client.GlobalVariables.GetGlobalVariableByID(ctx, *request)
```

**Response Type:** `GetGlobalVariableByIDResponseBody`
- `Items []map[string]any` - Generic map, not typed struct! ❌

The provider then has to manually extract fields from `items[0]` map, but it doesn't - it just stores the whole map in `items` field.

## Comparison

### Our Import Tool (✅ Correct)
```go
// Uses GetGlobalVariable - returns []shared.GlobalVar
for _, gv := range res.Object.Items {
    // gv is already a typed GlobalVar struct
    id := gv.ID                    // ✅ Direct access
    desc := gv.Description         // ✅ Direct access
    lib := gv.Lib                 // ✅ Direct access
    // ... etc
}
```

### Provider's Read Function (❌ Issue)
```go
// Uses GetGlobalVariableByID - returns []map[string]any
resp := res.Object.Items  // []map[string]any
// Should extract: resp[0]["id"], resp[0]["description"], etc.
// But it doesn't - just stores whole map in items field
```

## Solution: Match OpenAPI Schema in Provider

The provider's `RefreshFromOperationsGetGlobalVariableByIDResponseBody` should:

1. Extract the first item from `items` array
2. Parse it to `shared.GlobalVar` struct (or extract fields manually)
3. Populate individual fields in the model

```go
// What it should do:
if len(resp.Items) > 0 {
    item := resp.Items[0]  // map[string]any
    
    // Extract fields from map
    if id, ok := item["id"].(string); ok {
        r.ID = types.StringValue(id)
    }
    if desc, ok := item["description"].(string); ok {
        r.Description = types.StringValue(desc)
    }
    // ... etc for all fields
}
```

## Our Tool Status

✅ **Our import tool already correctly matches the OpenAPI schema:**
- Uses the correct endpoint (`GetGlobalVariable`)
- Gets properly typed responses (`[]shared.GlobalVar`)
- Extracts all fields correctly
- Generates accurate Terraform config

The updates you see after import are due to the **provider's Read function not matching the schema** - it doesn't extract individual fields from the response.

