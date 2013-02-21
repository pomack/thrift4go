package thrift

import (
	"log"
)

// Debug TProtocol interface
type TDebugProtocol struct {
	delegate TProtocol
}

func NewTDebugProtocol(delegate TProtocol) *TDebugProtocol {
	return &TDebugProtocol{delegate: delegate}
}

func (_m *TDebugProtocol) Flush() (e TProtocolException) {
	log.Println("->", "Flush")
	e = _m.delegate.Flush()
	log.Println("<-", "Flush", e)
	return
}

func (_m *TDebugProtocol) ReadBinary() (ret []byte, e TProtocolException) {
	log.Println("->", "ReadBinary")
	ret, e = _m.delegate.ReadBinary()
	log.Println("<-", "ReadBinary", ret, e)
	return
}

func (_m *TDebugProtocol) ReadBool() (ret bool, e TProtocolException) {
	log.Println("->", "ReadBool")
	ret, e = _m.delegate.ReadBool()
	log.Println("<-", "ReadBool", ret, e)
	return
}

func (_m *TDebugProtocol) ReadByte() (ret int8, e TProtocolException) {
	log.Println("->", "ReadByte")
	ret, e = _m.delegate.ReadByte()
	log.Println("<-", "ReadByte", ret, e)
	return
}

func (_m *TDebugProtocol) ReadDouble() (ret float64, e TProtocolException) {
	log.Println("->", "ReadDouble")
	ret, e = _m.delegate.ReadDouble()
	log.Println("<-", "ReadDouble", ret, e)
	return
}

func (_m *TDebugProtocol) ReadFieldBegin() (name string, t TType, tag int16, e TProtocolException) {
	log.Println("->", "ReadFieldBegin")
	name, t, tag, e = _m.delegate.ReadFieldBegin()
	log.Println("<-", "ReadFieldBegin", name, t, tag, e)
	return
}

func (_m *TDebugProtocol) ReadFieldEnd() (e TProtocolException) {
	log.Println("->", "ReadFieldEnd")
	e = _m.delegate.ReadFieldEnd()
	log.Println("<-", "ReadFieldEnd", e)
	return
}

func (_m *TDebugProtocol) ReadI16() (ret int16, e TProtocolException) {
	log.Println("->", "ReadI16")
	ret, e = _m.delegate.ReadI16()
	log.Println("<-", "ReadI16", ret, e)
	return
}

func (_m *TDebugProtocol) ReadI32() (ret int32, e TProtocolException) {
	log.Println("->", "ReadI32")
	ret, e = _m.delegate.ReadI32()
	log.Println("<-", "ReadI32", ret, e)
	return
}

func (_m *TDebugProtocol) ReadI64() (ret int64, e TProtocolException) {
	log.Println("->", "ReadI64")
	ret, e = _m.delegate.ReadI64()
	log.Println("<-", "ReadI64", ret, e)
	return
}

func (_m *TDebugProtocol) ReadListBegin() (t TType, sz int, e TProtocolException) {
	log.Println("->", "ReadListBegin")
	t, sz, e = _m.delegate.ReadListBegin()
	log.Println("<-", "ReadListBegin", t, sz, e)
	return
}

func (_m *TDebugProtocol) ReadListEnd() (e TProtocolException) {
	log.Println("->", "ReadListEnd")
	e = _m.delegate.ReadListEnd()
	log.Println("<-", "ReadListEnd", e)
	return
}

func (_m *TDebugProtocol) ReadMapBegin() (t TType, t2 TType, sz int, e TProtocolException) {
	log.Println("->", "ReadMapBegin")
	t, t2, sz, e = _m.delegate.ReadMapBegin()
	log.Println("<-", "ReadMapBegin", t, sz, e)
	return
}

func (_m *TDebugProtocol) ReadMapEnd() (e TProtocolException) {
	log.Println("->", "ReadMapEnd")
	e = _m.delegate.ReadMapEnd()
	log.Println("<-", "ReadMapEnd", e)
	return
}

