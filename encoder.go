package marshaler

import (
	"reflect"
	"strings"
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Encoder struct {
	name string
}

type Field struct {
	Index int
	Name  string
	Tags  []string
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new decoder object with 'name' used as struct tag for interpreting
// the field name
func NewEncoder(name string) *Encoder {
	return &Encoder{name}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Reflect on a structure (or pointer to structure) returns field names and
// their tags or nil if any field is ignored
func (this *Encoder) Reflect(v interface{}) []*Field {
	rv := reflect.ValueOf(v)
	// Fudge pointers
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	// Enumerate struct fields
	var result []*Field
	for i := 0; i < rv.Type().NumField(); i++ {
		result = append(result, reflectField(rv.Type().Field(i), this.name))
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func reflectField(field reflect.StructField, name string) *Field {
	var result Field

	// Private or anonymous fields not supported
	if field.Anonymous || unicode.IsLower(rune(field.Name[0])) {
		return nil
	}

	// Set index
	result.Index = field.Index[0]

	// Set the field name
	tags := strings.Split(field.Tag.Get(name), ",")
	if tags[0] == "-" {
		return nil
	} else if tags[0] == "" {
		result.Name = field.Name
	} else {
		result.Name = tags[0]
	}

	// Set the tags
	result.Tags = tags[1:]

	// Return the result
	return &result
}
