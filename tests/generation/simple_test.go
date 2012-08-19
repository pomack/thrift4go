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
		t.Errorf("NewContainerOfEnums emitted nil, not the struct.")
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
	emission.DefaultSeventh = UndefinedTwo
	emission.DefaultEighth = DefinedTwo
	emission.DefaultNineth = HeterogeneousTwo
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

	if emission.DefaultSeventh.String() != "UndefinedOne" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultSeventh, emission.DefaultSeventh.String(), "UndefinedOne")
	}

	if emission.DefaultEighth.String() != "DefinedOne" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultEighth, emission.DefaultEighth.String(), "DefinedOne")
	}

	if emission.DefaultNineth.String() != "HeterogeneousOne" {
		t.Errorf("%q.String() => %q, want %q", emission.DefaultNineth, emission.DefaultNineth.String(), "HeterogeneousOne")
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

	emission.OptionalFourth = UndefinedOne

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFourth = UndefinedTwo

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFourth = UndefinedThree

	if emission.IsSetOptionalFourth() != true {
		t.Errorf("emission.OptionalFourth = %q; emission.IsSetOptionalFourth() => %s, want %s", emission.OptionalFourth, emission.IsSetOptionalFourth(), true)
	}

	emission.OptionalFifth = DefinedOne

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedTwo

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalFifth = DefinedThree

	if emission.IsSetOptionalFifth() != true {
		t.Errorf("emission.OptionalFifth = %q; emission.IsSetOptionalFifth() => %s, want %s", emission.OptionalFifth, emission.IsSetOptionalFifth(), true)
	}

	emission.OptionalSixth = HeterogeneousOne

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.OptionalSixth = HeterogeneousTwo

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.OptionalSixth = HeterogeneousThree

	if emission.IsSetOptionalSixth() != true {
		t.Errorf("emission.OptionalSixth = %q; emission.IsSetOptionalSixth() => %s, want %s", emission.OptionalSixth, emission.IsSetOptionalSixth(), true)
	}

	emission.DefaultSeventh = UndefinedOne

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultSeventh = UndefinedTwo

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultSeventh = UndefinedThree

	if emission.IsSetDefaultSeventh() != true {
		t.Errorf("emission.DefaultSeventh = %q; emission.IsSetDefaultSeventh() => %s, want %s", emission.DefaultSeventh, emission.IsSetDefaultSeventh(), true)
	}

	emission.DefaultEighth = DefinedOne

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultEighth = DefinedTwo

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultEighth = DefinedThree

	if emission.IsSetDefaultEighth() != true {
		t.Errorf("emission.DefaultEighth = %q; emission.IsSetDefaultEighth() => %s, want %s", emission.DefaultEighth, emission.IsSetDefaultEighth(), true)
	}

	emission.DefaultNineth = HeterogeneousOne

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}

	emission.DefaultNineth = HeterogeneousTwo

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}

	emission.DefaultNineth = HeterogeneousThree

	if emission.IsSetDefaultNineth() != true {
		t.Errorf("emission.DefaultNineth = %q; emission.IsSetDefaultNineth() => %s, want %s", emission.DefaultNineth, emission.IsSetDefaultNineth(), true)
	}
}