func (_m *TDebugProtocol) ReadMessageBegin() (s string, t TMessageType, tag int32, e TProtocolException) {
	log.Println("->", "ReadMessageBegin")
	s, t, tag, e = _m.delegate.ReadMessageBegin()
	log.Println("<-", "ReadMessageBegin", s, t, tag, e)
	return
}

func (_m *TDebugProtocol) ReadMessageEnd() (e TProtocolException) {
	log.Println("->", "ReadMessageEnd")
	e = _m.delegate.ReadMessageEnd()
	log.Println("<-", "ReadMessageEnd", e)
	return
}

func (_m *TDebugProtocol) ReadSetBegin() (t TType, sz int, e TProtocolException) {
	log.Println("->", "ReadSetBegin")
	t, sz, e = _m.delegate.ReadSetBegin()
	log.Println("<-", "ReadSetBegin", t, sz, e)
	return
}

func (_m *TDebugProtocol) ReadSetEnd() (e TProtocolException) {
	log.Println("->", "ReadSetEnd")
	e = _m.delegate.ReadSetEnd()
	log.Println("<-", "ReadSetEnd", e)
	return
}

func (_m *TDebugProtocol) ReadString() (ret string, e TProtocolException) {
	log.Println("->", "ReadString")
	ret, e = _m.delegate.ReadString()
	log.Println("<-", "ReadString", ret, e)
	return
}

func (_m *TDebugProtocol) ReadStructBegin() (s string, e TProtocolException) {
	log.Println("->", "ReadStructBegin")
	s, e = _m.delegate.ReadStructBegin()
	log.Println("<-", "ReadStructBegin", s, e)
	return
}

func (_m *TDebugProtocol) ReadStructEnd() (e TProtocolException) {
	log.Println("->", "ReadStructEnd")
	e = _m.delegate.ReadStructEnd()
	log.Println("<-", "ReadStructEnd", e)
	return
}

func (_m *TDebugProtocol) Skip(_param0 TType) (e TProtocolException) {
	log.Println("->", "Skip", _param0)
	e = _m.delegate.Skip(_param0)
	log.Println("<-", "Skip", e)
	return
}

func (_m *TDebugProtocol) Transport() (r TTransport) {
	log.Println("->", "Transport")
	r = _m.delegate.Transport()
	log.Println("<-", "Transport", r)
	return
}

func (_m *TDebugProtocol) WriteBinary(_param0 []byte) (e TProtocolException) {
	log.Println("->", "WriteBinary", _param0)
	e = _m.delegate.WriteBinary(_param0)
	log.Println("<-", "WriteBinary", e)
	return
}

func (_m *TDebugProtocol) WriteBool(_param0 bool) (e TProtocolException) {
	log.Println("->", "WriteBool", _param0)
	e = _m.delegate.WriteBool(_param0)
	log.Println("<-", "WriteBool", e)
	return
}

func (_m *TDebugProtocol) WriteByte(_param0 int8) (e TProtocolException) {
	log.Println("->", "WriteByte", _param0)
	e = _m.delegate.WriteByte(_param0)
	log.Println("<-", "WriteByte", e)
	return
}

func (_m *TDebugProtocol) WriteDouble(_param0 float64) (e TProtocolException) {
	log.Println("->", "WriteDouble", _param0)
	e = _m.delegate.WriteDouble(_param0)
	log.Println("<-", "WriteDouble", e)
	return
}

func (_m *TDebugProtocol) WriteFieldBegin(_param0 string, _param1 TType, _param2 int16) (e TProtocolException) {
	log.Println("->", "WriteFieldBegin", _param0, _param1, _param2)
	e = _m.delegate.WriteFieldBegin(_param0, _param1, _param2)
	log.Println("<-", "WriteFieldBegin", e)
	return
}

func (_m *TDebugProtocol) WriteFieldEnd() (e TProtocolException) {
	log.Println("->", "WriteFieldEnd")
	e = _m.delegate.WriteFieldEnd()
	log.Println("<-", "WriteFieldEnd", e)
	return
}

func (_m *TDebugProtocol) WriteFieldStop() (e TProtocolException) {
	log.Println("->", "WriteFieldStop")
	e = _m.delegate.WriteFieldStop()
	log.Println("<-", "WriteFieldStop", e)
	return
}

