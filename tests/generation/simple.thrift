enum UndefinedValues {
  UndefinedOne,
  UndefinedTwo,
  UndefinedThree,
}

enum DefinedValues {
  DefinedOne = 1,
  DefinedTwo = 2,
  DefinedThree = 3,
}

enum HeterogeneousValues {
  HeterogeneousOne,
  HeterogeneousTwo = 2,
  HeterogeneousThree,
  HeterogeneousFour = 4,
}

struct ContainerOfEnums {
  1: UndefinedValues first,
  2: DefinedValues second,
  3: HeterogeneousValues third,
  4: optional UndefinedValues optional_fourth,
  5: optional DefinedValues optional_fifth,
  6: optional HeterogeneousValues optional_sixth,
  7: optional UndefinedValues default_seventh = UndefinedValues.UndefinedOne,
  8: optional DefinedValues default_eighth = DefinedValues.DefinedOne,
  9: optional HeterogeneousValues default_nineth = HeterogeneousValues.HeterogeneousOne,
}
