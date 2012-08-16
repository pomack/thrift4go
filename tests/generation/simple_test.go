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
