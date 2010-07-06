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
 *
 */

package protocol

import (
    "../base/thrift";
    "../transport/ttransport";
)

type TProtocolException struct {
    Message string;
    ErrorCode int32;
}

func (self *TProtocolException) String() string {
    return self.Message;
}

func (self *TProtocolException) Code() int32 {
    return self.ErrorCode;
}

var (
  UNKNOWN *TProtocolException = &TProtocolException{ErrorCode:0, Message:"Unknown"};
  INVALID_DATA *TProtocolException = &TProtocolException{ErrorCode:1, Message:"Invalid Data"};
  NEGATIVE_SIZE *TProtocolException = &TProtocolException{ErrorCode:2, Message:"Negative Size"};
  SIZE_LIMIT *TProtocolException = &TProtocolException{ErrorCode:3, Message:"Size Limit"};
  BAD_VERSION *TProtocolException = &TProtocolException{ErrorCode:4, Message:"Bad Version"};
  NOT_IMPLEMENTED *TProtocolException = &TProtocolException{ErrorCode:5, Message:"Not Implemented"};
)

type TProtocol interface {
    WriteMessageBegin(name string, messageType base.TMessageType, seqid int32) *TProtocolException;
    WriteMessageEnd() *TProtocolException;
    WriteStructBegin(name string) *TProtocolException;
    WriteStructEnd() *TProtocolException;
    WriteFieldBegin(name string, fieldType base.TType, seqid int16) *TProtocolException;
    WriteFieldEnd() *TProtocolException;
    WriteFieldStop() *TProtocolException;
    WriteMapBegin(ktype, vtype base.TType, size int32) *TProtocolException;
    WriteMapEnd() *TProtocolException;
    WriteListBegin(etype base.TType, size int32) *TProtocolException;
    WriteListEnd() *TProtocolException;
    WriteSetBegin(etype base.TType, size int32) *TProtocolException;
    WriteSetEnd() *TProtocolException;
    WriteBool(value bool) *TProtocolException;
    WriteByte(value byte) *TProtocolException;
    WriteI16(value int16) *TProtocolException;
    WriteI32(value int32) *TProtocolException;
    WriteI64(value int64) *TProtocolException;
    WriteDouble(value float64) *TProtocolException;
    WriteString(value string) *TProtocolException;

    ReadMessageBegin() (name string, messageType base.TMessageType, seqid int32, err *TProtocolException);
    ReadMessageEnd() *TProtocolException;
    ReadStructBegin() (name string, err *TProtocolException);
    ReadStructEnd() *TProtocolException;
    ReadFieldBegin() (name string, fieldType base.TType, seqid int16, err *TProtocolException);
    ReadFieldEnd() *TProtocolException;
    ReadMapBegin() (ktype, vtype base.TType, size int32, err *TProtocolException);
    ReadMapEnd() *TProtocolException;
    ReadListBegin() (etype base.TType, size int32, err *TProtocolException);
    ReadListEnd() *TProtocolException;
    ReadSetBegin() (etype base.TType, size int32, err *TProtocolException);
    ReadSetEnd() *TProtocolException;
    ReadBool() (value bool, err *TProtocolException);
    ReadByte() (value byte, err *TProtocolException);
    ReadI16() (value int16, err *TProtocolException);
    ReadI32() (value int32, err *TProtocolException);
    ReadI64() (value int64, err *TProtocolException);
    ReadDouble() (value float64, err *TProtocolException);
    ReadString() (value string, err *TProtocolException);

    Skip(fieldType base.TType) (data interface{}, err *TProtocolException);    
}

/*Base class for Thrift protocol driver.*/
type TProtocolBase struct {
    trans transport.TTransport;
}

type EmptyInterface interface {}

