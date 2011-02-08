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
	"bufio"
	"bytes"
	"container/vector"
	"encoding/base64"
	"fmt"
	"io"
	"json"
	"math"
	"os"
	"strconv"
	"strings"
)

type _ParseContext int
const (
  _CONTEXT_IN_TOPLEVEL _ParseContext = 1
  _CONTEXT_IN_LIST_FIRST _ParseContext = 2
  _CONTEXT_IN_LIST _ParseContext = 3
  _CONTEXT_IN_OBJECT_FIRST _ParseContext = 4
  _CONTEXT_IN_OBJECT_NEXT_KEY _ParseContext = 5
  _CONTEXT_IN_OBJECT_NEXT_VALUE _ParseContext = 6
)

func (p _ParseContext) String() string {
  switch p {
  case _CONTEXT_IN_TOPLEVEL:
    return "TOPLEVEL"
  case _CONTEXT_IN_LIST_FIRST:
    return "LIST-FIRST"
  case _CONTEXT_IN_LIST:
    return "LIST"
  case _CONTEXT_IN_OBJECT_FIRST:
    return "OBJECT-FIRST"
  case _CONTEXT_IN_OBJECT_NEXT_KEY:
    return "OBJECT-NEXT-KEY"
  case _CONTEXT_IN_OBJECT_NEXT_VALUE:
    return "OBJECT-NEXT-VALUE"
  }
  return "UNKNOWN-PARSE-CONTEXT"
}

/**
 * JSON protocol implementation for thrift.
 *
 * This protocol produces/consumes a simple output format
 * suitable for parsing by scripting languages.  It should not be
 * confused with the full-featured TJSONProtocol.
 *
 */
type TSimpleJSONProtocol struct {
	//TProtocolBase;
	trans TTransport

	/**
	 * Stack of nested contexts that we may be in.
	 */
	parseContextStack vector.IntVector
	/**
	 * Stack of nested contexts that we may be in.
	 */
	dumpContext vector.IntVector

	/**
	 * Current context that we are in
	 */
	writer      TTransport
	reader      *bufio.Reader
}

/**
 * Constructor
 */
func NewTSimpleJSONProtocol(t TTransport) *TSimpleJSONProtocol {
	v := &TSimpleJSONProtocol{trans: t,
		writer:            t,
		reader:            bufio.NewReader(t),
	}
	v.parseContextStack.Push(int(_CONTEXT_IN_TOPLEVEL))
	v.dumpContext.Push(int(_CONTEXT_IN_TOPLEVEL))
	return v
}

/**
 * Factory
 */
type TSimpleJSONProtocolFactory struct{}

func (p *TSimpleJSONProtocolFactory) GetProtocol(trans TTransport) TProtocol {
	return NewTSimpleJSONProtocol(trans)
}

func NewTSimpleJSONProtocolFactory() *TSimpleJSONProtocolFactory {
	return &TSimpleJSONProtocolFactory{}
}

var (
	JSON_COMMA                   []byte
	JSON_COLON                   []byte
	JSON_LBRACE                  []byte
	JSON_RBRACE                  []byte
	JSON_LBRACKET                []byte
	JSON_RBRACKET                []byte
	JSON_QUOTE                   byte
	JSON_QUOTE_BYTES             []byte
	JSON_NULL                    []byte
	JSON_TRUE                    []byte
	JSON_FALSE                   []byte
	JSON_INFINITY                string
	JSON_NEGATIVE_INFINITY       string
	JSON_NAN                     string
	JSON_INFINITY_BYTES          []byte
	JSON_NEGATIVE_INFINITY_BYTES []byte
	JSON_NAN_BYTES               []byte
	json_nonbase_map_elem_bytes  []byte
)

func init() {
	JSON_COMMA = []byte{','}
	JSON_COLON = []byte{':'}
	JSON_LBRACE = []byte{'{'}
	JSON_RBRACE = []byte{'}'}
	JSON_LBRACKET = []byte{'['}
	JSON_RBRACKET = []byte{']'}
	JSON_QUOTE = '"'
	JSON_QUOTE_BYTES = []byte{'"'}
	JSON_NULL = []byte{'n', 'u', 'l', 'l'}
	JSON_TRUE = []byte{'t', 'r', 'u', 'e'}
	JSON_FALSE = []byte{'f', 'a', 'l', 's', 'e'}
	JSON_INFINITY = "Infinity"
	JSON_NEGATIVE_INFINITY = "-Infinity"
	JSON_NAN = "NaN"
	JSON_INFINITY_BYTES = []byte{'I', 'n', 'f', 'i', 'n', 'i', 't', 'y'}
	JSON_NEGATIVE_INFINITY_BYTES = []byte{'-', 'I', 'n', 'f', 'i', 'n', 'i', 't', 'y'}
	JSON_NAN_BYTES = []byte{'N', 'a', 'N'}
	json_nonbase_map_elem_bytes = []byte{']', ',', '['}
}

