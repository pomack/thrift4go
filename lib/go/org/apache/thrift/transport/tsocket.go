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
    "io";
    "net";
    "os";
    "../base/thrift";
)

type TSocketBase struct {
    TTransportBase;
    addr net.Addr;
    conn net.Conn;
    readTimeout int64;
    writeTimeout int64;
}

func (self *TSocketBase) Close() (err os.Error) {
    if self.conn != nil {
        err = self.conn.Close();
        self.conn = nil;
    }
    return err;
}

type TSocket struct {
    TSocketBase;
}

func NewTSocket(addr net.Addr, readTimeout int64, writeTimeout int64) (sock *TSocket) {
    sock = new(TSocket);
    sock.readTimeout = readTimeout;
    sock.writeTimeout = writeTimeout;
    sock.addr = addr;
    sock.conn = nil;
    return sock;
}

func (self *TSocket) IsOpen() bool {
    return self.conn != nil;
}

func (self *TSocket) Open() (err os.Error) {
    self.conn, err = net.Dial(self.addr.Network(), "", self.addr.String());
    if err == nil && self.readTimeout > 0 {
        err = self.conn.SetReadTimeout(self.readTimeout);
    }
    if err == nil && self.writeTimeout > 0 {
        err = self.conn.SetWriteTimeout(self.writeTimeout);
    }
    if err != nil {
        if self.conn != nil {
            self.conn.Close();
            self.conn = nil;
        }
    }
    return err;
}

func (self *TSocket) Read(b []byte) (n int, err os.Error) {
    if !self.IsOpen() {
        return 0, NOT_OPEN;
    }
    n, err = self.conn.Read(b);
    if err != nil { return n, err; }
    if err == os.EOF { return n, END_OF_FILE; }
    return n, err;
}


func (self *TSocket) Write(b []byte) (n int, err os.Error) {
    if !self.IsOpen() {
        return 0, NOT_OPEN;
    }
    n, err = self.conn.Write(b);
    if err != nil { return n, err; }
    if err == os.EOF { return n, END_OF_FILE; }
    return n, err;
}

func (self *TSocket) Flush() (err os.Error) {
    return nil;
}

func (self *TSocket) SetReadTimeout(msec float64) (err os.Error) {
    nsec := int64(msec/1000.0);
    self.readTimeout = nsec;
    if self.conn != nil {
        return self.conn.SetReadTimeout(nsec);
    }
    return nil;
}

func (self *TSocket) SetWriteTimeout(msec float64) (err os.Error) {
    nsec := int64(msec/1000.0);
    self.writeTimeout = nsec;
    if self.conn != nil {
        return self.conn.SetWriteTimeout(nsec);
    }
    return nil;
}

func (self *TSocket) Addr() (net.Addr) {
    return self.addr;
}

type TProcessor interface {
    Process(fd io.ReadWriter, l net.Listener, done chan<- int);
}

type TServerSocket struct {
    TSocketBase;
}

func NewTServerSocket(addr net.Addr, readTimeout int64, writeTimeout int64) (sock *TServerSocket) {
    sock = new(TServerSocket);
    sock.readTimeout = readTimeout;
    sock.writeTimeout = writeTimeout;
    sock.addr = addr;
    sock.conn = nil;
    return sock;
}

func (self *TServerSocket) Serve(processor TProcessor, done chan int) (os.Error) {
    l, err := net.Listen(self.TSocketBase.addr.Network(), self.TSocketBase.addr.String());
    if err != nil {
        return err;
    }
    for {
        fd, err := l.Accept();
        if err != nil {
            break;
        }
        if self.TSocketBase.readTimeout > 0 {
            fd.SetReadTimeout(self.TSocketBase.readTimeout);
        }
        if self.TSocketBase.writeTimeout > 0 {
            fd.SetWriteTimeout(self.TSocketBase.writeTimeout);
        }
        go processor.Process(fd, l, done);
        v, ok := <-done;
        if ok && v == 1 {
            l.Close();
            break;
        }
    }
    return nil;
}

