package fixtures

// Foo is a fixture used to test code generation for imported packages
type Foo struct {
	A string
	B map[string]string
}

type bar struct {
	Z string
}

// Baz is also just a fixture, trying to make it's way in the world.
type Baz struct {
	B *bar
}

// Quux is a fixture that wishes it was part of Foo.
type Quux struct {
	a string
}

// Banana is a fixture that is yellow.
type Banana struct {
	a map[string]string
}

// Apple is a fixture that doesn't taste as yummy as Banana.
type Apple struct {
	A string
}

// Copy is the function Microsoft uses to make it's products :trollface:
func (a *Apple) Copy() *Apple {
	var copy Apple
	copy = *a
	return &copy
}

// Apricot is actually really good, but not very common.
type Apricot struct {
	A string
}

// Copy is the function that Microsoft accidentally used instead of Apple.Copy
// once and out came Microsoft Bob.
func (a Apricot) Copy() Apricot {
	return a
}

// StrSlice is a slice of strings
type StrSlice []string

type Unsettable struct {
	a string
}
