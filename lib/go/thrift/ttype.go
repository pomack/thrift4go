/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


package thrift

import (
  "bytes"
  "container/list"
  "container/vector"
  "strconv"
)

/**
 * Type constants in the Thrift protocol.
 */
type TType interface {
    ThriftTypeId() byte
    IsBinary() bool
    IsUTF8() bool
    String() string
    IsBaseType() bool
    IsEmptyType() bool
    IsEnum() bool
    IsNumericType() bool
    IsStringType() bool
    IsContainer() bool
    IsStruct() bool
    IsMap() bool
    IsList() bool
    IsSet() bool
    IsInt() bool
    LessType(other interface{}) bool
    Less(i, j interface{}) bool
    Compare(i, j interface{}) (int, bool)
    CompareValueArrays(li, lj []interface{}) (int, bool)
    Equals(other interface{}) bool
    Coerce(other interface{}) TType
    CoerceData(data interface{}) (interface{}, bool)
}

type tType struct {
    thriftTypeId byte
    isBinary bool
    isEnum bool
    isUTF8 bool
}

const (
    iSTOP   = 0
    iVOID   = 1
    iBOOL   = 2
    iBYTE   = 3
    iI08    = 3
    iDOUBLE = 4
    iI16    = 6
    iI32    = 8
    iI64    = 10
    iSTRING = 11
    iUTF7   = 11
    iBINARY = 11
    iSTRUCT = 12
    iMAP    = 13
    iSET    = 14
    iLIST   = 15
    iENUM   = 16
    iUTF8   = 16
    iUTF16  = 17
    iGENERIC = 127
)

var (
  STOP    = newStandardTType(iSTOP)
  VOID    = newStandardTType(iVOID)
  BOOL    = newStandardTType(iBOOL)
  BYTE    = newStandardTType(iBYTE)
  I08     = newStandardTType(iI08)
  DOUBLE  = newStandardTType(iDOUBLE)
  I16     = newStandardTType(iI16)
  I32     = newStandardTType(iI32)
  I64     = newStandardTType(iI64)
  STRING  = newStandardTType(iSTRING)
  UTF7    = newStandardTType(iUTF7)
  BINARY  = newTType(iBINARY, true, false, false)
  STRUCT  = newStandardTType(iSTRUCT)
  MAP     = newStandardTType(iMAP)
  SET     = newStandardTType(iSET)
  LIST    = newStandardTType(iLIST)
  ENUM    = newTType(iENUM, false, true, false)
  UTF8    = newTType(iUTF8, false, false, true)
  UTF16   = newStandardTType(iUTF16)
  GENERIC   = newStandardTType(iGENERIC)
)


func newTType(thriftTypeId byte, isBinary, isEnum, isUTF8 bool) TType {
    return &tType{
        thriftTypeId: thriftTypeId, 
        isBinary: isBinary, 
        isEnum: isEnum, 
        isUTF8: isUTF8,
    }
}

func newStandardTType(thriftTypeId byte) TType {
    return &tType{
        thriftTypeId: thriftTypeId, 
        isBinary: false, 
        isEnum: false, 
        isUTF8: false,
    }
}

func TTypeFromThriftTypeId(thriftTypeId byte) TType {
    var theType TType
    switch thriftTypeId {
    default:
        theType = STOP
    case iSTOP:
        theType = STOP
    case iVOID: 
        theType = VOID
    case iBOOL: 
        theType = BOOL
    case iBYTE:
        theType = BYTE
    case iDOUBLE:
        theType = DOUBLE
    case iI16:
        theType = I16
    case iI32:
        theType = I32
    case iI64:
        theType = I64
    case iSTRING:
        theType = STRING
    case iSTRUCT:
        theType = STRUCT
    case iMAP:
        theType = MAP
    case iSET:
        theType = SET
    case iLIST:
        theType = LIST
    case iENUM:
        theType = ENUM
    case iUTF16:
        theType = UTF16
    case iGENERIC:
        theType = GENERIC
    }
    return theType
}

