package marshaler_test

import (
	"testing"
	"time"

	"github.com/djthorpe/go-marshaler"
)

func Test_Encoder_001(t *testing.T) {
	if enc := marshaler.NewEncoder("yaml"); enc == nil {
		t.Fatal("Unexpected nil return from NewEncoder")
	}
}

func Test_Encoder_002(t *testing.T) {

	type ts struct {
		Time   time.Time `yaml:"timestamp"`
		Int    int       `yaml:"-"`
		String string    `yaml:",a,b,c"`
		Bool   bool      `yaml:"bool,a,b,c"`
		Struct struct {
			A string `yaml:"a"`
			B string
		} `yaml:"struct"`
	}

	fields := marshaler.NewEncoder("yaml").Reflect(ts{})
	if fields == nil {
		t.Fatal("Unexpected nil return from NewEncoder")
	}
	for _, field := range fields {
		t.Log(field)
	}
}
