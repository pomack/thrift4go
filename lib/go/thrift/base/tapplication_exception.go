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

package base;

import (
  "os"
  "thrift/protocol"
  "thrift"
)

const (
  UNKNOWN              = 0
  UNKNOWN_METHOD       = 1
  INVALID_MESSAGE_TYPE = 2
  WRONG_METHOD_NAME    = 3
  BAD_SEQUENCE_ID      = 4
  MISSING_RESULT       = 5
  INTERNAL_ERROR       = 6
)


/**
 * Application level exception
 *
 */
type TApplicationException interface {
  thrift.TException;
  TypeId() int32;
  Read(iprot protocol.TProtocol) (TApplicationException, os.Error);
  Write(oprot protocol.TProtocol) (os.Error);
}

type tApplicationException struct {
  thrift.TException;
  type_ int32;
}

func NewTApplicationExceptionDefault() TApplicationException {
  return NewTApplicationException(UNKNOWN, "UNKNOWN");
}

func NewTApplicationExceptionType(type_ int32) TApplicationException {
  return NewTApplicationException(type_, "UNKNOWN");
}

func NewTApplicationException(type_ int32, message string) TApplicationException {
  return &tApplicationException{TException:thrift.NewTException(message), type_:type_};
}

func NewTApplicationExceptionMessage(message string) TApplicationException {
  return NewTApplicationException(UNKNOWN, message);
}

func (p *tApplicationException) TypeId() int32 {
  return p.type_;
}

func (p *tApplicationException) Read(iprot protocol.TProtocol) (error TApplicationException, err os.Error) {
  _, err = iprot.ReadStructBegin();
  if err != nil { return; }
  
  message := "";
  type_ := int32(UNKNOWN);

  for {
    _, ttype, id, err := iprot.ReadFieldBegin();
    if err != nil { return; }
    if ttype == protocol.STOP {
      break;
    }
    switch id {
    case 1:
      if ttype == protocol.STRING {
        message, err = iprot.ReadString();
        if err != nil { return; }
      } else {
        err = protocol.SkipDefaultDepth(iprot, ttype);
        if err != nil { return; }
      }
      break;
    case 2:
      if ttype == protocol.I32 {
        type_, err = iprot.ReadI32();
        if err != nil { return; }
      } else {
        err = protocol.SkipDefaultDepth(iprot, ttype);
        if err != nil { return; }
      }
      break;
    default:
      err = protocol.SkipDefaultDepth(iprot, ttype);
      if err != nil { return; }
      break;
    }
    err = iprot.ReadFieldEnd();
    if err != nil { return; }
  }
  err = iprot.ReadStructEnd();
  error = NewTApplicationException(type_, message);
  return;
}

func (p *tApplicationException) Write(oprot protocol.TProtocol) (err os.Error) {
  err = oprot.WriteStructBegin("TApplicationException");
  if len(p.String()) > 0 {
    err = oprot.WriteFieldBegin("message", protocol.STRING, 1);
    if err != nil { return; }
    err = oprot.WriteString(p.String());
    if err != nil { return; }
    err = oprot.WriteFieldEnd();
    if err != nil { return; }
  }
  err = oprot.WriteFieldBegin("type", protocol.I32, 2);
  if err != nil { return; }
  err = oprot.WriteI32(p.type_);
  if err != nil { return; }
  err = oprot.WriteFieldEnd();
  if err != nil { return; }
  err = oprot.WriteFieldStop();
  if err != nil { return; }
  err = oprot.WriteStructEnd();
  return;
}


