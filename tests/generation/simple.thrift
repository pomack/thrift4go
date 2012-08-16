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
}
