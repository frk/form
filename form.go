package form

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// Transform takes the values of src and stores them into the value
// pointed to by dst which must be a non-nil pointer to a struct.
func Transform(src url.Values, dst interface{}) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidArgumentError{reflect.TypeOf(dst), "Transform"}
	}

	rv = rv.Elem()
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return &InvalidArgumentError{reflect.TypeOf(dst), "Transform"}
	}

	done := make(map[string]bool)
	return decodeValues(src, rv, done)
}

type InvalidArgumentError struct {
	Type     reflect.Type
	FuncName string
}

func (e *InvalidArgumentError) Error() string {
	if e.Type == nil {
		return "form: " + e.FuncName + "(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "form: " + e.FuncName + "(non-pointer " + e.Type.String() + ")"
	}
	return "form: " + e.FuncName + "(nil " + e.Type.String() + ")"
}

func decodeValues(uv url.Values, sv reflect.Value, done map[string]bool) error {
	st := sv.Type()
	n := sv.NumField()
	embedded := []reflect.Value{}

	for i := 0; i < n; i++ {
		sf := st.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}
		tag := sf.Tag.Get(DefaultTagName)
		if tag == "-" {
			continue
		}
		key, _ := parseTag(tag)
		if key == "" {
			key = sf.Name
		}
		if done[key] {
			continue
		}

		fv := sv.Field(i)
		fk := fv.Kind()
		if fk == reflect.Struct {
			if sf.Anonymous {
				// embedded struct values are handled outside of this loop
				embedded = append(embedded, fv)
			}
			// nested structs are ignored
			continue
		}

		vals := uv[key]
		ln := len(vals)
		if ln == 0 {
			// if no value is associated with the key we continue with the next field
			continue
		}

		if fk == reflect.Slice {
			sl := reflect.MakeSlice(fv.Type(), ln, ln)
			for j := 0; j < ln; j++ {
				str := vals[j]
				if err := decodeString(sl.Index(j), key, str); err != nil {
					return err
				}
			}
			fv.Set(sl)
		} else {
			if err := decodeString(fv, key, vals[0]); err != nil {
				return err
			}
		}

		done[key] = true
	}

	// embedded structs are handled last to ensure decoding into the first available field
	for _, fv := range embedded {
		if err := decodeValues(uv, fv, done); err != nil {
			return err
		}
	}

	return nil
}

func decodeString(rv reflect.Value, key, value string) error {
	if len(value) == 0 {
		return nil
	}

	k := rv.Kind()
	if k == reflect.Ptr && rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}

	switch k {
	case reflect.String:
		rv.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return &ValueError{Key: key, Value: value, Type: k.String()}
		}
		rv.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, int(rv.Type().Size())*8)
		if err != nil {
			return &ValueError{Key: key, Value: value, Type: k.String()}
		}
		rv.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, int(rv.Type().Size())*8)
		if err != nil {
			return &ValueError{Key: key, Value: value, Type: k.String()}
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, int(rv.Type().Size())*8)
		if err != nil {
			return &ValueError{Key: key, Value: value, Type: k.String()}
		}
		rv.SetUint(u)
	case reflect.Ptr:
		return decodeString(rv.Elem(), key, value)
	}
	return nil
}

type ValueError struct {
	Key   string
	Value string
	Type  string
}

func (err *ValueError) Error() string {
	return fmt.Sprintf("form: %q value %q could not be parsed into %q", err.Key, err.Value, err.Type)
}