package simple

import (
	"testing"
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
