package deepcopy

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func TestGenerate(t *testing.T) {
	type run struct {
		explain string
		test    interface{}
		x       []byte
		err     error
		ignore  []interface{}
	}
	type mapType map[string]string
	cases := []run{
		{"A string type", stringType(""), stringTypeX, nil, nil},
		{"An array type", arrayType{}, arrayTypeX, nil, nil},
		{"An array of arrays type", arrayOfArray{}, arrayOfArrayX, nil, nil},
		{"A simple slice type", sliceType{}, sliceTypeX, nil, nil},
		{"A 2-D slice type", doubleSliceType{}, doubleSliceTypeX, nil, nil},
		{"A 2-D slice with struct ptr", doubleSliceWithStructPtr{}, doubleSliceWithStructPtrX, nil, nil},
		{"A simple map type", mapType{}, mapTypeX, nil, nil},
		{"A map of slices", mapOfSlices{}, mapOfSlicesX, nil, nil},
		{"A map of maps", mapOfMaps{}, mapOfMapsX, nil, nil},
		{"A simple struct", simpleStruct{}, simpleStructX, nil, nil},
		{"A struct with an embedded struct pointer", structWithEmbeddedPointer{}, structWithEmbeddedPointerX, nil, nil},
		{"A struct pointer", &simpleStruct{}, structPointerX, nil, nil},
		{"A complex struct with mixed reference types", complexStruct{}, complexStructX, nil, nil},
		{"A struct which imports from another package", structWithImports{}, structWithImportsX, nil, nil},
		{"A struct with imports that are unexported in another pkg", structWithUnexportedImportTypes{}, nil, ErrUnexportedType, nil},
		{"A struct which imports from another package with unexported but simple fields", structWithImportsAndSimpleFields{}, structWithImportsAndSimpleFieldsX, nil, nil},
		{"A struct which imports from another package with unsettable fields", structWithImportsAndUnsettableFields{}, nil, ErrUnsettableField, nil},
		{"A struct which imports from another package with unsettable fields that are ignored", structWithImportsAndUnsettableFields{}, structWithImportsAndUnsettableFieldsX, nil, []interface{}{fixtures.Banana{}}},
		{"A struct which implements DeepCopy", structWithDeepCopy{}, structWithDeepCopyX, nil, nil},
		{"A struct with an imported ptr struct which implements DeepCopy", structPtrWithCopyMethod{}, structPtrWithCopyMethodX, nil, nil},
		{"A struct with an imported struct which implements DeepCopy", structWithCopyMethod{}, structWithCopyMethodX, nil, nil},
		{"A struct with an imported struct that does not require an import statement", structWithImportButNotNeeded{}, structWithImportButNotNeededX, nil, nil},
		{"A struct with an imported struct in a map that needs an import statement", structWithImportNeededMap{}, structWithImportNeededMapX, nil, nil},
		{"A struct with an imported struct in a slice that needs an import statement", structWithImportNeededSlice{}, structWithImportNeededSliceX, nil, nil},
		{"A struct that uses an imported custom slice type", structWithImportedCustomSliceType{}, structWithImportedCustomSliceTypeX, nil, nil},
		{"A struct type with a channel", structWithChannel{}, nil, ErrUnsupportedType, nil},
	}

	for _, c := range cases {
		t.Run(c.explain, func(t *testing.T) {
			imports, copyFunc, err := Generate("o", c.test, c.ignore)
			if err := cause(err); err != c.err {
				t.Fatalf("%s: expected '%v', got: %v", c.explain, c.err, err)
			}

			actual, err := format.Source(append(imports, copyFunc...))
			if err != nil {
				t.Fatal(err.Error() + "\n" + string(copyFunc))
			}

			xFmt, err := format.Source(c.x)
			if err != nil {
				t.Fatalf("%s: %v\n\n%s", c.explain, err, string(c.x))
			}
			if !bytes.Equal(bytes.TrimSpace(actual), bytes.TrimSpace(xFmt)) {
				t.Fatalf("%s: expected: \n%s\n\ngot: \n%s\n\n", c.explain, string(xFmt), string(actual))
			}
		})
	}
}