func readNonSignificantWhitespace(r *bufio.Reader) os.Error {
	for {
		c, err := r.ReadByte()
		if err == os.EOF  {
			return nil
		} else if e, ok := err.(TTransportException); ok && e.TypeId() == END_OF_FILE {
		  return nil
		}
		if err != nil {
			return err
		}
		switch c {
		case ' ', '\n', '\t':
			continue
		default:
			r.UnreadByte()
			break
		}
		break
	}
	return nil
}

func JsonQuote(s string) string {
	b, _ := json.Marshal(s)
	s1 := string(b)
	return s1
}

func JsonUnquote(s string) (string, bool) {
	s1 := new(string)
	err := json.Unmarshal([]byte(s), s1)
	return *s1, err == nil
}


func (p *TSimpleJSONProtocol) WriteMessageBegin(name string, typeId TMessageType, seqId int32) TProtocolException {
  if e := writeListBegin(p, p.writer); e != nil {
    return e
  }
  if e := p.WriteString(name); e != nil {
    return e
  }
  if e := p.WriteByte(byte(typeId)); e != nil {
    return e
  }
	if e := p.WriteI32(seqId); e != nil {
    return e
  }
	return nil
}

func (p *TSimpleJSONProtocol) WriteMessageEnd() TProtocolException {
  return writeListEnd(p, p.writer)
}

func (p *TSimpleJSONProtocol) WriteStructBegin(name string) TProtocolException {
  if e := writeObjectBegin(p, p.writer); e != nil {
    return e
  }
	return nil
}

func (p *TSimpleJSONProtocol) WriteStructEnd() TProtocolException {
  return writeObjectEnd(p, p.writer)
}

func (p *TSimpleJSONProtocol) WriteFieldBegin(name string, typeId TType, id int16) TProtocolException {
  if e := p.WriteString(name); e != nil {
    return e
  }
  return nil
  /*
	if e := writeListBegin(p, p.writer); e != nil {
    return e
  }
  if e := p.WriteByte(byte(typeId)); e != nil {
    return e
  }
  return p.WriteI16(id)
  */
}

func (p *TSimpleJSONProtocol) WriteFieldEnd() TProtocolException {
  //return writeListEnd(p, p.writer)
  return nil
}

func (p *TSimpleJSONProtocol) WriteFieldStop() TProtocolException { return nil }

func (p *TSimpleJSONProtocol) WriteMapBegin(keyType TType, valueType TType, size int) TProtocolException {
  if e := writeListBegin(p, p.writer); e != nil {
    return e
  }
  if e := p.WriteByte(byte(keyType)); e != nil {
    return e
  }
  if e := p.WriteByte(byte(valueType)); e != nil {
    return e
  }
  return p.WriteI32(int32(size))
}

func (p *TSimpleJSONProtocol) WriteMapEnd() TProtocolException {
	return writeListEnd(p, p.writer)
}

func (p *TSimpleJSONProtocol) WriteListBegin(elemType TType, size int) TProtocolException {
  return writeElemListBegin(p, p.writer, elemType, size)
}

func (p *TSimpleJSONProtocol) WriteListEnd() TProtocolException {
	return writeListEnd(p, p.writer)
}

func (p *TSimpleJSONProtocol) WriteSetBegin(elemType TType, size int) TProtocolException {
	return writeElemListBegin(p, p.writer, elemType, size)
}

func (p *TSimpleJSONProtocol) WriteSetEnd() TProtocolException {
	return writeListEnd(p, p.writer)
}

func (p *TSimpleJSONProtocol) WriteBool(b bool) TProtocolException {
	return writeBool(p, p.writer, b)
}

func (p *TSimpleJSONProtocol) WriteByte(b byte) TProtocolException {
	return p.WriteI32(int32(b))
}

func (p *TSimpleJSONProtocol) WriteI16(v int16) TProtocolException {
	return p.WriteI32(int32(v))
}

func (p *TSimpleJSONProtocol) WriteI32(v int32) TProtocolException {
  return writeI64(p, p.writer, int64(v))
}

func (p *TSimpleJSONProtocol) WriteI64(v int64) TProtocolException {
  return writeI64(p, p.writer, int64(v))
}

func (p *TSimpleJSONProtocol) WriteDouble(v float64) TProtocolException {
  return writeF64(p, p.writer, v)
}

func (p *TSimpleJSONProtocol) WriteString(v string) TProtocolException {
  return writeString(p, p.writer, v)
}

func (p *TSimpleJSONProtocol) WriteBinary(v []byte) TProtocolException {
	// JSON library only takes in a string, 
	// not an arbitrary byte array, to ensure bytes are transmitted
	// efficiently we must convert this into a valid JSON string
	// therefore we use base64 encoding to avoid excessive escaping/quoting
	if e := writePreValue(p, p.writer); e != nil {
    return e
  }
  if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	writer := base64.NewEncoder(base64.StdEncoding, p.writer)
	if _, e := writer.Write(v); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
  if e := writer.Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return writePostValue(p, p.writer)
}

