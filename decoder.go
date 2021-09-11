package marshaler

import (
	"fmt"
	"math"
	"net/url"
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
	nilValue           = reflect.ValueOf(nil)
	timeType           = reflect.TypeOf(time.Time{})
	durationType       = reflect.TypeOf(time.Duration(0))
	intType            = reflect.TypeOf(int(0))
	uintType           = reflect.TypeOf(uint(0))
	int64Type          = reflect.TypeOf(int64(0))
	uint64Type         = reflect.TypeOf(uint64(0))
	stringSliceType    = reflect.TypeOf([]string{})
	interfaceSliceType = reflect.TypeOf([]interface{}{})
	mapInterfaceType   = reflect.TypeOf(map[string]interface{}{})
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

// Decode decodes a map[string]interface{} type
func (this *Decoder) Decode(src, dest interface{}) error {
	if src == nil {
		return ErrBadParameter.With("Decode: nil value")
	}
	switch kind := reflect.ValueOf(src).Kind(); kind {
	case reflect.Map:
		return UnmarshalStruct(src, dest, this.name, this.unmarshalscalar)
	case reflect.Slice:
		return UnmarshalSlice(src, dest, this.unmarshalscalar)
	default:
		return ErrBadParameter.With("Decode: unable to decode ", kind)
	}
}

// DecodeQuery decodes a url.Values type
func (this *Decoder) DecodeQuery(src url.Values, dest interface{}) error {
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

// ConvertQueryValues returns a value from a []string
func ConvertQueryValues(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	// Skip this hook if type is not string slice type
	if v.Type() != stringSliceType {
		return nilValue, nil
	}
	// If source has length of zero return zero value
	if v.Len() == 0 {
		return reflect.Zero(dest), nil
	}
	// Support conversions to scalars and slices
	if dest.Kind() == reflect.Slice {
		return v, nil
	} else if v.Len() == 1 {
		return v.Index(0), nil
	}
	// Cannot convert
	return nilValue, fmt.Errorf("cannot convert %q to %q", v, dest)
}

// ConvertIntUint allows conversion from int value to a different int value,
// and uint value to a different uint value
func ConvertIntUint(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	// Skip this hook if type is not int or uint
	if v.Type() != intType && v.Type() != uintType {
		return nilValue, nil
	}
	// No conversion needed if destination is int or uint
	if v.Type() == dest {
		return v, nil
	}
	// Skip if can't convert
	if v.CanConvert(dest) == false {
		return nilValue, nil
	}
	// Check for bounds
	switch dest.Kind() {
	case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
		if v.CanConvert(dest) {
			return v.Convert(dest), nil
		}
	case reflect.Int8:
		if v.Int() >= math.MinInt8 && v.Int() <= math.MaxInt8 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	case reflect.Uint8:
		if v.Uint() <= math.MaxUint8 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	case reflect.Int16:
		if v.Int() >= math.MinInt16 && v.Int() <= math.MaxInt16 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	case reflect.Uint16:
		if v.Uint() <= math.MaxUint16 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	case reflect.Int32:
		if v.Int() >= math.MinInt32 && v.Int() <= math.MaxInt32 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	case reflect.Uint32:
		if v.Uint() <= math.MaxUint32 {
			return v.Convert(dest), nil
		} else {
			return nilValue, fmt.Errorf("value %v out of bounds for %v", v, dest)
		}
	}

	// Cannot convert
	return nilValue, fmt.Errorf("cannot convert %q to %q", v.Type(), dest)
}

// ConvertStringToNumber returns int, uint,float or bool from string
func ConvertStringToNumber(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	// Pass value through
	if v.Type() == dest {
		return v, nil
	}
	// Skip this hook if source is not string
	if v.Kind() != reflect.String {
		return nilValue, nil
	}
	// Convert to int, uint, float or bool
	switch dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, err := strconv.ParseInt(v.String(), 0, 64); err == nil {
			return reflect.ValueOf(value).Convert(dest), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value, err := strconv.ParseUint(v.String(), 0, 64); err == nil {
			return reflect.ValueOf(value).Convert(dest), nil
		}
	case reflect.Float32, reflect.Float64:
		if value, err := strconv.ParseFloat(v.String(), 64); err == nil {
			return reflect.ValueOf(value).Convert(dest), nil
		}
	case reflect.Bool:
		if value, err := strconv.ParseBool(v.String()); err == nil {
			return reflect.ValueOf(value).Convert(dest), nil
		}
	}
	// Skip
	return nilValue, nil
}

// ConvertMapInterface returns map[string]<type> from map[string]interface{} when all types
// within the interface match the destination type
func ConvertMapInterface(v reflect.Value, dest reflect.Type) (reflect.Value, error) {
	// Pass value through
	if v.Type() == dest {
		return v, nil
	}
	// Skip this hook if source is not map[string]interface{}
	if v.Type() != mapInterfaceType {
		return nilValue, nil
	}
	// Iterate through types in source map, skip if any type is not the same as destination type
	d := reflect.MakeMap(dest)
	for _, key := range v.MapKeys() {
		elem := v.MapIndex(key)
		if elem.Kind() == reflect.Interface && elem.CanInterface() {
			elem = reflect.ValueOf(elem.Interface())
		}
		if elem.Type() != dest.Elem() {
			return nilValue, fmt.Errorf("value of type %v in map cannot be converted to %v", elem.Type(), dest.Elem())
		} else {
			d.SetMapIndex(key, elem)
		}
	}

	// Return converted map
	return d, nil
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
