struct intstruct {
  1: i32 ifield;
}

struct manylists {
  1: list<byte> bytelist /* = [1, 2, 3] */,
  2: optional list<i16> i16list = [1, 2, 3],
  3: list<i32> i32list = [1, 2, 3],
  4: list<string> stringlist = ["one", "two", "three"],
  5: list<intstruct> structlist = [{'ifield':1}, {'ifield':2}, {'ifield':3}],
  6: list<list<i32>> listoflist,
  7: optional i32 optionalint = 2,
}

const manylists MANYLIST_ITEM = {
 'bytelist': [5, 6, 7],
 'i16list': [5, 6, 7],
 'i32list': [5, 6, 7],
 'stringlist': ["five", "six", "seven"],
 'structlist': [{'ifield':5}, {'ifield':6}, {'ifield':7}],
}
