package form

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

//
// multipart/form-data; boundary="foobar"
//
// TODO test tags

func boolp(b bool) *bool { return &b }

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

var boolValues = url.Values{"Bool": {"false"}, "Boolp": {"true"}, "Bools": {"true", "false", "false"}, "Boolps": {"true", "true"}}

const boolValString = `Bool=false&Boolp=true&Bools=true&Bools=false&Bools=false&Boolps=true&Boolps=true`
const boolValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Bool"` + "\n\n" + `false
--foobar` + "\n" + `Content-Disposition: form-data; name="Boolp"` + "\n\n" + `true
--foobar` + "\n" + `Content-Disposition: form-data; name="Bools"` + "\n\n" + `true
--foobar` + "\n" + `Content-Disposition: form-data; name="Bools"` + "\n\n" + `false
--foobar` + "\n" + `Content-Disposition: form-data; name="Boolps"` + "\n\n" + `true
--foobar` + "\n" + `Content-Disposition: form-data; name="Boolps"` + "\n\n" + `true
--foobar--
`

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

var intValues = url.Values{"Int": {"1"}, "Int8": {"2"}, "Int16": {"3"}, "Int32": {"4"}, "Int64": {"5"}}

const intValString = `Int=1&Int8=2&Int16=3&Int32=4&Int64=5`
const intValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Int"` + "\n\n" + `1
--foobar` + "\n" + `Content-Disposition: form-data; name="Int8"` + "\n\n" + `2
--foobar` + "\n" + `Content-Disposition: form-data; name="Int16"` + "\n\n" + `3
--foobar` + "\n" + `Content-Disposition: form-data; name="Int32"` + "\n\n" + `4
--foobar` + "\n" + `Content-Disposition: form-data; name="Int64"` + "\n\n" + `5
--foobar--
`

func i16p(i int16) *int16 { return &i }
func i32p(i int32) *int32 { return &i }
func i64p(i int64) *int64 { return &i }
func i8p(i int8) *int8    { return &i }
func intp(i int) *int     { return &i }

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

var intpValues = url.Values{"Intp": {"6"}, "Int8p": {"7"}, "Int16p": {"8"}, "Int32p": {"9"}, "Int64p": {"10"}}

const intpValString = `Intp=6&Int8p=7&Int16p=8&Int32p=9&Int64p=10`
const intpValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Intp"` + "\n\n" + `6
--foobar` + "\n" + `Content-Disposition: form-data; name="Intp8"` + "\n\n" + `7
--foobar` + "\n" + `Content-Disposition: form-data; name="Intp16"` + "\n\n" + `8
--foobar` + "\n" + `Content-Disposition: form-data; name="Intp32"` + "\n\n" + `9
--foobar` + "\n" + `Content-Disposition: form-data; name="Intp64"` + "\n\n" + `10
--foobar--
`

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

var intsValues = url.Values{"Ints": {"11"}, "Int8s": {"12", "13"}, "Int16s": {"14"}, "Int32s": {"15"}, "Int64s": {"16", "17"}}

const intsValString = `Ints=11&Int8s=12&Int8s=13&Int16s=14&Int32s=15&Int64s=16&Int64s=17`
const intsValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints"` + "\n\n" + `11
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints8"` + "\n\n" + `12
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints8"` + "\n\n" + `13
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints16"` + "\n\n" + `14
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints32"` + "\n\n" + `15
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints64"` + "\n\n" + `16
--foobar` + "\n" + `Content-Disposition: form-data; name="Ints64"` + "\n\n" + `17
--foobar--
`

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

var intpsValues = url.Values{"Intps": {"18", "19"}, "Int8ps": {"20"}, "Int16ps": {"21", "22"}, "Int32ps": {"23", "24"}, "Int64ps": {"25"}}

const intpsValString = `Intps=18&Intps=19&Int8ps=20&Int16ps=21&Int16ps=22&Int32ps=23&Int32ps=24&Int64ps=25`
const intpsValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps"` + "\n\n" + `18
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps"` + "\n\n" + `19
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps8"` + "\n\n" + `20
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps16"` + "\n\n" + `21
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps16"` + "\n\n" + `22
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps32"` + "\n\n" + `23
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps32"` + "\n\n" + `24
--foobar` + "\n" + `Content-Disposition: form-data; name="Intps64"` + "\n\n" + `25
--foobar--
`

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

var uintValues = url.Values{"Uint": {"26"}, "Uint8": {"27"}, "Uint16": {"28"}, "Uint32": {"29"}, "Uint64": {"30"}}

const uintValString = `Uint=26&Uint8=27&Uint16=28&Uint32=29&Uint64=30`
const uintValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Uint"` + "\n\n" + `26
--foobar` + "\n" + `Content-Disposition: form-data; name="Uint8"` + "\n\n" + `27
--foobar` + "\n" + `Content-Disposition: form-data; name="Uint16"` + "\n\n" + `28
--foobar` + "\n" + `Content-Disposition: form-data; name="Uint32"` + "\n\n" + `29
--foobar` + "\n" + `Content-Disposition: form-data; name="Uint64"` + "\n\n" + `30
--foobar--
`

