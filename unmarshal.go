package marshaler

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/go-multierror"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Custom function for converting a scalar value, the first argument is
// the source value and the second argument is the type of the destination
type UnmarshalScalarFunc func(reflect.Value, reflect.Type) (reflect.Value, error)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// UnmarshalStruct will decode src into dest field names identified by tag
func UnmarshalStruct(src, dst interface{}, name string, fn UnmarshalScalarFunc) error {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dst)

	// Destination should be a pointer to a struct
	if d.Kind() != reflect.Ptr || d.Elem().Kind() != reflect.Struct {
		return ErrBadParameter.With("destination should be ptr to struct")
	} else {
		d = d.Elem()
	}

	// Source should be map[string]
	if s.Kind() != reflect.Map || s.Type().Key().Kind() != reflect.String {
		return ErrBadParameter.With("source should be map[string]...")
	}

	// Unmarshal into each field
	var result error
	for i := 0; i < d.NumField(); i++ {
		if tag := tagName(d.Type().Field(i), name); tag != "" {
			if v := s.MapIndex(reflect.ValueOf(tag)); v.IsValid() && !v.IsNil() {
				if err := unmarshalValue(v, d.Field(i), fn); err != nil {
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
func unmarshalValue(src, dest reflect.Value, fn UnmarshalScalarFunc) error {
	switch src.Kind() {
	case reflect.Ptr:
		src := src.Elem()
		if src.IsValid() {
			dest.Set(reflect.New(src.Type()))
			return unmarshalValue(src, dest.Elem(), fn)
		}
	case reflect.Interface:
		src := src.Elem()
		recursive := true
		if fn != nil {
			if v, err := fn(src, dest.Type()); err != nil {
				return err
			} else if v.IsValid() {
				recursive = false
				src = v
			}
		}

		// Check appropriate type
		if src.Type() != dest.Type() {
			return ErrBadParameter.With("destination is ", dest.Type(), " but expected ", src.Type())
		}

		// Make copy of src if recursive, or set otherwise
		if recursive {
			copyValue := reflect.New(src.Type()).Elem()
			if err := unmarshalValue(src, copyValue, fn); err != nil {
				return err
			}
			dest.Set(copyValue)
		} else {
			dest.Set(src)
		}

	case reflect.Map:
		// Make a new map
		dest.Set(reflect.MakeMap(src.Type()))

		// Unmarshal each key/value pair
		for _, key := range src.MapKeys() {
			v := src.MapIndex(key)
			if !v.IsNil() {
				copy := reflect.New(v.Type()).Elem()
				if err := unmarshalValue(v, copy, fn); err != nil {
					return err
				}
				dest.SetMapIndex(key, copy)
			}
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i += 1 {
			if err := unmarshalValue(src.Field(i), dest.Field(i), fn); err != nil {
				return err
			}
		}
	case reflect.Slice:
		if fn != nil {
			if v, err := fn(src, dest.Type()); err != nil {
				return err
			} else if v.IsValid() && v.Kind() != reflect.Slice {
				return unmarshalValue(v, dest, fn)
			}
		}

		// Check for both slices, source can be []interface{}
		if src.Kind() != dest.Kind() && src.Type() != interfaceSliceType {
			return ErrBadParameter.With("destination is ", dest.Kind(), " but expected ", src.Kind())
		}

		// Make a new slice
		dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

		// Copy source elements
		for i := 0; i < src.Len(); i++ {
			if err := unmarshalValue(src.Index(i), dest.Index(i), fn); err != nil {
				return err
			}
		}
	default:
		if fn != nil {
			if v, err := fn(src, dest.Type()); err != nil {
				return err
			} else if v.IsValid() {
				src = v
			}
		}
		// Check appropriate type
		if src.Kind() != dest.Kind() {
			return ErrBadParameter.With("Unmarshal: ", "Destination is ", dest.Kind(), " but expected ", src.Kind())
		}

		// Set scalar
		dest.Set(src)
	}

	// Return success
	return nil
}
