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

Will output `map[string]interface{}{ a: 0, b: 0.0 }`.


