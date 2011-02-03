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
	"container/list"
	"encoding/base64"
	"fmt"
	"io"
	"json"
	"math"
	"os"
	"strconv"
	"strings"
)

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

	_BASE_CONTEXT Context

	/**
	 * Stack of nested contexts that we may be in.
	 */
	writeContextStack *list.List

	/**
	 * Current context that we are in
	 */
	writeContext Context

	/**
	 * Stack of nested contexts that we may be in.
	 */
	readContextStack *list.List

	/**
	 * Current context that we are in
	 */
	readContext Context
	writer      TTransport
	reader      *bufio.Reader
}

/**
 * Constructor
 */
func NewTSimpleJSONProtocol(t TTransport) *TSimpleJSONProtocol {
	cxt := NewContextWriter(t)
	cxt2 := NewContext()
	l := list.New()
	l.PushFront(cxt)
	l2 := list.New()
	l2.PushFront(cxt)
	v := &TSimpleJSONProtocol{trans: t,
		writeContext:      cxt,
		_BASE_CONTEXT:     cxt,
		writeContextStack: l,
		readContext:       cxt2,
		readContextStack:  l2,
		writer:            t,
		reader:            bufio.NewReader(t),
	}
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
		if err == os.EOF {
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

type Context interface {
	Read(r *bufio.Reader) os.Error
	Close() os.Error
	ReadClose(r *bufio.Reader) os.Error
	Populate(r *bufio.Reader) os.Error
	Value() interface{}
	TType() TType
	WriteNull() os.Error
	WriteBool(value bool) os.Error
	WriteI32(value int32) os.Error
	WriteI64(value int64) os.Error
	WriteDouble(value float64) os.Error
	WriteString(value string) os.Error
	WriteBinary(value []byte) os.Error
	WriteStartSubContext() os.Error
	WriteEndSubContext() os.Error
}

type context struct {
	value  interface{}
	writer io.Writer
}

func NewContext() Context {
	return &context{}
}

func NewContextWriter(w io.Writer) Context {
	return &context{writer: w}
}

func (p *context) Read(r *bufio.Reader) os.Error {
	return readNonSignificantWhitespace(r)
}

func (p *context) Close() os.Error {
	if p.writer != nil {
		if c, ok := p.writer.(io.Closer); ok {
			e := c.Close()
			p.writer = nil
			return e
		}
		p.writer = nil
	}
	return nil
}

func (p *context) ReadClose(r *bufio.Reader) os.Error {
	return readNonSignificantWhitespace(r)
}

func (p *context) Numeric() Numeric {
	n, ok := p.value.(Numeric)
	if ok {
		return n
	}
	return nil
}

func (p *context) WriteNull() (err os.Error) {
	if p.writer == nil {
		return os.EOF
	}
	p.value = nil
	_, err = p.writer.Write(JSON_NULL)
	return err
}

func (p *context) WriteBool(value bool) (err os.Error) {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	if value {
		_, err = p.writer.Write(JSON_TRUE)
	} else {
		_, err = p.writer.Write(JSON_FALSE)
	}
	return err
}

func (p *context) WriteI32(value int32) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	return p.writeStringData(strconv.Itoa(int(value)))
}

func (p *context) WriteI64(value int64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	return p.writeStringData(strconv.Itoa64(value))
}

func (p *context) WriteDouble(value float64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	return p.writeStringData(strconv.Ftoa64(value, 'g', 10))
}

func (p *context) WriteString(value string) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	return p.writeStringData(JsonQuote(value))
}

func (p *context) WriteBinary(value []byte) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	p.value = value
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return e
	}
	writer := base64.NewEncoder(base64.StdEncoding, p.writer)
	if _, e := writer.Write(value); e != nil {
		return e
	}
	if e := writer.Close(); e != nil {
		return e
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return e
	}
	return nil
}

func (p *context) writeStringData(s string) (err os.Error) {
	_, err = io.Copyn(p.writer, strings.NewReader(s), int64(len(s)))
	return err
}

func (p *context) WriteStartSubContext() os.Error { return nil }

func (p *context) WriteEndSubContext() os.Error { return nil }

func (p *context) Int64() int64 {
	n := p.Numeric()
	if n != nil {
		return 0
	}
	return n.Int64()
}

func (p *context) Int32() int32 {
	n := p.Numeric()
	if n != nil {
		return 0
	}
	return n.Int32()
}

func (p *context) Int() int {
	n := p.Numeric()
	if n != nil {
		return 0
	}
	return n.Int()
}

func (p *context) Float64() float64 {
	n := p.Numeric()
	if n != nil {
		return 0
	}
	return n.Float64()
}

func (p *context) Float32() float32 {
	n := p.Numeric()
	if n != nil {
		return 0
	}
	return n.Float32()
}

func (p *context) String() string {
	s, ok := p.value.(string)
	if ok {
		return s
	}
	return ""
}

func (p *context) isNull() bool {
	return p.value == nil
}

