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
		{UndefinedOne, "UndefinedOne"},
		{UndefinedTwo, "UndefinedTwo"},
		{UndefinedThree, "UndefinedThree"},
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
		{DefinedOne, "DefinedOne"},
		{DefinedTwo, "DefinedTwo"},
		{DefinedThree, "DefinedThree"},
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
		{HeterogeneousOne, "HeterogeneousOne"},
		{HeterogeneousTwo, "HeterogeneousTwo"},
		{HeterogeneousThree, "HeterogeneousThree"},
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
		{UndefinedOne, 0},
		{UndefinedTwo, 1},
		{UndefinedThree, 2},
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
		{DefinedOne, 1},
		{DefinedTwo, 2},
		{DefinedThree, 3},
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
		{HeterogeneousOne, 0},
		{HeterogeneousTwo, 2},
		{HeterogeneousThree, 3},
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
		t.Fatalf("NewContainerOfEnums emitted nil, not the struct.")
	}
}

func TestContainerOfEnumsFieldsSet(t *testing.T) {
	emission := NewContainerOfEnums()

	emission.First = UndefinedOne
	emission.Second = DefinedOne
	emission.Third = HeterogeneousOne
	emission.OptionalFourth = UndefinedTwo
	emission.OptionalFifth = DefinedTwo
	emission.OptionalSixth = HeterogeneousTwo
}

func TestContainerOfEnumsDefaultFieldsGet(t *testing.T) {
	emission := NewContainerOfEnums()

	if emission.First.String() != "UndefinedOne" {
		t.Fatalf("emission.First = %q, want %q", emission.First.String(), "UndefinedOne")
	}

	if emission.Second.String() != "" {
		t.Fatalf("emission.Second = %q, want %q", emission.Second.String(), "")
	}

	if emission.Third.String() != "HeterogeneousOne" {
		t.Fatalf("emission.Third = %q, want %q", emission.Third.String(), "HeterogeneousOne")
	}

	if emission.OptionalFourth.String() != "UndefinedOne" {
		t.Fatalf("emission.OptionalFourth = %q, want %q", emission.OptionalFourth.String(), "UndefinedOne")
	}

	if emission.OptionalFifth.String() != "" {
		t.Fatalf("emission.OptionalFifth = %q, want %q", emission.OptionalFifth.String(), "")
	}

	if emission.OptionalSixth.String() != "HeterogeneousOne" {
		t.Fatalf("emission.OptionalSixth = %q, want %q", emission.OptionalSixth.String(), "HeterogeneousOne")
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

	emission.OptionalFourth = UndefinedOne

	if emission.IsSetOptionalFourth() != false {
		t.Fatalf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), false)
	}

	emission.OptionalFourth = UndefinedTwo

	if emission.IsSetOptionalFourth() != true {
		t.Fatalf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFourth = UndefinedThree

	if emission.IsSetOptionalFourth() != true {
		t.Fatalf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFifth = DefinedOne

	if emission.IsSetOptionalFifth() != true {
		t.Fatalf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedTwo

	if emission.IsSetOptionalFifth() != true {
		t.Fatalf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedThree

	if emission.IsSetOptionalFifth() != true {
		t.Fatalf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalSixth = HeterogeneousOne

	if emission.IsSetOptionalSixth() != false {
		t.Fatalf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), false)
	}

	emission.OptionalSixth = HeterogeneousTwo

	if emission.IsSetOptionalSixth() != true {
		t.Fatalf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.OptionalSixth = HeterogeneousThree

	if emission.IsSetOptionalSixth() != true {
		t.Fatalf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}
}

type protocolBuilder func() thrift.TProtocol

func TestWireFormatWithDefaultPayload(t *testing.T) {
	var transport thrift.TTransport

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
			t.Fatalf("%d (%s): Could not read from buffer: %q\n", i, name, err)
		}

		if emission.First != incoming.First {
			t.Errorf("%d (%s) emission.First (%q) != incoming.First (%q)\n", i, name, emission.First, incoming.First)
		}

		if emission.Second != incoming.Second {
			t.Errorf("%d (%s) emission.Second (%q) != incoming.Second (%q)\n", i, name, emission.Second, incoming.Second)
		}

		if emission.Third != incoming.Third {
			t.Errorf("%d (%s) emission.Third (%q) != incoming.Third (%q)\n", i, name, emission.Third, incoming.Third)
		}

		if emission.OptionalFourth != incoming.OptionalFourth {
			t.Errorf("%d (%s) emission.OptionalFourth (%q) != incoming.OptionalFourth (%q)\n", i, name, emission.OptionalFourth, incoming.OptionalFourth)
		}

		if emission.OptionalFifth != incoming.OptionalFifth {
			t.Errorf("%d (%s) emission.OptionalFifth (%q) != incoming.OptionalFifth (%q)\n", i, name, emission.OptionalFifth, incoming.OptionalFifth)
		}

		if emission.OptionalSixth != incoming.OptionalSixth {
			t.Errorf("%d (%s) emission.OptionalSixth (%q) != incoming.OptionalSixth (%q)\n", i, name, emission.OptionalSixth, incoming.OptionalSixth)
		}

		if emission.IsSetOptionalFourth() != incoming.IsSetOptionalFourth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFourth (%q) != incoming.IsSetOptionalFourth (%q)\n", i, name, emission.IsSetOptionalFourth(), incoming.IsSetOptionalFourth())
		}

		if emission.IsSetOptionalFifth() != incoming.IsSetOptionalFifth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFifth (%q) != incoming.IsSetOptionalFifth (%q)\n", i, name, emission.IsSetOptionalFifth(), incoming.IsSetOptionalFifth())
		}

		if emission.IsSetOptionalSixth() != incoming.IsSetOptionalSixth() {
			t.Errorf("%d (%s) emission.IsSetOptionalSixth (%q) != incoming.IsSetOptionalSixth (%q)\n", i, name, emission.IsSetOptionalSixth(), incoming.IsSetOptionalSixth())
		}
	}
}

