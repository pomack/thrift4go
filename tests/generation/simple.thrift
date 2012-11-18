namespace java thrift4go.generated

enum UndefinedValues {
  One,
  Two,
  Three,
}

enum DefinedValues {
  One = 1,
  Two = 2,
  Three = 3,
}

enum HeterogeneousValues {
  One,
  Two = 2,
  Three,
  Four = 4,
}

struct ContainerOfEnums {
  1: UndefinedValues first,
  2: DefinedValues second,
  3: HeterogeneousValues third,
  4: optional UndefinedValues optional_fourth,
  5: optional DefinedValues optional_fifth,
  6: optional HeterogeneousValues optional_sixth,
  7: optional UndefinedValues default_seventh = UndefinedValues.One,
  8: optional DefinedValues default_eighth = DefinedValues.One,
  9: optional HeterogeneousValues default_nineth = HeterogeneousValues.One,
}

service ContainerOfEnumsTestService {
  ContainerOfEnums echo(1: ContainerOfEnums message);
}
