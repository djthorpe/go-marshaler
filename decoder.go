package marshaler

import (
	"fmt"
	"reflect"
	"strconv"
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
	nilValue     = reflect.ValueOf(nil)
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
	int64Type    = reflect.TypeOf(int64(0))
	uint64Type   = reflect.TypeOf(uint64(0))
	stringType   = reflect.TypeOf(string(""))
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

// ConvertDuration returns time.Duration from integer, string or time.Duration
func ConvertDuration(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	fmt.Println("ConvertDuration", v, v.Type(), "=>", dest)
	// Skip this hook if type is not time type
	if dest != durationType {
		return nilValue, nil
	}
	// Return value is source is already time type
	if v.Type() == durationType {
		return v, nil
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := v.Convert(int64Type).Int()
		return reflect.ValueOf(time.Duration(v) * time.Second), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := v.Convert(uint64Type).Uint()
		return reflect.ValueOf(time.Duration(v) * time.Second), nil
	case reflect.String:
		if v_, err := time.ParseDuration(v.String()); err == nil {
			return reflect.ValueOf(v_), nil
		} else if v_, err := strconv.ParseUint(v.String(), 0, 64); err == nil {
			return reflect.ValueOf(time.Duration(v_) * time.Second), nil
		} else {
			return nilValue, fmt.Errorf("cannot convert %q to time.Duration", v.String())
		}
	}
	return nilValue, fmt.Errorf("cannot convert %q to time.Duration", v.Kind())
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
