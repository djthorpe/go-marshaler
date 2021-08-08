# go-marshaler

This utility module currently decodes a map of values into a go struct.
For example,

```go
package main

import (
  marshaler "github.com/djthorpe/go-marshaler"
)

func main() {
  var a struct {
    A int `test:"a"`
    B float32 `test:"b"`
  }
  b := map[string]interface{}{}
  if err := marshaler.UnmarshalStruct(a,b,"test",nil); err != nil {
    panic(err)
  }
  fmt.Println(b)
}
```

Will output `map[string]interface{}{ a: 0, b: 0.0 }`. You can define a custom
scalar decoding function which can transform a source value into a destination type.
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

