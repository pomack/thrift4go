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

package server;

import (
  "thrift/base";
  "thrift/protocol";
  "thrift/transport";
  "os";
)


/**
 * Simple singlethreaded server for testing.
 *
 */
type TSimpleServer struct {
  stopped bool;
  
  processorFactory base.TProcessorFactory;
  serverTransport transport.TServerTransport;
  inputTransportFactory transport.TTransportFactory;
  outputTransportFactory transport.TTransportFactory;
  inputProtocolFactory protocol.TProtocolFactory;
  outputProtocolFactory protocol.TProtocolFactory;
  
  seqId int32;
}

func NewTSimpleServer2(processor base.TProcessor, 
                       serverTransport transport.TServerTransport,
                      ) *TSimpleServer {
  return NewTSimpleServerFactory2(base.NewTProcessorFactory(processor), serverTransport);
}

func NewTSimpleServer4(processor base.TProcessor, 
                       serverTransport transport.TServerTransport, 
                       transportFactory transport.TTransportFactory, 
                       protocolFactory protocol.TProtocolFactory,
                      ) *TSimpleServer {
  return NewTSimpleServerFactory4(base.NewTProcessorFactory(processor), 
                                  serverTransport, 
                                  transportFactory, 
                                  protocolFactory,
                                 );
}

func NewTSimpleServer6(processor base.TProcessor, 
                       serverTransport transport.TServerTransport,
                       inputTransportFactory transport.TTransportFactory,
                       outputTransportFactory transport.TTransportFactory,
                       inputProtocolFactory protocol.TProtocolFactory,
                       outputProtocolFactory protocol.TProtocolFactory,
                      ) *TSimpleServer {
  return NewTSimpleServerFactory6(base.NewTProcessorFactory(processor),
                                  serverTransport,
                                  inputTransportFactory,
                                  outputTransportFactory,
                                  inputProtocolFactory,
                                  outputProtocolFactory,
                                 );
}

func NewTSimpleServerFactory2(processorFactory base.TProcessorFactory,
                              serverTransport transport.TServerTransport,
                              ) *TSimpleServer {
  return NewTSimpleServerFactory6(processorFactory,
                                  serverTransport,
                                  transport.NewTTransportFactory(),
                                  transport.NewTTransportFactory(),
                                  protocol.NewTBinaryProtocolFactoryDefault(),
                                  protocol.NewTBinaryProtocolFactoryDefault(),
                                 );
}

func NewTSimpleServerFactory4(processorFactory base.TProcessorFactory, 
                              serverTransport transport.TServerTransport, 
                              transportFactory transport.TTransportFactory, 
                              protocolFactory protocol.TProtocolFactory,
                             ) *TSimpleServer {
  return NewTSimpleServerFactory6(processorFactory,
                                  serverTransport,
                                  transportFactory,
                                  transportFactory,
                                  protocolFactory,
                                  protocolFactory,
                                 );
}

func NewTSimpleServerFactory6(processorFactory base.TProcessorFactory, 
                              serverTransport transport.TServerTransport,
                              inputTransportFactory transport.TTransportFactory,
                              outputTransportFactory transport.TTransportFactory,
                              inputProtocolFactory protocol.TProtocolFactory,
                              outputProtocolFactory protocol.TProtocolFactory,
                             ) *TSimpleServer {
  return &TSimpleServer{processorFactory:processorFactory,
                        serverTransport:serverTransport,
                        inputTransportFactory:inputTransportFactory,
                        outputTransportFactory:outputTransportFactory,
                        inputProtocolFactory:inputProtocolFactory,
                        outputProtocolFactory:outputProtocolFactory,
                       };
}

func (p *TSimpleServer) ProcessorFactory() base.TProcessorFactory {
  return p.processorFactory;
}

func (p *TSimpleServer) ServerTransport() transport.TServerTransport {
  return p.serverTransport;
}

func (p *TSimpleServer) InputTransportFactory() transport.TTransportFactory {
  return p.inputTransportFactory;
}

func (p *TSimpleServer) OutputTransportFactory() transport.TTransportFactory {
  return p.outputTransportFactory;
}

func (p *TSimpleServer) InputProtocolFactory() protocol.TProtocolFactory {
  return p.inputProtocolFactory;
}

func (p *TSimpleServer) OutputProtocolFactory() protocol.TProtocolFactory {
  return p.outputProtocolFactory;
}

func (p *TSimpleServer) Serve() os.Error {
  p.stopped = false;
  err := p.serverTransport.Listen();
  if err != nil { return err; }
  for !p.stopped {
    client, err := p.serverTransport.Accept();
    if err != nil { return err; }
    if client != nil {
      p.processRequest(client);
    }
  }
  return nil;
}

func (p *TSimpleServer) Stop() os.Error {
  p.stopped = true;
  p.serverTransport.Interrupt();
  return nil;
}

func (p *TSimpleServer) processRequest(client transport.TTransport) {
  processor := p.processorFactory.GetProcessor(client);
  inputTransport := p.inputTransportFactory.GetTransport(client);
  outputTransport := p.outputTransportFactory.GetTransport(client);
  inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport);
  outputProtocol := p.outputProtocolFactory.GetProtocol(outputTransport);
  for {
    ok, e := processor.Process(inputProtocol, outputProtocol, p.seqId);
    p.seqId++;
    if e != nil {
      if !p.stopped {
        // TODO(pomack) log error
        break;
      }
    }
    if !ok {
      break;
    }
  }
  if inputTransport != nil {
    e2 := inputTransport.Close();
    if e2 != nil {
      // TODO(pomack) log error
      return;
    }
  }
  if outputTransport != nil {
    e2 := outputTransport.Close();
    if e2 != nil {
      // TODO(pomack) log error
      return;
    }
  }
}