/**
 * Reading methods.
 */

func (p *TSimpleJSONProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err TProtocolException) {
  if isNull, err := readListBegin(p, p.reader); isNull || err != nil {
    return name, typeId, seqId, err
  }
  if name, err = p.ReadString(); err != nil {
    return name, typeId, seqId, err
  }
  bTypeId, err := p.ReadByte()
  typeId = TMessageType(bTypeId)
  if err != nil {
    return name, typeId, seqId, err
  }
  if seqId, err = p.ReadI32(); err != nil {
    return name, typeId, seqId, err
  }
  return name, typeId, seqId, nil
}

func (p *TSimpleJSONProtocol) ReadMessageEnd() TProtocolException {
  return readListEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadStructBegin() (name string, err TProtocolException) {
  _, err = readObjectStart(p, p.reader)
  return "", err
}

func (p *TSimpleJSONProtocol) ReadStructEnd() TProtocolException {
  return readObjectEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadFieldBegin() (string, TType, int16, TProtocolException) {
  if err := readPreValue(p, p.reader); err != nil {
    return "", STOP, 0, err
  }
  b, _ := p.reader.Peek(1)
  if len(b) > 0 {
    switch b[0] {
    case JSON_RBRACE[0]:
      return "", STOP, 0, nil
    case JSON_QUOTE:
      p.reader.ReadByte()
      name, err := ReadStringBody(p.reader)
      if err != nil {
        return name, STOP, 0, err
      }
      return name, GENERIC, -1, readPostValue(p, p.reader)
      /*
      if err = readPostValue(p, p.reader); err != nil {
        return name, STOP, 0, err
      }
      if isNull, err := readListBegin(p, p.reader); isNull || err != nil {
        return name, STOP, 0, err
      }
      bType, err := p.ReadByte()
      thetype := TType(bType)
      if err != nil {
        return name, thetype, 0, err
      }
      id, err := p.ReadI16()
      return name, thetype, id, err
      */
    }
    return "", STOP, 0, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected \"}\" or '\"', but found: '", string(b), "'"))
  }
  return "", STOP, 0, NewTProtocolExceptionFromOsError(os.EOF)
}

func (p *TSimpleJSONProtocol) ReadFieldEnd() TProtocolException {
  return nil
  //return readListEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, e TProtocolException) {
  if isNull, e := readListBegin(p, p.reader); isNull || e != nil {
    return VOID, VOID, 0, e
  }
  
  // read keyType
  bKeyType, e := p.ReadByte()
  keyType = TType(bKeyType)
  if e != nil {
    return keyType, valueType, size, e
  }
  
  // read valueType
  bValueType, e := p.ReadByte()
  valueType = TType(bValueType)
  if e != nil {
    return keyType, valueType, size, e
  }
  
  // read size
  iSize, err := p.ReadI64()
  size = int(iSize)
  return keyType, valueType, size, err
}

