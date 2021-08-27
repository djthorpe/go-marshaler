package marshaler_test

import (
	"testing"
	"time"

	"github.com/djthorpe/go-marshaler"
)

type ts struct {
	Time     time.Time     `yaml:"timestamp"`
	Duration time.Duration `yaml:"duration"`
}

func Test_Decoder_001(t *testing.T) {
	if dec := marshaler.NewDecoder("yaml"); dec == nil {
		t.Fatal("Unexpected nil return from NewDecoder")
	}
}

func Test_Decoder_002(t *testing.T) {
	dest := ts{}
	src := map[string]interface{}{
		"timestamp": "2016-01-01T00:00:00Z",
	}
	if err := marshaler.NewDecoder("yaml", marshaler.ConvertTime).Decode(src, &dest); err != nil {
		t.Fatal(err)
	} else if dest.Time.Equal(time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)) == false {
		t.Fatal("Unexpected time value")
	}
}
func Test_Decoder_003(t *testing.T) {
	dest := ts{}
	src := map[string]interface{}{
		"timestamp": "   ",
	}
	if err := marshaler.NewDecoder("yaml", marshaler.ConvertTime).Decode(src, &dest); err != nil {
		t.Fatal(err)
	} else if dest.Time.IsZero() == false {
		t.Fatal("Unexpected non-zero time")
	}
}

func Test_Decoder_004(t *testing.T) {
	dest := ts{}
	src := map[string]interface{}{
		"timestamp": "2006",
	}
	if err := marshaler.NewDecoder("yaml", marshaler.ConvertTime).Decode(src, &dest); err == nil {
		t.Fatal("Expected error")
	}
}

func Test_Decoder_005(t *testing.T) {
	dest := ts{}
	src := map[string]interface{}{}
	tests := []struct {
		src  interface{}
		dest time.Duration
	}{
		{"100ms", 100 * time.Millisecond},
		{"10", 10 * time.Second},
		{"50h", 50 * time.Hour},
		{int(314), 314 * time.Second},
	}
	for _, test := range tests {
		src["duration"] = test.src
		if err := marshaler.NewDecoder("yaml", marshaler.ConvertTime, marshaler.ConvertDuration).Decode(src, &dest); err != nil {
			t.Fatal(err)
		} else if dest.Duration != test.dest {
			t.Fatal("Unexpected value", dest.Duration, " expected ", test.dest)
		}
	}
}
