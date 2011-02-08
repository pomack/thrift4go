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

package thrift_test

import (
  . "thrift"
  "testing"
  "http"
  "math"
  "net"
  "io"
  "bytes"
  "fmt"
)

const PROTOCOL_BINARY_DATA_SIZE = 155

var (
  data           string // test data for writing
  protocol_bdata []byte // test data for writing; same as data
  BOOL_VALUES    []bool
  BYTE_VALUES    []byte
  INT16_VALUES   []int16
  INT32_VALUES   []int32
  INT64_VALUES   []int64
  DOUBLE_VALUES  []float64
  STRING_VALUES  []string
)


func init() {
  protocol_bdata = make([]byte, PROTOCOL_BINARY_DATA_SIZE)
  for i := 0; i < PROTOCOL_BINARY_DATA_SIZE; i++ {
    protocol_bdata[i] = byte((i + 'a') % 255)
  }
  data = string(protocol_bdata)
  BOOL_VALUES = []bool{false, true, false, false, true}
  BYTE_VALUES = []byte{117, 0, 1, 32, 127, 128, 255}
  INT16_VALUES = []int16{459, 0, 1, -1, -128, 127, 32767, -32768}
  INT32_VALUES = []int32{459, 0, 1, -1, -128, 127, 32767, 2147483647, -2147483535}
  INT64_VALUES = []int64{459, 0, 1, -1, -128, 127, 32767, 2147483647, -2147483535, 34359738481, -35184372088719, -9223372036854775808, 9223372036854775807}
  DOUBLE_VALUES = []float64{459.3, 0.0, -1.0, 1.0, 0.5, 0.3333, 3.14159, 1.537e-38, 1.673e25, 6.02214179e23, -6.02214179e23, INFINITY.Float64(), NEGATIVE_INFINITY.Float64(), NAN.Float64()}
  STRING_VALUES = []string{"", "a", "st[uf]f", "st,u:ff with spaces", "stuff\twith\nescape\\characters'...\"lots{of}fun</xml>"}
}

type HTTPEchoServer struct{}

func (p *HTTPEchoServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.Copy(w, req.Body)
}

func HttpClientSetupForTest(t *testing.T) (net.Listener, net.Addr) {
  addr, err := FindAvailableTCPServerPort(40000)
  if err != nil {
    t.Fatalf("Unable to find available tcp port addr: %s", err)
  }
  l, err := net.Listen(addr.Network(), addr.String())
  if err != nil {
    t.Fatalf("Unable to setup tcp listener on %s: %s", addr.String(), err)
  }
  go http.Serve(l, &HTTPEchoServer{})
  return l, addr
}