func (p *TSimpleJSONProtocol) ReadMapEnd() TProtocolException {
  return readListEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadListBegin() (elemType TType, size int, e TProtocolException) {
  return readElemListBegin(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadListEnd() TProtocolException {
  return readListEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadSetBegin() (elemType TType, size int, e TProtocolException) {
  return readElemListBegin(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadSetEnd() TProtocolException {
  return readListEnd(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadBool() (bool, TProtocolException) {
  var value bool
  if err := readPreValue(p, p.reader); err != nil {
    return value, err
  }
  b, _ := p.reader.Peek(len(JSON_FALSE))
  if len(b) > 0 {
    switch b[0] {
    case JSON_TRUE[0]:
      if string(b[0:len(JSON_TRUE)]) == string(JSON_TRUE) {
        p.reader.Read(b[0:len(JSON_TRUE)])
        value = true
      } else {
        return value, NewTProtocolException(INVALID_DATA, "Expected \"true\" but found: " + string(b))
      }
      break
    case JSON_FALSE[0]:
      if string(b[0:len(JSON_FALSE)]) == string(JSON_FALSE) {
        p.reader.Read(b[0:len(JSON_FALSE)])
        value = false
      } else {
        return value, NewTProtocolException(INVALID_DATA, "Expected \"false\" but found: " + string(b))
      }
      break
    case JSON_NULL[0]:
      if string(b[0:len(JSON_NULL)]) == string(JSON_NULL) {
        p.reader.Read(b[0:len(JSON_NULL)])
        value = false
      } else {
        return value, NewTProtocolException(INVALID_DATA, "Expected \"null\" but found: " + string(b))
      }
    default:
      return value, NewTProtocolException(INVALID_DATA, "Expected \"true\", \"false\", or \"null\" but found: " + string(b))
    }
  }
  return value, readPostValue(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadByte() (byte, TProtocolException) {
	v, err := p.ReadI64()
	return byte(v), err
}

func (p *TSimpleJSONProtocol) ReadI16() (int16, TProtocolException) {
	v, err := p.ReadI64()
	return int16(v), err
}

func (p *TSimpleJSONProtocol) ReadI32() (int32, TProtocolException) {
	v, err := p.ReadI64()
	return int32(v), err
}

func (p *TSimpleJSONProtocol) ReadI64() (int64, TProtocolException) {
  if err := readPreValue(p, p.reader); err != nil {
    return 0, err
  }
  v, err := readNumeric(p.reader)
  value := v.Int64()
  if err != nil {
    return value, err
  }
  return value, readPostValue(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadDouble() (float64, TProtocolException) {
  if err := readPreValue(p, p.reader); err != nil {
    return 0, err
  }
  v, err := readNumeric(p.reader)
  value := v.Float64()
  if err != nil {
    return value, err
  }
  return value, readPostValue(p, p.reader)
}

func (p *TSimpleJSONProtocol) ReadString() (string, TProtocolException) {
  var v string
  if err := readPreValue(p, p.reader); err != nil {
    return v, err
  }
  b, _ := p.reader.Peek(len(JSON_NULL))
  if len(b) > 0 && b[0] == JSON_QUOTE {
    p.reader.ReadByte()
    value, err := ReadStringBody(p.reader)
    v = value
    if err != nil {
      return v, err
    }
  } else if len(b) >= len(JSON_NULL) && string(b[0:len(JSON_NULL)]) == string(JSON_NULL) {
    _, err := p.reader.Read(b[0:len(JSON_NULL)])
    if err != nil {
      return v, NewTProtocolExceptionFromOsError(err)
    }
  } else {
    return v, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected a JSON string, found ", string(b)))
  }
  return v, readPostValue(p, p.reader)
}

func ReadStringBody(reader *bufio.Reader) (string, TProtocolException) {
	line, err := reader.ReadString(JSON_QUOTE)
	if err != nil {
		return "", NewTProtocolExceptionFromOsError(err)
	}
	l := len(line)
	// count number of escapes to see if we need to keep going
	i := 1
	for ; i < l; i++ {
		if line[l-i-1] != '\\' {
			break
		}
	}
	if i&0x01 == 1 {
		v, ok := JsonUnquote(string(JSON_QUOTE) + line)
		if !ok {
			return "", NewTProtocolExceptionFromOsError(err)
		}
		return v, nil
	}
	s, err := readQuotedStringBody(reader)
	if err != nil {
		return "", NewTProtocolExceptionFromOsError(err)
	}
	str := string(JSON_QUOTE) + line + s
	v, ok := JsonUnquote(str)
	if !ok {
		return "", NewTProtocolException(INVALID_DATA, "Unable to parse as JSON string "+str)
	}
	return v, nil
}

func readQuotedStringBody(reader *bufio.Reader) (string, TProtocolException) {
	line, err := reader.ReadString(JSON_QUOTE)
	if err != nil {
		return "", NewTProtocolExceptionFromOsError(err)
	}
	l := len(line)
	// count number of escapes to see if we need to keep going
	i := 1
	for ; i < l; i++ {
		if line[l-i-1] != '\\' {
			break
		}
	}
	if i&0x01 == 1 {
		return line, nil
	}
	s, err := readQuotedStringBody(reader)
	if err != nil {
		return "", NewTProtocolExceptionFromOsError(err)
	}
	v := line + s
	return v, nil
}

func readBase64EncodedBody(reader *bufio.Reader) ([]byte, TProtocolException) {
	line, err := reader.ReadBytes(JSON_QUOTE)
	if err != nil {
		return line, NewTProtocolExceptionFromOsError(err)
	}
	line2 := line[0 : len(line)-1]
	l := len(line2)
	output := make([]byte, base64.StdEncoding.DecodedLen(l))
	n, err := base64.StdEncoding.Decode(output, line2)
	return output[0:n], NewTProtocolExceptionFromOsError(err)
}

func (p *TSimpleJSONProtocol) ReadBinary() ([]byte, TProtocolException) {
  var v []byte
  if err := readPreValue(p, p.reader); err != nil {
    return nil, err
  }
  b, _ := p.reader.Peek(len(JSON_NULL))
  if len(b) > 0 && b[0] == JSON_QUOTE {
    p.reader.ReadByte()
    value, err := readBase64EncodedBody(p.reader)
    v = value
    if err != nil {
      return v, err
    }
  } else if len(b) >= len(JSON_NULL) && string(b[0:len(JSON_NULL)]) == string(JSON_NULL) {
    _, err := p.reader.Read(b[0:len(JSON_NULL)])
    if err != nil {
      return v, NewTProtocolExceptionFromOsError(err)
    }
  } else {
    return v, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected a JSON string, found ", string(b)))
  }
  return v, readPostValue(p, p.reader)
}

func (p *TSimpleJSONProtocol) Flush() (err TProtocolException) {
	return NewTProtocolExceptionFromOsError(p.writer.Flush())
}

func (p *TSimpleJSONProtocol) Skip(fieldType TType) (err TProtocolException) {
	return SkipDefaultDepth(p, fieldType)
}

func (p *TSimpleJSONProtocol) Transport() TTransport {
	return p.trans
}


func writePreValue(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  cxt := _ParseContext(p.dumpContext.Last())
  switch cxt {
  case _CONTEXT_IN_LIST, _CONTEXT_IN_OBJECT_NEXT_KEY:
    if _, e := w.Write(JSON_COMMA); e != nil {
      return NewTProtocolExceptionFromOsError(e)
    }
    break
  case _CONTEXT_IN_OBJECT_NEXT_VALUE:
    if _, e := w.Write(JSON_COLON); e != nil {
      return NewTProtocolExceptionFromOsError(e)
    }
    break
  }
  return nil
}

func writePostValue(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  cxt := _ParseContext(p.dumpContext.Last())
  switch cxt {
  case _CONTEXT_IN_LIST_FIRST:
    p.dumpContext.Pop()
    p.dumpContext.Push(int(_CONTEXT_IN_LIST))
    break
  case _CONTEXT_IN_OBJECT_FIRST:
    p.dumpContext.Pop()
    p.dumpContext.Push(int(_CONTEXT_IN_OBJECT_NEXT_VALUE))
    break
  case _CONTEXT_IN_OBJECT_NEXT_KEY:
    p.dumpContext.Pop()
    p.dumpContext.Push(int(_CONTEXT_IN_OBJECT_NEXT_VALUE))
    break
  case _CONTEXT_IN_OBJECT_NEXT_VALUE:
    p.dumpContext.Pop()
    p.dumpContext.Push(int(_CONTEXT_IN_OBJECT_NEXT_KEY))
    break
  }
  return nil
}

func writeBool(p *TSimpleJSONProtocol, w TTransport, value bool) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  var v []byte
  if value {
    v = JSON_TRUE
  } else {
    v = JSON_FALSE
  }
  if _, e := w.Write(v); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  return writePostValue(p, w)
}

func writeNull(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  if _, e := w.Write(JSON_NULL); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  return writePostValue(p, w)
}

func writeF64(p *TSimpleJSONProtocol, w TTransport, value float64) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  var v string
  if math.IsNaN(value) {
    v = string(JSON_QUOTE) + JSON_NAN + string(JSON_QUOTE)
  } else if math.IsInf(value, 1) {
    v = string(JSON_QUOTE) + JSON_INFINITY + string(JSON_QUOTE)
  } else if math.IsInf(value, -1) {
    v = string(JSON_QUOTE) + JSON_NEGATIVE_INFINITY + string(JSON_QUOTE)
  } else {
    v = strconv.Ftoa64(value, 'g', -1)
  }
  if e := writeStringData(w, v); e != nil {
    return e
  }
  return writePostValue(p, w)
}

func writeI64(p *TSimpleJSONProtocol, w TTransport, value int64) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  if e := writeStringData(w, strconv.Itoa64(value)); e != nil {
    return e
  }
  return writePostValue(p, w)
}

func writeString(p *TSimpleJSONProtocol, w TTransport, s string) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  if e := writeStringData(p.writer, JsonQuote(s)); e != nil {
    return e
  }
  return writePostValue(p, w)
}

func writeStringData(w TTransport, s string) TProtocolException {
	_, e := io.Copyn(w, strings.NewReader(s), int64(len(s)))
	return NewTProtocolExceptionFromOsError(e)
}

func writeObjectBegin(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  if _, e := w.Write(JSON_LBRACE); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  p.dumpContext.Push(int(_CONTEXT_IN_OBJECT_FIRST))
  return nil
}

func writeObjectEnd(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  if _, e := w.Write(JSON_RBRACE); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  p.dumpContext.Pop()
  if e := writePostValue(p, w); e != nil {
    return e
  }
  return nil
}

func writeListBegin(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  if e := writePreValue(p, w); e != nil {
    return e
  }
  if _, e := w.Write(JSON_LBRACKET); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  p.dumpContext.Push(int(_CONTEXT_IN_LIST_FIRST))
  return nil
}

func writeListEnd(p *TSimpleJSONProtocol, w TTransport) TProtocolException {
  if _, e := w.Write(JSON_RBRACKET); e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  p.dumpContext.Pop()
  if e := writePostValue(p, w); e != nil {
    return e
  }
  return nil
}

func writeElemListBegin(p *TSimpleJSONProtocol, w TTransport, elemType TType, size int) TProtocolException {
  if e := writeListBegin(p, w); e != nil {
    return e
  }
  if e := p.WriteByte(byte(elemType)); e != nil {
    return e
  }
  if e := p.WriteI64(int64(size)); e != nil {
    return e
  }
  return nil
}

func readPreValue(p *TSimpleJSONProtocol, r *bufio.Reader) TProtocolException {
  if e := readNonSignificantWhitespace(r); e != nil {
  	return NewTProtocolExceptionFromOsError(e)
  }
  cxt := _ParseContext(p.parseContextStack.Last())
  b, _ := r.Peek(1)
  switch cxt {
  case _CONTEXT_IN_LIST:
    if len(b) > 0 {
      switch b[0] {
      case JSON_RBRACKET[0]:
        return nil
      case JSON_COMMA[0]:
        r.ReadByte()
        if e := readNonSignificantWhitespace(r); e != nil {
        	return NewTProtocolExceptionFromOsError(e)
        }
        return nil
      default:
        return NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected \"]\" or \",\" in list context, but found \"", string(b), "\""))
      }
    }
    break
  case _CONTEXT_IN_OBJECT_NEXT_KEY:
    if len(b) > 0 {
      switch b[0] {
      case JSON_RBRACE[0]:
        return nil
      case JSON_COMMA[0]:
        r.ReadByte()
        if e := readNonSignificantWhitespace(r); e != nil {
        	return NewTProtocolExceptionFromOsError(e)
        }
        return nil
      default:
        return NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected \"}\" or \",\" in object context, but found \"", string(b), "\""))
      }
    }
    break
  case _CONTEXT_IN_OBJECT_NEXT_VALUE:
    if len(b) > 0 {
      switch b[0] {
      case JSON_COLON[0]:
        r.ReadByte()
        if e := readNonSignificantWhitespace(r); e != nil {
        	return NewTProtocolExceptionFromOsError(e)
        }
        return nil
      default:
        return NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected \":\" in object context, but found \"", string(b), "\""))
      }
    }
    break
  }
  return nil
}

func readPostValue(p *TSimpleJSONProtocol, r *bufio.Reader) TProtocolException {
  if e := readNonSignificantWhitespace(r); e != nil {
  	return NewTProtocolExceptionFromOsError(e)
  }
  cxt := _ParseContext(p.parseContextStack.Last())
  switch cxt {
  case _CONTEXT_IN_LIST_FIRST:
    p.parseContextStack.Pop()
    p.parseContextStack.Push(int(_CONTEXT_IN_LIST))
    break
  case _CONTEXT_IN_OBJECT_FIRST, _CONTEXT_IN_OBJECT_NEXT_KEY:
    p.parseContextStack.Pop()
    p.parseContextStack.Push(int(_CONTEXT_IN_OBJECT_NEXT_VALUE))
    break
  case _CONTEXT_IN_OBJECT_NEXT_VALUE:
    p.parseContextStack.Pop()
    p.parseContextStack.Push(int(_CONTEXT_IN_OBJECT_NEXT_KEY))
    break
  }
  return nil
}

func readObjectStart(p *TSimpleJSONProtocol, r *bufio.Reader) (bool, TProtocolException) {
  if err := readPreValue(p, r); err != nil {
    return false, err
  }
  b, _ := r.Peek(len(JSON_NULL))
  if len(b) > 0 && b[0] == JSON_LBRACE[0] {
    r.ReadByte()
    p.parseContextStack.Push(int(_CONTEXT_IN_OBJECT_FIRST))
    return false, nil
  } else if len(b) >= len(JSON_NULL) && string(b[0:len(JSON_NULL)]) == string(JSON_NULL) {
    return true, nil
  }
  return false, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected '{' or null, but found '", string(b), "'"))
}

func readObjectEnd(p *TSimpleJSONProtocol, r *bufio.Reader) TProtocolException {
  if e := readNonSignificantWhitespace(r); e != nil {
  	return NewTProtocolExceptionFromOsError(e)
  }
  if isNull, err := readIfNull(r); isNull || err != nil {
    return err
  }
  cxt := _ParseContext(p.parseContextStack.Last())
  if cxt != _CONTEXT_IN_OBJECT_FIRST && cxt != _CONTEXT_IN_OBJECT_NEXT_KEY {
    return NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected to be in the Object Context, but not in Object Context"))
  }
  b, _ := r.Peek(1)
  if len(b) > 0 && b[0] == JSON_RBRACE[0] {
    _, e := r.ReadByte()
    p.parseContextStack.Pop()
    return NewTProtocolExceptionFromOsError(e)
  }
  return readPostValue(p, r)
}

func readListBegin(p *TSimpleJSONProtocol, r *bufio.Reader) (bool, TProtocolException) {
  if e := readPreValue(p, r); e != nil {
    return false, e
  }
  b, e := r.Peek(len(JSON_NULL))
  if e == nil && len(b) >= 1 && b[0] == JSON_LBRACKET[0] {
    p.parseContextStack.Push(int(_CONTEXT_IN_LIST_FIRST))
    r.ReadByte()
    return false, nil
  } else if e == nil && len(b) >= len(JSON_NULL) && string(b) == string(JSON_NULL) {
    return true, nil
  }
  return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected 'null' or '{', received '%q'", b))
}

func readElemListBegin(p *TSimpleJSONProtocol, r *bufio.Reader) (elemType TType, size int, e TProtocolException) {
  if isNull, e := readListBegin(p, r); isNull || e != nil {
  	return VOID, 0, e
  }
  bElemType, err := p.ReadByte()
  elemType = TType(bElemType)
  if err != nil {
    return elemType, size, err
  }
  nSize, err2 := p.ReadI64()
  size = int(nSize)
  return elemType, size, err2
}

func readListEnd(p *TSimpleJSONProtocol, r *bufio.Reader) TProtocolException {
  if e := readNonSignificantWhitespace(r); e != nil {
  	return NewTProtocolExceptionFromOsError(e)
  }
  if isNull, err := readIfNull(r); isNull || err != nil {
    return err
  }
  if _ParseContext(p.parseContextStack.Last()) != _CONTEXT_IN_LIST {
    return NewTProtocolException(INVALID_DATA, "Expected to be in the List Context, but not in List Context")
  }
  b, e := r.Peek(1)
  if e != nil {
    return NewTProtocolExceptionFromOsError(e)
  }
  if len(b) >= 1 && b[0] == JSON_RBRACKET[0] {
    r.ReadByte()
    p.parseContextStack.Pop()
    return readPostValue(p, r)
  }
  return NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected \"", string(JSON_RBRACKET[0]), "\" but found \"", string(b), "\""))
}

