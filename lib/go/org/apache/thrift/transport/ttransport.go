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

package transport

import (
    "bufio";
    "io";
    "os";
    "encoding/binary";
    "bytes";
)

type TTransportException struct {
    Message string;
    ErrorCode int32;
};

var (
    UNKNOWN *TTransportException = &TTransportException{ErrorCode:0, Message:"Unknown"};
    NOT_OPEN *TTransportException = &TTransportException{ErrorCode:1, Message:"Not Open"};
    ALREADY_OPEN *TTransportException = &TTransportException{ErrorCode:2, Message:"Already Open"};
    TIMED_OUT *TTransportException = &TTransportException{ErrorCode:3, Message:"Timed Out"};
    END_OF_FILE *TTransportException = &TTransportException{ErrorCode:4, Message:"End of File"};
)

func (self *TTransportException) String() string {
    return self.Message;
}

func (self *TTransportException) Code() int32 {
    return self.ErrorCode;
}

type TTransport interface {
    IsOpen() bool;
    Open() (err os.Error);
    Close() (err os.Error);
    Read(buf []byte) (n int, err os.Error);
    ReadAll(size int) (b []byte, err os.Error);
    Write(buf []byte) (n int, err os.Error);
    Flush() (err os.Error);
}

//Base class for Thrift transport layer.
type TTransportBase struct {
}

func (self *TTransportBase) IsOpen() bool {
    return false;
}

func (self *TTransportBase) Open() os.Error {
    return nil;
}

func (self *TTransportBase) Close() os.Error {
    return nil;
}

func (self *TTransportBase) Read(buf []byte) (n int, err os.Error) {
    return 0, nil;
}

func (self *TTransportBase) ReadAll(size int) (buf []byte, err os.Error) {
    buf = make([]byte, size);
    for have := 0; have < size; {
        var n int;
        n, err = self.Read(buf[have:]);
        have += n;
        if err == nil {
            continue;
        }
        if err == os.EOF {
            return buf[0:have], err;
        }
    }
    return buf, nil;
}

func (self *TTransportBase) Write(buf []byte) (n int, err os.Error) {
    return 0, nil;
}

func (self *TTransportBase) Flush() (err os.Error) {
    return nil;
}

type TServerTransport interface {
    Listen();
    Accept();
    Close();
}

/*Base class for Thrift server transports.*/
type TServerTransportBase struct {};

func (self *TServerTransportBase) Listen() {
    
}


func (self *TServerTransportBase) Accept() {
    
}

func (self *TServerTransportBase) Close() {
    
}

type TTransportFactory interface {
    GetTransport(trans *TTransport) (*TTransport);
}

/*Base class for a Transport Factory*/
type TTransportFactoryBase struct {};

func (self *TTransportFactoryBase) GetTransport(trans *TTransport) (*TTransport) {
    return trans;
}

/*Factory transport that builds buffered transports*/
type TBufferedTransportFactory struct {};
func (self *TBufferedTransportFactory) GetTransport(trans *TTransport) (*TBufferedTransport) {
    return &TBufferedTransport{trans, &bytes.Buffer{}, &bytes.Buffer{}};
}

const DEFAULT_BUFFER = 4096;

/*Class that wraps another transport and buffers its I/O.*/
type TBufferedTransport struct {
    trans *TTransport;
    writeBuf *bytes.Buffer;
    readBuf *bytes.Buffer;
}

func (self *TBufferedTransport) IsOpen() bool {
    return self.trans.IsOpen();
}

func (self *TBufferedTransport) Open() os.Error {
    return self.trans.Open();
}

func (self *TBufferedTransport) Close() os.Error {
    return self.trans.Close()
}

func (self *TBufferedTransport) Read(buf []byte) (n int, err os.Error) {
    n, err = self.readBuf.Read(buf);
    if n != 0 {
        return n, err;
    }
    self.readBuf.Reset();
    n2, err2 := self.readBuf.ReadFrom((*self.trans).(io.Reader));
    if err2 != nil { return int(n2), err2; }
    return self.readBuf.Read(buf);
}