func ReadWriteProtocolTest(t *testing.T, protocolFactory TProtocolFactory) {
  buf := bytes.NewBuffer(make([]byte, 0, 1024))
  l, addr := HttpClientSetupForTest(t)
  transports := []TTransportFactory{
    NewTMemoryBufferTransportFactory(1024),
    NewTIOStreamTransportFactory(buf, buf, true),
    NewTFramedTransportFactory(NewTMemoryBufferTransportFactory(1024)),
    NewTHttpPostClientTransportFactory("http://" + addr.String()),
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteBool(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteByte(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteI16(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteI32(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteI64(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteDouble(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteString(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteBinary(t, p, trans)
    trans.Close()
  }
  for _, tf := range transports {
    trans := tf.GetTransport(nil)
    p := protocolFactory.GetProtocol(trans)
    ReadWriteWork(t, p, trans)
    trans.Close()
  }

  // this test doesn't work in all cases due to EOF issues between
  // buffer read and buffer write when using the same bufio for both
  //for _, tf := range transports {
  //  trans := tf.GetTransport(nil)
  //  p := GetProtocol(trans);
  //  ReadWriteI64(t, p, trans);
  //  ReadWriteDouble(t, p, trans);
  //  ReadWriteBinary(t, p, trans);
  //  ReadWriteByte(t, p, trans);
  //  trans.Close()
  //}

  l.Close()
}

func ReadWriteBool(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(BOOL)
  thelen := len(BOOL_VALUES)
  err := p.WriteListBegin(thetype, thelen)
  if err != nil {
    t.Errorf("%s: %T %T %q Error writing list begin: %q", "ReadWriteBool", p, trans, err, thetype)
  }
  for k, v := range BOOL_VALUES {
    err = p.WriteBool(v)
    if err != nil {
      t.Errorf("%s: %T %T %q Error writing bool in list at index %d: %q", "ReadWriteBool", p, trans, err, k, v)
    }
  }
  p.WriteListEnd()
  if err != nil {
    t.Errorf("%s: %T %T %q Error writing list end: %q", "ReadWriteBool", p, trans, err, BOOL_VALUES)
  }
  p.Flush()
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteBool", p, trans, err, BOOL_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteBool", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteBool", p, trans, thelen, thelen2)
    }
  }
  for k, v := range BOOL_VALUES {
    value, err := p.ReadBool()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading bool at index %d: %q", "ReadWriteBool", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: index %d %q %q %q != %q", "ReadWriteBool", k, p, trans, v, value)
    }
  }
  err = p.ReadListEnd()
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteBool", p, trans, err)
  }
}

func ReadWriteByte(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(BYTE)
  thelen := len(BYTE_VALUES)
  err := p.WriteListBegin(thetype, thelen)
  if err != nil {
    t.Errorf("%s: %T %T %q Error writing list begin: %q", "ReadWriteByte", p, trans, err, thetype)
  }
  for k, v := range BYTE_VALUES {
    err = p.WriteByte(v)
    if err != nil {
      t.Errorf("%s: %T %T %q Error writing byte in list at index %d: %q", "ReadWriteByte", p, trans, err, k, v)
    }
  }
  err = p.WriteListEnd()
  if err != nil {
    t.Errorf("%s: %T %T %q Error writing list end: %q", "ReadWriteByte", p, trans, err, BYTE_VALUES)
  }
  err = p.Flush()
  if err != nil {
    t.Errorf("%s: %T %T %q Error flushing list of bytes: %q", "ReadWriteByte", p, trans, err, BYTE_VALUES)
  }
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteByte", p, trans, err, BYTE_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteByte", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteByte", p, trans, thelen, thelen2)
    }
  }
  for k, v := range BYTE_VALUES {
    value, err := p.ReadByte()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading byte at index %d: %q", "ReadWriteByte", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: %T %T %d != %d", "ReadWriteByte", p, trans, v, value)
    }
  }
  err = p.ReadListEnd()
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteByte", p, trans, err)
  }
}

func ReadWriteI16(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(I16)
  thelen := len(INT16_VALUES)
  p.WriteListBegin(thetype, thelen)
  for _, v := range INT16_VALUES {
    p.WriteI16(v)
  }
  p.WriteListEnd()
  p.Flush()
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteI16", p, trans, err, INT16_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteI16", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteI16", p, trans, thelen, thelen2)
    }
  }
  for k, v := range INT16_VALUES {
    value, err := p.ReadI16()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading int16 at index %d: %q", "ReadWriteI16", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: %T %T %d != %d", "ReadWriteI16", p, trans, v, value)
    }
  }
  err = p.ReadListEnd()
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteI16", p, trans, err)
  }
}

func ReadWriteI32(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(I32)
  thelen := len(INT32_VALUES)
  p.WriteListBegin(thetype, thelen)
  for _, v := range INT32_VALUES {
    p.WriteI32(v)
  }
  p.WriteListEnd()
  p.Flush()
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteI32", p, trans, err, INT32_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteI32", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteI32", p, trans, thelen, thelen2)
    }
  }
  for k, v := range INT32_VALUES {
    value, err := p.ReadI32()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading int32 at index %d: %q", "ReadWriteI32", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: %T %T %d != %d", "ReadWriteI32", p, trans, v, value)
    }
  }
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteI32", p, trans, err)
  }
}