func (self *TProtocolBase) WriteMessageBegin(name string, messageType base.TMessageType, seqid int32) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteMessageEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteStructBegin(name string) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteStructEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteFieldBegin(name string, fieldType base.TType, seqid int16) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteFieldEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteFieldStop() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteMapBegin(ktype, vtype base.TType, size int32) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteMapEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteListBegin(etype base.TType, size int32) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteListEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteSetBegin(etype base.TType, size int32) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteSetEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) WriteBool(value bool) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteByte(value byte) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteI16(value int16) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteI32(value int32) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteI64(value int64) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteDouble(value float64) *TProtocolException { return nil; }
func (self *TProtocolBase) WriteString(value string) *TProtocolException { return nil; }

func (self *TProtocolBase) ReadMessageBegin() (name string, messageType base.TMessageType, seqid int32, err *TProtocolException) {
    return "", 0, 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadMessageEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) ReadStructBegin() (name string, err *TProtocolException) { return "", nil; }
func (self *TProtocolBase) ReadStructEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) ReadFieldBegin() (name string, fieldType base.TType, seqid int16, err *TProtocolException) {
    return "", 0, 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadFieldEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) ReadMapBegin() (ktype, vtype base.TType, size int32, err *TProtocolException) {
    return 0, 0, 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadMapEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) ReadListBegin() (etype base.TType, size int32, err *TProtocolException) {
    return 0, 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadListEnd() *TProtocolException { return nil }
func (self *TProtocolBase) ReadSetBegin() (etype base.TType, size int32, err *TProtocolException) {
    return 0, 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadSetEnd() *TProtocolException { return nil; }
func (self *TProtocolBase) ReadBool() (value bool, err *TProtocolException) { 
    return false, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadByte() (value byte, err *TProtocolException) { 
    return 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadI16() (value int16, err *TProtocolException) { 
    return 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadI32() (value int32, err *TProtocolException) { 
    return 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadI64() (value int64, err *TProtocolException) { 
    return 0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadDouble() (value float64, err *TProtocolException) { 
    return 0.0, NOT_IMPLEMENTED;
}
func (self *TProtocolBase) ReadString() (value string, err *TProtocolException) { 
    return "", NOT_IMPLEMENTED;
}

func (self *TProtocolBase) Skip(fieldType base.TType) (data interface{}, err *TProtocolException) {
    switch(fieldType) {
    case base.STOP: return nil, nil;
    case base.BOOL: return self.ReadBool();
    case base.BYTE: return self.ReadByte();
    case base.I16: return self.ReadI16();
    case base.I32: return self.ReadI32();
    case base.I64: return self.ReadI64();
    case base.DOUBLE: return self.ReadDouble();
    case base.STRING: return self.ReadString();
    case base.STRUCT: {
        _, err = self.ReadStructBegin();
        if err != nil {
            return nil, err;
        }
        for {
            _, fieldType2, _, _ := self.ReadFieldBegin();
            if fieldType2 == base.STOP {
              break;
            }
            self.Skip(fieldType2);
            self.ReadFieldEnd()
        }
        return nil, self.ReadStructEnd();
    }
    case base.MAP: {
        ktype, vtype, size, err := self.ReadMapBegin();
        if err != nil {
            return nil, err;
        }
        for i:=int32(0); i<size; i++ {
            self.Skip(ktype);
            self.Skip(vtype);
        }
        return nil, self.ReadMapEnd();
    }
    case base.SET: {
        etype, size, err := self.ReadSetBegin();
        if err != nil {
            return nil, err;
        }
        for i:=int32(0); i<size; i++ {
            self.Skip(etype);
        }
        return nil, self.ReadSetEnd();
    }
    case base.LIST: {
        etype, size, err := self.ReadListBegin();
        if err != nil {
            return nil, err;
        }
        for i:=int32(0); i<size; i++ {
            self.Skip(etype);
        }
        return nil, self.ReadListEnd();
    }
    }
    return nil, nil;
}

type TProtocolFactory interface {
    GetProtocol(trans transport.TTransport) (TProtocol);
};
