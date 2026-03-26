package provider

import (
	"reflect"
	"strings"
)

// PromoteFirstInputUnionItemToTopLevel copies the non-nil branch from Items[0] onto the model's
// top-level Input* fields. Speakeasy RefreshFrom fills Items only; this hoists like destination's
// RefreshFromSharedOutput. Called from source_resource / packsource_resource after RefreshFrom
// (those files are genignored) and from import-cli export after converter.Convert.
//
// It then normalizes empty Go slices to nil for optional list attributes. RefreshFrom always
// materializes missing API arrays as make([]T, 0), which Terraform treats as distinct from null
// (omitted in config). Destination's field-by-field RefreshFromSharedOutput often leaves nil
// slices when the API sends nothing; source’s generated loop always fills empty slices.
func PromoteFirstInputUnionItemToTopLevel(model interface{}) {
	if model == nil {
		return
	}
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = rv.Elem()
	itemsField := rv.FieldByName("Items")
	if !itemsField.IsValid() || itemsField.Kind() != reflect.Slice || itemsField.Len() == 0 {
		return
	}
	first := itemsField.Index(0)
	if first.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < first.NumField(); i++ {
		sf := first.Type().Field(i)
		name := sf.Name
		if !strings.HasPrefix(name, "Input") {
			continue
		}
		fv := first.Field(i)
		if fv.Kind() != reflect.Ptr || fv.IsNil() {
			continue
		}
		dest := rv.FieldByName(name)
		if !dest.IsValid() || !dest.CanSet() || dest.Kind() != reflect.Ptr {
			continue
		}
		if !fv.Type().AssignableTo(dest.Type()) {
			continue
		}
		dest.Set(fv)
	}
	for i := 0; i < first.NumField(); i++ {
		if !strings.HasPrefix(first.Type().Field(i).Name, "Input") {
			continue
		}
		fv := first.Field(i)
		if fv.Kind() != reflect.Ptr || fv.IsNil() {
			continue
		}
		normalizeEmptyCollectionsToNilRecursive(fv.Elem())
	}
}

// normalizeEmptyCollectionsToNilRecursive sets len-0 slices and maps to nil so optional attrs match
// omitted/null in config (avoids [] / {} vs null plan drift after RefreshFrom).
func normalizeEmptyCollectionsToNilRecursive(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		normalizeEmptyCollectionsToNilRecursive(v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if !fv.CanSet() {
				continue
			}
			switch fv.Kind() {
			case reflect.Slice:
				if fv.Len() == 0 {
					fv.Set(reflect.Zero(fv.Type()))
					continue
				}
				// Recurse into elements (e.g. auth_tokens[].allowed_indexes_at_token).
				for j := 0; j < fv.Len(); j++ {
					el := fv.Index(j)
					switch el.Kind() {
					case reflect.Struct:
						normalizeEmptyCollectionsToNilRecursive(el)
					case reflect.Ptr:
						if !el.IsNil() {
							normalizeEmptyCollectionsToNilRecursive(el.Elem())
						}
					}
				}
			case reflect.Map:
				if fv.Len() == 0 {
					fv.Set(reflect.Zero(fv.Type()))
					continue
				}
				for _, mk := range fv.MapKeys() {
					mv := fv.MapIndex(mk)
					if mv.Kind() != reflect.Struct {
						continue
					}
					c := reflect.New(mv.Type()).Elem()
					c.Set(mv)
					normalizeEmptyCollectionsToNilRecursive(c)
					fv.SetMapIndex(mk, c)
				}
			case reflect.Ptr:
				if !fv.IsNil() {
					normalizeEmptyCollectionsToNilRecursive(fv.Elem())
				}
			case reflect.Struct:
				normalizeEmptyCollectionsToNilRecursive(fv)
			}
		}
	}
}