func ReadWriteI64(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(I64)
  thelen := len(INT64_VALUES)
  p.WriteListBegin(thetype, thelen)
  for _, v := range INT64_VALUES {
    p.WriteI64(v)
  }
  p.WriteListEnd()
  p.Flush()
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteI64", p, trans, err, INT64_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteI64", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteI64", p, trans, thelen, thelen2)
    }
  }
  for k, v := range INT64_VALUES {
    value, err := p.ReadI64()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading int64 at index %d: %q", "ReadWriteI64", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: %T %T %q != %q", "ReadWriteI64", p, trans, v, value)
    }
  }
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteI64", p, trans, err)
  }
}

func ReadWriteDouble(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(DOUBLE)
  thelen := len(DOUBLE_VALUES)
  p.WriteListBegin(thetype, thelen)
  for _, v := range DOUBLE_VALUES {
    p.WriteDouble(v)
  }
  p.WriteListEnd()
  p.Flush()
  wrotebuffer := ""
  if memtrans, ok := trans.(*TMemoryBuffer); ok {
    wrotebuffer = memtrans.String()
  }
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q, wrote: %v", "ReadWriteDouble", p, trans, err, DOUBLE_VALUES, wrotebuffer)
  }
  if thetype != thetype2 {
    t.Errorf("%s: %T %T type %s != type %s", "ReadWriteDouble", p, trans, thetype, thetype2)
  }
  if thelen != thelen2 {
    t.Errorf("%s: %T %T len %s != len %s", "ReadWriteDouble", p, trans, thelen, thelen2)
  }
  for k, v := range DOUBLE_VALUES {
    value, err := p.ReadDouble()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading double at index %d: %q", "ReadWriteDouble", p, trans, err, k, v)
    }
    if math.IsNaN(v) {
      if !math.IsNaN(value) {
        t.Errorf("%s: %T %T math.IsNaN(%q) != math.IsNaN(%q)", "ReadWriteDouble", p, trans, v, value)
      }
    } else if v != value {
      t.Errorf("%s: %T %T %v != %q", "ReadWriteDouble", p, trans, v, value)
    }
  }
  err = p.ReadListEnd()
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteDouble", p, trans, err)
  }
}

func ReadWriteString(t *testing.T, p TProtocol, trans TTransport) {
  thetype := TType(STRING)
  thelen := len(STRING_VALUES)
  p.WriteListBegin(thetype, thelen)
  for _, v := range STRING_VALUES {
    p.WriteString(v)
  }
  p.WriteListEnd()
  p.Flush()
  thetype2, thelen2, err := p.ReadListBegin()
  if err != nil {
    t.Errorf("%s: %T %T %q Error reading list: %q", "ReadWriteString", p, trans, err, STRING_VALUES)
  }
  _, ok := p.(*TSimpleJSONProtocol)
  if !ok {
    if thetype != thetype2 {
      t.Errorf("%s: %T %T type %s != type %s", "ReadWriteString", p, trans, thetype, thetype2)
    }
    if thelen != thelen2 {
      t.Errorf("%s: %T %T len %s != len %s", "ReadWriteString", p, trans, thelen, thelen2)
    }
  }
  for k, v := range STRING_VALUES {
    value, err := p.ReadString()
    if err != nil {
      t.Errorf("%s: %T %T %q Error reading string at index %d: %q", "ReadWriteString", p, trans, err, k, v)
    }
    if v != value {
      t.Errorf("%s: %T %T %d != %d", "ReadWriteString", p, trans, v, value)
    }
  }
  if err != nil {
    t.Errorf("%s: %T %T Unable to read list end: %q", "ReadWriteString", p, trans, err)
  }
}


func ReadWriteBinary(t *testing.T, p TProtocol, trans TTransport) {
  v := protocol_bdata
  p.WriteBinary(v)
  p.Flush()
  value, err := p.ReadBinary()
  if err != nil {
    t.Errorf("%s: %T %T Unable to read binary: %s", "ReadWriteBinary", p, trans, err.String())
  }
  if len(v) != len(value) {
    t.Errorf("%s: %T %T len(v) != len(value)... %d != %d", "ReadWriteBinary", p, trans, len(v), len(value))
  } else {
    for i := 0; i < len(v); i++ {
      if v[i] != value[i] {
        t.Errorf("%s: %T %T %s != %s", "ReadWriteBinary", p, trans, v, value)
      }
    }
  }
}