func TestWireFormatWithSetPayload(t *testing.T) {
	var transport thrift.TTransport

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
		emission.First = UndefinedTwo
		emission.Second = DefinedTwo
		emission.Third = HeterogeneousTwo
		emission.OptionalFourth = UndefinedThree
		emission.OptionalFifth = DefinedThree
		emission.OptionalSixth = HeterogeneousThree

		if err := emission.Write(protocol); err != nil {
			t.Fatalf("%d (%s): Could not emit simple %q to JSON.", i, name, emission)
		}

		if err := protocol.Flush(); err != nil {
			t.Fatalf("%d (%s): Could not flush emission.", i, name)
		}

		incoming := NewContainerOfEnums()

		if err := incoming.Read(protocol); err != nil {
			t.Fatalf("%d (%s): Could not read from buffer: %q\n", i, name, err)
		}

		if emission.First != incoming.First {
			t.Errorf("%d (%s) emission.First (%q) != incoming.First (%q)\n", i, name, emission.First, incoming.First)
		}

		if emission.Second != incoming.Second {
			t.Errorf("%d (%s) emission.Second (%q) != incoming.Second (%q)\n", i, name, emission.Second, incoming.Second)
		}

		if emission.Third != incoming.Third {
			t.Errorf("%d (%s) emission.Third (%q) != incoming.Third (%q)\n", i, name, emission.Third, incoming.Third)
		}

		if emission.OptionalFourth != incoming.OptionalFourth {
			t.Errorf("%d (%s) emission.OptionalFourth (%q) != incoming.OptionalFourth (%q)\n", i, name, emission.OptionalFourth, incoming.OptionalFourth)
		}

		if emission.OptionalFifth != incoming.OptionalFifth {
			t.Errorf("%d (%s) emission.OptionalFifth (%q) != incoming.OptionalFifth (%q)\n", i, name, emission.OptionalFifth, incoming.OptionalFifth)
		}

		if emission.OptionalSixth != incoming.OptionalSixth {
			t.Errorf("%d (%s) emission.OptionalSixth (%q) != incoming.OptionalSixth (%q)\n", i, name, emission.OptionalSixth, incoming.OptionalSixth)
		}

		if emission.IsSetOptionalFourth() != incoming.IsSetOptionalFourth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFourth (%q) != incoming.IsSetOptionalFourth (%q)\n", i, name, emission.IsSetOptionalFourth(), incoming.IsSetOptionalFourth())
		}

		if emission.IsSetOptionalFifth() != incoming.IsSetOptionalFifth() {
			t.Errorf("%d (%s) emission.IsSetOptionalFifth (%q) != incoming.IsSetOptionalFifth (%q)\n", i, name, emission.IsSetOptionalFifth(), incoming.IsSetOptionalFifth())
		}

		if emission.IsSetOptionalSixth() != incoming.IsSetOptionalSixth() {
			t.Errorf("%d (%s) emission.IsSetOptionalSixth (%q) != incoming.IsSetOptionalSixth (%q)\n", i, name, emission.IsSetOptionalSixth(), incoming.IsSetOptionalSixth())
		}
	}

}