func readSingleValue(r *bufio.Reader) (interface{}, TType, TProtocolException) {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return nil, VOID, NewTProtocolExceptionFromOsError(e)
	}
	b, e := r.Peek(10)
	fmt.Fprint(os.Stderr, "[JSON]: Read bytes \"", b, "\": \"", string(b), "\"")
	if len(b) > 0 {
	  c := b[0]
  	switch c {
  	case JSON_NULL[0]:
  		buf := make([]byte, len(JSON_NULL))
  		_, e := r.Read(buf)
  		if e != nil {
  			return nil, VOID, NewTProtocolExceptionFromOsError(e)
  		}
  		if string(JSON_NULL) != string(buf) {
  			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_NULL)+"' but found '"+string(buf)+"' while parsing JSON.")
  			return nil, VOID, e
  		}
  		return nil, VOID, nil
  	case JSON_QUOTE:
  	  r.ReadByte()
  		v, e := ReadStringBody(r)
  		if e != nil {
  			return v, UTF8, NewTProtocolExceptionFromOsError(e)
  		}
  		if v == JSON_INFINITY {
  			return INFINITY, DOUBLE, nil
  		} else if v == JSON_NEGATIVE_INFINITY {
  			return NEGATIVE_INFINITY, DOUBLE, nil
  		} else if v == JSON_NAN {
  			return NAN, DOUBLE, nil
  		}
  		return v, UTF8, nil
  	case JSON_TRUE[0]:
  		buf := make([]byte, len(JSON_TRUE))
  		_, e := r.Read(buf)
  		if e != nil {
  			return true, BOOL, NewTProtocolExceptionFromOsError(e)
  		}
  		if string(JSON_TRUE) != string(buf) {
  			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_TRUE)+"' but found '"+string(buf)+"' while parsing JSON.")
  			return true, BOOL, NewTProtocolExceptionFromOsError(e)
  		}
  		return true, BOOL, nil
  	case JSON_FALSE[0]:
  		buf := make([]byte, len(JSON_FALSE))
  		_, e := r.Read(buf)
  		if e != nil {
  			return false, BOOL, NewTProtocolExceptionFromOsError(e)
  		}
  		if string(JSON_FALSE) != string(buf) {
  			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_FALSE)+"' but found '"+string(buf)+"' while parsing JSON.")
  			return false, BOOL, NewTProtocolExceptionFromOsError(e)
  		}
  		return false, BOOL, nil
  	case JSON_LBRACKET[0]:
  		_, e := r.ReadByte()
			return make([]interface{}, 0), LIST, NewTProtocolExceptionFromOsError(e)
  	case JSON_LBRACE[0]:
  		_, e := r.ReadByte()
			return make(map[string]interface{}), STRUCT, NewTProtocolExceptionFromOsError(e)
  	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'e', 'E', '.', '+', '-', JSON_INFINITY[0], JSON_NAN[0]:
  		// assume numeric
  		v, e := readNumeric(r)
  		return v, DOUBLE, e
  	default:
  		return nil, VOID, NewTProtocolException(INVALID_DATA, "Expected element in list but found '"+string(c)+"' while parsing JSON.")
  	}
	}
	return nil, VOID, NewTProtocolException(INVALID_DATA, "Cannot read a single element while parsing JSON.")
  
}