func (p *context) ListContext() *ListContext {
	l, ok := p.value.(*ListContext)
	if ok {
		return l
	}
	return nil
}

func (p *context) ObjectContext() *ObjectContext {
	m, ok := p.value.(*ObjectContext)
	if ok {
		return m
	}
	return nil
}

func (p *context) Value() interface{} {
	return p.value
}

func (p *context) TType() TType {
	if p.value == nil {
		return STOP
	}
	_, ok := p.value.(string)
	if ok {
		return UTF8
	}
	_, ok = p.value.(Numeric)
	if ok {
		return DOUBLE
	}
	_, ok = p.value.(bool)
	if ok {
		return BOOL
	}
	_, ok = p.value.(*ListContext)
	if ok {
		return LIST
	}
	_, ok = p.value.(*ObjectContext)
	if ok {
		return STRUCT
	}
	return STOP
}

func (p *context) Populate(r *bufio.Reader) os.Error {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return e
	}
	c, e := r.ReadByte()
	switch c {
	case JSON_NULL[0]:
		buf := make([]byte, len(JSON_NULL))
		buf[0] = c
		_, e := r.Read(buf[1:])
		if e != nil {
			return e
		}
		if string(JSON_NULL) != string(buf) {
			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_NULL)+"' but found '"+string(buf)+"' while parsing JSON.")
			return e
		}
		p.value = nil
	case JSON_QUOTE:
		v, e := ReadStringBody(r)
		if e != nil {
			return e
		}
		if v == JSON_INFINITY {
			p.value = INFINITY
		} else if v == JSON_NEGATIVE_INFINITY {
			p.value = NEGATIVE_INFINITY
		} else if v == JSON_NAN {
			p.value = NAN
		} else {
			p.value = v
		}
	case JSON_TRUE[0]:
		buf := make([]byte, len(JSON_TRUE))
		buf[0] = c
		_, e := r.Read(buf[1:])
		if e != nil {
			return e
		}
		if string(JSON_TRUE) != string(buf) {
			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_TRUE)+"' but found '"+string(buf)+"' while parsing JSON.")
			return e
		}
		p.value = true
	case JSON_FALSE[0]:
		buf := make([]byte, len(JSON_FALSE))
		buf[0] = c
		_, e := r.Read(buf[1:])
		if e != nil {
			return e
		}
		if string(JSON_FALSE) != string(buf) {
			e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_FALSE)+"' but found '"+string(buf)+"' while parsing JSON.")
			return e
		}
		p.value = false
	case JSON_LBRACKET[0]:
		cxt := NewListContext()
		e = cxt.Populate(r)
		if e != nil {
			return NewTProtocolExceptionFromOsError(e)
		}
		p.value = cxt
	case JSON_LBRACE[0]:
		cxt := NewObjectContext()
		e = cxt.Populate(r)
		if e != nil {
			return NewTProtocolExceptionFromOsError(e)
		}
		p.value = cxt
	case JSON_INFINITY[0]:
		buf := make([]byte, len(JSON_INFINITY))
		buf[0] = c
		_, e := r.Read(buf[1:])
		if e != nil {
			return e
		}
		if JSON_INFINITY != string(buf) {
			e := NewTProtocolException(INVALID_DATA, "Expected '"+JSON_INFINITY+"' but found '"+string(buf)+"' while parsing JSON.")
			return e
		}
		p.value = INFINITY
	case JSON_NAN[0]:
		buf := make([]byte, len(JSON_NAN))
		buf[0] = c
		_, e := r.Read(buf[1:])
		if e != nil {
			return e
		}
		if JSON_NAN != string(buf) {
			e := NewTProtocolException(INVALID_DATA, "Expected '"+JSON_NAN+"' but found '"+string(buf)+"' while parsing JSON.")
			return e
		}
		p.value = NAN
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'e', 'E', '.', '+', '-':
		// assume numeric
		r.UnreadByte()
		v, e := readNumeric(r)
		p.value = v
		return e
	default:
		e := NewTProtocolException(INVALID_DATA, "Expected element in list but found '"+string(c)+"' while parsing JSON.")
		return e
	}
	return nil
}


type ListContext struct {
	first    bool
	readList *list.List
	elemType TType
	next     *list.Element
	reader   *bufio.Reader
	writer   io.Writer
}

func NewListContext() *ListContext {
	return &ListContext{first: true, readList: list.New(), elemType: TType(STOP)}
}

func NewListContextEType(elemType TType) *ListContext {
	return &ListContext{first: true, readList: list.New(), elemType: elemType}
}

func NewListContextETypeWriter(elemType TType, w io.Writer) *ListContext {
	return &ListContext{first: true, readList: list.New(), elemType: elemType, writer: w}
}

func (p *ListContext) Len() int {
	return p.readList.Len()
}

func (p *ListContext) TType() TType {
	return LIST
}

func (p *ListContext) ElementType() TType {
	return p.elemType
}

func (p *ListContext) Next() *list.Element {
	if p.next == nil {
		p.next = p.readList.Front()
	} else {
		p.next = p.next.Next()
	}
	return p.next
}

