package marshaler_test

import (
	"fmt"
	"math/rand"
	"net/url"
	"testing"

	marshaler "github.com/djthorpe/go-marshaler"
)

func Test_Unmarshall_001(t *testing.T) {
	t.Log(t.Name())
}

func Test_Unmarshall_002(t *testing.T) {
	var a = map[string]interface{}{
		"int":    int(0),
		"uint":   uint(0),
		"int8":   int8(0),
		"uint8":  uint8(0),
		"int16":  int16(0),
		"uint16": uint16(0),
		"int32":  int32(0),
		"uint32": uint32(0),
		"int64":  int64(0),
		"uint64": uint64(0),
	}
	var b struct {
		A int    `yaml:"int"`
		B uint   `yaml:"uint"`
		C int8   `yaml:"int8"`
		D uint8  `yaml:"uint8"`
		E int16  `yaml:"int16"`
		F uint16 `yaml:"uint16"`
		G int32  `yaml:"int32"`
		H uint32 `yaml:"uint32"`
		I int64  `yaml:"int64"`
		J uint64 `yaml:"uint64"`
	}
	for i := 0; i < 1000; i++ {
		a["int"] = i
		a["uint"] = uint(i)
		a["int8"] = int8(i)
		a["uint8"] = uint8(i)
		a["int16"] = int16(i)
		a["uint16"] = uint16(i)
		a["int32"] = int32(i)
		a["uint32"] = uint32(i)
		a["int64"] = int64(i)
		a["uint64"] = uint64(i)
		if err := marshaler.UnmarshalStruct(a, &b, "yaml"); err != nil {
			t.Fatal(err)
		}
		if b.A != i {
			t.Errorf("int %v != %v", b.A, i)
		}
		if b.B != uint(i) {
			t.Errorf("uint %v != %v", b.B, i)
		}
		if b.C != int8(i) {
			t.Errorf("int8 %v != %v", b.C, i)
		}
		if b.D != uint8(i) {
			t.Errorf("uint8 %v != %v", b.D, i)
		}
		if b.E != int16(i) {
			t.Errorf("int16 %v != %v", b.E, i)
		}
		if b.F != uint16(i) {
			t.Errorf("uint16 %v != %v", b.F, i)
		}
		if b.G != int32(i) {
			t.Errorf("int32 %v != %v", b.G, i)
		}
		if b.H != uint32(i) {
			t.Errorf("uint32 %v != %v", b.H, i)
		}
		if b.I != int64(i) {
			t.Errorf("int64 %v != %v", b.I, i)
		}
		if b.J != uint64(i) {
			t.Errorf("uint64 %v != %v", b.J, i)
		}
	}
}

func Test_Unmarshall_003(t *testing.T) {
	var a = map[string]interface{}{
		"true":  true,
		"false": false,
	}
	var b struct {
		T bool `yaml:"true"`
		F bool `yaml:"false"`
	}
	if err := marshaler.UnmarshalStruct(a, &b, "yaml"); err != nil {
		t.Fatal(err)
	}
	if b.T != true {
		t.Errorf("true %v != %v", b.T, true)
	}
	if b.F != false {
		t.Errorf("false %v != %v", b.T, false)
	}
}

func Test_Unmarshall_004(t *testing.T) {
	var a = map[string]interface{}{
		"float32": float32(0),
		"float64": float64(0),
	}
	var b struct {
		Float32 float32
		Float64 float64
	}
	for i := 0; i < 100; i++ {
		a["float32"] = rand.Float32()
		a["float64"] = rand.Float64()
		if err := marshaler.UnmarshalStruct(a, &b, "yaml"); err != nil {
			t.Fatal(err)
		}
		if b.Float32 != a["float32"] {
			t.Errorf("float32 %v != %v", b.Float32, a["float32"])
		}
		if b.Float64 != a["float64"] {
			t.Errorf("float64 %v != %v", b.Float64, a["float64"])
		}
	}
}

func Test_Unmarshall_005(t *testing.T) {
	var a = map[string]interface{}{
		"strings": []string{"a", "b", "c"},
		"nil":     []string{},
	}
	var b struct {
		Strings []string
		Nil     []string
	}
	if err := marshaler.UnmarshalStruct(a, &b, "yaml"); err != nil {
		t.Fatal(err)
	}
	if !stringSliceEqual(b.Strings, a["strings"].([]string)) {
		t.Errorf("strings %v != %v", b.Strings, a["strings"])
	}
	if !stringSliceEqual(b.Nil, a["nil"].([]string)) {
		t.Errorf("strings %v != %v", b.Strings, a["strings"])
	}
}

func Test_Unmarshall_006(t *testing.T) {
	a := url.Values{}
	var b struct {
		Q []string
	}
	for i := 0; i < 100; i++ {
		a.Add("q", fmt.Sprint(i))
		if err := marshaler.UnmarshalStruct(a, &b, ""); err != nil {
			t.Fatal(err)
		} else if stringSliceEqual(a["q"], b.Q) == false {
			t.Error("a != b")
		}
	}
}

func Test_Unmarshall_007(t *testing.T) {
	a := url.Values{}
	var b struct {
		Q string
	}
	for i := 0; i < 100; i++ {
		a.Add("q", fmt.Sprint(i))
		if err := marshaler.UnmarshalStruct(a, &b, ""); err != nil {
			t.Fatal(err)
		} else if stringSliceEqual(a["q"], []string{b.Q}) == false {
			t.Error("a != b")
		}
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}