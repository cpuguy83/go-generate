## DeepCopy

DeepCopy is a library that can be used to generate a function on any given object
which will perform a deep copy of that object.

No CLI is provided as you will need to supply your own types to the copy generator,
which cannot come from CLI arguments directly, however you can create your own
CLI tailored for your use-case. See the example usage below.

Some types are unsupported, such as `chan` types, since it does not make sense
to copy these.

Some types will use unexported types from another package, which this tool cannot
generate a copy function for, in such a case it will generate an error. You can
choose to ignore these errors by passing in the list of types you wish to ignore
these kinds of errors from.  
One good example of this is `time.Time`, which is a struct that has a field `loc`
which is a pointer to another struct. Naturally we can't make a copy of this since
the field is unexported so we need to either move the field to another object not
being copied, or ignore errors from `time.Time{}` (shown in the example below)

### Example Usage

```go
package main

import (
	"fmt"
	"go/format"
	"os"
	"time"

	"github.com/cpuguy83/go-generate/deepcopy"
)

type Foo struct {
	a string
	B string
	C map[string]string
	D time.Time
}

func main() {
	imports, fn, err := deepcopy.Generate("o", &Foo{}, []interface{}{time.Time{}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := []byte("\n\npackage hello\n\n")
	b = append(b, imports...)
	b = append(b, fn...)
	code, _ := format.Source(b)
	os.Stdout.Write(code)
}
```

```
$ go run main.go
package hello

func (o *Foo) Copy() *Foo {
	var oCopy Foo
	oCopy = *o
	if o.C != nil {
		oCopy.C = make(map[string]string, len(o.C))
		for i0, v0 := range o.C {
			oCopy.C[i0] = v0
		}

	}

	oCopy.D = o.D

	return &oCopy
}
```


### TODO

