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
	"./gen-go/librarytest"
	"./gen-go/tutorial"
	"bytes"
	. "github.com/apesternikov/thrift4go/lib/go/src/thrift"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"reflect"
	"testing"
)

const PROTOCOL_BINARY_DATA_SIZE = 155

var (
	data           string // test data for writing
	protocol_bdata []byte // test data for writing; same as data
	BOOL_VALUES    []bool
	BYTE_VALUES    []int8
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
	BYTE_VALUES = []int8{117, 0, 1, 32, 127, -1, -127}
	INT16_VALUES = []int16{459, 0, 1, -1, -128, 127, 32767, -32768}
	INT32_VALUES = []int32{459, 0, 1, -1, -128, 127, 32767, 2147483647, -2147483535}
	INT64_VALUES = []int64{459, 0, 1, -1, -128, 127, 32767, 2147483647, -2147483535, 34359738481, -35184372088719, -9223372036854775808, 9223372036854775807}
	DOUBLE_VALUES = []float64{459.3, 0.0, -1.0, 1.0, 0.5, 0.3333, 3.14159, 1.537e-38, 1.673e25, 6.02214179e23, -6.02214179e23, INFINITY.Float64(), NEGATIVE_INFINITY.Float64(), NAN.Float64()}
	STRING_VALUES = []string{"", "a", "st[uf]f", "st,u:ff with spaces", "stuff\twith\nescape\\characters'...\"lots{of}fun</xml>"}
}

type HTTPEchoServer struct{}

func (p *HTTPEchoServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	}
}

func HttpClientSetupForTest(t *testing.T) (net.Listener, net.Addr) {
	addr, err := FindAvailableTCPServerPort(40000)
	if err != nil {
		t.Fatalf("Unable to find available tcp port addr: %s", err)
		return nil, addr
	}
	l, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		t.Fatalf("Unable to setup tcp listener on %s: %s", addr.String(), err)
		return l, addr
	}
	go http.Serve(l, &HTTPEchoServer{})
	return l, addr
}

func ReadWriteProtocolTest(t *testing.T, protocolFactory TProtocolFactory) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	l, addr := HttpClientSetupForTest(t)
	defer l.Close()
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
	for _, tf := range transports {
		trans := tf.GetTransport(nil)
		p := protocolFactory.GetProtocol(trans)
		ReadWriteCalculate(t, p, trans)
		trans.Close()
	}
	for _, tf := range transports {
		trans := tf.GetTransport(nil)
		p := protocolFactory.GetProtocol(trans)
		ReadWriteCompactProtoTestStruct(t, p, trans)
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
		t.Errorf("%s: %T %T Unable to read binary: %s", "ReadWriteBinary", p, trans, err.Error())
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
	orig := tutorial.NewWork()
	orig.Num1 = 25
	orig.Num2 = 102
	orig.Op = tutorial.Operation_ADD
	orig.Comment = "Add: 25 + 102"
	return //TODO: seriously?
	if e := orig.Write(p); e != nil {
		t.Fatalf("Unable to write %s value %#v due to error: %s", thetype, orig, e.Error())
	}
	read := tutorial.NewWork()
	e := read.Read(p)
	if e != nil {
		t.Fatalf("Unable to read %s due to error: %s", thetype, e.Error())
	}
	if !reflect.DeepEqual(orig, read) {
		t.Fatalf("Original Write != Read: %#v != %#v ", orig, read)
	}
}

func ReadWriteCompactProtoTestStruct(t *testing.T, p TProtocol, trans TTransport) {
	return
	orig := librarytest.CompactProtoTestStruct_COMPACTTEST
	read := librarytest.NewCompactProtoTestStruct()
	if e := orig.Write(p); e != nil {
		t.Fatalf("Unable to write CompactProtoTestStruct value %#v due to error: %s", orig, e.Error())
	}
	if e := read.Read(p); e != nil {
		t.Fatalf("Unable to read CompactProtoTestStruct due to error: %s", e.Error())
	}
	if !reflect.DeepEqual(orig, read) {
		t.Fatalf("Original Write != Read: \nOrig: %#v \nRead: %#v ", orig, read)
	}
}

func ReadWriteCalculate(t *testing.T, p TProtocol, trans TTransport) {
	messageName := "calculate"
	logid := int32(12)
	seqId := int32(35)
	w := tutorial.NewWork()
	w.Num1 = 25
	w.Num2 = 102
	w.Op = tutorial.Operation_ADD
	w.Comment = "Add: 25 + 102"

	args31 := tutorial.NewCalculateArgs()
	args31.Logid = logid
	args31.W = w
	p.WriteMessageBegin(messageName, CALL, seqId)
	if err := args31.Write(p); err != nil {
		t.Fatalf("%s: %T %T Unable to write message: %s", messageName, p, trans, err.Error())
	}
	p.WriteMessageEnd()
	p.Transport().Flush()

	name, ttype, seqid, err1 := p.ReadMessageBegin()
	if err1 != nil {
		t.Fatalf("%s: %T %T Unable to read message begin: %s", messageName, p, trans, err1.Error())
	}
	if name != messageName {
		t.Errorf("%s: %T %T Expected message named \"%s\", but was: \"%s\"", messageName, p, trans, messageName, name)
	}
	if ttype != CALL {
		t.Errorf("%s: %T %T Expected message type \"%s\", but was: \"%s\"", messageName, p, trans, CALL, ttype)
	}
	if seqid != seqId {
		t.Errorf("%s: %T %T Expected message type \"%s\", but was: \"%s\"", messageName, p, trans, seqId, seqid)
	}
	calcArgs := tutorial.NewCalculateArgs()
	err2 := calcArgs.Read(p)
	if !reflect.DeepEqual(args31, calcArgs) {
		t.Errorf("%s: %T %T Calculate args not as expected, %+v vs %+v", messageName, p, trans, args31, calcArgs)
	}
	if err2 != nil {
		t.Fatalf("%s: %T %T Unable to read message end: %s", messageName, p, trans, err2.Error())
	}
	err3 := p.ReadMessageEnd()
	if err3 != nil {
		t.Fatalf("%s: %T %T Unable to read message end: %s", messageName, p, trans, err3.Error())
	}
}
