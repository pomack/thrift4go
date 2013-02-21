struct intstruct {
  1: i32 ifield;
}

const map<string,string> MAPCONSTANT = {'hello':'world', 'goodnight':'moon'}

struct manymaps {
  1: map<byte, string> bytemap ,
  2: optional map<i16, string> i16map,
  3: map<i32, string> i32map,
  4: map<string,string> stringmap,
  5: map<binary,string> binarymap,
  6: map<binary,intstruct> binarstructymap,
}


/*
const manylists MANYLIST_ITEM = {
 'bytelist': [5, 6, 7],
 'i16list': [5, 6, 7],
 'i32list': [5, 6, 7],
 'stringlist': ["five", "six", "seven"],
 'structlist': [{'ifield':5}, {'ifield':6}, {'ifield':7}],
}
*/

