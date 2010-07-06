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


package base

type TType byte;
const (
    STOP   = 0;
    VOID   = 1;
    BOOL   = 2;
    BYTE   = 3;
    I08    = 3;
    DOUBLE = 4;
    I16    = 6;
    I32    = 8;
    I64    = 10;
    STRING = 11;
    UTF7   = 11;
    STRUCT = 12;
    MAP    = 13;
    SET    = 14;
    LIST   = 15;
    ENUM   = 16;
    UTF8   = 16;
    UTF16  = 17;
)

type TMessageType int32;
const (
    CALL TMessageType  = 1;
    REPLY TMessageType = 2;
    EXCEPTION TMessageType = 3;
    ONEWAY TMessageType = 4;
)

type TProcessor interface {
    Process(iprot, oprot TProtocol);
}

type TException interface {
    String() string;
};

type TApplicationException interface {
    String() string;
    Code() int32;
}

type TProtocol interface {
    WriteMessageBegin(name string, messageType TMessageType, seqid byte) TException;
    WriteMessageEnd() TException;
    WriteStructBegin(name string) TException;
    WriteStructEnd() TException;
    WriteFieldBegin(name string, fieldType TType, seqid int) TException;
    WriteFieldEnd() TException;
    WriteFieldStop() TException;
    WriteMapBegin(ktype, vtype TType, size int) TException;
    WriteMapEnd() TException;
    WriteListBegin(etype TType, size int) TException;
    WriteListEnd() TException;
    WriteSetBegin(etype TType, size int) TException;
    WriteSetEnd() TException;
    WriteBool(value bool) TException;
    WriteByte(value byte) TException;
    WriteI16(value int16) TException;
    WriteI32(value int32) TException;
    WriteI64(value int64) TException;
    WriteDouble(value float64) TException;
    WriteString(value string) TException;

    ReadMessageBegin() (name string, messageType TMessageType, seqid byte, err TException);
    ReadMessageEnd() TException;
    ReadStructBegin() (name string, err TException);
    ReadStructEnd() TException;
    ReadFieldBegin() (name string, fieldType TType, seqid byte, err TException);
    ReadFieldEnd() TException;
    ReadMapBegin() (ktype, vtype TType, size int, err TException);
    ReadMapEnd() TException;
    ReadListBegin() (etype TType, size int, err TException);
    ReadListEnd() TException;
    ReadSetBegin() (etype TType, size int, err TException);
    ReadSetEnd() TException;
    ReadBool() (value bool, err TException);
    ReadByte() (value byte, err TException);
    ReadI16() (value int16, err TException);
    ReadI32() (value int32, err TException);
    ReadI64() (value int64, err TException);
    ReadDouble() (value float64, err TException);
    ReadString() (value string, err TException);

    Skip(fieldType TType) (data interface{}, err TException);
}

type TApplicationError struct {
    Message string;
    ErrorCode int32;
}

func (self *TApplicationError) String() string {
    return self.Message;
}

func (self *TApplicationError) Code() int32 {
    return self.ErrorCode;
}

var (
    UNKNOWN *TApplicationError = &TApplicationError{ErrorCode:0, Message:"Unknown"};
    UNKNOWN_METHOD *TApplicationError = &TApplicationError{ErrorCode:1, Message:"Unknown method"};
    INVALID_MESSAGE_TYPE *TApplicationError = &TApplicationError{ErrorCode:2, Message:"Invalid message type"};
    WRONG_METHOD_NAME *TApplicationError = &TApplicationError{ErrorCode:3, Message:"Wrong method name"};
    BAD_SEQUENCE_ID *TApplicationError = &TApplicationError{ErrorCode:4, Message:"Bad Sequence ID"};
    MISSING_RESULT *TApplicationError = &TApplicationError{ErrorCode:5, Message:"Missing Result"};
)

func (self *TApplicationError) Read(iprot *TProtocol) {
  iprot.ReadStructBegin();
  for {
    _, ftype, fid, _ := iprot.ReadFieldBegin();
    if ftype == STOP {
      break;
    }
    if fid == 1 {
      if ftype == STRING {
        self.Message, _ = iprot.ReadString();
      } else {
        iprot.Skip(ftype);
      }
    } else if fid == 2 {
      if ftype == I32 {
        self.ErrorCode, _ = iprot.ReadI32();
      } else {
        iprot.Skip(ftype);
      }
    } else {
      iprot.Skip(ftype);
    }
    iprot.ReadFieldEnd();
  }
  iprot.ReadStructEnd();
}

func (self *TApplicationError) Write(oprot *TProtocol) {
  oprot.WriteStructBegin("TApplicationException");
  if len(self.Message) != 0 {
    oprot.WriteFieldBegin("message", STRING, 1);
    oprot.WriteString(self.Message);
    oprot.WriteFieldEnd();
  }
  if self.ErrorCode != -1 {
    oprot.WriteFieldBegin("type", I32, 2);
    oprot.WriteI32(int32(self.ErrorCode));
    oprot.WriteFieldEnd();
  }
  oprot.WriteFieldStop();
  oprot.WriteStructEnd();
}