func (p *ListContext) Front() *list.Element {
	return p.readList.Front()
}

func (p *ListContext) Back() *list.Element {
	return p.readList.Back()
}

func (p *ListContext) List() *list.List {
	return p.readList
}

func (p *ListContext) Value() interface{} {
	return p.readList
}

func (p *ListContext) WriteNull() os.Error {
	if p.writer == nil {
		return os.EOF
	}
	if e := p.writeContext(); e != nil {
		return e
	}
	_, err := p.writer.Write(JSON_NULL)
	return err
}

func (p *ListContext) WriteBool(value bool) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	s := JSON_FALSE
	if value {
		s = JSON_TRUE
	}
	if e := p.writeContext(); e != nil {
		return e
	}
	_, err := p.writer.Write(s)
	return err
}

func (p *ListContext) WriteI32(value int32) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	return p.writeStringData(strconv.Itoa(int(value)))
}

func (p *ListContext) WriteI64(value int64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	return p.writeStringData(strconv.Itoa64(value))
}

func (p *ListContext) WriteDouble(value float64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	return p.writeStringData(strconv.Ftoa64(value, 'g', 10))
}

func (p *ListContext) WriteString(value string) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	return p.writeStringData(JsonQuote(value))
}

func (p *ListContext) WriteBinary(value []byte) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	if e := p.writeContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	writer := base64.NewEncoder(base64.StdEncoding, p.writer)
	if _, e := writer.Write(value); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if e := writer.Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return nil
}

func (p *ListContext) writeContext() os.Error {
	if p.first {
		p.first = false
		_, err := p.writer.Write(JSON_LBRACKET)
		return err
	}
	_, err := p.writer.Write(JSON_COMMA)
	return err
}

func (p *ListContext) writeStringData(s string) os.Error {
	if e := p.writeContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	_, e := io.Copyn(p.writer, strings.NewReader(s), int64(len(s)))
	return NewTProtocolExceptionFromOsError(e)
}

func (p *ListContext) WriteStartSubContext() os.Error {
	return p.writeContext()
}

func (p *ListContext) WriteEndSubContext() os.Error {
	return nil
}

func (p *ListContext) Populate(r *bufio.Reader) os.Error {
	isFirst := p.first
	p.readList = list.New()
	for {
		e := readNonSignificantWhitespace(r)
		if e != nil {
			return e
		}
		c, e := r.ReadByte()
		if c == JSON_RBRACKET[0] {
			break
		}
		if isFirst {
			isFirst = false
			r.UnreadByte()
		} else if c == JSON_COMMA[0] {
			e = readNonSignificantWhitespace(r)
			if e != nil {
				return e
			}
		} else {
			e := NewTProtocolException(INVALID_DATA, "Expected ',' or ']' but found '"+string(c)+"' while parsing JSON.")
			return e
		}
		cxt := NewContext()
		e = cxt.Populate(r)
		if e != nil {
			return NewTProtocolExceptionFromOsError(e)
		}
		p.readList.PushBack(cxt.Value())
		t := cxt.TType()
		if t != STOP {
			p.elemType = t
		}
	}
	return nil
}

func (p *ListContext) Read(r *bufio.Reader) os.Error {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return e
	}
	if p.first {
		p.first = false
	} else {
		_, err := r.ReadBytes(JSON_COMMA[0])
		if err != nil {
			return err
		}
	}
	return readNonSignificantWhitespace(r)
}

func (p *ListContext) Close() os.Error {
	if p.first {
		if e := p.writeContext(); e != nil {
			return e
		}
	}
	_, e := p.writer.Write(JSON_RBRACKET)
	p.writer = nil
	return e
}

func (p *ListContext) ReadClose(r *bufio.Reader) os.Error {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return e
	}
	_, err := r.ReadBytes(JSON_RBRACKET[0])
	if err != nil {
		return err
	}
	return readNonSignificantWhitespace(r)
}

type StructEntry struct {
	key   interface{}
	value interface{}
}

type ObjectContext struct {
	first     bool
	onKey     bool
	onValue   bool
	readList  *list.List
	keyType   TType
	valueType TType
	next      *list.Element
	writer    io.Writer
}

func NewObjectContext() *ObjectContext {
	return &ObjectContext{first: true, onKey: true, onValue: false, readList: list.New(), keyType: TType(STOP), valueType: TType(STOP)}
}

func NewObjectContextKVType(keyType, valueType TType) *ObjectContext {
	return &ObjectContext{first: true, onKey: true, onValue: false, readList: list.New(), keyType: keyType, valueType: valueType}
}

func NewObjectContextKVTypeWriter(keyType, valueType TType, w io.Writer) *ObjectContext {
	return &ObjectContext{first: true, onKey: true, onValue: false, readList: list.New(), keyType: keyType, valueType: valueType, writer: w}
}

func (p *ObjectContext) KeyType() TType {
	return p.keyType
}

func (p *ObjectContext) ValueType() TType {
	return p.valueType
}

