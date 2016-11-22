package deepcopy

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// reflect type wrapts reflect.Type with functionality for traversing the type
// and links it to its parent.
type reflectType struct {
	parent *reflectType
	reflect.Type

	// index is a value that is incremented at each walk step for the next type
	// when the current type is indexable.
	// This is done to keep track of indexes in maps/slices/arrays.
	index int

	// fieldIndex is used by struct field types to get the correct field information
	// (e.g. field name) from the parent
	fieldIndex int
}

// Next returns traverses the item to the next thing
// e.g., if this is a struct, it would return the next struct field
func (t *reflectType) Next() *reflectType {
	next := &reflectType{
		parent: t,
		index:  t.index,
	}

	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		next.Type = t.Elem()
	case reflect.Struct:
		if t.NumField() > 0 {
			field := t.Field(0)
			next.Type = field.Type
		}
	default:
		if t.parent != nil {
			switch t.parent.Kind() {
			case reflect.Struct:
				if t.fieldIndex < t.parent.NumField()-1 {
					next.fieldIndex = t.fieldIndex + 1
					field := t.parent.Field(next.fieldIndex)
					next.Type = field.Type
				}
			}
		}
	}

	if isKind(t, reflect.Map, reflect.Slice, reflect.Array) {
		next.index++
	}

	return next
}