func uintp(u uint) *uint    { return &u }
func u16p(u uint16) *uint16 { return &u }
func u32p(u uint32) *uint32 { return &u }
func u64p(u uint64) *uint64 { return &u }
func u8p(u uint8) *uint8    { return &u }

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

var uintpValues = url.Values{"Uintp": {"31"}, "Uint8p": {"32"}, "Uint16p": {"33"}, "Uint32p": {"34"}, "Uint64p": {"35"}}

const uintpValString = `Uintp=31&Uint8p=32&Uint16p=33&Uint32p=34&Uint64p=35`
const uintpValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Uintp"` + "\n\n" + `31
--foobar` + "\n" + `Content-Disposition: form-data; name="Uintp8"` + "\n\n" + `32
--foobar` + "\n" + `Content-Disposition: form-data; name="Uintp16"` + "\n\n" + `33
--foobar` + "\n" + `Content-Disposition: form-data; name="Uintp32"` + "\n\n" + `34
--foobar` + "\n" + `Content-Disposition: form-data; name="Uintp64"` + "\n\n" + `35
--foobar--
`

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

var uintsValues = url.Values{"Uints": {"36"}, "Uint8s": {"37", "38"}, "Uint16s": {"39"}, "Uint32s": {"40"}, "Uint64s": {"41", "42"}}

const uintsValString = `Uints=36&Uint8s=37&Uint8s=38&Uint16s=39&Uint32s=40&Uint64s=41&Uint64s=42`
const uintsValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints"` + "\n\n" + `36
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints8"` + "\n\n" + `37
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints8"` + "\n\n" + `38
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints16"` + "\n\n" + `39
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints32"` + "\n\n" + `40
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints64"` + "\n\n" + `41
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints64"` + "\n\n" + `42
--foobar--
`

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

var uintpsValues = url.Values{"Uintps": {"43", "44"}, "Uint8ps": {"45"}, "Uint16ps": {"46", "47"}, "Uint32ps": {"48", "49"}, "Uint64ps": {"50"}}

const uintpsValString = `Uintps=43&Uintps=44&Uint8ps=45&Uint16ps=46&Uint16ps=47&Uint32ps=48&Uint32ps=49&Uint64ps=50`
const uintpsValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints"` + "\n\n" + `43
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints"` + "\n\n" + `44
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints8"` + "\n\n" + `45
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints16"` + "\n\n" + `46
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints16"` + "\n\n" + `47
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints32"` + "\n\n" + `48
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints32"` + "\n\n" + `49
--foobar` + "\n" + `Content-Disposition: form-data; name="Uints64"` + "\n\n" + `50
--foobar--
`

func strp(s string) *string { return &s }

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

var stringValues = url.Values{"String": {"51"}, "Stringp": {"foo"}, "Strings": {"foo", "bar", "baz"}, "Stringps": {"baz", "bar", "foo"}}

const stringValString = `String=51&Stringp=foo&Strings=foo&Strings=bar&Strings=baz&Stringps=baz&Stringps=bar&Stringps=foo`
const stringValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="String"` + "\n\n" + `51
--foobar` + "\n" + `Content-Disposition: form-data; name="Stringp"` + "\n\n" + `foo
--foobar` + "\n" + `Content-Disposition: form-data; name="Strings"` + "\n\n" + `foo
--foobar` + "\n" + `Content-Disposition: form-data; name="Strings"` + "\n\n" + `bar
--foobar` + "\n" + `Content-Disposition: form-data; name="Strings"` + "\n\n" + `baz
--foobar` + "\n" + `Content-Disposition: form-data; name="Stringps"` + "\n\n" + `baz
--foobar` + "\n" + `Content-Disposition: form-data; name="Stringps"` + "\n\n" + `bar
--foobar` + "\n" + `Content-Disposition: form-data; name="Stringps"` + "\n\n" + `foo
--foobar--
`

