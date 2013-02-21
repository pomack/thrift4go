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

/**
 * Type constants in the Thrift protocol.
 */
type TType int8

const (
	STOP   = 0
	VOID   = 1
	BOOL   = 2
	BYTE   = 3
	DOUBLE = 4
	I16    = 6
	I32    = 8
	I64    = 10
	STRING = 11
	STRUCT = 12
	MAP    = 13
	SET    = 14
	LIST   = 15
	ENUM   = 16

//	GENERIC = 127
)

func (p TType) ThriftTypeId() int8 {
	return int8(p)
}

func (p TType) String() string {
	switch p {
	case STOP:
		return "STOP"
	case VOID:
		return "VOID"
	case BOOL:
		return "BOOL"
	case BYTE:
		return "BYTE"
	case DOUBLE:
		return "DOUBLE"
	case I16:
		return "I16"
	case I32:
		return "I32"
	case I64:
		return "I64"
	case STRING:
		return "STRING"
	case STRUCT:
		return "STRUCT"
	case MAP:
		return "MAP"
	case SET:
		return "SET"
	case LIST:
		return "LIST"
	case ENUM:
		return "ENUM"
	}
	return "Unknown"
}