func (p *ObjectContext) TType() TType {
	return STRUCT
}

func (p *ObjectContext) Len() int {
	return p.readList.Len()
}

func (p *ObjectContext) Next() *list.Element {
	if p.next == nil {
		p.next = p.readList.Front()
	} else {
		p.next = p.next.Next()
	}
	return p.next
}

func (p *ObjectContext) Front() *list.Element {
	return p.readList.Front()
}

func (p *ObjectContext) Back() *list.Element {
	return p.readList.Back()
}

func (p *ObjectContext) List() *list.List {
	return p.readList
}

func (p *ObjectContext) Value() interface{} {
	return p.readList
}

func (p *ObjectContext) WriteNull() os.Error {
	if p.writer == nil {
		return os.EOF
	}
	if e := p.writeContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	_, err := p.writer.Write(JSON_NULL)
	p.onKey, p.onValue = p.onValue, p.onKey
	return err
}

func (p *ObjectContext) WriteBool(value bool) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	s, v := JSON_FALSE, []byte{'"', 'f', 'a', 'l', 's', 'e', '"'}
	if value {
		s, v = JSON_TRUE, []byte{'"', 't', 'r', 'u', 'e', '"'}
	}
	if p.onKey {
		if e := p.writeContext(); e != nil {
			return e
		}
		if _, e := p.writer.Write(v); e != nil {
			return e
		}
		p.onKey, p.onValue = p.onValue, p.onKey
		return nil
	}
	if e := p.writeContext(); e != nil {
		return e
	}
	_, err := p.writer.Write(s)
	p.onKey, p.onValue = p.onValue, p.onKey
	return err
}

func (p *ObjectContext) WriteI32(value int32) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	s := strconv.Itoa(int(value))
	if p.onKey {
		s = strconv.Quote(s)
	}
	return p.writeStringData(s)
}

func (p *ObjectContext) WriteI64(value int64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	s := strconv.Itoa64(value)
	if p.onKey {
		s = strconv.Quote(s)
	}
	return p.writeStringData(s)
}

func (p *ObjectContext) WriteDouble(value float64) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	s := strconv.Ftoa64(value, 'g', 10)
	if p.onKey {
		s = strconv.Quote(s)
	}
	return p.writeStringData(s)
}

func (p *ObjectContext) WriteString(value string) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	return p.writeStringData(JsonQuote(value))
}

func (p *ObjectContext) WriteBinary(value []byte) os.Error {
	if p.writer == nil {
		return os.EOF
	}
	if e := p.writeContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	writer := base64.NewEncoder(base64.StdEncoding, p.writer)
	if _, e := writer.Write(value); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if e := writer.Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.onKey, p.onValue = p.onValue, p.onKey
	return nil
}

func (p *ObjectContext) writeContext() os.Error {
	if p.first {
		p.first = false
		p.onKey, p.onValue = true, false
		if !p.keyType.IsBaseType() {
			s := fmt.Sprint("{\"__type__\":\"map\",\"__keytype__\":\"", p.keyType.String(), "\",\"__valuetype__\":\"", p.valueType.String(), "\",\"__values__\":[")
			return p.writeStringData(s)
		}
		_, err := p.writer.Write(JSON_LBRACE)
		return err
	}
	if !p.keyType.IsBaseType() {
		v := JSON_COMMA
		if p.onValue {
			v = json_nonbase_map_elem_bytes
		}
		_, err := p.writer.Write(v)
		return err
	}
	v := JSON_COMMA
	if p.onValue {
		v = JSON_COLON
	}
	_, err := p.writer.Write(v)
	return err
}

func (p *ObjectContext) writeStringData(s string) os.Error {
	if e := p.writeContext(); e != nil {
		return e
	}
	_, err := io.Copyn(p.writer, strings.NewReader(s), int64(len(s)))
	p.onKey, p.onValue = p.onValue, p.onKey
	return err
}

func (p *ObjectContext) WriteStartSubContext() os.Error {
	return p.writeContext()
}

func (p *ObjectContext) WriteEndSubContext() os.Error {
	return nil
}

func (p *ObjectContext) Populate(r *bufio.Reader) os.Error {
	isFirst := p.first
	isKey := true
	p.readList = list.New()
	entry := &StructEntry{}
	for {
		e := readNonSignificantWhitespace(r)
		if e != nil {
			return e
		}
		c, e := r.ReadByte()
		if isKey {
			if c == JSON_RBRACE[0] {
				break
			}
			if isFirst {
				isFirst = false
			} else if c == JSON_COMMA[0] {
				e = readNonSignificantWhitespace(r)
				if e != nil {
					return e
				}
				c, e = r.ReadByte()
			} else {
				e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_COMMA[0])+"' or '"+string(JSON_RBRACE[0])+"' but found '"+string(c)+"' while parsing JSON map.")
				return e
			}
			entry = &StructEntry{}
		} else {
			if c == JSON_COLON[0] {
				e = readNonSignificantWhitespace(r)
				if e != nil {
					return e
				}
				c, e = r.ReadByte()
			} else {
				e := NewTProtocolException(INVALID_DATA, "Expected '"+string(JSON_COLON[0])+"' but found '"+string(c)+"' while parsing JSON map.")
				return e
			}
		}
		cxt := NewContext()
		e = cxt.Populate(r)
		if e != nil {
			return NewTProtocolExceptionFromOsError(e)
		}
		t := cxt.TType()
		if isKey {
			entry.key = cxt.Value()
			if t != STOP {
				p.keyType = t
			}
			isKey = false
		} else {
			entry.value = cxt.Value()
			if t != STOP {
				p.valueType = t
			}
			p.readList.PushBack(entry)
			isKey = true
		}
	}
	return nil
}