func f32p(f float32) *float32 { return &f }
func f64p(f float64) *float64 { return &f }

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

var floatValues = url.Values{"Float32": {"52.00001"}, "Float32p": {"52.1234"}, "Float64": {"52.64"}, "Float64p": {"52"}}

const floatValString = `Float32=52.00001&Float32p=52.1234&Float64=52.64&Float64p=52`
const floatValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Float32"` + "\n\n" + `52.00001
--foobar` + "\n" + `Content-Disposition: form-data; name="Float32p"` + "\n\n" + `52.1234
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64"` + "\n\n" + `52.64
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64p"` + "\n\n" + `52
--foobar--
`

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

var floatsValues = url.Values{"Float32s": {"53.01"}, "Float32ps": {"53.1", "53.2"}, "Float64s": {"53.03", "53.04"}, "Float64ps": {"53.0005", "53.6"}}

const floatsValString = `Float32s=53.01&Float32ps=53.1&Float32ps=53.2&Float64s=53.03&Float64s=53.04&Float64ps=53.0005&Float64ps=53.6`
const floatsValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Float32s"` + "\n\n" + `53.01
--foobar` + "\n" + `Content-Disposition: form-data; name="Float32ps"` + "\n\n" + `53.1
--foobar` + "\n" + `Content-Disposition: form-data; name="Float32ps"` + "\n\n" + `53.2
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64s"` + "\n\n" + `53.03
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64s"` + "\n\n" + `52.04
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64ps"` + "\n\n" + `53.0005
--foobar` + "\n" + `Content-Disposition: form-data; name="Float64ps"` + "\n\n" + `53.6
--foobar--
`

type embed0 struct {
	Field string
	embed1
}

type embed1 struct {
	Field int
	embed2
}

type embed2 struct {
	Field float64
}

var embedVal = embed0{
	embed1: embed1{
		embed2: embed2{34.67},
		Field:  3467,
	},
	Field: "string",
}

var embedValues = url.Values{"Field": {"string", "3467", "34.67"}}

const embedValString = `Field=string&Field=3467&Field=34.67`
const embedValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="Field"` + "\n\n" + `string
--foobar` + "\n" + `Content-Disposition: form-data; name="Field"` + "\n\n" + `3467
--foobar` + "\n" + `Content-Disposition: form-data; name="Field"` + "\n\n" + `34.67
--foobar--
`

type marshalSlice []string

func (s *marshalSlice) MarshalText() ([]byte, error) {
	return []byte(strings.Join(*s, ",")), nil
}

func (s *marshalSlice) UnmarshalText(text []byte) error {
	texts := strings.Split(string(text), ",")
	for _, txt := range texts {
		*s = append(*s, txt)
	}
	return nil
}

type marshalType struct {
	M *marshalSlice
	N marshalSlice
}

var marshalVal = marshalType{
	M: &marshalSlice{"foo", "bar", "baz"},
	N: marshalSlice{"foo", "bar", "baz"},
}

var marshalValues = url.Values{"M": {"foo", "bar", "baz"}, "N": {"foo", "bar", "baz"}}

const marshalValString = `M=foo%2Cbar%2Cbaz&N=foo&N=bar&N=baz`
const marshalValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="M"` + "\n\n" + `foo,bar,baz
--foobar` + "\n" + `Content-Disposition: form-data; name="N"` + "\n\n" + `foo
--foobar` + "\n" + `Content-Disposition: form-data; name="N"` + "\n\n" + `bar
--foobar` + "\n" + `Content-Disposition: form-data; name="N"` + "\n\n" + `baz
--foobar--
`

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

var ifaceValues = url.Values{"IString": {"foo"}, "IInt": {"32"}, "IBool": {"true"}, "ISlice": {"abc", "32.123455"}}

const ifaceValString = `IString=foo&IInt=32&IBool=true&ISlice=abc&ISlice=32.123455`
const ifaceValStringMultipart = `
--foobar` + "\n" + `Content-Disposition: form-data; name="IString"` + "\n\n" + `foo
--foobar` + "\n" + `Content-Disposition: form-data; name="IInt"` + "\n\n" + `32
--foobar` + "\n" + `Content-Disposition: form-data; name="IBool"` + "\n\n" + `true
--foobar` + "\n" + `Content-Disposition: form-data; name="ISlice"` + "\n\n" + `abc
--foobar` + "\n" + `Content-Disposition: form-data; name="ISlice"` + "\n\n" + `32.123455
--foobar--
`