func (_m *TDebugProtocol) WriteI16(_param0 int16) (e TProtocolException) {
	log.Println("->", "WriteI16", _param0)
	e = _m.delegate.WriteI16(_param0)
	log.Println("<-", "WriteI16", e)
	return
}

func (_m *TDebugProtocol) WriteI32(_param0 int32) (e TProtocolException) {
	log.Println("->", "WriteI32", _param0)
	e = _m.delegate.WriteI32(_param0)
	log.Println("<-", "WriteI32", e)
	return
}

func (_m *TDebugProtocol) WriteI64(_param0 int64) (e TProtocolException) {
	log.Println("->", "WriteI64", _param0)
	e = _m.delegate.WriteI64(_param0)
	log.Println("<-", "WriteI64", e)
	return
}

func (_m *TDebugProtocol) WriteListBegin(_param0 TType, _param1 int) (e TProtocolException) {
	log.Println("->", "WriteListBegin", _param0, _param1)
	e = _m.delegate.WriteListBegin(_param0, _param1)
	log.Println("<-", "WriteListBegin", e)
	return
}

func (_m *TDebugProtocol) WriteListEnd() (e TProtocolException) {
	log.Println("->", "WriteListEnd")
	e = _m.delegate.WriteListEnd()
	log.Println("<-", "WriteListEnd", e)
	return
}

func (_m *TDebugProtocol) WriteMapBegin(_param0 TType, _param1 TType, _param2 int) (e TProtocolException) {
	log.Println("->", "WriteMapBegin", _param0, _param1, _param2)
	e = _m.delegate.WriteMapBegin(_param0, _param1, _param2)
	log.Println("<-", "WriteMapBegin", e)
	return
}

func (_m *TDebugProtocol) WriteMapEnd() (e TProtocolException) {
	log.Println("->", "WriteMapEnd")
	e = _m.delegate.WriteMapEnd()
	log.Println("<-", "WriteMapEnd", e)
	return
}

func (_m *TDebugProtocol) WriteMessageBegin(_param0 string, _param1 TMessageType, _param2 int32) (e TProtocolException) {
	log.Println("->", "WriteMessageBegin", _param0, _param1, _param2)
	e = _m.delegate.WriteMessageBegin(_param0, _param1, _param2)
	log.Println("<-", "WriteMessageBegin", e)
	return
}

func (_m *TDebugProtocol) WriteMessageEnd() (e TProtocolException) {
	log.Println("->", "WriteMessageEnd")
	e = _m.delegate.WriteMessageEnd()
	log.Println("<-", "WriteMessageEnd", e)
	return
}

func (_m *TDebugProtocol) WriteSetBegin(_param0 TType, _param1 int) (e TProtocolException) {
	log.Println("->", "WriteSetBegin", _param0, _param1)
	e = _m.delegate.WriteSetBegin(_param0, _param1)
	log.Println("<-", "WriteSetBegin", e)
	return
}

func (_m *TDebugProtocol) WriteSetEnd() (e TProtocolException) {
	log.Println("->", "WriteSetEnd")
	e = _m.delegate.WriteSetEnd()
	log.Println("<-", "WriteSetEnd", e)
	return
}

func (_m *TDebugProtocol) WriteString(_param0 string) (e TProtocolException) {
	log.Println("->", "WriteString", _param0)
	e = _m.delegate.WriteString(_param0)
	log.Println("<-", "WriteString", e)
	return
}

func (_m *TDebugProtocol) WriteStructBegin(_param0 string) (e TProtocolException) {
	log.Println("->", "WriteStructBegin", _param0)
	e = _m.delegate.WriteStructBegin(_param0)
	log.Println("<-", "WriteStructBegin", e)
	return
}

func (_m *TDebugProtocol) WriteStructEnd() (e TProtocolException) {
	log.Println("->", "WriteStructEnd")
	e = _m.delegate.WriteStructEnd()
	log.Println("<-", "WriteStructEnd", e)
	return
}