func (self *TBufferedTransport) Write(buf []byte) (n int, err os.Error) {
    return self.writeBuf.Write(buf);
}

func (self *TBufferedTransport) Flush() (err os.Error) {
    return nil;
}

type TMemoryBuffer struct {
    TTransportBase;
    buffer *bytes.Buffer;
}

func (self *TMemoryBuffer) IsOpen() bool {
    return self.buffer != nil;
}

func (self *TMemoryBuffer) Open() os.Error {
    self.buffer = &bytes.Buffer{};
    return nil;
}

func (self *TMemoryBuffer) Close() os.Error {
    self.buffer = nil;
    return nil;
}

func (self *TMemoryBuffer) Read(buf []byte) (n int, err os.Error) {
    return self.buffer.Read(buf);
}

func (self *TMemoryBuffer) Write(buf []byte) (n int, err os.Error) {
    return self.buffer.Write(buf);
}

func (self *TMemoryBuffer) Flush() os.Error {
    return nil;
}

func (self *TMemoryBuffer) Bytes() []byte {
    return self.buffer.Bytes();
}

/*Class that wraps another transport and frames its I/O when writing.*/
type TFramedTransport struct {
    trans *TTransport;
    readBuf, writeBuf *bytes.Buffer;
}

/*Factory transport that builds framed transports*/
type TFramedTransportFactory struct {};
func (self *TFramedTransportFactory) GetTransport(trans *TTransport) *TFramedTransport {
    return &TFramedTransport{trans, &bytes.Buffer{}, &bytes.Buffer{}};
}

func (self *TFramedTransport) IsOpen() bool {
    return self.trans.IsOpen();
}

func (self *TFramedTransport) Open() os.Error {
    return self.trans.Open();
}

func (self *TFramedTransport) Close() os.Error {
    return self.trans.Close();
}

func (self *TFramedTransport) Read(buf []byte) (n int, err os.Error) {
    n, err = self.readBuf.Read(buf);
    if n != 0 {
        return n, err;
    }
    self.readFrame(buf);
    return self.readBuf.Read(buf);
}

func (self *TFramedTransport) readFrame(buf []byte) (err os.Error) {
    var inbuf []byte;
    inbuf, err = self.trans.ReadAll(4);
    if len(inbuf) > 3 {
        size := int(binary.BigEndian.Uint32(inbuf));
        if err == nil {
            self.readBuf.Reset();
            inbuf, err = self.trans.ReadAll(size);
            self.readBuf = bytes.NewBuffer(inbuf);
            if err == nil {
                self.readBuf.Read(buf);
            }
        }
    }
    return err;
}

func (self *TFramedTransport) Write(buf []byte) (n int, err os.Error) {
    return self.writeBuf.Write(buf);
}

func (self *TFramedTransport) Flush() (err os.Error) {
    size := int32(self.writeBuf.Len());
    writer := (*self.trans).(io.Writer);
    binary.Write(writer, binary.BigEndian, size);
    _, err = self.writeBuf.WriteTo(writer);
    return err;
}

/*Wraps a file-like object to make it work as a Thrift transport.*/
type TFileObjectTransport struct {
    TTransportBase;
    FileObj *os.File;
    buffer *bufio.ReadWriter;
}

func (self *TFileObjectTransport) IsOpen() bool {
    return true;
}

func (self *TFileObjectTransport) Open() os.Error {
    /* FIXME: figure out why type assertion is failing */
    /*
    reader, _ := (*self.FileObj).(io.Reader);
    writer, _ := (*self.FileObj).(io.Writer);
    self.buffer = bufio.NewReadWriter(bufio.NewReader(reader), bufio.NewWriter(writer));
    */
    return nil;
}

func (self *TFileObjectTransport) Close() os.Error {
    self.buffer = nil;
    return self.FileObj.Close();
}

func (self *TFileObjectTransport) Read(buf []byte) (int, os.Error) {
    return self.buffer.Read(buf);
}

func (self *TFileObjectTransport) Write(buf []byte) (int, os.Error) {
    return self.buffer.Write(buf);
}

func (self *TFileObjectTransport) Flush() os.Error {
    return self.buffer.Flush();
}