func (p *ObjectContext) Read(r *bufio.Reader) os.Error {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return e
	}
	if p.first {
		p.first = false
		p.onKey, p.onValue = true, false
	} else {
		var b byte
		if p.onKey {
			b = JSON_COLON[0]
		} else {
			b = JSON_COMMA[0]
		}
		_, err := r.ReadBytes(b)
		if err != nil {
			return err
		}
	}
	return readNonSignificantWhitespace(r)
}

func (p *ObjectContext) Close() os.Error {
	if p.writer == nil {
		return os.EOF
	}
	if p.first {
		v := []byte{JSON_LBRACE[0], JSON_RBRACE[0]}
		_, e := p.writer.Write(v)
		p.writer = nil
		return e
	}
	if p.onValue {
		if e := p.WriteNull(); e != nil {
			return e
		}
	}
	if !p.keyType.IsBaseType() {
		v := []byte{JSON_RBRACKET[0], JSON_RBRACE[0]}
		_, e := p.writer.Write(v)
		p.writer = nil
		return e
	}
	_, e := p.writer.Write(JSON_RBRACE)
	p.writer = nil
	return e
}

func (p *ObjectContext) ReadClose(r *bufio.Reader) os.Error {
	e := readNonSignificantWhitespace(r)
	if e != nil {
		return e
	}
	_, err := r.ReadBytes(JSON_RBRACE[0])
	if err != nil {
		return err
	}
	return readNonSignificantWhitespace(r)
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


/**
 * Push a new write context onto the stack.
 */
func (p *TSimpleJSONProtocol) pushWriteContext(c Context) {
	p.writeContextStack.PushFront(c)
	p.writeContext = c
}

/**
 * Pop the last write context off the stack
 */
func (p *TSimpleJSONProtocol) popWriteContext() Context {
	c := p.writeContextStack.Front()
	p.writeContextStack.Remove(c)
	p.writeContext = p.writeContextStack.Front().Value.(Context)
	return c.Value.(Context)
}

/**
 * Push a new read context onto the stack.
 */
func (p *TSimpleJSONProtocol) pushReadContext(c Context) {
	p.readContextStack.PushFront(c)
	p.readContext = c
}

/**
 * Pop the last read context off the stack
 */
func (p *TSimpleJSONProtocol) popReadContext() Context {
	c := p.readContextStack.Front()
	if c == nil {
		panic("Cannot read front of context stack")
		return nil
	}
	p.readContextStack.Remove(c)
	f := p.readContextStack.Front()
	if f == nil {
		panic("Cannot read front of context stack")
		return nil
	}
	p.readContext = f.Value.(Context)
	return c.Value.(Context)
}

func (p *TSimpleJSONProtocol) WriteMessageBegin(name string, typeId TMessageType, seqId int32) TProtocolException {
	if e := p.writeContext.WriteStartSubContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.pushWriteContext(NewListContextETypeWriter(TType(VOID), p.writer))
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
	if e := p.popWriteContext().Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteEndSubContext())
}

func (p *TSimpleJSONProtocol) WriteStructBegin(name string) TProtocolException {
	if e := p.writeContext.WriteStartSubContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.pushWriteContext(NewObjectContextKVTypeWriter(TType(STRING), TType(STOP), p.writer))
	return nil
}

func (p *TSimpleJSONProtocol) WriteStructEnd() TProtocolException {
	if e := p.popWriteContext().Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteEndSubContext())
}

func (p *TSimpleJSONProtocol) WriteFieldBegin(name string, typeId TType, id int16) TProtocolException {
	// Note that extra type information is omitted in JSON!
	return p.WriteString(name)
}

func (p *TSimpleJSONProtocol) WriteFieldEnd() TProtocolException { return nil }

func (p *TSimpleJSONProtocol) WriteFieldStop() TProtocolException { return nil }

func (p *TSimpleJSONProtocol) WriteMapBegin(keyType TType, valueType TType, size int) TProtocolException {
	if e := p.writeContext.WriteStartSubContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.pushWriteContext(NewObjectContextKVTypeWriter(keyType, valueType, p.writer))
	// No metadata!
	return nil
}

func (p *TSimpleJSONProtocol) WriteMapEnd() TProtocolException {
	if e := p.popWriteContext().Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteEndSubContext())
}

