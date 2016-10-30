package form

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
)

/*
TODO?

- documentation
- decoder tests
- stream implementation of the Decoder a la the json.Decoder
- Marshal/Encoder
- benchmarks
- add support for "multipart/form-data"

*/

type Decoder struct {
	r       io.Reader
	tagName string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, tagName: DefaultTagName}
}

func (d *Decoder) Decode(v interface{}) error {
	rv, ok := structValueOf(v)
	if !ok {
		return &InvalidArgumentError{reflect.TypeOf(v)}
	}

	data, err := ioutil.ReadAll(d.r)
	if err != nil {
		return err
	}
	m, err := parseData(data)
	if err != nil {
		return err
	}

	done := make(map[string]bool)
	return decodeValues(url.Values(m), rv, done, d.tagName)
}

func (d *Decoder) UseTagName(tagName string) *Decoder {
	d.tagName = tagName
	return d
}

// Unmarshal
func Unmarshal(data []byte, v interface{}) error {
	rv, ok := structValueOf(v)
	if !ok {
		return &InvalidArgumentError{reflect.TypeOf(v)}
	}

	m, err := parseData(data)
	if err != nil {
		return err
	}

	done := make(map[string]bool)
	return decodeValues(url.Values(m), rv, done, DefaultTagName)
}

func parseData(data []byte) (map[string][]string, error) {
	m := make(map[string][]string)
	for len(data) != 0 {
		pair := data
		if i := bytes.IndexAny(pair, "&;"); i >= 0 {
			pair, data = pair[:i], pair[i+1:]
		} else {
			data = nil
		}
		if len(pair) == 0 {
			continue
		}

		var key, value []byte
		if i := bytes.IndexByte(pair, '='); i >= 0 {
			key, value = pair[:i], pair[i+1:]
		}

		k, err := url.QueryUnescape(string(key))
		if err != nil {
			return nil, err
		}
		v, err := url.QueryUnescape(string(value))
		if err != nil {
			return nil, err
		}

		m[k] = append(m[k], v)
	}
	return m, nil
}

// Transform
func Transform(vv url.Values, v interface{}) error {
	rv, ok := structValueOf(v)
	if !ok {
		return &InvalidArgumentError{reflect.TypeOf(v)}
	}

	done := make(map[string]bool)
	return decodeValues(vv, rv, done, DefaultTagName)
}

func structValueOf(v interface{}) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return rv, false
	}
	rv = rv.Elem()
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return rv, false
	}
	return rv, true
}

type InvalidArgumentError struct {
	Type reflect.Type
}

func (e *InvalidArgumentError) Error() string {
	s := "form: the v interface{} argument must be a non-nil pointer to a struct, instead got "
	if e.Type == nil {
		s += "nil"
	} else if e.Type.Kind() != reflect.Ptr {
		s += "non-pointer " + e.Type.String()
	} else {
		s += e.Type.String()
	}

	return s
}

func decodeValues(uv url.Values, sv reflect.Value, done map[string]bool, tagName string) error {
	st := sv.Type()
	n := sv.NumField()
	embedded := []reflect.Value{}

	for i := 0; i < n; i++ {
		sf := st.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}
		tag := sf.Tag.Get(tagName)
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

		vals := uv[key]
		ln := len(vals)
		if ln == 0 {
			// if no value is associated with the key we continue with the next field
			if fk == reflect.Struct && sf.Anonymous {
				embedded = append(embedded, fv)
			}
			continue
		}

		var pv reflect.Value
		if fk != reflect.Ptr && fv.Type().Name() != "" && fv.CanAddr() {
			pv = fv.Addr()
		} else if fk == reflect.Ptr {
			pv = fv
		}

		if pv.IsValid() && pv.Type().NumMethod() > 0 {
			if pv.IsNil() {
				pv.Set(reflect.New(pv.Type().Elem()))
			}
			if tu, ok := pv.Interface().(encoding.TextUnmarshaler); ok {
				for _, s := range vals {
					if err := tu.UnmarshalText([]byte(s)); err != nil {
						return err
					}
				}
				continue
			}
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
		if err := decodeValues(uv, fv, done, tagName); err != nil {
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