// getCopyName generates the var declaration for the given type
func getCopyName(root, copyRoot string, t *reflectType) (copyStr, copyVal, varStr string) {
	if t == nil {
		return copyRoot, root, copyRoot
	}

	copyStr, copyVal, varStr = getCopyName(root, copyRoot, t.parent)

	if t.parent == nil {
		return
	}

	switch t.parent.Kind() {
	case reflect.Struct:
		if t.fieldIndex <= t.parent.NumField()-1 {
			copyStr += "." + t.parent.Field(t.fieldIndex).Name
			copyVal += "." + t.parent.Field(t.fieldIndex).Name
			varStr += "_" + t.parent.Field(t.fieldIndex).Name
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		copyStr += "[i" + strconv.Itoa(t.parent.index) + "]"
		copyVal = "v" + strconv.Itoa(t.parent.index)
		varStr += strconv.Itoa(t.parent.index)
	case reflect.Ptr:
		varStr += strconv.Itoa(t.index)
	}

	return
}

// Generate is used to generate a function for any object that will create a
// deep copy of that object.
// it takes:
//   - `ref` is a string which is the name of the function reciever should use (e.g. `o` would be `func(o MyObject)`
//   - `o` is the object which you want to generate a copy function for
//   - `ignorePkgErrs` takes a list of objects that you would like to ignore errors related to non-accessible types in other packages
// It returns the neccessary import statements and the generated copy function to use.
func Generate(ref string, o interface{}, ignorePkgErrs []interface{}) (importsBuf []byte, copyFnBuf []byte, err error) {
	imports := make(map[string]struct{})
	buf := bytes.NewBuffer(nil)
	root := &reflectType{parent: nil, Type: reflect.TypeOf(o)}
	rootPkg := getPkgName(root.Type)
	baseCopy := ref + "Copy"

	ignored := make(map[reflect.Type]bool, len(ignorePkgErrs))
	for _, i := range ignorePkgErrs {
		ignored[reflect.TypeOf(i)] = true
	}

	var generate func(t *reflectType) error
	generate = func(t *reflectType) error {
		if t == nil {
			return nil
		}

		copyStr, copyVal, varStr := getCopyName(ref, baseCopy, t)
		getPkgName(t)

		if root != t && hasCopyMethod(t.Type) {
			copyStr, copyVal, _ := getCopyName(ref, baseCopy, t)
			_, err := buf.Write([]byte(copyStr + " = " + copyVal + ".Copy()\n"))
			return err
		}

		switch t.Kind() {
		case reflect.Chan:
			return wrapErr(ErrUnsupportedType, "cannot make copy of channel types")
		case reflect.Struct:
			equals := ":="
			if t.parent != nil {
				equals = "="
			}

			if !isKind(t.parent, reflect.Ptr) {
				buf.Write([]byte(fmt.Sprintf("%s %s %s\n", copyStr, equals, copyVal)))
			}

			// go through each struct field and generate copies for that type
			//
			// TODO(cpuguy83): Why doesn't this work properly in Next()?
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if dcTagVal, ok := field.Tag.Lookup("deepcopy"); ok && dcTagVal == "skip" {
					continue
				}

				next := &reflectType{
					parent:     t,
					Type:       field.Type,
					fieldIndex: i,
					index:      t.index,
				}
				if nextPkg := getPkgName(next.Type); nextPkg != "" && nextPkg != rootPkg {
					name := getName(next.Type, rootPkg)
					if ln := strings.ToLower(name[0:]); ln == name[0:] {
						if ignored[next.Type] {
							continue
						}
						return wrapErr(ErrUnexportedType, fmt.Sprintf("cannot use type: %s", next.Type))
					}
				}

				if curPkg := getPkgName(t.Type); curPkg != rootPkg {
					if field.Name[0:] == strings.ToLower(field.Name[0:]) && isKind(next, reflect.Map, reflect.Ptr, reflect.Array, reflect.Slice) {
						if ignored[t.Type] {
							continue
						}
						_, copyVal, _ = getCopyName(ref, baseCopy, next)
						return wrapErr(ErrUnsettableField, fmt.Sprintf("cannot make copy of type '%v' with unexported field in another package: %s", t.Type, copyVal))
					}
				}

				if err := generate(next); err != nil {
					return err
				}
			}
			return nil
		case reflect.Ptr:
			if t.parent == nil {
				_, err := buf.Write([]byte(fmt.Sprintf("var %s %s\n", varStr, getName(t.Type.Elem(), rootPkg))))
				if err != nil {
					return err
				}
			} else {
				buf.Write([]byte(fmt.Sprintf("if %s != nil {", copyVal)))
				_, err := buf.Write([]byte(fmt.Sprintf("var %s %s\n", varStr, getName(t.Type.Elem(), rootPkg))))
				if err != nil {
					return err
				}
			}
			equals := ":="
			if t.parent != nil || copyStr == varStr {
				equals = "="
			}
			buf.Write([]byte(fmt.Sprintf("%s %s *%s\n", varStr, equals, copyVal)))
			if t.parent != nil {
				buf.Write([]byte(fmt.Sprintf("%s = &%s\n", copyStr, varStr)))
			}

			next := t.Next()
			addImport(next, rootPkg, imports)
			if err := generate(next); err != nil {
				return err
			}
			if t.parent != nil {
				_, err = buf.Write([]byte{'}', '\n', '\n'})
			}
			return err
		case reflect.Array:
			if t.parent == nil {
				buf.Write([]byte(fmt.Sprintf("var %s %s\n", varStr, getName(t.Type, rootPkg))))
			}
			s := fmt.Sprintf("for i%d, v%d := range %s {\n", t.index, t.index, copyVal)
			_, err := buf.Write([]byte(s))
			if err != nil {
				return err
			}
			if err := generate(t.Next()); err != nil {
				return err
			}
			_, err = buf.Write([]byte{'}', '\n'})
			return err
		case reflect.Map, reflect.Slice:
			next := t.Next()
			addImport(t, rootPkg, imports)
			addImport(next, rootPkg, imports)
			var s string
			name := getName(t.Type, rootPkg)
			if t.parent == nil {
				s = fmt.Sprintf("%s := make(%s, len(%s))\n", copyStr, name, copyVal)
			} else {
				s = fmt.Sprintf(`if %s != nil {
				%s = make(%s, len(%s))
				`, copyVal, copyStr, name, copyVal)
			}
			_, err := buf.Write([]byte(s))
			if err != nil {
				return err
			}

			s = fmt.Sprintf("for i%d, v%d := range %s {\n", t.index, t.index, copyVal)
			_, err = buf.Write([]byte(s))
			if err != nil {
				return err
			}

			if err := generate(next); err != nil {
				return err
			}
			_, err = buf.Write([]byte{'}', '\n'})
			if err != nil {
				return err
			}
			if t.parent != nil {
				_, err = buf.Write([]byte{'\n', '}', '\n', '\n'})
			}
			return err
		default:
			if isKind(t.parent, reflect.Struct, reflect.Ptr) {
				return nil
			}
			equals := ":="
			if t.parent != nil {
				equals = "="
			}

			_, err := buf.Write([]byte(fmt.Sprintf("%s %s %s\n", copyStr, equals, copyVal)))
			return err
		}
	}

	name := getName(root.Type, rootPkg)
	_, err = buf.Write([]byte("func(" + ref + " " + name + ") Copy() " + name + " {\n"))

	if err != nil {
		return nil, nil, err
	}

	if err := generate(root); err != nil {
		return nil, nil, err
	}

	buf.Write([]byte("\nreturn "))
	if root.Kind() == reflect.Ptr {
		buf.Write([]byte{'&'})
	}
	buf.Write([]byte(baseCopy + "\n"))
	buf.Write([]byte{'}', '\n'})

	importsW := bytes.NewBuffer(nil)
	if len(imports) > 0 {
		importsW.Write([]byte("import (\n"))
	}
	for i := range imports {
		alias := getPkgAlias(i)
		importsW.Write([]byte(alias + " " + `"` + i + `"` + "\n"))
	}
	if len(imports) > 0 {
		importsW.Write([]byte{'\n', ')', '\n'})
	}
	return importsW.Bytes(), buf.Bytes(), nil
}

