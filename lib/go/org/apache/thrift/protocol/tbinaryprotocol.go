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

package protocol

import (
    "fmt";
    "os";
    "encoding/binary";
    "math";
    "strings";
    "../base/thrift";
    "../transport/ttransport";
)


const VERSION_MASK uint32 = 0xffff0000;
const VERSION_1 uint32 = 0x80010000;
const TYPE_MASK uint32 = 0x000000ff;

type TBinaryProtocol struct {
    TProtocolBase;
    StrictRead bool;
    StrictWrite bool;
}

func (self *TBinaryProtocol) handleWriteError(oserr os.Error) (err *TProtocolException) {
    if oserr == nil {
        return nil;
    } else if oserr == os.EOF {
        err = &TProtocolException{ErrorCode:int32(oserr.(os.Errno)), Message:oserr.String()};
        return err;
    }
    err = &TProtocolException{ErrorCode:int32(oserr.(os.Errno)), Message:oserr.String()};
    return err;
}

func (self *TBinaryProtocol) handleReadError(oserr os.Error) (err *TProtocolException) {
    if oserr == nil {
        return nil;
    } else if oserr == os.EOF {
        err = &TProtocolException{ErrorCode:int32(oserr.(os.Errno)), Message:oserr.String()};
        return err;
    }
    err = &TProtocolException{ErrorCode:int32(oserr.(os.Errno)), Message:oserr.String()};
    return err;
}

func (self *TBinaryProtocol) WriteMessageBegin(name string, messageType base.TMessageType, seqid int32) (err *TProtocolException) {
    if self.StrictWrite {
        err = self.WriteI32(int32(uint32(VERSION_1) | uint32(messageType)));
        if err != nil { err = self.WriteString(name); }
        if err != nil { err = self.WriteI32(seqid); }
    } else {
        err = self.WriteString(name);
        if err != nil { err = self.WriteByte(byte(messageType)); }
        if err != nil { err = self.WriteI32(seqid); }
    }
    return err;
}

func (self *TBinaryProtocol) WriteFieldBegin(name string, fieldType base.TType, seqid int16) (err *TProtocolException) {
    err = self.WriteByte(byte(fieldType));
    if err != nil { self.WriteI16(seqid); }
    return err;
}

func (self *TBinaryProtocol) WriteFieldStop() (err *TProtocolException) {
    return self.WriteByte(byte(base.STOP));
}

func (self *TBinaryProtocol) WriteMapBegin(ktype, vtype base.TType, size int32) (err *TProtocolException) {
    err = self.WriteByte(byte(ktype));
    if err != nil { err = self.WriteByte(byte(vtype)); }
    if err != nil { err = self.WriteI32(size); }
    return err;
}

func (self *TBinaryProtocol) WriteListBegin(etype base.TType, size int32) (err *TProtocolException) {
    err = self.WriteByte(byte(etype));
    if err != nil { err = self.WriteI32(size); }
    return err;
}

func (self *TBinaryProtocol) WriteSetBegin(etype base.TType, size int32) (err *TProtocolException) {
    err = self.WriteByte(byte(etype));
    if err != nil { err = self.WriteI32(size); }
    return err;
}

func (self *TBinaryProtocol) WriteBool(value bool) (err *TProtocolException) {
    if value {
        err = self.WriteByte(1);
    } else {
        err = self.WriteByte(0);
    }
    return err;
}

func (self *TBinaryProtocol) WriteByte(value byte) (err *TProtocolException) {
    _, oserr := self.trans.Write([]byte{value});
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) WriteI16(value int16) (err *TProtocolException) {
    oserr := binary.Write(self.trans, binary.BigEndian, value);
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) WriteI32(value int32) (err *TProtocolException) {
    oserr := binary.Write(self.trans, binary.BigEndian, value);
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) WriteI64(value int64) (err *TProtocolException) {
    oserr := binary.Write(self.trans, binary.BigEndian, value);
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) WriteDouble(value float64) (err *TProtocolException) {
    /* FIXME: This may need to be fixed */
    oserr := binary.Write(self.trans, binary.BigEndian, value);
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) WriteString(value string) (err *TProtocolException) {
    var l int32 = int32(len(value));
    binary.Write(self.trans, binary.BigEndian, l);
    _, oserr := self.trans.Write(strings.Bytes(value));
    return self.handleWriteError(oserr);
}

