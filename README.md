# go-marshaler

This utility module currently decodes a map of values into a go struct. You
could use this for converting "flat" wire formats (for example, from an SQL 
database or HTTP response into go structures.

## UnmarshalStruct

The simple method is to use "UnmarshalStruct" to perform the transformation. For example,

```go
package main

import (
  marshaler "github.com/djthorpe/go-marshaler"
)

func main() {
  var dest struct {
    A int     `test:"a"`
    B float32 `test:"b"`
    C []int   `test:"c"`
  }
  src := map[string]interface{}{
    "a": int(10),
    "b": float32(3.1415),
    "c": []int{ 1, 2, 3 },
  }
  if err := marshaler.UnmarshalStruct(src, &dest, "test", nil); err != nil {
    panic(err)
  }
  fmt.Println(dest)
}
```

Will output `map[string]interface{}{ a: 10, b: 3.1415, c: [ 1, 2, 3 ] }`. You can define 
a custom scalar decoding function which can transform a source value into a destination type.
Pass this custom function as the last argument to `UnmarshalStruct`. For example,

```go
func CustomTransformer(src reflect.Value, dest reflect.Type) (reflect.Value,error) {
  if dest != /* type I am interested in transforming to */ {
		return reflect.ValueOf(nil), nil
	}
  // Return type if it's already converted
  if src.Type() == dest {
		return v, nil
	}
  // Do transformation here to dest type, return error if the
  // source cannot be transformed....
}
```

Your function should return an invalid value to skip the transformation, and an error
if you want to return an error out of the `UnmarshalStruct` function.

## Decode

Decoding can also be performed as follows:

```go
package main

import (
  marshaler "github.com/djthorpe/go-marshaler"
)

func main() {
  var dest struct {
    A int       `test:"a"`
    B float32   `test:"b"`
    C []int     `test:"c"`
    D time.Time `test:"d"`
  }
  src := map[string]interface{}{
    "a": int(10),
    "b": float32(3.1415),
    "c": []int{ 1, 2, 3 },
    "d": "2016-01-01T00:00:00Z",
  }

  dec := marshaler.NewDecoder("test",marshaler.ConvertTime)

  if err := dec.Decode(src, &dest); err != nil {
    panic(err)
  }
  fmt.Println(dest)
}
```

The `NewDecoder` method takes one or more custom functions which can be used for decoding
fields. In this example, the `ConvertTime` decoder will convert a string into a `time.Time`
structure. All the custom functions provided are:

  * `marshaler.ConvertTime` Converts strings formatted as RFC3339 into `time.Time`. Empty
    strings are converted into zero-time.

