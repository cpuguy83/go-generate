package deepcopy

import "github.com/cpuguy83/go-generate/deepcopy/fixtures"

type stringType string

var stringTypeX = []byte(`
func(o stringType) Copy() stringType {
	oCopy := o

	return oCopy
}
`)

type sliceType []string

type arrayType [3]string

var arrayTypeX = []byte(`
func(o arrayType) Copy() arrayType {
	var oCopy arrayType
	for i0, v0 := range o {
		oCopy[i0] = v0
	}

	return oCopy
}
`)

type arrayOfArray [4][2]string

var arrayOfArrayX = []byte(`
func(o arrayOfArray) Copy() arrayOfArray {
	var oCopy arrayOfArray
	for i0, v0 := range o {
		for i1, v1 := range v0 {
			oCopy[i0][i1] = v1
		}
	}

	return oCopy
}
`)

var sliceTypeX = []byte(`
func(o sliceType) Copy() sliceType {
	oCopy := make(sliceType, len(o))
	for i0, v0 := range o {
		oCopy[i0] = v0
	}

	return oCopy
}
`)

type doubleSliceType [][]string

var doubleSliceTypeX = []byte(`
func(o doubleSliceType) Copy() doubleSliceType {
	oCopy := make(doubleSliceType, len(o))
	for i0, v0 := range o {
		if v0 != nil {
			oCopy[i0] = make([]string, len(v0))
			for i1, v1 := range v0 {
				oCopy[i0][i1] = v1
			}

		}

	}

	return oCopy
}
`)

type doubleSliceWithStructPtr [][]*simpleStruct

var doubleSliceWithStructPtrX = []byte(`
func(o doubleSliceWithStructPtr) Copy() doubleSliceWithStructPtr {
	oCopy := make(doubleSliceWithStructPtr, len(o))
	for i0, v0 := range o {
		if v0 != nil {
			oCopy[i0] = make([]*simpleStruct, len(v0))
			for i1, v1 := range v0 {
				if v1 != nil {
					var oCopy01 simpleStruct
					oCopy01 = *v1
					oCopy[i0][i1] = &oCopy01
				}

			}

		}

	}

	return oCopy
}
`)

type mapType map[string]string

var mapTypeX = []byte(`
func(o mapType) Copy() mapType {
	oCopy := make(mapType, len(o))
	for i0, v0 := range o {
		oCopy[i0] = v0
	}

	return oCopy
}
`)

type mapOfSlices map[string][]string

var mapOfSlicesX = []byte(`
func(o mapOfSlices) Copy() mapOfSlices {
	oCopy := make(mapOfSlices, len(o))
	for i0, v0 := range o {
		if v0 != nil {
			oCopy[i0] = make([]string, len(v0))
			for i1, v1 := range v0 {
				oCopy[i0][i1] = v1
			}

		}

	}

	return oCopy
}
`)

type mapOfMaps map[string]map[string]struct{}

var mapOfMapsX = []byte(`
func(o mapOfMaps) Copy() mapOfMaps {
	oCopy := make(mapOfMaps, len(o))
	for i0, v0 := range o {
		if v0 != nil {
			oCopy[i0] = make(map[string]struct{}, len(v0))
			for i1, v1 := range v0 {
				oCopy[i0][i1] = v1
			}

		}

	}

	return oCopy
}
`)

type simpleStruct struct {
	A string
	b string
}

var simpleStructX = []byte(`
func(o simpleStruct) Copy() simpleStruct {
	oCopy := o

	return oCopy
}
`)

type structWithEmbeddedPointer struct {
	A *struct{ B string }
}

var structWithEmbeddedPointerX = []byte(`
func(o structWithEmbeddedPointer) Copy() structWithEmbeddedPointer {
	oCopy := o
	if o.A != nil {
		var oCopy_A struct{ B string }
		oCopy_A = *o.A
		oCopy.A = &oCopy_A
	}

	return oCopy
}
`)

var structPointerX = []byte(`
func(o *simpleStruct) Copy() *simpleStruct {
	var oCopy simpleStruct
	oCopy = *o

	return &oCopy
}
`)

type anotherStruct struct {
	simpleStruct
	X map[string]*struct{ A *string }
	Y map[string]struct{ A *string }
	Z map[string]*string
}

type complexStruct struct {
	A string
	B map[string]int
	C []*simpleStruct
	D map[string]*simpleStruct
	E [][]*simpleStruct
	F [2][]*simpleStruct
	G [2]simpleStruct
	H *anotherStruct
}