func ReadWriteWork(t *testing.T, p TProtocol, trans TTransport) {
  thetype := "struct"
  orig := NewWork()
  orig.Num1 = 25
  orig.Num2 = 102
  orig.Op = ADD
  orig.Comment = "Add: 25 + 102"
  return
  if e := orig.Write(p); e != nil {
    t.Fatalf("Unable to write %s value %#v due to error: %s", thetype, orig, e.String())
  }
  read := NewWork()
  e := read.Read(p)
  if e != nil {
    t.Fatalf("Unable to read %s due to error: %s", thetype, e.String())
  }
  if !orig.Equals(read) {
    t.Fatalf("Original Write != Read: %#v != %#v ", orig, read)
  }
}


/**
 *You can define enums, which are just 32 bit integers. Values are optional
 *and start at 1 if not supplied, C style again.
 */
type Operation int

const (
  ADD      Operation = 1
  SUBTRACT Operation = 2
  MULTIPLY Operation = 3
  DIVIDE   Operation = 4
)

func (p Operation) String() string {
  switch p {
  case ADD:
    return "ADD"
  case SUBTRACT:
    return "SUBTRACT"
  case MULTIPLY:
    return "MULTIPLY"
  case DIVIDE:
    return "DIVIDE"
  }
  return ""
}

func FromOperationString(s string) Operation {
  switch s {
  case "ADD":
    return ADD
  case "SUBTRACT":
    return SUBTRACT
  case "MULTIPLY":
    return MULTIPLY
  case "DIVIDE":
    return DIVIDE
  }
  return Operation(-10000)
}

func (p Operation) Value() int {
  return int(p)
}

func (p Operation) IsEnum() bool {
  return true
}

/**
 *Thrift lets you do typedefs to get pretty names for your types. Standard
 *C style here.
 */
type MyInteger int32

const INT32CONSTANT = 9853

var MAPCONSTANT TMap
/**
 * Structs are the basic complex data structures. They are comprised of fields
 * which each have an integer identifier, a type, a symbolic name, and an
 * optional default value.
 * 
 * Fields can be declared "optional", which ensures they will not be included
 * in the serialized output if they aren't set.  Note that this requires some
 * manual management in some languages.
 * 
 * Attributes:
 *  - Num1
 *  - Num2
 *  - Op
 *  - Comment
 */
type Work struct {
  TStruct
  _       interface{} "num1"    // nil # 0
  Num1    int32       "num1"    // 1
  Num2    int32       "num2"    // 2
  Op      Operation   "op"      // 3
  Comment string      "comment" // 4
}

func NewWork() *Work {
  output := &Work{
    TStruct: NewTStruct("Work", []TField{
      NewTField("num1", I32, 1),
      NewTField("num2", I32, 2),
      NewTField("op", I32, 3),
      NewTField("comment", STRING, 4),
    }),
  }
  {
    output.Num1 = 0
  }
  return output
}