func readIfNull(reader *bufio.Reader) (bool, TProtocolException) {
	b, err := reader.ReadByte()
	if err != nil {
		if err == os.EOF {
			return false, nil
		}
		return false, NewTProtocolExceptionFromOsError(err)
	}
	if b != JSON_NULL[0] {
		return false, NewTProtocolExceptionFromOsError(reader.UnreadByte())
	}
	buf := make([]byte, len(JSON_NULL))
	buf[0] = b
	for i := 1; i < len(JSON_NULL); i++ {
		b, err = reader.ReadByte()
		buf[i] = b
		if b != JSON_NULL[i] {
			reader.UnreadByte()
			return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expecting 'null', found '%s'", string(buf)))
		}
	}
	return true, nil
}

func readQuoteIfNext(reader *bufio.Reader) {
  b, _ := reader.Peek(1)
  if len(b) > 0 && b[0] == JSON_QUOTE {
    reader.ReadByte()
  }
}

func readNumeric(reader *bufio.Reader) (Numeric, TProtocolException) {
	isNull, err := readIfNull(reader)
	if isNull || err != nil {
		return NUMERIC_NULL, err
	}
	hasDecimalPoint := false
	nextCanBeSign := true
	hasE := false
	MAX_LEN := 40
	buf := bytes.NewBuffer(make([]byte, 0, MAX_LEN))
	continueFor := true
	inQuotes := false
	for continueFor {
		c, err := reader.ReadByte()
		if err != nil {
			if err == os.EOF {
				break
			}
			return NUMERIC_NULL, NewTProtocolExceptionFromOsError(err)
		}
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buf.WriteByte(c)
			nextCanBeSign = false
		case '.':
			if hasDecimalPoint {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number with multiple decimal points '%s.'", buf.String()))
			}
			if hasE {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number with decimal points in the exponent '%s.'", buf.String()))
			}
			buf.WriteByte(c)
			hasDecimalPoint, nextCanBeSign = true, false
		case 'e', 'E':
			if hasE {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number with multiple exponents '%s%c'", buf.String(), c))
			}
			buf.WriteByte(c)
			hasE, nextCanBeSign = true, true
		case '-', '+':
			if !nextCanBeSign {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprint("Negative sign within number"))
			}
			buf.WriteByte(c)
			nextCanBeSign = false
		case ' ', 0, '\t', '\n', '\r', JSON_RBRACE[0], JSON_RBRACKET[0], JSON_COMMA[0], JSON_COLON[0]:
			reader.UnreadByte()
			continueFor = false
		case JSON_NAN[0]:
			if buf.Len() == 0 {
				buffer := make([]byte, len(JSON_NAN))
				buffer[0] = c
				_, e := reader.Read(buffer[1:])
				if e != nil {
					return NUMERIC_NULL, NewTProtocolExceptionFromOsError(e)
				}
				if JSON_NAN != string(buffer) {
					e := NewTProtocolException(INVALID_DATA, "Expected '"+JSON_NAN+"' but found '"+string(buffer)+"' while parsing JSON.")
					return NUMERIC_NULL, e
				}
  			if inQuotes {
  			  readQuoteIfNext(reader)
  			}
				return NAN, nil
			} else {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number starting with character '%c'", c))
			}
		case JSON_INFINITY[0]:
			if buf.Len() == 0 || (buf.Len() == 1 && buf.Bytes()[0] == '+') {
				buffer := make([]byte, len(JSON_INFINITY))
				buffer[0] = c
				_, e := reader.Read(buffer[1:])
				if e != nil {
					return NUMERIC_NULL, NewTProtocolExceptionFromOsError(e)
				}
				if JSON_INFINITY != string(buffer) {
					e := NewTProtocolException(INVALID_DATA, "Expected '"+JSON_INFINITY+"' but found '"+string(buffer)+"' while parsing JSON.")
					return NUMERIC_NULL, e
				}
  			if inQuotes {
  			  readQuoteIfNext(reader)
  			}
				return INFINITY, nil
			} else if buf.Len() == 1 && buf.Bytes()[0] == JSON_NEGATIVE_INFINITY[0] {
				buffer := make([]byte, len(JSON_NEGATIVE_INFINITY))
				buffer[0] = JSON_NEGATIVE_INFINITY[0]
				buffer[1] = c
				_, e := reader.Read(buffer[2:])
				if e != nil {
					return NUMERIC_NULL, NewTProtocolExceptionFromOsError(e)
				}
				if JSON_NEGATIVE_INFINITY != string(buffer) {
					e := NewTProtocolException(INVALID_DATA, "Expected '"+JSON_NEGATIVE_INFINITY+"' but found '"+string(buffer)+"' while parsing JSON.")
					return NUMERIC_NULL, e
				}
				if inQuotes {
				  readQuoteIfNext(reader)
				}
				return NEGATIVE_INFINITY, nil
			} else {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number starting with character '%c' due to existing buffer %s", c, buf.String()))
			}
		case JSON_QUOTE:
		  if !inQuotes {
		    inQuotes = true
	    } else {
	      break
	    }
		default:
			return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number starting with character '%c'", c))
		}
	}
	if buf.Len() == 0 {
		return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprint("Unable to parse number from empty string ''"))
	}
	return NewNumericFromJSONString(buf.String(), false), nil
}