var complexStructX = []byte(`
func (o complexStruct) Copy() complexStruct {
	oCopy := o
	if o.B != nil {
		oCopy.B = make(map[string]int, len(o.B))
		for i0, v0 := range o.B {
			oCopy.B[i0] = v0
		}

	}

	if o.C != nil {
		oCopy.C = make([]*simpleStruct, len(o.C))
		for i0, v0 := range o.C {
			if v0 != nil {
				var oCopy_C0 simpleStruct
				oCopy_C0 = *v0
				oCopy.C[i0] = &oCopy_C0
			}

		}

	}

	if o.D != nil {
		oCopy.D = make(map[string]*simpleStruct, len(o.D))
		for i0, v0 := range o.D {
			if v0 != nil {
				var oCopy_D0 simpleStruct
				oCopy_D0 = *v0
				oCopy.D[i0] = &oCopy_D0
			}

		}

	}

	if o.E != nil {
		oCopy.E = make([][]*simpleStruct, len(o.E))
		for i0, v0 := range o.E {
			if v0 != nil {
				oCopy.E[i0] = make([]*simpleStruct, len(v0))
				for i1, v1 := range v0 {
					if v1 != nil {
						var oCopy_E01 simpleStruct
						oCopy_E01 = *v1
						oCopy.E[i0][i1] = &oCopy_E01
					}

				}

			}

		}

	}

	for i0, v0 := range o.F {
		if v0 != nil {
			oCopy.F[i0] = make([]*simpleStruct, len(v0))
			for i1, v1 := range v0 {
				if v1 != nil {
					var oCopy_F01 simpleStruct
					oCopy_F01 = *v1
					oCopy.F[i0][i1] = &oCopy_F01
				}

			}

		}

	}
	for i0, v0 := range o.G {
		oCopy.G[i0] = v0
	}
	if o.H != nil {
		var oCopy_H anotherStruct
		oCopy_H = *o.H
		oCopy.H = &oCopy_H
		oCopy.H.simpleStruct = o.H.simpleStruct
		if o.H.X != nil {
			oCopy.H.X = make(map[string]*struct{ A *string }, len(o.H.X))
			for i0, v0 := range o.H.X {
				if v0 != nil {
					var oCopy_H0_X0 struct{ A *string }
					oCopy_H0_X0 = *v0
					oCopy.H.X[i0] = &oCopy_H0_X0
					if v0.A != nil {
						var oCopy_H0_X01_A string
						oCopy_H0_X01_A = *v0.A
						oCopy.H.X[i0].A = &oCopy_H0_X01_A
					}

				}

			}

		}

		if o.H.Y != nil {
			oCopy.H.Y = make(map[string]struct{ A *string }, len(o.H.Y))
			for i0, v0 := range o.H.Y {
				oCopy.H.Y[i0] = v0
				if v0.A != nil {
					var oCopy_H0_Y0_A string
					oCopy_H0_Y0_A = *v0.A
					oCopy.H.Y[i0].A = &oCopy_H0_Y0_A
				}

			}

		}

		if o.H.Z != nil {
			oCopy.H.Z = make(map[string]*string, len(o.H.Z))
			for i0, v0 := range o.H.Z {
				if v0 != nil {
					var oCopy_H0_Z0 string
					oCopy_H0_Z0 = *v0
					oCopy.H.Z[i0] = &oCopy_H0_Z0
				}

			}

		}

	}

	return oCopy
}
`)

type structWithImports struct {
	A *fixtures.Foo
}

var structWithImportsX = []byte(`
import (
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func(o structWithImports) Copy() structWithImports {
	oCopy := o
	if o.A != nil {
		var oCopy_A github_com_cpuguy83_go_generate_deepcopy_fixtures.Foo
		oCopy_A = *o.A
		oCopy.A = &oCopy_A
		if o.A.B != nil {
			oCopy.A.B = make(map[string]string, len(o.A.B))
			for i0, v0 := range o.A.B {
				oCopy.A.B[i0] = v0
			}

		}

	}

	return oCopy
}
`)

// structWithImportButNotNeeded has field from another package, but we don't need to actually
// call into that package.
// It ensures that we don't have unused imports in our import block
type structWithImportButNotNeeded struct {
	A fixtures.Foo
	B [1]fixtures.Foo
}

var structWithImportButNotNeededX = []byte(`
func(o structWithImportButNotNeeded) Copy() structWithImportButNotNeeded {
	oCopy := o
	oCopy.A = o.A
	if o.A.B != nil {
		oCopy.A.B = make(map[string]string, len(o.A.B))
		for i0, v0 := range o.A.B {
			oCopy.A.B[i0] = v0
		}

	}

	for i0, v0 := range o.B {
			oCopy.B[i0] = v0
			if v0.B != nil {
				oCopy.B[i0].B = make(map[string]string, len(v0.B))
				for i1, v1 := range v0.B {
					oCopy.B[i0].B[i1] = v1
				}

			}

		}

	return oCopy
}
`)