func (p *tType) ThriftTypeId() byte {
    return p.thriftTypeId
}

func (p *tType) IsBinary() bool {
    return p.isBinary
}

func (p *tType) IsEnum() bool {
    return p.isEnum
}

func (p *tType) IsUTF8() bool {
    return p.isUTF8
}

func (p *tType) String() string {
  switch p.ThriftTypeId() {
  case iSTOP:
    return "STOP"
  case iVOID:
    return "VOID"
  case iBOOL:
    return "BOOL"
  case iBYTE:
    return "BYTE"
  case iDOUBLE:
    return "DOUBLE"
  case iI16:
    return "I16"
  case iI32:
    return "I32"
  case iI64:
    return "I64"
  case iSTRING:
    if p.IsBinary() {
      return "BINARY"
    }
    return "STRING"
  case iSTRUCT:
    return "STRUCT"
  case iMAP:
    return "MAP"
  case iSET:
    return "SET"
  case iLIST:
    return "LIST"
  case iENUM:
    if p.IsUTF8() {
      return "UTF8"
    }
    return "ENUM"
  case iUTF16:
    return "UTF16"
  case iGENERIC:
    return "GENERIC"
  }
  return "Unknown"
}

func (p *tType) IsBaseType() bool {
  switch p.ThriftTypeId() {
  case iBOOL, iBYTE, iDOUBLE, iI16, iI32, iI64, iSTRING, iUTF8, iUTF16:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsEmptyType() bool {
  switch p.ThriftTypeId() {
  case iSTOP, iVOID:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsNumericType() bool {
  switch p.ThriftTypeId() {
  case iENUM, iBOOL, iBYTE, iDOUBLE, iI16, iI32, iI64:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsStringType() bool {
  switch p.ThriftTypeId() {
  case iSTRING, iUTF8, iUTF16:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsContainer() bool {
  switch p.ThriftTypeId() {
  case iMAP, iSET, iLIST:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsStruct() bool {
  switch p.ThriftTypeId() {
  case iSTRUCT:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsMap() bool {
  switch p.ThriftTypeId() {
  case iMAP:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsList() bool {
  switch p.ThriftTypeId() {
  case iLIST:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsSet() bool {
  switch p.ThriftTypeId() {
  case iSET:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) IsInt() bool {
  switch p.ThriftTypeId() {
  case iBYTE, iI16, iI32, iI64:
    return true
  default:
    return false
  }
  return false
}

func (p *tType) Coerce(other interface{}) TType {
  if other == nil {
    return STOP
  }
  switch b := other.(type) {
  default:
    return STOP
  case nil:
    return STOP
  case TType:
    return b
  case byte:
    return TTypeFromThriftTypeId(b)
  case int:
    return TTypeFromThriftTypeId(byte(b))
  case int8:
    return TTypeFromThriftTypeId(byte(b))
  case int32:
    return TTypeFromThriftTypeId(byte(b))
  case int64:
    return TTypeFromThriftTypeId(byte(b))
  case uint:
    return TTypeFromThriftTypeId(byte(b))
  case uint32:
    return TTypeFromThriftTypeId(byte(b))
  case uint64:
    return TTypeFromThriftTypeId(byte(b))
  case float32:
    return TTypeFromThriftTypeId(byte(int(b)))
  case float64:
    return TTypeFromThriftTypeId(byte(int(b)))
  }
  return STOP
}

func (p *tType) LessType(other interface{}) bool {
  return p.ThriftTypeId() < p.Coerce(other).ThriftTypeId()
}

func (p *tType) Less(i, j interface{}) bool {
  cmp, ok := p.Compare(i, j)
  return ok && cmp > 0
}


func (p *tType) Compare(i, j interface{}) (int, bool) {
  if i == j {
    return 0, true
  }
  if i == nil {
    if j == nil {
      return 0, true
    }
    return -1, true
  }
  if j == nil {
    return 1, true
  }
  ci, iok := p.CoerceData(i)
  cj, jok := p.CoerceData(j)
  if iok && !jok {
    return 1, true
  }
  if !iok && jok {
    return -1, true
  }
  // hopefully this doesn't happen as Compare() would continuously return 0, false
  if !iok && !jok {
    return 0, false
  }
  if ci == cj {
    return 0, true
  }
  if ci == nil {
    if cj == nil {
      return 0, true
    }
    return -1, true
  }
  if cj == nil {
    return 1, true
  }
  switch p.ThriftTypeId() {
  case iSTOP, iVOID:
    // hopefully this doesn't happen as Compare() would continuously return 0, false
    return 0, false
  case iBOOL:
    vi, iok := ci.(bool)
    vj, jok := cj.(bool)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi == false {
      return -1, true
    }
    return 1, true
  case iBYTE:
    vi, iok := ci.(byte)
    vj, jok := cj.(byte)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi < vj {
      return -1, true
    }
    return 1, true
  case iDOUBLE:
    vi, iok := ci.(float64)
    vj, jok := cj.(float64)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi < vj {
      return -1, true
    }
    return 1, true
  case iI16:
    vi, iok := ci.(int16)
    vj, jok := cj.(int16)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi < vj {
      return -1, true
    }
    return 1, true
  case iI32:
    vi, iok := ci.(int32)
    vj, jok := cj.(int32)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi < vj {
      return -1, true
    }
    return 1, true
  case iI64:
    vi, iok := ci.(int64)
    vj, jok := cj.(int64)
    if !iok || !jok {
      return 0, false
    }
    if vi == vj {
      return 0, true
    }
    if vi < vj {
      return -1, true
    }
    return 1, true
  case iSTRING, iUTF8, iUTF16:
    if !p.IsBinary() {
      vi, iok := ci.(string)
      vj, jok := cj.(string)
      if !iok || !jok {
        return 0, false
      }
      if vi == vj {
        return 0, true
      }
      if vi < vj {
        return -1, true
      }
      return 1, true
    } else {
      vi, iok := ci.([]byte)
      vj, jok := cj.([]byte)
      if !iok || !jok {
        return 0, false
      }
      return bytes.Compare(vi, vj), true
    }
  case iSTRUCT:
    si, iok := ci.(TStruct)
    sj, jok := cj.(TStruct)
    if !iok || !jok {
      return 0, false
    }
    if cmp := CompareString(si.ThriftName(), sj.ThriftName()); cmp != 0 {
      return cmp, true
    }
    if cmp, ok := si.TStructFields().CompareTo(sj.TStructFields()); !ok || cmp != 0 {
      return cmp, ok
    }
    for field := range si.TStructFields().Iter() {
      a := si.AttributeFromFieldId(field.Id())
      b := sj.AttributeFromFieldId(field.Id())
      if cmp, ok := field.TypeId().Compare(a, b); !ok || cmp != 0 {
        return cmp, ok
      }
    }
    return 0, true
  case iMAP:
    mi, iok := ci.(TMap)
    mj, jok := cj.(TMap)
    if !iok || !jok {
      return 0, false
    }
    ei := mi.KeyType()
    if ej := mj.KeyType(); ei != ej {
      return CompareInt(int(ei.ThriftTypeId()), int(ej.ThriftTypeId())), true
    }
    if size := mi.Len(); size != mj.Len() {
      return CompareInt(size, mj.Len()), true
    }
    if c, cok := ei.Compare(mi.Keys(), mj.Keys()); c != 0 || !cok {
      return c, cok
    }
    return ei.Compare(mi.Values(), mj.Values())
  case iLIST:
    li, iok := ci.(TList)
    lj, jok := cj.(TList)
    if !iok || !jok {
      return 0, false
    }
    ei := li.ElemType()
    ej := lj.ElemType()
    if ei != ej {
      return CompareInt(int(ei.ThriftTypeId()), int(ej.ThriftTypeId())), true
    }
    size := li.Len()
    if size != lj.Len() {
      return CompareInt(size, lj.Len()), true
    }
    for k := 0; k < size; k++ {
      vi := li.At(k)
      vj := lj.At(k)
      c, cok := ei.Compare(vi, vj)
      if c != 0 || !cok {
        return c, cok
      }
    }
    return 0, true
  case iSET:
    li, iok := ci.(TSet)
    lj, jok := cj.(TSet)
    if !iok || !jok {
      return 0, false
    }
    ei := li.ElemType()
    ej := lj.ElemType()
    if ei != ej {
      return CompareInt(int(ei.ThriftTypeId()), int(ej.ThriftTypeId())), true
    }
    size := li.Len()
    if size != lj.Len() {
      return CompareInt(size, lj.Len()), true
    }
    return ei.Compare(li.Values(), lj.Values())
  default:
    panic("Invalid thrift type to coerce")
  }
  return 0, false
}

func (p *tType) CompareValueArrays(li, lj []interface{}) (int, bool) {
  size := len(li)
  if cmp := CompareInt(size, len(lj)); cmp != 0 {
    return cmp, true
  }
  for i := 0; i < size; i++ {
    vi := li[i]
    vj := lj[i]
    c, cok := p.Compare(vi, vj)
    if c != 0 || !cok {
      return c, cok
    }
  }
  return 0, true
}

func (p *tType) Equals(other interface{}) bool {
  return p == p.Coerce(other)
}

type Stringer interface {
  String() string
}

type Enumer interface {
  String() string
  Value() int
  IsEnum() bool
}

func TypeFromValue(data interface{}) TType {
  switch i := data.(type) {
  default:
    return STOP
  case nil:
    return VOID
  case bool:
    return BOOL
  case float32, float64:
    return DOUBLE
  case int, int32:
    return I32
  case byte:
    return BYTE
  case int8:
    return I08
  case int16:
    return I16
  case int64:
    return I64
  case string:
    return STRING
  case []byte:
    return BINARY
  case TStruct:
    return STRUCT
  case TMap:
    return MAP
  case TSet:
    return SET
  case []interface{}, *list.List, *vector.Vector, TList:
    return LIST
  }
  return STOP
}

func (p *tType) CoerceData(data interface{}) (interface{}, bool) {
  if data == nil {
    switch p.ThriftTypeId() {
    case iSTOP:
      return nil, true
    case iVOID:
      return nil, true
    case iBOOL:
      return false, true
    case iBYTE:
      return byte(0), true
    case iDOUBLE:
      return float64(0), true
    case iI16:
      return int16(0), true
    case iI32:
      return int32(0), true
    case iI64:
      return int64(0), true
    case iSTRING:
      if p.IsBinary() {
          return nil, true
      }
      return "", true
    case iSTRUCT:
      return NewTStructEmpty(""), true
    case iMAP:
      return NewTMapDefault(), true
    case iLIST:
      return NewTListDefault(), true
    case iSET:
      return NewTSetDefault(), true
    case iENUM:
        if p.IsUTF8() {
            return "", true
        }
        return int32(0), true
    case iUTF16:
        return "", true
    default:
      panic("Invalid thrift type to coerce")
    }
  }
  switch p.ThriftTypeId() {
  case iSTOP:
    return nil, true
  case iVOID:
    return nil, true
  case iBOOL:
    switch b := data.(type) {
    default:
      return false, false
    case bool:
      return b, true
    case Numeric:
      return bool(b.Int() != 0), true
    case int:
      return b != 0, true
    case byte:
      return b != 0, true
    case int8:
      return b != 0, true
    case int16:
      return b != 0, true
    case int32:
      return b != 0, true
    case int64:
      return b != 0, true
    case uint:
      return b != 0, true
    case uint16:
      return b != 0, true
    case uint32:
      return b != 0, true
    case uint64:
      return b != 0, true
    case float32:
      return b != 0, true
    case float64:
      return b != 0, true
    case []byte:
      return len(b) > 1 || (len(b) == 1 && b[0] != 0), true
    case Stringer:
      v := b.String()
      if v == "false" || v == "0" || len(v) == 0 {
        return false, true
      }
      return true, true
    case string:
      if b == "false" || b == "0" || len(b) == 0 {
        return false, true
      }
      return true, true
    }
  case iBYTE:
    if b, ok := data.(byte); ok {
      return b, true
    }
    if b, ok := data.(Numeric); ok {
      return b.Byte(), true
    }
    if b, ok := data.(bool); ok {
      if b {
        return byte(1), true
      }
      return byte(0), true
    }
    if b, ok := data.(int); ok {
      return byte(b), true
    }
    if b, ok := data.(int8); ok {
      return byte(b), true
    }
    if b, ok := data.(int16); ok {
      return byte(b), true
    }
    if b, ok := data.(int32); ok {
      return byte(b), true
    }
    if b, ok := data.(int64); ok {
      return byte(b), true
    }
    if b, ok := data.(uint); ok {
      return byte(b), true
    }
    if b, ok := data.(uint8); ok {
      return byte(b), true
    }
    if b, ok := data.(uint16); ok {
      return byte(b), true
    }
    if b, ok := data.(uint32); ok {
      return byte(b), true
    }
    if b, ok := data.(uint64); ok {
      return byte(b), true
    }
    if b, ok := data.(float32); ok {
      return byte(int(b)), true
    }
    if b, ok := data.(float64); ok {
      return byte(int(b)), true
    }
    if b, ok := data.([]byte); ok {
      if len(b) > 0 {
        return b[0], true
      }
      return byte(0), true
    }
    if b, ok := data.(Stringer); ok {
      data = b.String()
    }
    if b, ok := data.(string); ok {
      i1, err := strconv.Atoi(b)
      if err != nil {
        return byte(int(i1)), true
      }
    }
    return byte(0), false
  case iDOUBLE:
    if b, ok := data.(float32); ok {
      return float64(b), true
    }
    if b, ok := data.(float64); ok {
      return b, true
    }
    if b, ok := data.(Numeric); ok {
      return bool(b.Float64() != 0.0), true
    }
    if b, ok := data.(byte); ok {
      return float64(b), true
    }
    if b, ok := data.(bool); ok {
      if b {
        return float64(1.0), true
      }
      return float64(0.0), true
    }
    if b, ok := data.(int); ok {
      return float64(b), true
    }
    if b, ok := data.(int8); ok {
      return float64(b), true
    }
    if b, ok := data.(int16); ok {
      return float64(b), true
    }
    if b, ok := data.(int32); ok {
      return float64(b), true
    }
    if b, ok := data.(int64); ok {
      return float64(b), true
    }
    if b, ok := data.(uint); ok {
      return float64(b), true
    }
    if b, ok := data.(uint8); ok {
      return float64(b), true
    }
    if b, ok := data.(uint16); ok {
      return float64(b), true
    }
    if b, ok := data.(uint32); ok {
      return float64(b), true
    }
    if b, ok := data.(uint64); ok {
      return float64(b), true
    }
    if b, ok := data.([]byte); ok {
      if len(b) > 0 {
        return float64(b[0]), true
      }
      return float64(0), true
    }
    if b, ok := data.(Stringer); ok {
      data = b.String()
    }
    if b, ok := data.(string); ok {
      d1, err := strconv.Atof64(b)
      if err != nil {
        return d1, true
      }
    }
    return float64(0), false
  case iI16:
    if b, ok := data.(int16); ok {
      return b, true
    }
    if b, ok := data.(int); ok {
      return int16(b), true
    }
    if b, ok := data.(Numeric); ok {
      return bool(b.Int16() != 0), true
    }
    if b, ok := data.(byte); ok {
      return int16(b), true
    }
    if b, ok := data.(bool); ok {
      if b {
        return int16(1.0), true
      }
      return int16(0.0), true
    }
    if b, ok := data.(int8); ok {
      return int16(b), true
    }
    if b, ok := data.(int32); ok {
      return int16(b), true
    }
    if b, ok := data.(int64); ok {
      return int16(b), true
    }
    if b, ok := data.(uint); ok {
      return int16(b), true
    }
    if b, ok := data.(uint8); ok {
      return int16(b), true
    }
    if b, ok := data.(uint16); ok {
      return int16(b), true
    }
    if b, ok := data.(uint32); ok {
      return int16(b), true
    }
    if b, ok := data.(uint64); ok {
      return int16(b), true
    }
    if b, ok := data.(float32); ok {
      return int16(int(b)), true
    }
    if b, ok := data.(float64); ok {
      return int16(int(b)), true
    }
    if b, ok := data.(Stringer); ok {
      data = b.String()
    }
    if b, ok := data.(string); ok {
      i1, err := strconv.Atoi(b)
      if err != nil {
        return int16(i1), true
      }
    }
    return int16(0), false
  case iI32:
    if b, ok := data.(int32); ok {
      return b, true
    }
    if b, ok := data.(int); ok {
      return int32(b), true
    }
    if b, ok := data.(Numeric); ok {
      return bool(b.Int32() != 0), true
    }
    if b, ok := data.(byte); ok {
      return int32(b), true
    }
    if b, ok := data.(bool); ok {
      if b {
        return int32(1.0), true
      }
      return int32(0.0), true
    }
    if b, ok := data.(int8); ok {
      return int32(b), true
    }
    if b, ok := data.(int16); ok {
      return int32(b), true
    }
    if b, ok := data.(int64); ok {
      return int32(b), true
    }
    if b, ok := data.(uint); ok {
      return int32(b), true
    }
    if b, ok := data.(uint8); ok {
      return int32(b), true
    }
    if b, ok := data.(uint16); ok {
      return int32(b), true
    }
    if b, ok := data.(uint32); ok {
      return int32(b), true
    }
    if b, ok := data.(uint64); ok {
      return int32(b), true
    }
    if b, ok := data.(float32); ok {
      return int32(int(b)), true
    }
    if b, ok := data.(float64); ok {
      return int32(int(b)), true
    }
    if b, ok := data.(Stringer); ok {
      data = b.String()
    }
    if b, ok := data.(string); ok {
      i1, err := strconv.Atoi(b)
      if err != nil {
        return int32(i1), true
      }
    }
    return int32(0), false
  case iI64:
    if b, ok := data.(int64); ok {
      return b, true
    }
    if b, ok := data.(int32); ok {
      return int64(b), true
    }
    if b, ok := data.(int); ok {
      return int64(b), true
    }
    if b, ok := data.(Numeric); ok {
      return bool(b.Int64() != 0), true
    }
    if b, ok := data.(byte); ok {
      return int64(b), true
    }
    if b, ok := data.(bool); ok {
      if b {
        return int64(1.0), true
      }
      return int64(0.0), true
    }
    if b, ok := data.(int8); ok {
      return int64(b), true
    }
    if b, ok := data.(int16); ok {
      return int64(b), true
    }
    if b, ok := data.(uint); ok {
      return int64(b), true
    }
    if b, ok := data.(uint8); ok {
      return int64(b), true
    }
    if b, ok := data.(uint16); ok {
      return int64(b), true
    }
    if b, ok := data.(uint32); ok {
      return int64(b), true
    }
    if b, ok := data.(uint64); ok {
      return int64(b), true
    }
    if b, ok := data.(float32); ok {
      return int64(b), true
    }
    if b, ok := data.(float64); ok {
      return int64(b), true
    }
    if b, ok := data.(Stringer); ok {
      data = b.String()
    }
    if b, ok := data.(string); ok {
      i1, err := strconv.Atoi64(b)
      if err != nil {
        return i1, true
      }
    }
    return int64(0), false
  case iSTRING, iUTF8, iUTF16:
    if !p.IsBinary() {
      if b, ok := data.([]byte); ok {
        return string(b), true
      }
      if b, ok := data.(Enumer); ok {
        if i1, ok := data.(int); ok {
          return string(i1), true
        }
        return b.String(), true
      }
      if b, ok := data.(Stringer); ok {
        return b.String(), true
      }
      if b, ok := data.(string); ok {
        return b, true
      }
      if b, ok := data.(int); ok {
        return string(b), true
      }
      if b, ok := data.(byte); ok {
        return string(b), true
      }
      if b, ok := data.(bool); ok {
        if b {
          return "true", true
        }
        return "false", true
      }
      if b, ok := data.(int8); ok {
        return string(b), true
      }
      if b, ok := data.(int16); ok {
        return string(b), true
      }
      if b, ok := data.(int32); ok {
        return string(b), true
      }
      if b, ok := data.(int64); ok {
        return string(b), true
      }
      if b, ok := data.(uint); ok {
        return string(b), true
      }
      if b, ok := data.(uint8); ok {
        return string(b), true
      }
      if b, ok := data.(uint16); ok {
        return string(b), true
      }
      if b, ok := data.(uint32); ok {
        return string(b), true
      }
      if b, ok := data.(uint64); ok {
        return string(b), true
      }
      if b, ok := data.(float32); ok {
        return strconv.Ftoa32(b, 'g', -1), true
      }
      if b, ok := data.(float64); ok {
        return strconv.Ftoa64(b, 'g', -1), true
      }
      return "", false
    } else {
      if b, ok := data.([]byte); ok {
        return b, true
      }
      if b, ok := data.(Enumer); ok {
        if i1, ok := data.(int); ok {
          return []byte(string(i1)), true
        }
        return []byte(b.String()), true
      }
      if b, ok := data.(Stringer); ok {
        return []byte(b.String()), true
      }
      if b, ok := data.(string); ok {
        return []byte(b), true
      }
      if b, ok := data.(int); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(byte); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(bool); ok {
        if b {
          return []byte("true"), true
        }
        return []byte("false"), true
      }
      if b, ok := data.(int8); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(int16); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(int32); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(int64); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(uint); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(uint8); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(uint16); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(uint32); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(uint64); ok {
        return []byte(string(b)), true
      }
      if b, ok := data.(float32); ok {
        return []byte(strconv.Ftoa32(b, 'g', -1)), true
      }
      if b, ok := data.(float64); ok {
        return []byte(strconv.Ftoa64(b, 'g', -1)), true
      }
      return "", false
    }
  case iSTRUCT:
    if b, ok := data.(TStruct); ok {
      return b, true
    }
    return NewTStructEmpty(""), true
  case iMAP:
    if b, ok := data.(TMap); ok {
      return b, true
    }
    return NewTMapDefault(), false
  case iLIST:
    if b, ok := data.(TList); ok {
      return b, true
    }
    return NewTListDefault(), false
  case iSET:
    if b, ok := data.(TSet); ok {
      return b, true
    }
    return NewTSetDefault(), false
  default:
    panic("Invalid thrift type to coerce")
  }
  return nil, false
}

type EqualsOtherInterface interface {
  Equals(other interface{}) bool
}

type EqualsMap interface {
  Equals(other TMap) bool
}

type EqualsSet interface {
  Equals(other TSet) bool
}

type EqualsList interface {
  Equals(other TList) bool
}

type EqualsStruct interface {
  Equals(other TStruct) bool
}