func (self *TBinaryProtocol) ReadMessageBegin() (name string, messageType base.TMessageType, seqid int32, err *TProtocolException) {
    var sz int32;
    sz, err = self.ReadI32();
    if err != nil { return "", 0, 0, err; }
    if sz < 0 {
        version := uint32(sz) & uint32(VERSION_MASK);
        if version != VERSION_1 {
            err = &TProtocolException{ErrorCode:BAD_VERSION.ErrorCode};
            err.Message = fmt.Sprintf("Bad version in readMessageBegin: %d", sz);
            return "", 0, 0, err;
        }
        messageType = base.TMessageType(byte(uint32(sz) & uint32(TYPE_MASK)));
        name, err = self.ReadString();
        if err != nil { return "", 0, 0, err; }
        seqid, err = self.ReadI32();
        if err != nil { return "", 0, 0, err; }
    } else {
        if self.StrictRead {
            err = &TProtocolException{ErrorCode:BAD_VERSION.ErrorCode};
            err.Message = "No protocol version header";
            return "", 0, 0, err;
        }
        b, oserr := self.trans.ReadAll(int(sz));
        if oserr != nil { return "", 0, 0, self.handleReadError(oserr); }
        name = string(b);
        var b1 byte;
        b1, err = self.ReadByte();
        messageType = base.TMessageType(b1);
        if err != nil { return "", 0, 0, err; }
        seqid, err = self.ReadI32();
        if err != nil { return "", 0, 0, err; }
    }
    return name, messageType, seqid, err;
}

func (self *TBinaryProtocol) ReadMessageEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadStructBegin() (name string, err *TProtocolException) {
    return "", nil;
}

func (self *TBinaryProtocol) ReadStructEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadFieldBegin() (name string, fieldType base.TType, seqid int16, err *TProtocolException) {
    var b1 byte;
    b1, err = self.ReadByte();
    fieldType = base.TType(b1);
    if err != nil { return "", 0, 0, err; }
    if fieldType == base.STOP { return "", fieldType, 0, nil; }
    seqid, err = self.ReadI16();
    if err != nil { return "", 0, 0, err; }
    return "", fieldType, seqid, err;
}

func (self *TBinaryProtocol) ReadFieldEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadMapBegin() (ktype, vtype base.TType, size int32, err *TProtocolException) {
    var b1 byte;
    b1, err = self.ReadByte();
    ktype = base.TType(b1);
    if err != nil { return ktype, vtype, size, err; }
    b1, err = self.ReadByte();
    vtype = base.TType(b1);
    if err != nil { return ktype, vtype, size, err; }
    size, err = self.ReadI32();
    return ktype, vtype, size, err;
}

func (self *TBinaryProtocol) ReadMapEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadListBegin() (etype base.TType, size int32, err *TProtocolException) {
    var b1 byte;
    b1, err = self.ReadByte();
    etype = base.TType(b1);
    if err != nil { return etype, size, err; }
    size, err = self.ReadI32();
    return etype, size, err;
}

func (self *TBinaryProtocol) ReadListEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadSetBegin() (etype base.TType, size int32, err *TProtocolException) {
    var b1 byte;
    b1, err = self.ReadByte();
    etype = base.TType(b1);
    if err != nil { return etype, size, err; }
    size, err = self.ReadI32();
    return etype, size, err;
}

func (self *TBinaryProtocol) ReadSetEnd() (err *TProtocolException) {
    return nil;
}

func (self *TBinaryProtocol) ReadBool() (value bool, err *TProtocolException) {
    var b byte;
    b, err = self.ReadByte();
    if b == 0 { value = false; } else { value = true; }
    return value, err;
}

func (self *TBinaryProtocol) ReadByte() (value byte, err *TProtocolException) {
    b, oserr := self.trans.ReadAll(1);
    if len(b) > 0 { value = b[0]; }
    return value, self.handleReadError(oserr);
}

func (self *TBinaryProtocol) ReadI16() (value int16, err *TProtocolException) {
    b, oserr := self.trans.ReadAll(2);
    if len(b) > 2 { value = int16(binary.BigEndian.Uint16(b)); }
    return value, self.handleReadError(oserr);
}

func (self *TBinaryProtocol) ReadI32() (value int32, err *TProtocolException) {
    b, oserr := self.trans.ReadAll(4);
    if len(b) > 3 { value = int32(binary.BigEndian.Uint32(b)); }
    return value, self.handleReadError(oserr);
}

func (self *TBinaryProtocol) ReadI64() (value int64, err *TProtocolException) {
    b, oserr := self.trans.ReadAll(8);
    if len(b) > 7 { value = int64(binary.BigEndian.Uint64(b)); }
    return value, self.handleReadError(oserr);
}

func (self *TBinaryProtocol) ReadDouble() (value float64, err *TProtocolException) {
    b, oserr := self.trans.ReadAll(8);
    if len(b) > 7 { value = math.Float64frombits(binary.BigEndian.Uint64(b)); }
    return value, self.handleReadError(oserr);
}

func (self *TBinaryProtocol) ReadString() (value string, err *TProtocolException) {
    var sz int32;
    sz, err = self.ReadI32();
    if err != nil { return value, err; }
    b, oserr := self.trans.ReadAll(int(sz));
    value = string(b);
    return value, self.handleReadError(oserr);
}

type TBinaryProtocolFactory struct {
    StrictRead bool;
    StrictWrite bool;
}

func (self *TBinaryProtocolFactory) GetProtocol(trans transport.TTransport) (TProtocol) {
    value := &TBinaryProtocol{StrictRead:self.StrictRead, StrictWrite:self.StrictWrite};
    value.trans = trans;
    return value;
}