type structWithImportNeededMap struct {
	A map[string]fixtures.Foo
}

var structWithImportNeededMapX = []byte(`
import (
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func (o structWithImportNeededMap) Copy() structWithImportNeededMap {
	oCopy := o
	if o.A != nil {
		oCopy.A = make(map[string]github_com_cpuguy83_go_generate_deepcopy_fixtures.Foo, len(o.A))
		for i0, v0 := range o.A {
			oCopy.A[i0] = v0
			if v0.B != nil {
				oCopy.A[i0].B = make(map[string]string, len(v0.B))
				for i1, v1 := range v0.B {
					oCopy.A[i0].B[i1] = v1
				}

			}

		}

	}

	return oCopy
}`)

type structWithImportNeededSlice struct {
	A []fixtures.Foo
}

var structWithImportNeededSliceX = []byte(`
import (
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func (o structWithImportNeededSlice) Copy() structWithImportNeededSlice {
	oCopy := o
	if o.A != nil {
		oCopy.A = make([]github_com_cpuguy83_go_generate_deepcopy_fixtures.Foo, len(o.A))
		for i0, v0 := range o.A {
			oCopy.A[i0] = v0
			if v0.B != nil {
				oCopy.A[i0].B = make(map[string]string, len(v0.B))
				for i1, v1 := range v0.B {
					oCopy.A[i0].B[i1] = v1
				}

			}

		}

	}

	return oCopy
}
`)

type structWithImportedCustomSliceType struct {
	A fixtures.StrSlice
}

var structWithImportedCustomSliceTypeX = []byte(`
import (
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func (o structWithImportedCustomSliceType) Copy() structWithImportedCustomSliceType {
	oCopy := o
	if o.A != nil {
		oCopy.A = make(github_com_cpuguy83_go_generate_deepcopy_fixtures.StrSlice, len(o.A))
		for i0, v0 := range o.A {
			oCopy.A[i0] = v0
		}

	}

	return oCopy
}`)

type structWithUnexportedImportTypes struct {
	A *fixtures.Baz
}

type structWithImportsAndSimpleFields struct {
	A *fixtures.Quux
}

var structWithImportsAndSimpleFieldsX = []byte(`
import(
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func (o structWithImportsAndSimpleFields) Copy() structWithImportsAndSimpleFields {
	oCopy := o
	if o.A != nil {
		var oCopy_A github_com_cpuguy83_go_generate_deepcopy_fixtures.Quux
		oCopy_A = *o.A
		oCopy.A = &oCopy_A
	}

	return oCopy
}
`)

type structWithImportsAndUnsettableFields struct {
	A *fixtures.Banana
}

var structWithImportsAndUnsettableFieldsX = []byte(`
import (
	github_com_cpuguy83_go_generate_deepcopy_fixtures "github.com/cpuguy83/go-generate/deepcopy/fixtures"
)

func (o structWithImportsAndUnsettableFields) Copy() structWithImportsAndUnsettableFields {
	oCopy := o
	if o.A != nil {
		var oCopy_A github_com_cpuguy83_go_generate_deepcopy_fixtures.Banana
		oCopy_A = *o.A
		oCopy.A = &oCopy_A
	}

	return oCopy
}
`)

type structPtrWithCopyMethod struct {
	A *fixtures.Apple
}

var structPtrWithCopyMethodX = []byte(`
func(o structPtrWithCopyMethod) Copy() structPtrWithCopyMethod {
	oCopy := o
	oCopy.A = o.A.Copy()

	return oCopy
}
`)

type structWithCopyMethod struct {
	A fixtures.Apricot
}

var structWithCopyMethodX = []byte(`
func(o structWithCopyMethod) Copy() structWithCopyMethod {
	oCopy := o
	oCopy.A = o.A.Copy()

	return oCopy
}
`)

type structWithDeepCopy struct{}

func (s structWithDeepCopy) Copy() structWithDeepCopy {
	return s
}

// This is essentially a no-op since there are no fields on the struct
// The important part is that even though the struct already implements `DeepCopy`
// we need to be able to replace it since it is the top-level thing being generated.
var structWithDeepCopyX = []byte(`
func (o structWithDeepCopy) Copy() structWithDeepCopy {
	oCopy := o

	return oCopy
}
`)

type structWithChannel struct {
	c chan struct{}
}
