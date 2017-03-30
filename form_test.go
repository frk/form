package form

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// helpers
func boolp(b bool) *bool      { return &b }
func f32p(f float32) *float32 { return &f }
func f64p(f float64) *float64 { return &f }
func intp(i int) *int         { return &i }
func i16p(i int16) *int16     { return &i }
func i32p(i int32) *int32     { return &i }
func i64p(i int64) *int64     { return &i }
func i8p(i int8) *int8        { return &i }
func strp(s string) *string   { return &s }
func uintp(u uint) *uint      { return &u }
func u16p(u uint16) *uint16   { return &u }
func u32p(u uint32) *uint32   { return &u }
func u64p(u uint64) *uint64   { return &u }
func u8p(u uint8) *uint8      { return &u }

// returns a pointer to a shallow copy of a value pointed to by x.
func pcopy(x interface{}) interface{} {
	if x == nil {
		return nil
	}

	v1 := reflect.ValueOf(x)
	if v1.Kind() != reflect.Ptr {
		return x
	}

	v1 = v1.Elem()
	v2 := reflect.New(v1.Type()).Elem()
	v2.Set(v1)
	return v2.Addr().Interface()
}

type stringParams struct {
	Str   string    `form:"str"`
	Strp  *string   `form:"str_p"`
	Strs  []string  `form:"strs"`
	Strps []*string `form:"str_ps"`
}

type boolParams struct {
	Bool   bool    `form:"bool"`
	Boolp  *bool   `form:"bool_p"`
	Bools  []bool  `form:"bools"`
	Boolps []*bool `form:"bool_ps"`
}

type floatParams struct {
	Float32   float32    `form:"fl32"`
	Float32p  *float32   `form:"fl32_p"`
	Float32s  []float32  `form:"fl32s"`
	Float32ps []*float32 `form:"fl32_ps"`

	Float64   float64    `form:"fl64"`
	Float64p  *float64   `form:"fl64_p"`
	Float64s  []float64  `form:"fl64s"`
	Float64ps []*float64 `form:"fl64_ps"`
}

type intParams struct {
	Int   int    `form:"int"`
	Intp  *int   `form:"int_p"`
	Ints  []int  `form:"ints"`
	Intps []*int `form:"int_ps"`

	Int8   int8    `form:"int8"`
	Int8p  *int8   `form:"int8_p"`
	Int8s  []int8  `form:"int8s"`
	Int8ps []*int8 `form:"int8_ps"`

	Int16   int16    `form:"int16"`
	Int16p  *int16   `form:"int16_p"`
	Int16s  []int16  `form:"int16s"`
	Int16ps []*int16 `form:"int16_ps"`

	Int32   int32    `form:"int32"`
	Int32p  *int32   `form:"int32_p"`
	Int32s  []int32  `form:"int32s"`
	Int32ps []*int32 `form:"int32_ps"`

	Int64   int64    `form:"int64"`
	Int64p  *int64   `form:"int64_p"`
	Int64s  []int64  `form:"int64s"`
	Int64ps []*int64 `form:"int64_ps"`
}

type uintParams struct {
	Uint   uint    `form:"uint"`
	Uintp  *uint   `form:"uint_p"`
	Uints  []uint  `form:"uints"`
	Uintps []*uint `form:"uint_ps"`

	Uint8   uint8    `form:"uint8"`
	Uint8p  *uint8   `form:"uint8_p"`
	Uint8s  []uint8  `form:"uint8s"`
	Uint8ps []*uint8 `form:"uint8_ps"`

	Uint16   uint16    `form:"uint16"`
	Uint16p  *uint16   `form:"uint16_p"`
	Uint16s  []uint16  `form:"uint16s"`
	Uint16ps []*uint16 `form:"uint16_ps"`

	Uint32   uint32    `form:"uint32"`
	Uint32p  *uint32   `form:"uint32_p"`
	Uint32s  []uint32  `form:"uint32s"`
	Uint32ps []*uint32 `form:"uint32_ps"`

	Uint64   uint64    `form:"uint64"`
	Uint64p  *uint64   `form:"uint64_p"`
	Uint64s  []uint64  `form:"uint64s"`
	Uint64ps []*uint64 `form:"uint64_ps"`
}