func (p *TSimpleJSONProtocol) WriteListBegin(elemType TType, size int) TProtocolException {
	if e := p.writeContext.WriteStartSubContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.pushWriteContext(NewListContextETypeWriter(elemType, p.writer))
	// No metadata!
	return nil
}

func (p *TSimpleJSONProtocol) WriteListEnd() TProtocolException {
	if e := p.popWriteContext().Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteEndSubContext())
}

func (p *TSimpleJSONProtocol) WriteSetBegin(elemType TType, size int) TProtocolException {
	if e := p.writeContext.WriteStartSubContext(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	p.pushWriteContext(NewListContextETypeWriter(elemType, p.writer))
	// No metadata!
	return nil
}

func (p *TSimpleJSONProtocol) WriteSetEnd() TProtocolException {
	if e := p.popWriteContext().Close(); e != nil {
		return NewTProtocolExceptionFromOsError(e)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteEndSubContext())
}

func (p *TSimpleJSONProtocol) WriteBool(b bool) TProtocolException {
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteBool(b))
}

func (p *TSimpleJSONProtocol) WriteByte(b byte) TProtocolException {
	return p.WriteI32(int32(b))
}

func (p *TSimpleJSONProtocol) WriteI16(v int16) TProtocolException {
	return p.WriteI32(int32(v))
}

func (p *TSimpleJSONProtocol) WriteI32(v int32) TProtocolException {
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteI32(v))
}

func (p *TSimpleJSONProtocol) WriteI64(v int64) TProtocolException {
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteI64(v))
}

func (p *TSimpleJSONProtocol) WriteDouble(v float64) TProtocolException {
	if math.IsNaN(v) {
		return p.WriteString(JSON_NAN)
	} else if math.IsInf(v, 1) {
		return p.WriteString(JSON_INFINITY)
	} else if math.IsInf(v, -1) {
		return p.WriteString(JSON_NEGATIVE_INFINITY)
	}
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteDouble(v))
}

func (p *TSimpleJSONProtocol) WriteString(v string) TProtocolException {
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteString(v))
}

func (p *TSimpleJSONProtocol) WriteBinary(v []byte) TProtocolException {
	// JSON library only takes in a string, 
	// not an arbitrary byte array, to ensure bytes are transmitted
	// efficiently we must convert this into a valid JSON string
	// therefore we use base64 encoding to avoid excessive escaping/quoting
	return NewTProtocolExceptionFromOsError(p.writeContext.WriteBinary(v))
	/*
	  p.writeContext.Write(p.writer)
	  if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil { return NewTProtocolExceptionFromOsError(e) }
	  writer := base64.NewEncoder(base64.StdEncoding, p.writer)
	  if _, e := writer.Write(v); e != nil { return NewTProtocolExceptionFromOsError(e) }
	  if e := writer.Close(); e != nil { return NewTProtocolExceptionFromOsError(e) }
	  if _, e := p.writer.Write(JSON_QUOTE_BYTES); e != nil { return NewTProtocolExceptionFromOsError(e) }
	  return nil
	*/
	/*
	  escape.WriteByte(JSON_QUOTE);
	  for i := 0; i<length; i++ {
	    c := v[i];
	    switch c {
	    case '"', '\\':
	      escape.WriteByte('\\');
	      escape.WriteByte(byte(c));
	    case '\b':
	      escape.WriteByte('\\');
	      escape.WriteByte('b');
	    case '\f':
	      escape.WriteByte('\\');
	      escape.WriteByte('f');
	    case '\n':
	      escape.WriteByte('\\');
	      escape.WriteByte('n');
	    case '\r':
	      escape.WriteByte('\\');
	      escape.WriteByte('r');
	    case '\t':
	      escape.WriteByte('\\');
	      escape.WriteByte('t');
	    default:
	       if c < ' ' || c > '~' {
	        // Control characters! According to JSON RFC u0020 (space)
	        h := strconv.Uitob(uint(c), 16);
	        escape.WriteByte('\\');
	        escape.WriteByte('u');
	        for j := 4; j > len(h); j-- {
	          escape.WriteByte('0');
	        }
	        io.Copyn(escape, strings.NewReader(h), int64(len(h)));
	      } else {
	        escape.WriteByte(byte(c));
	      }
	    }
	  }
	  escape.WriteByte(JSON_QUOTE);
	  return p.writeStringData(escape.String());
	*/
}

/**
 * Reading methods.
 */

func (p *TSimpleJSONProtocol) ReadMessageBegin() (name string, typeId TMessageType, seqId int32, err TProtocolException) {
	// TODO(mcslee): implement
	e := p.readContext.Read(p.reader)
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	_, e = p.reader.ReadString(JSON_LBRACKET[0])
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	cxt := NewListContext()
	p.pushReadContext(cxt)
	e = cxt.Populate(p.reader)
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	name, err = p.ReadString()
	if err != nil {
		return
	}
	t, err := p.ReadByte()
	typeId = TMessageType(t)
	if err != nil {
		return
	}
	seqId, err = p.ReadI32()
	return
}

