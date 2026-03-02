// Required Terraform identifier fields that RefreshFrom* methods do not set;
// they are populated from the import ID or Get request params.
package converter

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IdentifierParamNames are the request/import param names that map to Terraform
// model fields (e.g. GroupID, ID, Pack, LakeID). Used to inject required identifiers
// after conversion so models are valid for Terraform and HCL generation.
var IdentifierParamNames = []string{"GroupID", "ID", "Pack", "LakeID"}

// InjectRequiredIdentifiers sets required Terraform identifier fields on the
// converted model (e.g. id, group_id) from the given identifiers map. Keys should
// match request param names (e.g. "ID", "GroupID", "Pack"). Only fields that
// exist on the model and are types.String are set. This ensures downstream HCL
// generation receives valid Terraform models with required identifiers populated.
func InjectRequiredIdentifiers(model interface{}, identifiers map[string]string) error {
	if model == nil || len(identifiers) == 0 {
		return nil
	}
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct or pointer to struct, got %s", reflect.ValueOf(model).Kind())
	}
	for _, paramName := range IdentifierParamNames {
		v, ok := identifiers[paramName]
		if !ok || v == "" {
			continue
		}
		f := val.FieldByName(paramName)
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		if !isTypesString(f.Type()) {
			continue
		}
		f.Set(reflect.ValueOf(types.StringValue(v)))
	}
	return nil
}

func isTypesString(t reflect.Type) bool {
	return t == reflect.TypeOf(types.String{})
}
