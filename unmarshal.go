package marshaler

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/go-multierror"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// UnmarshalStruct will decode src into dest field names identified by tag
func UnmarshalStruct(src, dst interface{}, name string) error {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dst)

	// Destination should be a pointer to a struct
	if d.Kind() != reflect.Ptr || d.Elem().Kind() != reflect.Struct {
		return ErrBadParameter.With("UnmarshalStruct: Destination should be ptr to struct")
	} else {
		d = d.Elem()
	}

	// Source should be map[string]
	if s.Kind() != reflect.Map || s.Type().Key().Kind() != reflect.String {
		return ErrBadParameter.With("UnmarshalStruct: Source should be map[string]...")
	}

	// Unmarshal into each field
	var result error
	for i := 0; i < d.NumField(); i++ {
		if tag := tagName(d.Type().Field(i), name); tag != "" {
			if v := s.MapIndex(reflect.ValueOf(tag)); v.IsValid() && !v.IsNil() {
				if err := unmarshalValue(v, d.Field(i)); err != nil {
					result = multierror.Append(result, err)
				}
			}
		}
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// tagName returns the name of the field based on tag or field name
// and returns empty string if the field should be ignored (not assignable)
func tagName(field reflect.StructField, tagName string) string {
	// Check for private field
	if field.Name != "" && unicode.IsLower(rune(field.Name[0])) {
		return ""
	}
	tags := strings.Split(field.Tag.Get(tagName), ",")
	if tags[0] == "-" {
		return ""
	} else if tags[0] == "" {
		return field.Name
	} else {
		return tags[0]
	}
}

// unmarshalValue recursively unmarshals src into dest and returns any errors if src is
// not assignable into dest
func unmarshalValue(src, dest reflect.Value) error {
	switch src.Kind() {
	case reflect.Ptr:
		src := src.Elem()
		if src.IsValid() {
			dest.Set(reflect.New(src.Type()))
			return unmarshalValue(src, dest.Elem())
		}
	case reflect.Interface:
		src := src.Elem()
		copyValue := reflect.New(src.Type()).Elem()
		if err := unmarshalValue(src, copyValue); err != nil {
			return err
		}
		dest.Set(copyValue)
	case reflect.Map:
		// Make a new map
		dest.Set(reflect.MakeMap(src.Type()))

		// Unmarshal each key/value pair
		for _, key := range src.MapKeys() {
			v := src.MapIndex(key)
			if !v.IsNil() {
				copy := reflect.New(v.Type()).Elem()
				if err := unmarshalValue(v, copy); err != nil {
					return err
				}
				dest.SetMapIndex(key, copy)
			}
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i += 1 {
			if err := unmarshalValue(src.Field(i), dest.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Slice:
		// Check for both slices
		if src.Kind() != dest.Kind() {
			return ErrBadParameter.With("Unmarshal: ", "Destination is ", dest.Kind(), " but expected ", src.Kind())
		}

		// Make a new slice
		dest.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))

		// Copy source elements
		for i := 0; i < src.Len(); i++ {
			if err := unmarshalValue(src.Index(i), dest.Index(i)); err != nil {
				return err
			}
		}
	default:
		if src.Kind() != dest.Kind() {
			return ErrBadParameter.With("Unmarshal: ", "Destination is ", dest.Kind(), " but expected ", src.Kind())
		}
		dest.Set(src)
	}

	// Return success
	return nil
}
