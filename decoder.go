package marshaler

import (
	"reflect"
	"strings"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Decoder struct {
	name  string
	hooks []UnmarshalScalarFunc
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	nilValue = reflect.ValueOf(nil)
	timeType = reflect.TypeOf(time.Time{})
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new decoder object with 'name' used as struct tag for interpreting
// the field name
func NewDecoder(name string, hooks ...UnmarshalScalarFunc) *Decoder {
	return &Decoder{name, hooks}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Decoder) Decode(src map[string]interface{}, dest interface{}) error {
	return UnmarshalStruct(src, dest, this.name, this.unmarshalscalar)
}

///////////////////////////////////////////////////////////////////////////////
// TIME

// ConvertTime returns time.Time and converts a ISO8601 string to a time.Time
// or empty string to time.Time{}
func ConvertTime(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	// Skip this hook if type is not time type
	if dest != timeType {
		return nilValue, nil
	}
	// Return value is source is already time type
	if v.Type() == timeType {
		return v, nil
	}
	// Skip if source is not a string
	if v.Kind() != reflect.String {
		return nilValue, nil
	}
	// Check for empty string which returns a time.Time{}
	if strings.TrimSpace(v.String()) == "" {
		return reflect.ValueOf(time.Time{}), nil
	}

	// Parse RFC3339 string to time.Time
	if t, err := time.Parse(time.RFC3339Nano, v.String()); err == nil {
		return reflect.ValueOf(t), nil
	} else {
		return nilValue, err
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Decoder) unmarshalscalar(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	for _, hook := range this.hooks {
		if value, err := hook(v, dest); err != nil {
			return nilValue, err
		} else if value.IsValid() {
			return value, nil
		}
	}
	return v, nil
}