// isKind is a helper function that returns true if the passed in type matches
// any of the passed in kinds.
func isKind(t *reflectType, kinds ...reflect.Kind) bool {
	if t == nil || t.Type == nil {
		return false
	}

	kind := t.Kind()
	for _, k := range kinds {
		if k == kind {
			return true
		}
	}

	return false
}

// addImport adds the package of the passed in type to the list of imports if
// the type is in a different package then the root object.
func addImport(t reflect.Type, rootPkg string, imports map[string]struct{}) {
	if name := getPkgName(t); name != rootPkg {
		if pkgPath := t.PkgPath(); pkgPath != "" {
			imports[pkgPath] = struct{}{}
		}
	}
}

// hasCopyMethod determines if the passed in type has a deep copy method generated
// for it.
// In practice, it checks the type for a function signature like so:
//    func(<T>) Copy() <T>
// It's used to determine if the generator needs to generate it's own copy code
// in-line with the root object, or if it can just rely on the the `Copy()` fn to
// create a deep-copy.
func hasCopyMethod(t reflect.Type) bool {
	m, ok := t.MethodByName("Copy")
	if !ok || m.Type.NumIn() != 1 || m.Type.NumOut() != 1 {
		return false
	}
	has := m.Type.In(0) == t && m.Type.Out(0) == t
	if !has {
		panic(fmt.Sprintf("%v - %v", m.Type.Out(0), t))
	}
	return has
}

// canonicalPkgName gets the alias used to refer to the package that the passed
// in type belongs to.
// If the type belongs to the same package as the root object, then an empty string
// is returned.
func canonicalPkgName(t reflect.Type, rootPkg string) string {
	pkgName := getPkgName(t)
	if pkgName == rootPkg || pkgName == "" {
		pkgName = ""
	}
	return getPkgAlias(pkgName)
}

// getPkgAlias converts an package path to an alias name.
// This makes it easier to refer to types in other packages by providing a
// consistent naming scheme.
func getPkgAlias(s string) string {
	alias := strings.Replace(s, "/", "_", -1)
	alias = strings.Replace(alias, "-", "_", -1)
	alias = strings.Replace(alias, ".", "_", -1)
	return alias
}

// getName is a recursive function that generates the name of the given type
// It traverses maps, slices, and pointers as needed.
func getName(t reflect.Type, rootPkg string) string {
	pkgName := canonicalPkgName(t, rootPkg)
	if n := t.Name(); n != "" {
		if pkgName != "" {
			pkgName += "."
		}
		return pkgName + n
	}
	switch t.Kind() {
	case reflect.Slice:
		return "[]" + getName(t.Elem(), rootPkg)
	case reflect.Map:
		var key, elem string
		key = getName(t.Key(), rootPkg)
		elem = getName(t.Elem(), rootPkg)
		return "map[" + key + "]" + elem
	case reflect.Ptr:
		return "*" + getName(t.Elem(), rootPkg)
	default:
		return t.String()
	}
}

// getPkgName gets the package name that the passed in type belongs to
func getPkgName(t reflect.Type) string {
	for {
		if name := t.PkgPath(); name != "" {
			return name
		}
		if t.Kind() != reflect.Slice && t.Kind() != reflect.Array && t.Kind() != reflect.Ptr && t.Kind() != reflect.Map {
			if t.Name() != "" || !strings.Contains(t.String(), ".") {
				return ""
			}
			panic(fmt.Sprintf("got unexpected type %v", t))
		}
		t = t.Elem()
	}
}
