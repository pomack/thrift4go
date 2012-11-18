package simple

import (
	"testing"
	"thrift"
)

func TestUndefinedValuesString(t *testing.T) {
	var inputAndExpected = []struct {
		in  UndefinedValues
		out string
	}{
		{UndefinedValues_One, "UndefinedValues_One"},
		{UndefinedValues_Two, "UndefinedValues_Two"},
		{UndefinedValues_Three, "UndefinedValues_Three"},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.String()

		if actual != definition.out {
			t.Errorf("%d. %q.String() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestDefinedValuesString(t *testing.T) {
	var inputAndExpected = []struct {
		in  DefinedValues
		out string
	}{
		{DefinedValues_One, "DefinedValues_One"},
		{DefinedValues_Two, "DefinedValues_Two"},
		{DefinedValues_Three, "DefinedValues_Three"},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.String()

		if actual != definition.out {
			t.Errorf("%d. %q.String() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestHeterogeneousValuesString(t *testing.T) {
	var inputAndExpected = []struct {
		in  HeterogeneousValues
		out string
	}{
		{HeterogeneousValues_One, "HeterogeneousValues_One"},
		{HeterogeneousValues_Two, "HeterogeneousValues_Two"},
		{HeterogeneousValues_Three, "HeterogeneousValues_Three"},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.String()

		if actual != definition.out {
			t.Errorf("%d. %q.String() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestUndefinedValuesValue(t *testing.T) {
	var inputAndExpected = []struct {
		in  UndefinedValues
		out int
	}{
		{UndefinedValues_One, 0},
		{UndefinedValues_Two, 1},
		{UndefinedValues_Three, 2},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.Value()

		if actual != definition.out {
			t.Errorf("%d. %q.Value() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestDefinedValuesValue(t *testing.T) {
	var inputAndExpected = []struct {
		in  DefinedValues
		out int
	}{
		{DefinedValues_One, 1},
		{DefinedValues_Two, 2},
		{DefinedValues_Three, 3},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.Value()

		if actual != definition.out {
			t.Errorf("%d. %q.Value() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestHeterogeneousValuesValue(t *testing.T) {
	var inputAndExpected = []struct {
		in  HeterogeneousValues
		out int
	}{
		{HeterogeneousValues_One, 0},
		{HeterogeneousValues_Two, 2},
		{HeterogeneousValues_Three, 3},
	}

	for i, definition := range inputAndExpected {
		actual := definition.in.Value()

		if actual != definition.out {
			t.Errorf("%d. %q.Value() => %q, want %q", i, definition.in, actual, definition.out)
		}
	}
}

func TestContainerOfEnumsNew(t *testing.T) {
	emission := NewContainerOfEnums()

	if emission == nil {
		t.Errorf("NewContainerOfEnums emitted nil, not the struct.")
	}
}

func TestContainerOfEnumsFieldsSet(t *testing.T) {
	emission := NewContainerOfEnums()

	emission.First = UndefinedValues_One
	emission.Second = DefinedValues_One
	emission.Third = HeterogeneousValues_One
	emission.OptionalFourth = UndefinedValues_Two
	emission.OptionalFifth = DefinedValues_Two
	emission.OptionalSixth = HeterogeneousValues_Two
	emission.DefaultSeventh = UndefinedValues_Two
	emission.DefaultEighth = DefinedValues_Two
	emission.DefaultNineth = HeterogeneousValues_Two
}

func TestContainerOfEnumsDefaultFieldsGet(t *testing.T) {
	emission := NewContainerOfEnums()

	definitions := []thrift.Enumer{
		emission.First,
		emission.Second,
		emission.Third,
		emission.OptionalFourth,
		emission.OptionalFifth,
		emission.OptionalSixth,
	}

	for i, definition := range definitions {
		actual := definition.String()

		if "<UNSET>" != actual {
			t.Errorf("%d. %q.String() => %q, want %q", i, definition, actual, "<UNSET>")
		}
	}

	if emission.DefaultSeventh.String() != "UndefinedValues_One" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultSeventh, emission.DefaultSeventh.String(), "UndefinedValues_One")
	}

	if emission.DefaultEighth.String() != "DefinedValues_One" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultEighth, emission.DefaultEighth.String(), "DefinedValues_One")
	}

	if emission.DefaultNineth.String() != "HeterogeneousValues_One" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultNineth, emission.DefaultNineth.String(), "HeterogeneousValues_One")
	}
}

// To validate https://github.com/pomack/thrift4go/issues/16.
func TestContainerOfEnumsOptionalFieldsAreSetStatusByDefault(t *testing.T) {
	emission := NewContainerOfEnums()

	var valueAndExpected = []struct {
		value    bool
		expected bool
	}{
		{emission.IsSetOptionalFourth(), false},
		{emission.IsSetOptionalFifth(), false},
		{emission.IsSetOptionalSixth(), false},
		{emission.IsSetDefaultSeventh(), true},
		{emission.IsSetDefaultEighth(), true},
		{emission.IsSetDefaultNineth(), true},
	}

	for i, definition := range valueAndExpected {
		actual := definition.value
		expected := definition.expected

		if actual != expected {
			t.Errorf("%d. %q.IsSet() => %q, want %q", i, definition, actual, expected)
		}
	}
}

// To validate https://github.com/pomack/thrift4go/issues/16.
func TestContainerOfEnumsOptionalFieldsAreSetStatusAfterSet(t *testing.T) {
	emission := NewContainerOfEnums()

	emission.OptionalFourth = UndefinedValues_One

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFourth = UndefinedValues_Two

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFourth = UndefinedValues_Three

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFifth = DefinedValues_One

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedValues_Two

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedValues_Three

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalSixth = HeterogeneousValues_One

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.OptionalSixth = HeterogeneousValues_Two

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.OptionalSixth = HeterogeneousValues_Three

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.DefaultSeventh = UndefinedValues_One

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultSeventh = UndefinedValues_Two

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultSeventh = UndefinedValues_Three

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultEighth = DefinedValues_One

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultEighth = DefinedValues_Two

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultEighth = DefinedValues_Three

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultNineth = HeterogeneousValues_One

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}

	emission.DefaultNineth = HeterogeneousValues_Two

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}

	emission.DefaultNineth = HeterogeneousValues_Three

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}
}

type protocolBuilder func() thrift.TProtocol

func TestWireFormatWithDefaultPayload(t *testing.T) {
	var transport *thrift.TMemoryBuffer

	var protocols = []struct {
		name    string
		builder protocolBuilder
	}{
		{
			"TBinaryProtocol",
			func() thrift.TProtocol {
				return thrift.NewTBinaryProtocolTransport(transport)
			},
		},
		{
			"TCompactProtocol",
			func() thrift.TProtocol {
				return thrift.NewTCompactProtocol(transport)
			},
		},
		{
			"TJSONProtocol",
			func() thrift.TProtocol {
				return thrift.NewTJSONProtocol(transport)
			},
		},
		{
			"TSimpleJSONProtocol",
			func() thrift.TProtocol {
				return thrift.NewTSimpleJSONProtocol(transport)
			},
		},
	}

	for i, definition := range protocols {
		transport = thrift.NewTMemoryBuffer()
		defer transport.Close()
		protocol := definition.builder()
		name := definition.name

		emission := NewContainerOfEnums()

		if err := emission.Write(protocol); err != nil {
			t.Fatalf("%d (%s): Could not emit simple %q to JSON.", i, name, emission)
		}

		if err := protocol.Flush(); err != nil {
			t.Fatalf("%d (%s): Could not flush emission.", i, name)
		}

		incoming := NewContainerOfEnums()

		if err := incoming.Read(protocol); err != nil {
			t.Fatalf("%d (%s): Could not read from buffer: %q", i, name, err)
		}

		if emission.First != incoming.First {
			t.Errorf("%d (%s) emission.First (%q) != incoming.First (%q)", i, name, emission.First, incoming.First)
		}

		if emission.Second != incoming.Second {
			t.Errorf("%d (%s) emission.Second (%q) != incoming.Second (%q)", i, name, emission.Second, incoming.Second)
		}

		if emission.Third != incoming.Third {
			t.Errorf("%d (%s) emission.Third (%q) != incoming.Third (%q)", i, name, emission.Third, incoming.Third)
		}

		if emission.OptionalFourth != incoming.OptionalFourth {
			t.Errorf("%d (%s) emission.OptionalFourth (%q) != incoming.OptionalFourth (%q)", i, name, emission.OptionalFourth, incoming.OptionalFourth)
		}

		if emission.OptionalFifth != incoming.OptionalFifth {
			t.Errorf("%d (%s) emission.OptionalFifth (%q) != incoming.OptionalFifth (%q)", i, name, emission.OptionalFifth, incoming.OptionalFifth)
		}

		if emission.OptionalSixth != incoming.OptionalSixth {
			t.Errorf("%d (%s) emission.OptionalSixth (%q) != incoming.OptionalSixth (%q)", i, name, emission.OptionalSixth, incoming.OptionalSixth)
		}

		if emission.DefaultSeventh != incoming.DefaultSeventh {
			t.Errorf("%d (%s) emission.DefaultSeventh (%q) != incoming.DefaultSeventh (%q)", i, name, emission.DefaultSeventh, incoming.DefaultSeventh)
		}

		if emission.DefaultEighth != incoming.DefaultEighth {
			t.Errorf("%d (%s) emission.DefaultEighth (%q) != incoming.DefaultEighth (%q)", i, name, emission.DefaultEighth, incoming.DefaultEighth)
		}

		if emission.DefaultNineth != incoming.DefaultNineth {
			t.Errorf("%d (%s) emission.DefaultNineth (%q) != incoming.DefaultNineth (%q)", i, name, emission.DefaultNineth, incoming.DefaultNineth)
		}

		if emission.IsSetOptionalFourth() != incoming.IsSetOptionalFourth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFourth (%q) != incoming.IsSetOptionalFourth (%q)", i, name, emission.IsSetOptionalFourth(), incoming.IsSetOptionalFourth())
		}

		if emission.IsSetOptionalFifth() != incoming.IsSetOptionalFifth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFifth (%q) != incoming.IsSetOptionalFifth (%q)", i, name, emission.IsSetOptionalFifth(), incoming.IsSetOptionalFifth())
		}

		if emission.IsSetOptionalSixth() != incoming.IsSetOptionalSixth() {
			t.Errorf("%d (%s) emission.IsSetOptionalSixth (%q) != incoming.IsSetOptionalSixth (%q)", i, name, emission.IsSetOptionalSixth(), incoming.IsSetOptionalSixth())
		}

		if emission.IsSetDefaultSeventh() != incoming.IsSetDefaultSeventh() {
			t.Errorf("%d (%s) emission.IsSetDefaultSeventh (%q) != incoming.IsSetDefaultSeventh (%q)", i, name, emission.IsSetDefaultSeventh(), incoming.IsSetDefaultSeventh())
		}

		if emission.IsSetDefaultEighth() != incoming.IsSetDefaultEighth() {
			t.Errorf("%d (%s) emission.IsSetDefaultEighth (%q) != incoming.IsSetDefaultEighth (%q)", i, name, emission.IsSetDefaultEighth(), incoming.IsSetDefaultEighth())
		}

		if emission.IsSetDefaultNineth() != incoming.IsSetDefaultNineth() {
			t.Errorf("%d (%s) emission.IsSetDefaultNineth (%q) != incoming.IsSetDefaultNineth (%q)", i, name, emission.IsSetDefaultNineth(), incoming.IsSetDefaultNineth())
		}
	}
}

func TestWireFormatWithSetPayload(t *testing.T) {
	var transport *thrift.TMemoryBuffer

	var protocols = []struct {
		name    string
		builder protocolBuilder
	}{
		{
			"TBinaryProtocol",
			func() thrift.TProtocol {
				return thrift.NewTBinaryProtocolTransport(transport)
			},
		},
		{
			"TCompactProtocol",
			func() thrift.TProtocol {
				return thrift.NewTCompactProtocol(transport)
			},
		},
		{
			"TJSONProtocol",
			func() thrift.TProtocol {
				return thrift.NewTJSONProtocol(transport)
			},
		},
		{
			"TSimpleJSONProtocol",
			func() thrift.TProtocol {
				return thrift.NewTSimpleJSONProtocol(transport)
			},
		},
	}

	for i, definition := range protocols {
		transport = thrift.NewTMemoryBuffer()
		defer transport.Close()
		protocol := definition.builder()
		name := definition.name

		emission := NewContainerOfEnums()
		emission.First = UndefinedValues_Two
		emission.Second = DefinedValues_Two
		emission.Third = HeterogeneousValues_Two
		emission.OptionalFourth = UndefinedValues_Three
		emission.OptionalFifth = DefinedValues_Three
		emission.OptionalSixth = HeterogeneousValues_Three

		if err := emission.Write(protocol); err != nil {
			t.Fatalf("%d (%s): Could not emit simple %q to JSON.", i, name, emission)
		}

		if err := protocol.Flush(); err != nil {
			t.Fatalf("%d (%s): Could not flush emission.", i, name)
		}

		incoming := NewContainerOfEnums()

		if err := incoming.Read(protocol); err != nil {
			t.Fatalf("%d (%s): Could not read from buffer: %q", i, name, err)
		}

		if emission.First != incoming.First {
			t.Errorf("%d (%s) emission.First (%q) != incoming.First (%q)", i, name, emission.First, incoming.First)
		}

		if emission.Second != incoming.Second {
			t.Errorf("%d (%s) emission.Second (%q) != incoming.Second (%q)", i, name, emission.Second, incoming.Second)
		}

		if emission.Third != incoming.Third {
			t.Errorf("%d (%s) emission.Third (%q) != incoming.Third (%q)", i, name, emission.Third, incoming.Third)
		}

		if emission.OptionalFourth != incoming.OptionalFourth {
			t.Errorf("%d (%s) emission.OptionalFourth (%q) != incoming.OptionalFourth (%q)", i, name, emission.OptionalFourth, incoming.OptionalFourth)
		}

		if emission.OptionalFifth != incoming.OptionalFifth {
			t.Errorf("%d (%s) emission.OptionalFifth (%q) != incoming.OptionalFifth (%q)", i, name, emission.OptionalFifth, incoming.OptionalFifth)
		}

		if emission.OptionalSixth != incoming.OptionalSixth {
			t.Errorf("%d (%s) emission.OptionalSixth (%q) != incoming.OptionalSixth (%q)", i, name, emission.OptionalSixth, incoming.OptionalSixth)
		}

		if emission.DefaultSeventh != incoming.DefaultSeventh {
			t.Errorf("%d (%s) emission.DefaultSeventh (%q) != incoming.DefaultSeventh (%q)", i, name, emission.DefaultSeventh, incoming.DefaultSeventh)
		}

		if emission.DefaultEighth != incoming.DefaultEighth {
			t.Errorf("%d (%s) emission.DefaultEighth (%q) != incoming.DefaultEighth (%q)", i, name, emission.DefaultEighth, incoming.DefaultEighth)
		}

		if emission.DefaultNineth != incoming.DefaultNineth {
			t.Errorf("%d (%s) emission.DefaultNineth (%q) != incoming.DefaultNineth (%q)", i, name, emission.DefaultNineth, incoming.DefaultNineth)
		}

		if emission.IsSetOptionalFourth() != incoming.IsSetOptionalFourth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFourth (%q) != incoming.IsSetOptionalFourth (%q)", i, name, emission.IsSetOptionalFourth(), incoming.IsSetOptionalFourth())
		}

		if emission.IsSetOptionalFifth() != incoming.IsSetOptionalFifth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFifth (%q) != incoming.IsSetOptionalFifth (%q)", i, name, emission.IsSetOptionalFifth(), incoming.IsSetOptionalFifth())
		}

		if emission.IsSetOptionalSixth() != incoming.IsSetOptionalSixth() {
			t.Errorf("%d (%s) emission.IsSetOptionalSixth (%q) != incoming.IsSetOptionalSixth (%q)", i, name, emission.IsSetOptionalSixth(), incoming.IsSetOptionalSixth())
		}

		if emission.IsSetDefaultSeventh() != incoming.IsSetDefaultSeventh() {
			t.Errorf("%d (%s) emission.IsSetDefaultSeventh (%q) != incoming.IsSetDefaultSeventh (%q)", i, name, emission.IsSetDefaultSeventh(), incoming.IsSetDefaultSeventh())
		}

		if emission.IsSetDefaultEighth() != incoming.IsSetDefaultEighth() {
			t.Errorf("%d (%s) emission.IsSetDefaultEighth (%q) != incoming.IsSetDefaultEighth (%q)", i, name, emission.IsSetDefaultEighth(), incoming.IsSetDefaultEighth())
		}

		if emission.IsSetDefaultNineth() != incoming.IsSetDefaultNineth() {
			t.Errorf("%d (%s) emission.IsSetDefaultNineth (%q) != incoming.IsSetDefaultNineth (%q)", i, name, emission.IsSetDefaultNineth(), incoming.IsSetDefaultNineth())
		}
	}
}
