package marshaler_test

import (
	"fmt"
	"math/rand"
	"net/url"
	"testing"
	"time"

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
		if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
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
	if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
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
		Float32 float32 `yaml:"float32"`
		Float64 float64 `yaml:"float64"`
	}
	for i := 0; i < 100; i++ {
		a["float32"] = rand.Float32()
		a["float64"] = rand.Float64()
		if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
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
		Strings []string `yaml:"strings"`
		Nil     []string `yaml:"nil"`
	}
	if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
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
		Q []string `yaml:"q"`
	}
	for i := 0; i < 100; i++ {
		a.Add("q", fmt.Sprint(i))
		if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
			t.Fatal(err)
		} else if stringSliceEqual(a["q"], b.Q) == false {
			t.Error("a != b", a, b)
		}
	}
}

func Test_Unmarshall_007(t *testing.T) {
	a := url.Values{}
	var b struct {
		Q []string `yaml:"q"`
	}
	for i := 0; i < 100; i++ {
		a.Add("q", fmt.Sprint(i))
		if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
			t.Fatal(err)
		} else if stringSliceEqual(a["q"], b.Q) == false {
			t.Error("a != b")
		}
	}
}

func Test_Unmarshall_008(t *testing.T) {
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
		"time":   time.Duration(0) * time.Second,
		"str1":   "0s",
		"str2":   "0",
	}
	var b struct {
		A time.Duration `yaml:"int"`
		B time.Duration `yaml:"uint"`
		C time.Duration `yaml:"int8"`
		D time.Duration `yaml:"uint8"`
		E time.Duration `yaml:"int16"`
		F time.Duration `yaml:"uint16"`
		G time.Duration `yaml:"int32"`
		H time.Duration `yaml:"uint32"`
		I time.Duration `yaml:"int64"`
		J time.Duration `yaml:"uint64"`
		K time.Duration `yaml:"time"`
		L time.Duration `yaml:"str1"`
		M time.Duration `yaml:"str2"`
	}
	for i := 0; i < 127; i++ {
		v := time.Duration(i) * time.Second
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
		a["time"] = v
		a["str1"] = fmt.Sprint(i, "s")
		a["str2"] = fmt.Sprint(i)
		if err := marshaler.UnmarshalStruct(a, &b, "yaml", marshaler.ConvertDuration); err != nil {
			t.Fatal(err)
		}
		if b.A != v {
			t.Errorf("int %v != %v", b.A, v)
		}
		if b.B != v {
			t.Errorf("uint %v != %v", b.B, v)
		}
		if b.C != v {
			t.Errorf("int8 %v != %v", b.C, v)
		}
		if b.D != v {
			t.Errorf("uint8 %v != %v", b.D, v)
		}
		if b.E != v {
			t.Errorf("int16 %v != %v", b.E, v)
		}
		if b.F != v {
			t.Errorf("uint16 %v != %v", b.F, v)
		}
		if b.G != v {
			t.Errorf("int32 %v != %v", b.G, v)
		}
		if b.H != v {
			t.Errorf("uint32 %v != %v", b.H, v)
		}
		if b.I != v {
			t.Errorf("int64 %v != %v", b.I, v)
		}
		if b.J != v {
			t.Errorf("uint64 %v != %v", b.J, v)
		}
		if b.K != v {
			t.Errorf("time %v != %v", b.K, v)
		}
		if b.L != v {
			t.Errorf("str1 %v != %v", b.L, v)
		}
		if b.M != v {
			t.Errorf("str2 %v != %v", b.M, v)
		}
	}
}

func Test_Unmarshall_009(t *testing.T) {
	var a = map[string]interface{}{
		"slice": []interface{}{"a", "b", "c"},
	}
	var b struct {
		Slice []string `yaml:"slice"`
	}
	if err := marshaler.UnmarshalStruct(a, &b, "yaml", nil); err != nil {
		t.Fatal(err)
	} else if stringSliceEqual([]string{"a", "b", "c"}, b.Slice) == false {
		t.Error("a != b")
	}
}

func Test_Unmarshall_010(t *testing.T) {
	var src = []string{"a", "b", "c"}
	var dest = []string{}
	if err := marshaler.UnmarshalSlice(src, &dest, nil); err != nil {
		t.Fatal(err)
	} else if stringSliceEqual(src, dest) == false {
		t.Error("src != dest")
	}
}

func Test_Unmarshall_011(t *testing.T) {
	var src = []string{"1", "2", "3", "4"}
	var dest = []uint{}
	if err := marshaler.UnmarshalSlice(src, &dest, marshaler.ConvertStringToNumber); err != nil {
		t.Fatal(err)
	} else if fmt.Sprint(src) != fmt.Sprint(dest) {
		t.Error("src != dest", src, dest)
	}
}

func Test_Unmarshall_012(t *testing.T) {
	var src = map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	var dest1 = map[string]string{}
	var dest2 = map[string]int{}
	var dest3 map[string]float32
	if err := marshaler.UnmarshalStruct(src, &dest1, "yaml", nil); err != nil {
		t.Fatal(err)
	} else if len(src) != len(dest1) {
		t.Error("src != dest1", src, dest1)
	} else {
		t.Log(src, dest1)
	}
	if err := marshaler.UnmarshalStruct(src, &dest2, "yaml", marshaler.ConvertStringToNumber); err != nil {
		t.Fatal(err)
	} else if len(src) != len(dest2) {
		t.Error("src != dest2", src, dest2)
	} else {
		t.Log(src, dest2)
	}
	if err := marshaler.UnmarshalStruct(src, &dest3, "yaml", marshaler.ConvertStringToNumber); err != nil {
		t.Fatal(err)
	} else if len(src) != len(dest3) {
		t.Error("src != dest3", src, dest3)
	} else {
		t.Log(src, dest3)
	}
}

func Test_Unmarshall_013(t *testing.T) {
	type A struct {
		FieldA int `yaml:"a"`
	}
	type B struct {
		FieldB int `yaml:"b"`
	}
	var dest struct {
		A
		B
		FieldC int `yaml:"c"`
	}
	var src = map[string]string{
		"a": "55",
		"b": "66",
		"c": "77",
	}
	if err := marshaler.UnmarshalStruct(src, &dest, "yaml", marshaler.ConvertStringToNumber); err != nil {
		t.Fatal(err)
	} else if dest.FieldA != 55 || dest.FieldB != 66 || dest.FieldC != 77 {
		t.Error("src != dest", src, dest)
	} else {
		t.Log(src, dest)
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