func TestMarshal(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
		str  string
		err  error
	}{{
		name: "bool values",
		val:  boolVal,
		str:  boolValString,
	}, {
		name: "int values",
		val:  intVal,
		str:  intValString,
	}, {
		name: "int pointer values",
		val:  intpVal,
		str:  intpValString,
	}, {
		name: "int slices",
		val:  intsVal,
		str:  intsValString,
	}, {
		name: "int pointer slices",
		val:  intpsVal,
		str:  intpsValString,
	}, {
		name: "uint values",
		val:  uintVal,
		str:  uintValString,
	}, {
		name: "uint pointer values",
		val:  uintpVal,
		str:  uintpValString,
	}, {
		name: "uint slices",
		val:  uintsVal,
		str:  uintsValString,
	}, {
		name: "uint pointer slices",
		val:  uintpsVal,
		str:  uintpsValString,
	}, {
		name: "string values",
		val:  stringVal,
		str:  stringValString,
	}, {
		name: "float values",
		val:  floatVal,
		str:  floatValString,
	}, {
		name: "float slices",
		val:  floatsVal,
		str:  floatsValString,
	}, {
		name: "embedded types",
		val:  embedVal,
		str:  embedValString,
	}, {
		name: "TextMarshaler type",
		val:  marshalVal,
		str:  marshalValString,
	}, {
		name: "interface values",
		val:  ifaceVal,
		str:  ifaceValString,
	}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if bgot, err := Marshal(tt.val); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: Marshal error got %v, want %v", i, err, tt.err)
			} else if got := string(bgot); got != tt.str {
				t.Errorf("#%d: got %q, want %q", i, got, tt.str)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		data string
		dst  interface{}
		want interface{}
		err  error
	}{{
		name: "bool values",
		data: boolValString,
		dst:  &boolType{},
		want: &boolVal,
	}, {
		name: "int values",
		data: intValString,
		dst:  &intType{},
		want: &intVal,
	}, {
		name: "int pointer values",
		data: intpValString,
		dst:  &intpType{},
		want: &intpVal,
	}, {
		name: "int slices",
		data: intsValString,
		dst:  &intsType{},
		want: &intsVal,
	}, {
		name: "int pointer slices",
		data: intpsValString,
		dst:  &intpsType{},
		want: &intpsVal,
	}, {
		name: "uint values",
		data: uintValString,
		dst:  &uintType{},
		want: &uintVal,
	}, {
		name: "uint pointer values",
		data: uintpValString,
		dst:  &uintpType{},
		want: &uintpVal,
	}, {
		name: "uint slices",
		data: uintsValString,
		dst:  &uintsType{},
		want: &uintsVal,
	}, {
		name: "uint pointer slices",
		data: uintpsValString,
		dst:  &uintpsType{},
		want: &uintpsVal,
	}, {
		name: "string values",
		data: stringValString,
		dst:  &stringType{},
		want: &stringVal,
	}, {
		name: "float values",
		data: floatValString,
		dst:  &floatType{},
		want: &floatVal,
	}, {
		name: "float slices",
		data: floatsValString,
		dst:  &floatsType{},
		want: &floatsVal,
	}, {
		name: "TextMarshaler type",
		data: marshalValString,
		dst:  &marshalType{},
		want: &marshalVal,
	}, {
		name: "embeded types",
		data: embedValString,
		dst:  &embed0{},
		want: &embedVal,
	}, {
		name: "interface values",
		data: ifaceValString,
		dst:  &ifaceType{},
		want: &ifaceVal,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal([]byte(tt.data), tt.dst); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: got %#v, want %#v", i, err, tt.err)
			} else if !reflect.DeepEqual(tt.want, tt.dst) {
				t.Errorf("#%d: got %+v, want %+v", i, tt.dst, tt.want)
			}
		})
	}
}