func (p *TSimpleJSONProtocol) ReadMessageEnd() TProtocolException {
	err := p.popReadContext().ReadClose(p.reader)
	if err != nil {
		return NewTProtocolExceptionFromOsError(err)
	}
	_, err = p.reader.ReadString(JSON_RBRACKET[0])
	return NewTProtocolExceptionFromOsError(err)
}

func (p *TSimpleJSONProtocol) ReadStructBegin() (name string, err TProtocolException) {
	e := p.readContext.Read(p.reader)
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	_, e = p.reader.ReadBytes(JSON_LBRACE[0])
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	cxt := NewObjectContext()
	p.pushReadContext(cxt)
	e = cxt.Populate(p.reader)
	if e != nil {
		err = NewTProtocolExceptionFromOsError(e)
		return
	}
	return
}

func (p *TSimpleJSONProtocol) ReadStructEnd() TProtocolException {
	e := p.popReadContext().ReadClose(p.reader)
	return NewTProtocolExceptionFromOsError(e)
}

func (p *TSimpleJSONProtocol) ReadFieldBegin() (string, TType, int16, TProtocolException) {

	return "", STOP, 0, nil
}

func (p *TSimpleJSONProtocol) ReadFieldEnd() TProtocolException { return nil }

func (p *TSimpleJSONProtocol) ReadMapBegin() (keyType TType, valueType TType, size int, e TProtocolException) {
	err := p.readContext.Read(p.reader)
	if err != nil {
		return STOP, STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	var b byte
	b, err = p.reader.ReadByte()
	if err != nil {
		return STOP, STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	if b == JSON_LBRACE[0] {
		cxt := NewObjectContext()
		p.pushReadContext(cxt)
		err = cxt.Populate(p.reader)
		return cxt.KeyType(), cxt.ValueType(), cxt.Len(), NewTProtocolExceptionFromOsError(err)
	}
	return STOP, STOP, 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected 'null' or '{', received '%q'", b))
}

func (p *TSimpleJSONProtocol) ReadMapEnd() TProtocolException {
	thecxt, ok := p.readContext.(*ObjectContext)
	if !ok {
		return NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected current context should be ObjectContext, but was %T", thecxt))
	}
	cxt := p.popReadContext()
	if cxt == nil {
		return NewTProtocolException(INVALID_DATA, "Popped a nil context.  Something's wrong in your program")
	}
	return nil
}

func (p *TSimpleJSONProtocol) ReadListBegin() (elemType TType, size int, e TProtocolException) {
	err := p.readContext.Read(p.reader)
	if err != nil {
		return STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	var b byte
	b, err = p.reader.ReadByte()
	if err != nil {
		return STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	if b == JSON_LBRACKET[0] {
		cxt := NewListContext()
		p.pushReadContext(cxt)
		err = cxt.Populate(p.reader)
		return cxt.TType(), cxt.Len(), NewTProtocolExceptionFromOsError(err)
	}
	return STOP, 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected 'null' or '[', received '%q'", b))
}

func (p *TSimpleJSONProtocol) ReadListEnd() TProtocolException {
	thecxt, ok := p.readContext.(*ListContext)
	if !ok {
		return NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected current context should be ListContext, but was %T", thecxt))
	}
	cxt := p.popReadContext()
	if cxt == nil {
		return NewTProtocolException(INVALID_DATA, "Popped a nil context.  Something's wrong in your program")
	}
	return nil
}

func (p *TSimpleJSONProtocol) ReadSetBegin() (elemType TType, size int, e TProtocolException) {
	err := p.readContext.Read(p.reader)
	if err != nil {
		return STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	isNull, e1 := readIfNull(p.reader)
	if e1 != nil {
		return STOP, 0, e1
	}
	if isNull {
		p.pushReadContext(NewContext())
		return STOP, 0, nil
	}
	var b byte
	b, err = p.reader.ReadByte()
	if err != nil {
		return STOP, 0, NewTProtocolExceptionFromOsError(err)
	}
	if b == JSON_LBRACKET[0] {
		cxt := NewListContext()
		p.pushReadContext(cxt)
		err = cxt.Populate(p.reader)
		return cxt.TType(), cxt.Len(), NewTProtocolExceptionFromOsError(err)
	}
	return STOP, 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected 'null' or '[', received '%q'", b))
}

func (p *TSimpleJSONProtocol) ReadSetEnd() TProtocolException {
	e := p.popReadContext().ReadClose(p.reader)
	return NewTProtocolExceptionFromOsError(e)
}

func (p *TSimpleJSONProtocol) ReadBool() (bool, TProtocolException) {
	l, ok := p.readContext.(*ListContext)
	if ok {
		elem := l.Next()
		if elem == nil {
			return false, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected bool in list, found nil element in list"))
		}
		v := elem.Value
		b, ok := v.(bool)
		if ok {
			return b, nil
		} else {
			return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected bool in list, found %q", v))
		}
	}
	s, ok := p.readContext.(*ObjectContext)
	if ok {
		v := s.Next().Value
		b, ok := v.(bool)
		if ok {
			return b, nil
		} else {
			return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected bool in struct, found %q", v))
		}
	}
	err := p.readContext.Read(p.reader)
	if err != nil {
		return false, NewTProtocolExceptionFromOsError(err)
	}
	c, e := p.reader.ReadByte()
	if e != nil {
		return false, NewTProtocolExceptionFromOsError(e)
	}
	var b []byte
	var value bool
	switch c {
	case JSON_TRUE[0]:
		b = JSON_TRUE
		value = true
	case JSON_FALSE[0]:
		b = JSON_FALSE
		value = false
	default:
		return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse boolean expression starting with %q", c))
	}
	bytes := make([]byte, len(b))
	bytes[0] = c
	_, e = p.reader.Read(bytes[1:])
	if e != nil {
		return false, NewTProtocolExceptionFromOsError(e)
	}
	for i := 1; i < len(bytes); i++ {
		if b[i] != bytes[i] {
			return false, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse boolean expression starting with %q", c))
		}
	}
	return value, nil
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
	l, ok := p.readContext.(*ListContext)
	if ok {
		elem := l.Next()
		if elem == nil {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected int64 in list, found nil element in list"))
		}
		v := elem.Value
		n, ok := v.(Numeric)
		if ok {
			return n.Int64(), nil
		} else {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected int64 in list, found %q", v))
		}
	}
	s, ok := p.readContext.(*ObjectContext)
	if ok {
		v := s.Next().Value
		n, ok := v.(Numeric)
		if ok {
			return n.Int64(), nil
		} else {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected int64 in struct, found %q", v))
		}
	}
	p.readContext.Read(p.reader)
	v, err := readNumeric(p.reader)
	return v.Int64(), NewTProtocolExceptionFromOsError(err)
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
				return NEGATIVE_INFINITY, nil
			} else {
				return NUMERIC_NULL, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Unable to parse number starting with character '%c' due to existing buffer %s", c, buf.String()))
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

func (p *TSimpleJSONProtocol) ReadDouble() (float64, TProtocolException) {
	l, ok := p.readContext.(*ListContext)
	if ok {
		elem := l.Next()
		if elem == nil {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected double in list, found nil element in list"))
		}
		v := elem.Value
		n, ok := v.(Numeric)
		if ok {
			return n.Float64(), nil
		} else {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected double in list, found %q", v))
		}
	}
	s, ok := p.readContext.(*ObjectContext)
	if ok {
		v := s.Next().Value
		n, ok := v.(Numeric)
		if ok {
			return n.Float64(), nil
		} else {
			return 0, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected double in struct, found %q", v))
		}
	}
	p.readContext.Read(p.reader)
	v, err := readNumeric(p.reader)
	return v.Float64(), err
}

func (p *TSimpleJSONProtocol) ReadString() (string, TProtocolException) {
	l, ok := p.readContext.(*ListContext)
	if ok {
		elem := l.Next()
		if elem == nil {
			return "", NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected string in list, found nil element in list"))
		}
		v := elem.Value
		n, ok := v.(string)
		if ok {
			return n, nil
		} else {
			return "", NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected string in list, found %q", v))
		}
	}
	s, ok := p.readContext.(*ObjectContext)
	if ok {
		v := s.Next().Value
		n, ok := v.(string)
		if ok {
			return n, nil
		} else {
			return "", NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected string in struct, found %q", v))
		}
	}
	p.readContext.Read(p.reader)
	_, err := p.reader.ReadString(JSON_QUOTE)
	if err != nil {
		return "", NewTProtocolExceptionFromOsError(err)
	}
	str, e := ReadStringBody(p.reader)
	return str, e
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
	l, ok := p.readContext.(*ListContext)
	if ok {
		elem := l.Next()
		if elem == nil {
			return nil, NewTProtocolException(INVALID_DATA, fmt.Sprint("Expected []byte in list, found nil element in list"))
		}
		v := elem.Value
		n, ok := v.([]byte)
		if ok {
			return n, nil
		} else {
			return nil, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected []byte in list, found %q", v))
		}
	}
	s, ok := p.readContext.(*ObjectContext)
	if ok {
		v := s.Next().Value
		n, ok := v.([]byte)
		if ok {
			return n, nil
		} else {
			return nil, NewTProtocolException(INVALID_DATA, fmt.Sprintf("Expected []byte in struct, found %q", v))
		}
	}
	// JSON library only takes in a string, 
	// not an arbitrary byte array, to ensure bytes are transmitted
	// efficiently we must convert this into a valid JSON string
	// but since it can contain invalid escape sequences for
	// strings, we must do it manually
	p.readContext.Read(p.reader)
	_, err := p.reader.ReadString(JSON_QUOTE)
	if err != nil {
		return []byte{}, NewTProtocolExceptionFromOsError(err)
	}
	return readBase64EncodedBody(p.reader)
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
