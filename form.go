// The form package decodes "application/x-www-form-urlencoded" type data into Go structs.
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
TODO:
- encoder
- support for "multipart/form-data"
*/

// The ArgumentError will be returned by one of the package's expored functions
// or methods if the argument passed to them is not a non-nil pointer to a struct.
type ArgumentError struct {
	Type reflect.Type
}

func (e *ArgumentError) Error() string {
	var t string
	if e.Type == nil {
		t = "nil"
	} else if e.Type.Kind() != reflect.Ptr {
		t = "non-pointer " + e.Type.String()
	} else {
		t = e.Type.String()
	}
	return "form: the v interface{} argument must be a non-nil pointer to a struct, instead got " + t
}

// A ValueError describes a URL-encoded value that was not
// appropriate for a value of a specific Go type.
type ValueError struct {
	Key   string
	Value string
	Type  string
}

func (err *ValueError) Error() string {
	return fmt.Sprintf("form: %q value %q could not be parsed into %q", err.Key, err.Value, err.Type)
}

// Unmarshal parses the URL-encoded data and stores the result in the value
// pointed to by v. The v argument must point to a struct value otherwhise
// an ArgumentError will be returned.
func Unmarshal(data []byte, v interface{}) error {
	src, err := parseBytes(data)
	if err != nil {
		return err
	}
	d := &Decoder{
		src:  src,
		done: make(map[string]bool),
	}
	return d.Decode(v)
}

// Transform takes the url.Values src and stores its elements into the value
// pointed to by dst. The dst argument must point to a struct value otherwhise
// an ArgumentError will be returned.
func Transform(src url.Values, dst interface{}) error {
	d := &Decoder{
		src:  map[string][]string(src),
		done: make(map[string]bool),
	}
	return d.Decode(dst)
}

// A Decoder reads and decodes URL-encoded values.
type Decoder struct {
	src    map[string][]string
	done   map[string]bool
	tagKey string
	err    error

	vals []string
	key  string
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return &Decoder{err: err}
	}

	src, err := parseBytes(data)
	if err != nil {
		return &Decoder{err: err}
	}

	return &Decoder{src: src, done: make(map[string]bool)}
}

// Decode reads the URL-encoded data from its input and stores it in the
// value pointed to by v. The v argument must point to a struct value
// otherwhise an ArgumentError will be returned.
func (d *Decoder) Decode(v interface{}) error {
	if d.err != nil {
		return d.err
	}
	if d.tagKey == "" {
		d.tagKey = DefaultTagKey
	}

	rv, ok := structValueOf(v)
	if !ok {
		return &ArgumentError{reflect.TypeOf(v)}
	}
	return d.decode(rv)
}

// The decode method decodes the Decoder's src values into the dst struct value.
func (d *Decoder) decode(dst reflect.Value) error {
	var (
		n        = dst.NumField()
		stype    = dst.Type()
		embedded = []reflect.Value{}
	)

	for i := 0; i < n; i++ {
		field := stype.Field(i)

		// Skip fields that are unexported and are not embedded.
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		tag := field.Tag.Get(d.tagKey)
		if tag == "-" {
			continue
		}
		key, _ := parseTag(tag)
		if key == "" {
			key = field.Name
		}

		// If a field with this key was already decoded,
		// continue to the next one.
		if d.done[key] {
			continue
		}

		d.vals = d.src[key]
		d.key = key

		fv := dst.Field(i)
		fk := fv.Kind()
		ln := len(d.vals)

		if ln == 0 {
			// If the field is a struct and it is embedded, "record"
			// it and decode its fields after the main loop's done.
			if fk == reflect.Struct && field.Anonymous {
				embedded = append(embedded, fv)
			}

			// If no value is associated with the key
			// continue to the next field.
			continue
		}

		// If the value implements encoding.TextUnmarshaler, loop over
		// the values and call its UnmarshalText method with each value.
		if fk != reflect.Ptr && fv.CanAddr() && fv.Type().Name() != "" {
			fv = fv.Addr()
		}
		if fv.IsValid() && fv.Type().NumMethod() > 0 {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			if tu, ok := fv.Interface().(encoding.TextUnmarshaler); ok {
				for _, s := range d.vals {
					if err := tu.UnmarshalText([]byte(s)); err != nil {
						return err
					}
				}
				d.done[d.key] = true
				continue
			}
		}

		// If the field is a slice, allocate a new slice with length
		// equal to the number of elements in values, loop over the
		// values and decode each one into its respective position.
		if fk == reflect.Slice {
			sl := reflect.MakeSlice(fv.Type(), ln, ln)
			for j := 0; j < ln; j++ {
				if err := decodeString(sl.Index(j), d.vals[j]); err != nil {
					return &ValueError{Key: key, Value: d.vals[j], Type: fk.String()}
				}
			}
			fv.Set(sl)
			d.done[key] = true
			continue
		}

		if err := decodeString(fv, d.vals[0]); err != nil {
			return &ValueError{Key: key, Value: d.vals[0], Type: fk.String()}
		}
		d.done[key] = true
	}

	// Loop over all of the embedded struct values, if there were any, and decode them.
	for _, v := range embedded {
		if err := d.decode(v); err != nil {
			return err
		}
	}

	return nil
}

// decodeString decodes the string src into the reflect.Value dst. If src
// cannot be decoded into the dst value, decodeString will return an error.
// If dst is not one of the supported kinds it will be ignored.
func decodeString(dst reflect.Value, src string) error {
	if len(src) == 0 {
		return nil
	}

	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()

	}

	switch k := dst.Kind(); k {
	case reflect.String:
		dst.SetString(src)
	case reflect.Bool:
		b, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		dst.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(src, int(dst.Type().Size())*8)
		if err != nil {
			return err
		}
		dst.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(src, 10, int(dst.Type().Size())*8)
		if err != nil {
			return err
		}
		dst.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(src, 10, int(dst.Type().Size())*8)
		if err != nil {
			return err
		}
		dst.SetUint(u)
	}
	return nil
}

// structValueOf returns a new reflect.Value initialized to the concrete
// struct value stored in the interface v. The ok return value reports
// whether the value stored in v is a non-nil pointer to a struct or not.
func structValueOf(v interface{}) (rv reflect.Value, ok bool) {
	rv = reflect.ValueOf(v)
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

// parseBytes parses the URL-encoded data and returns
// a map listing the values specified for each key.
func parseBytes(data []byte) (map[string][]string, error) {
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