type embed0 struct {
	embed1
	Name string `form:"-"`
}
type embed1 struct {
	embed2
	Name string `form:"name"`
}
type embed2 struct {
	Name string `form:"name"`
}

type nest0 struct {
	nest1 struct {
		Name string `form:"name"`
	}
}

type MarshSlice []string

func (s *MarshSlice) UnmarshalText(b []byte) error {
	*s = append(*s, strings.ToUpper(string(b)))
	return nil
}

type MarshStruct struct {
	First, Last string
}

func (s *MarshStruct) UnmarshalText(b []byte) error {
	ss := strings.Split(string(b), " ")
	if len(ss) > 0 {
		s.First = ss[0]
	}
	if len(ss) > 1 {
		s.Last = ss[1]
	}
	return nil
}

type marshParams struct {
	Slice   MarshSlice   `form:"slice"`
	Struct  MarshStruct  `form:"struct"`
	Slicep  *MarshSlice  `form:"slice_p"`
	Structp *MarshStruct `form:"struct_p"`
}

type testCase struct {
	name string
	vals url.Values
	dst  interface{}
	want interface{}
	err  error
}

var unmarshalTests = []testCase{
	{
		name: "string params should be set as is",
		vals: url.Values{
			"str": {"hello world"}, "str_p": {"bar"},
			"strs":   {"foo", "bar", "baz"},
			"str_ps": {"baz", "bar", "foo"},
		},
		dst: &stringParams{},
		want: &stringParams{
			Str: "hello world", Strp: strp("bar"),
			Strs:  []string{"foo", "bar", "baz"},
			Strps: []*string{strp("baz"), strp("bar"), strp("foo")},
		},
		err: nil,
	}, {
		name: "bool params should be parsed correctly",
		vals: url.Values{
			"bool": {"TRUE"}, "bool_p": {"FALSE"},
			"bools":   {"TRUE", "True", "true", "T", "t", "1", "FALSE", "False", "false", "F", "f", "0"},
			"bool_ps": {"FALSE", "False", "false", "F", "f", "0", "TRUE", "True", "true", "T", "t", "1"},
		},
		dst: &boolParams{},
		want: &boolParams{
			Bool: true, Boolp: boolp(false),
			Bools:  []bool{true, true, true, true, true, true, false, false, false, false, false, false},
			Boolps: []*bool{boolp(false), boolp(false), boolp(false), boolp(false), boolp(false), boolp(false), boolp(true), boolp(true), boolp(true), boolp(true), boolp(true), boolp(true)},
		},
		err: nil,
	}, {
		name: "float params should be parsed correctly",
		vals: url.Values{
			"fl32": {"42.0987654321"}, "fl32_p": {"3.14159265359"},
			"fl32s":   {"12", "-3.14", "1.61803"},
			"fl32_ps": {"-3.14", "1.61803", "12"},

			"fl64": {"42.0987654321"}, "fl64_p": {"3.14159265359"},
			"fl64s":   {"12", "-3.14", "1.61803"},
			"fl64_ps": {"-3.14", "1.61803", "12"},
		},
		dst: &floatParams{},
		want: &floatParams{
			Float32: 42.098766, Float32p: f32p(3.1415927),
			Float32s:  []float32{12, -3.14, 1.61803},
			Float32ps: []*float32{f32p(-3.14), f32p(1.61803), f32p(12.0)},

			Float64: 42.0987654321, Float64p: f64p(3.14159265359),
			Float64s:  []float64{12, -3.14, 1.61803},
			Float64ps: []*float64{f64p(-3.14), f64p(1.61803), f64p(12.0)},
		},
		err: nil,
	}, {
		name: "int params should be parsed correctly",
		vals: url.Values{
			"int": {"2147483647"}, "int_p": {"-2147483648"},
			"ints":   {"0", "-123", "999"},
			"int_ps": {"999", "0", "-123"},

			"int8": {"127"}, "int8_p": {"-128"},
			"int8s":   {"0", "-123", "99"},
			"int8_ps": {"99", "0", "-123"},

			"int16": {"32767"}, "int16_p": {"-32768"},
			"int16s":   {"0", "-1968", "9999"},
			"int16_ps": {"9999", "0", "-1968"},

			"int32": {"2147483647"}, "int32_p": {"-2147483648"},
			"int32s":   {"0", "-62976", "999999"},
			"int32_ps": {"999999", "0", "-62976"},

			"int64": {"9223372036854775807"}, "int64_p": {"-9223372036854775808"},
			"int64s":   {"0", "-4030464", "999999999"},
			"int64_ps": {"999999999", "0", "-4030464"},
		},
		dst: &intParams{},
		want: &intParams{
			Int: (1 << 31) - 1, Intp: intp(-1 << 31),
			Ints:  []int{0, -123, 999},
			Intps: []*int{intp(999), intp(0), intp(-123)},

			Int8: (1 << 7) - 1, Int8p: i8p(-1 << 7),
			Int8s:  []int8{0, -123, 99},
			Int8ps: []*int8{i8p(99), i8p(0), i8p(-123)},

			Int16: (1 << 15) - 1, Int16p: i16p(-1 << 15),
			Int16s:  []int16{0, -1968, 9999},
			Int16ps: []*int16{i16p(9999), i16p(0), i16p(-1968)},

			Int32: (1 << 31) - 1, Int32p: i32p(-1 << 31),
			Int32s:  []int32{0, -62976, 999999},
			Int32ps: []*int32{i32p(999999), i32p(0), i32p(-62976)},

			Int64: (1 << 63) - 1, Int64p: i64p(-1 << 63),
			Int64s:  []int64{0, -4030464, 999999999},
			Int64ps: []*int64{i64p(999999999), i64p(0), i64p(-4030464)},
		},
		err: nil,
	}, {
		name: "uint params should be parsed correctly",
		vals: url.Values{
			"uint": {"4294967295"}, "uint_p": {"4294967295"},
			"uints":   {"0", "123", "999"},
			"uint_ps": {"999", "0", "123"},

			"uint8": {"255"}, "uint8_p": {"255"},
			"uint8s":   {"0", "123", "99"},
			"uint8_ps": {"99", "0", "123"},

			"uint16": {"65535"}, "uint16_p": {"65535"},
			"uint16s":   {"0", "1968", "9999"},
			"uint16_ps": {"9999", "0", "1968"},

			"uint32": {"4294967295"}, "uint32_p": {"4294967295"},
			"uint32s":   {"0", "62976", "999999"},
			"uint32_ps": {"999999", "0", "62976"},

			"uint64": {"18446744073709551615"}, "uint64_p": {"18446744073709551615"},
			"uint64s":   {"0", "4030464", "999999999"},
			"uint64_ps": {"999999999", "0", "4030464"},
		},
		dst: &uintParams{},
		want: &uintParams{
			Uint: (1 << 32) - 1, Uintp: uintp((1 << 32) - 1),
			Uints:  []uint{0, 123, 999},
			Uintps: []*uint{uintp(999), uintp(0), uintp(123)},

			Uint8: (1 << 8) - 1, Uint8p: u8p((1 << 8) - 1),
			Uint8s:  []uint8{0, 123, 99},
			Uint8ps: []*uint8{u8p(99), u8p(0), u8p(123)},

			Uint16: (1 << 16) - 1, Uint16p: u16p((1 << 16) - 1),
			Uint16s:  []uint16{0, 1968, 9999},
			Uint16ps: []*uint16{u16p(9999), u16p(0), u16p(1968)},

			Uint32: (1 << 32) - 1, Uint32p: u32p((1 << 32) - 1),
			Uint32s:  []uint32{0, 62976, 999999},
			Uint32ps: []*uint32{u32p(999999), u32p(0), u32p(62976)},

			Uint64: (1 << 64) - 1, Uint64p: u64p((1 << 64) - 1),
			Uint64s:  []uint64{0, 4030464, 999999999},
			Uint64ps: []*uint64{u64p(999999999), u64p(0), u64p(4030464)},
		},
		err: nil,
	}, {
		name: "unmarshal text",
		vals: url.Values{
			"slice": {"John Doe", "Jane Doe"}, "struct": {"John Doe"},
			"slice_p": {"hello world", "hello 世界"}, "struct_p": {"Jane Doe"},
		},
		dst: &marshParams{},
		want: &marshParams{
			Slice:   MarshSlice{"JOHN DOE", "JANE DOE"},
			Struct:  MarshStruct{"John", "Doe"},
			Slicep:  &MarshSlice{"HELLO WORLD", "HELLO 世界"},
			Structp: &MarshStruct{"Jane", "Doe"},
		},
		err: nil,
	}, {
		name: "should transform to an embeded field but only once",
		vals: url.Values{"name": {"John Doe"}},
		dst:  &embed0{}, want: &embed0{embed1: embed1{Name: "John Doe"}},
		err: nil,
	}, {
		name: "should not transform to a nested struct",
		vals: url.Values{"name": {"John Doe"}},
		dst:  &nest0{}, want: &nest0{},
		err: nil,
	}, {
		name: "nil dst argument should err",
		vals: url.Values{},
		dst:  nil, want: nil,
		err: &ArgumentError{nil},
	}, {
		name: "nil dst should return error",
		dst:  nil, want: nil,
		err: &ArgumentError{nil},
	}, {
		name: "non-pointer-struct dst should return error",
		dst:  struct{}{}, want: struct{}{},
		err: &ArgumentError{reflect.TypeOf(struct{}{})},
	}, {
		name: "any dst that is not a pointer to a struct should return error",
		dst:  []int{}, want: []int{},
		err: &ArgumentError{reflect.TypeOf([]int{})},
	}, {
		name: "bool type err",
		vals: url.Values{"bool": {"Hello World"}},
		dst:  &boolParams{}, want: &boolParams{},
		err: &ValueError{Key: "bool", Value: "Hello World", Type: "bool"},
	}, {
		name: "int type err",
		vals: url.Values{"int": {"Hello World"}},
		dst:  &intParams{}, want: &intParams{},
		err: &ValueError{Key: "int", Value: "Hello World", Type: "int"},
	}, {
		name: "int8 value out of range",
		vals: url.Values{"int8": {"128"}},
		dst:  &intParams{}, want: &intParams{},
		err: &ValueError{Key: "int8", Value: "128", Type: "int8"},
	}, {
		name: "int16 value out of range",
		vals: url.Values{"int16": {"32768"}},
		dst:  &intParams{}, want: &intParams{},
		err: &ValueError{Key: "int16", Value: "32768", Type: "int16"},
	}, {
		name: "int32 value out of range",
		vals: url.Values{"int32": {"2147483648"}},
		dst:  &intParams{}, want: &intParams{},
		err: &ValueError{Key: "int32", Value: "2147483648", Type: "int32"},
	}, {
		name: "int64 value out of range",
		vals: url.Values{"int64": {"9223372036854775808"}},
		dst:  &intParams{}, want: &intParams{},
		err: &ValueError{Key: "int64", Value: "9223372036854775808", Type: "int64"},
	}, {
		name: "uint type err",
		vals: url.Values{"uint": {"Hello World"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint", Value: "Hello World", Type: "uint"},
	}, {
		name: "uint value out of range",
		vals: url.Values{"uint": {"-128"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint", Value: "-128", Type: "uint"},
	}, {
		name: "uint8 value out of range",
		vals: url.Values{"uint8": {"256"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint8", Value: "256", Type: "uint8"},
	}, {
		name: "uint16 value out of range",
		vals: url.Values{"uint16": {"65536"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint16", Value: "65536", Type: "uint16"},
	}, {
		name: "uint32 value out of range",
		vals: url.Values{"uint32": {"4294967296"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint32", Value: "4294967296", Type: "uint32"},
	}, {
		name: "uint64 value out of range",
		vals: url.Values{"uint64": {"18446744073709551616"}},
		dst:  &uintParams{}, want: &uintParams{},
		err: &ValueError{Key: "uint64", Value: "18446744073709551616", Type: "uint64"},
	},
}

func TestTransform(t *testing.T) {
	for i, tt := range unmarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			dst := pcopy(tt.dst)
			if err := Transform(tt.vals, dst); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: got %#v, want %#v", i, err, tt.err)
			}
			if !reflect.DeepEqual(tt.want, dst) {
				t.Errorf("#%d: got %+v, want %+v", i, dst, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range unmarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(tt.vals.Encode())
			dst := pcopy(tt.dst)
			if err := Unmarshal(data, dst); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: got %#v, want %#v", i, err, tt.err)
			}
			if !reflect.DeepEqual(tt.want, dst) {
				t.Errorf("#%d: got %+v, want %+v", i, dst, tt.want)
			}
		})
	}
}

func TestParseBytes(t *testing.T) {
	tests := []struct {
		in   string
		want map[string][]string
		err  error
	}{
		{
			in:   "a=1&b=2",
			want: map[string][]string{"a": []string{"1"}, "b": []string{"2"}},
		}, {
			in:   "a=1&a=2&a=banana",
			want: map[string][]string{"a": []string{"1", "2", "banana"}},
		}, {
			in:   "ascii=%3Ckey%3A+0x90%3E",
			want: map[string][]string{"ascii": []string{"<key: 0x90>"}},
		}, {
			in:   "a=1;b=2",
			want: map[string][]string{"a": []string{"1"}, "b": []string{"2"}},
		}, {
			in:   "a=1&a=2;a=banana",
			want: map[string][]string{"a": []string{"1", "2", "banana"}},
		}, {
			in:   "a=100%",
			want: nil, err: url.EscapeError("%"),
		},
	}

	for i, tt := range tests {
		got, err := parseBytes([]byte(tt.in))
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("#%d: got err %v, want %v", i, err, tt.err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("#%d: got %v, want %v", i, got, tt.want)
		}
	}
}

type boolType struct {
	Bool   bool
	Boolp  *bool
	Bools  []bool
	Boolps []*bool
}

var boolVal = boolType{
	Bool:   false,
	Boolp:  boolp(true),
	Bools:  []bool{true, false, false},
	Boolps: []*bool{boolp(true), boolp(true)},
}

var boolValString = `Bool=false&Boolp=true&Bools=true&Bools=false&Bools=false&Boolps=true&Boolps=true`

func TestUnmarshal2(t *testing.T) {
	var val boolType
	if err := Unmarshal([]byte(boolValString), &val); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(boolVal, val) {
		t.Errorf("got %v want %v", val, boolVal)
	}
}

type intType struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64
}

var intVal = intType{
	Int:   1,
	Int8:  2,
	Int16: 3,
	Int32: 4,
	Int64: 5,
}

var intValString = `Int=1&Int8=2&Int16=3&Int32=4&Int64=5`

type intpType struct {
	Intp   *int
	Int8p  *int8
	Int16p *int16
	Int32p *int32
	Int64p *int64
}

var intpVal = intpType{
	Intp:   intp(6),
	Int8p:  i8p(7),
	Int16p: i16p(8),
	Int32p: i32p(9),
	Int64p: i64p(10),
}

var intpValString = `Intp=6&Int8p=7&Int16p=8&Int32p=9&Int64p=10`

type intsType struct {
	Ints   []int
	Int8s  []int8
	Int16s []int16
	Int32s []int32
	Int64s []int64
}

var intsVal = intsType{
	Ints:   []int{11},
	Int8s:  []int8{12, 13},
	Int16s: []int16{14},
	Int32s: []int32{15},
	Int64s: []int64{16, 17},
}

var intsValString = `Ints=11&Int8s=12&Int8s=13&Int16s=14&Int32s=15&Int64s=16&Int64s=17`

type intpsType struct {
	Intps   []*int
	Int8ps  []*int8
	Int16ps []*int16
	Int32ps []*int32
	Int64ps []*int64
}

var intpsVal = intpsType{
	Intps:   []*int{intp(18), intp(19)},
	Int8ps:  []*int8{i8p(20)},
	Int16ps: []*int16{i16p(21), i16p(22)},
	Int32ps: []*int32{i32p(23), i32p(24)},
	Int64ps: []*int64{i64p(25)},
}

var intpsValString = `Intps=18&Intps=19&Int8ps=20&Int16ps=21&Int16ps=22&Int32ps=23&Int32ps=24&Int64ps=25`

type uintType struct {
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
}

var uintVal = uintType{
	Uint:   26,
	Uint8:  27,
	Uint16: 28,
	Uint32: 29,
	Uint64: 30,
}

var uintValString = `Uint=26&Uint8=27&Uint16=28&Uint32=29&Uint64=30`

type uintpType struct {
	Uintp   *uint
	Uint8p  *uint8
	Uint16p *uint16
	Uint32p *uint32
	Uint64p *uint64
}

var uintpVal = uintpType{
	Uintp:   uintp(31),
	Uint8p:  u8p(32),
	Uint16p: u16p(33),
	Uint32p: u32p(34),
	Uint64p: u64p(35),
}

var uintpValString = `Uintp=31&Uint8p=32&Uint16p=33&Uint32p=34&Uint64p=35`

type uintsType struct {
	Uints   []uint
	Uint8s  []uint8
	Uint16s []uint16
	Uint32s []uint32
	Uint64s []uint64
}

var uintsVal = uintsType{
	Uints:   []uint{36},
	Uint8s:  []uint8{37, 38},
	Uint16s: []uint16{39},
	Uint32s: []uint32{40},
	Uint64s: []uint64{41, 42},
}

var uintsValString = `Uints=36&Uint8s=37&Uint8s=38&Uint16s=39&Uint32s=40&Uint64s=41&Uint64s=42`

type uintpsType struct {
	Uintps   []*uint
	Uint8ps  []*uint8
	Uint16ps []*uint16
	Uint32ps []*uint32
	Uint64ps []*uint64
}

var uintpsVal = uintpsType{
	Uintps:   []*uint{uintp(43), uintp(44)},
	Uint8ps:  []*uint8{u8p(45)},
	Uint16ps: []*uint16{u16p(46), u16p(47)},
	Uint32ps: []*uint32{u32p(48), u32p(49)},
	Uint64ps: []*uint64{u64p(50)},
}

var uintpsValString = `Uintps=43&Uintps=44&Uint8ps=45&Uint16ps=46&Uint16ps=47&Uint32ps=48&Uint32ps=49&Uint64ps=50`

type stringType struct {
	String   string
	Stringp  *string
	Strings  []string
	Stringps []*string
}

var stringVal = stringType{
	String:   "51",
	Stringp:  strp("foo"),
	Strings:  []string{"foo", "bar", "baz"},
	Stringps: []*string{strp("baz"), strp("bar"), strp("foo")},
}

var stringValString = `String=51&Stringp=foo&Strings=foo&Strings=bar&Strings=baz&Stringps=baz&Stringps=bar&Stringps=foo`

type floatType struct {
	Float32  float32
	Float32p *float32
	Float64  float64
	Float64p *float64
}

var floatVal = floatType{
	Float32:  52.00001,
	Float32p: f32p(52.1234),
	Float64:  52.64,
	Float64p: f64p(52.0),
}

var floatValString = `Float32=52.00001&Float32p=52.1234&Float64=52.64&Float64p=52`

type floatsType struct {
	Float32s  []float32
	Float32ps []*float32
	Float64s  []float64
	Float64ps []*float64
}

var floatsVal = floatsType{
	Float32s:  []float32{53.01},
	Float32ps: []*float32{f32p(53.1), f32p(53.2)},
	Float64s:  []float64{53.03, 53.04},
	Float64ps: []*float64{f64p(53.0005), f64p(53.6)},
}

var floatsValString = `Float32s=53.01&Float32ps=53.1&Float32ps=53.2&Float64s=53.03&Float64s=53.04&Float64ps=53.0005&Float64ps=53.6`

type _embed0 struct {
	_embed1
	Field string
}

type _embed1 struct {
	_embed2
	Field int
}

type _embed2 struct {
	Field float64
}

var embedVal = _embed0{
	_embed1: _embed1{
		_embed2: _embed2{34.67},
		Field:   3467,
	},
	Field: "string",
}

var embedValString = `Field=34.67&Field=3467&Field=string`

type marshalSlice []string

func (s *marshalSlice) MarshalText() ([]byte, error) {
	return []byte(strings.Join(*s, ",")), nil
}

type marshalType struct {
	M *marshalSlice
	N marshalSlice
}

var marshalVal = marshalType{
	M: &marshalSlice{"foo", "bar", "baz"},
	N: marshalSlice{"foo", "bar", "baz"},
}

var marshalValString = `M=foo%2Cbar%2Cbaz&N=foo&N=bar&N=baz`

func ifacep(i interface{}) *interface{} { return &i }

type ifaceType struct {
	IString interface{}
	IInt    interface{}
	IBool   *interface{}
	ISlice  interface{}
}

var ifaceVal = ifaceType{
	IString: "foo",
	IInt:    intp(32),
	IBool:   ifacep(true),
	ISlice:  []interface{}{"abc", float32(32.1234567)},
}

var ifaceValString = `IString=foo&IInt=32&IBool=true&ISlice=abc&ISlice=32.123455`

// TODO test tags

func TestMarshal(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
		str  string
		dst  interface{}
	}{{
		name: "bool values",
		val:  boolVal,
		str:  boolValString,
		dst:  boolType{},
	}, {
		name: "int values",
		val:  intVal,
		str:  intValString,
		dst:  intType{},
	}, {
		name: "int pointer values",
		val:  intpVal,
		str:  intpValString,
		dst:  intpType{},
	}, {
		name: "int slices",
		val:  intsVal,
		str:  intsValString,
		dst:  intsType{},
	}, {
		name: "int pointer slices",
		val:  intpsVal,
		str:  intpsValString,
		dst:  intpsType{},
	}, {
		name: "uint values",
		val:  uintVal,
		str:  uintValString,
		dst:  uintType{},
	}, {
		name: "uint pointer values",
		val:  uintpVal,
		str:  uintpValString,
		dst:  uintpType{},
	}, {
		name: "uint slices",
		val:  uintsVal,
		str:  uintsValString,
		dst:  uintsType{},
	}, {
		name: "uint pointer slices",
		val:  uintpsVal,
		str:  uintpsValString,
		dst:  uintpsType{},
	}, {
		name: "string values",
		val:  stringVal,
		str:  stringValString,
		dst:  stringType{},
	}, {
		name: "float values",
		val:  floatVal,
		str:  floatValString,
		dst:  floatType{},
	}, {
		name: "float slices",
		val:  floatsVal,
		str:  floatsValString,
		dst:  floatsType{},
	}, {
		name: "embeded types",
		val:  embedVal,
		str:  embedValString,
		dst:  _embed0{},
	}, {
		name: "TextMarshaler type",
		val:  marshalVal,
		str:  marshalValString,
		dst:  marshalType{},
	}, {
		name: "interface values",
		val:  ifaceVal,
		str:  ifaceValString,
		dst:  ifaceType{},
	}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bgot, err := Marshal(tt.val)
			if err != nil {
				t.Error(err)
				return
			}

			if got := string(bgot); got != tt.str {
				t.Errorf("#%d: got %q, want %q", i, got, tt.str)
			}
		})
	}
}