func (p *Work) Read(iprot TProtocol) (err TProtocolException) {
  _, err = iprot.ReadStructBegin()
  if err != nil {
    return NewTProtocolExceptionReadStruct(p.ThriftName(), err)
  }
  for {
    fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
    if fieldId < 0 {
      fieldId = int16(p.FieldIdFromFieldName(fieldName))
    } else if fieldName == "" {
      fieldName = p.FieldNameFromFieldId(int(fieldId))
    }
    if fieldTypeId == GENERIC {
      fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
    }
    if err != nil {
      return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
    }
    if fieldTypeId == STOP {
      break
    }
    if fieldId == 1 || fieldName == "num1" {
      if fieldTypeId == I32 {
        err = p.ReadField1(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else if fieldTypeId == VOID {
        err = iprot.Skip(fieldTypeId)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else {
        err = p.ReadField1(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      }
    } else if fieldId == 2 || fieldName == "num2" {
      if fieldTypeId == I32 {
        err = p.ReadField2(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else if fieldTypeId == VOID {
        err = iprot.Skip(fieldTypeId)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else {
        err = p.ReadField2(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      }
    } else if fieldId == 3 || fieldName == "op" {
      if fieldTypeId == I32 {
        err = p.ReadField3(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else if fieldTypeId == VOID {
        err = iprot.Skip(fieldTypeId)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else {
        err = p.ReadField3(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      }
    } else if fieldId == 4 || fieldName == "comment" {
      if fieldTypeId == STRING {
        err = p.ReadField4(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else if fieldTypeId == VOID {
        err = iprot.Skip(fieldTypeId)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      } else {
        err = p.ReadField4(iprot)
        if err != nil {
          return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
        }
      }
    } else {
      err = iprot.Skip(fieldTypeId)
      if err != nil {
        return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
      }
    }
    err = iprot.ReadFieldEnd()
    if err != nil {
      return NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
    }
  }
  err = iprot.ReadStructEnd()
  if err != nil {
    return NewTProtocolExceptionReadStruct(p.ThriftName(), err)
  }
  return err
}

func (p *Work) ReadField1(iprot TProtocol) (err TProtocolException) {
  v4, err5 := iprot.ReadI32()
  if err5 != nil {
    return NewTProtocolExceptionReadField(1, "num1", p.ThriftName(), err5)
  }
  p.Num1 = v4
  return err
}

func (p *Work) ReadFieldNum1(iprot TProtocol) TProtocolException {
  return p.ReadField1(iprot)
}

func (p *Work) ReadField2(iprot TProtocol) (err TProtocolException) {
  v6, err7 := iprot.ReadI32()
  if err7 != nil {
    return NewTProtocolExceptionReadField(2, "num2", p.ThriftName(), err7)
  }
  p.Num2 = v6
  return err
}

func (p *Work) ReadFieldNum2(iprot TProtocol) TProtocolException {
  return p.ReadField2(iprot)
}

func (p *Work) ReadField3(iprot TProtocol) (err TProtocolException) {
  v8, err9 := iprot.ReadI32()
  if err9 != nil {
    return NewTProtocolExceptionReadField(3, "op", p.ThriftName(), err9)
  }
  p.Op = Operation(v8)
  return err
}

func (p *Work) ReadFieldOp(iprot TProtocol) TProtocolException {
  return p.ReadField3(iprot)
}

func (p *Work) ReadField4(iprot TProtocol) (err TProtocolException) {
  v10, err11 := iprot.ReadString()
  if err11 != nil {
    return NewTProtocolExceptionReadField(4, "comment", p.ThriftName(), err11)
  }
  p.Comment = v10
  return err
}

func (p *Work) ReadFieldComment(iprot TProtocol) TProtocolException {
  return p.ReadField4(iprot)
}

func (p *Work) Write(oprot TProtocol) (err TProtocolException) {
  err = oprot.WriteStructBegin("Work")
  if err != nil {
    return NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
  }
  err = p.WriteField1(oprot)
  if err != nil {
    return err
  }
  err = p.WriteField2(oprot)
  if err != nil {
    return err
  }
  err = p.WriteField3(oprot)
  if err != nil {
    return err
  }
  err = p.WriteField4(oprot)
  if err != nil {
    return err
  }
  err = oprot.WriteFieldStop()
  if err != nil {
    return NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
  }
  err = oprot.WriteStructEnd()
  if err != nil {
    return NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
  }
  return err
}

func (p *Work) WriteField1(oprot TProtocol) (err TProtocolException) {
  err = oprot.WriteFieldBegin("num1", I32, 1)
  if err != nil {
    return NewTProtocolExceptionWriteField(1, "num1", p.ThriftName(), err)
  }
  err = oprot.WriteI32(int32(p.Num1))
  if err != nil {
    return NewTProtocolExceptionWriteField(1, "num1", p.ThriftName(), err)
  }
  err = oprot.WriteFieldEnd()
  if err != nil {
    return NewTProtocolExceptionWriteField(1, "num1", p.ThriftName(), err)
  }
  return err
}

func (p *Work) WriteFieldNum1(oprot TProtocol) TProtocolException {
  return p.WriteField1(oprot)
}

func (p *Work) WriteField2(oprot TProtocol) (err TProtocolException) {
  err = oprot.WriteFieldBegin("num2", I32, 2)
  if err != nil {
    return NewTProtocolExceptionWriteField(2, "num2", p.ThriftName(), err)
  }
  err = oprot.WriteI32(int32(p.Num2))
  if err != nil {
    return NewTProtocolExceptionWriteField(2, "num2", p.ThriftName(), err)
  }
  err = oprot.WriteFieldEnd()
  if err != nil {
    return NewTProtocolExceptionWriteField(2, "num2", p.ThriftName(), err)
  }
  return err
}

func (p *Work) WriteFieldNum2(oprot TProtocol) TProtocolException {
  return p.WriteField2(oprot)
}

func (p *Work) WriteField3(oprot TProtocol) (err TProtocolException) {
  err = oprot.WriteFieldBegin("op", I32, 3)
  if err != nil {
    return NewTProtocolExceptionWriteField(3, "op", p.ThriftName(), err)
  }
  err = oprot.WriteI32(int32(p.Op))
  if err != nil {
    return NewTProtocolExceptionWriteField(3, "op", p.ThriftName(), err)
  }
  err = oprot.WriteFieldEnd()
  if err != nil {
    return NewTProtocolExceptionWriteField(3, "op", p.ThriftName(), err)
  }
  return err
}

func (p *Work) WriteFieldOp(oprot TProtocol) TProtocolException {
  return p.WriteField3(oprot)
}

func (p *Work) WriteField4(oprot TProtocol) (err TProtocolException) {
  err = oprot.WriteFieldBegin("comment", STRING, 4)
  if err != nil {
    return NewTProtocolExceptionWriteField(4, "comment", p.ThriftName(), err)
  }
  err = oprot.WriteString(string(p.Comment))
  if err != nil {
    return NewTProtocolExceptionWriteField(4, "comment", p.ThriftName(), err)
  }
  err = oprot.WriteFieldEnd()
  if err != nil {
    return NewTProtocolExceptionWriteField(4, "comment", p.ThriftName(), err)
  }
  return err
}

func (p *Work) WriteFieldComment(oprot TProtocol) TProtocolException {
  return p.WriteField4(oprot)
}

func (p *Work) TStructName() string {
  return "Work"
}

func (p *Work) ThriftName() string {
  return "Work"
}

func (p *Work) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Work(%+v)", *p)
}

func (p *Work) CompareTo(other interface{}) (int, bool) {
  if other == nil {
    return 1, true
  }
  data, ok := other.(Work)
  if !ok {
    return 0, false
  }
  if p.Num1 != data.Num1 {
    if p.Num1 < data.Num1 {
      return -1, true
    }
    return 1, true
  }
  if p.Num2 != data.Num2 {
    if p.Num2 < data.Num2 {
      return -1, true
    }
    return 1, true
  }
  if p.Op != data.Op {
    if p.Op < data.Op {
      return -1, true
    }
    return 1, true
  }
  if p.Comment != data.Comment {
    if p.Comment < data.Comment {
      return -1, true
    }
    return 1, true
  }
  return 0, true
}

func (p *Work) AttributeByFieldId(id int) interface{} {
  switch id {
  default:
    return nil
  case 1:
    return p.Num1
  case 2:
    return p.Num2
  case 3:
    return p.Op
  case 4:
    return p.Comment
  }
  return nil
}

func (p *Work) TStructFields() TFieldContainer {
  return NewTFieldContainer([]TField{
    NewTField("num1", I32, 1),
    NewTField("num2", I32, 2),
    NewTField("op", I32, 3),
    NewTField("comment", STRING, 4),
  })
}