func TestTransform(t *testing.T) {
	tests := []struct {
		name string
		vals url.Values
		dst  interface{}
		want interface{}
		err  error
	}{{
		name: "bool values",
		vals: boolValues,
		dst:  &boolType{},
		want: &boolVal,
	}, {
		name: "int values",
		vals: intValues,
		dst:  &intType{},
		want: &intVal,
	}, {
		name: "int pointer values",
		vals: intpValues,
		dst:  &intpType{},
		want: &intpVal,
	}, {
		name: "int slices",
		vals: intsValues,
		dst:  &intsType{},
		want: &intsVal,
	}, {
		name: "int pointer slices",
		vals: intpsValues,
		dst:  &intpsType{},
		want: &intpsVal,
	}, {
		name: "uint values",
		vals: uintValues,
		dst:  &uintType{},
		want: &uintVal,
	}, {
		name: "uint pointer values",
		vals: uintpValues,
		dst:  &uintpType{},
		want: &uintpVal,
	}, {
		name: "uint slices",
		vals: uintsValues,
		dst:  &uintsType{},
		want: &uintsVal,
	}, {
		name: "uint pointer slices",
		vals: uintpsValues,
		dst:  &uintpsType{},
		want: &uintpsVal,
	}, {
		name: "string values",
		vals: stringValues,
		dst:  &stringType{},
		want: &stringVal,
	}, {
		name: "float values",
		vals: floatValues,
		dst:  &floatType{},
		want: &floatVal,
	}, {
		name: "float slices",
		vals: floatsValues,
		dst:  &floatsType{},
		want: &floatsVal,
	}, {
		name: "TextMarshaler type",
		vals: marshalValues,
		dst:  &marshalType{},
		want: &marshalVal,
	}, {
		name: "embeded types",
		vals: embedValues,
		dst:  &embed0{},
		want: &embedVal,
	}, {
		name: "interface values",
		vals: ifaceValues,
		dst:  &ifaceType{},
		want: &ifaceVal,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Transform(tt.vals, tt.dst); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: got %#v, want %#v", i, err, tt.err)
			} else if !reflect.DeepEqual(tt.want, tt.dst) {
				t.Errorf("#%d: got %+v, want %+v", i, tt.dst, tt.want)
			}
		})
	}
}

func TestUnmarshal_errors(t *testing.T) {
	tests := []struct {
		name string
		val  string
		dst  interface{}
		err  error
	}{{
		name: "nil dst argument should err",
		val:  "foo=bar",
		dst:  nil,
		err:  &ArgumentError{nil},
	}, {
		name: "nil dst should return error",
		dst:  nil,
		err:  &ArgumentError{nil},
	}, {
		name: "non-pointer-struct dst should return error",
		dst:  struct{}{},
		err:  &ArgumentError{reflect.TypeOf(struct{}{})},
	}, {
		name: "any dst that is not a pointer to a struct should return error",
		dst:  []int{},
		err:  &ArgumentError{reflect.TypeOf([]int{})},
	}, {
		name: "bool type err",
		val:  "Bool=Hello World",
		dst:  &boolType{},
		err:  &ValueError{Key: "Bool", Value: "Hello World", Type: "bool"},
	}, {
		name: "int type err",
		val:  "Int=Hello World",
		dst:  &intType{},
		err:  &ValueError{Key: "Int", Value: "Hello World", Type: "int"},
	}, {
		name: "int8 value out of range",
		val:  "Int8=128",
		dst:  &intType{},
		err:  &ValueError{Key: "Int8", Value: "128", Type: "int8"},
	}, {
		name: "int16 value out of range",
		val:  "Int16=32768",
		dst:  &intType{},
		err:  &ValueError{Key: "Int16", Value: "32768", Type: "int16"},
	}, {
		name: "int32 value out of range",
		val:  "Int32=2147483648",
		dst:  &intType{},
		err:  &ValueError{Key: "Int32", Value: "2147483648", Type: "int32"},
	}, {
		name: "int64 value out of range",
		val:  "Int64=9223372036854775808",
		dst:  &intType{},
		err:  &ValueError{Key: "Int64", Value: "9223372036854775808", Type: "int64"},
	}, {
		name: "uint type err",
		val:  "Uint=Hello World",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint", Value: "Hello World", Type: "uint"},
	}, {
		name: "uint value out of range",
		val:  "Uint=-128",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint", Value: "-128", Type: "uint"},
	}, {
		name: "uint8 value out of range",
		val:  "Uint8=256",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint8", Value: "256", Type: "uint8"},
	}, {
		name: "uint16 value out of range",
		val:  "Uint16=65536",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint16", Value: "65536", Type: "uint16"},
	}, {
		name: "uint32 value out of range",
		val:  "Uint32=4294967296",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint32", Value: "4294967296", Type: "uint32"},
	}, {
		name: "uint64 value out of range",
		val:  "Uint64=18446744073709551616",
		dst:  &uintType{},
		err:  &ValueError{Key: "Uint64", Value: "18446744073709551616", Type: "uint64"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal([]byte(tt.val), tt.dst); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("error got %v, want %v", err, tt.err)
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
