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
	"sort"
)

/**
 * Helper class that encapsulates field metadata.
 *
 */
type TField interface {
	Name() string
	TypeId() TType
	Id() int
	String() string
}

type tField struct {
	name   string
	typeId TType
	id     int
}

func NewTFieldDefault() TField {
	return ANONYMOUS_FIELD
}

func NewTField(n string, t TType, i int) TField {
	return &tField{name: n, typeId: t, id: i}
}

func (p *tField) Name() string {
	if p == nil {
		return ""
	}
	return p.name
}

func (p *tField) TypeId() TType {
	if p == nil {
		return TType(VOID)
	}
	return p.typeId
}

func (p *tField) Id() int {
	if p == nil {
		return -1
	}
	return p.id
}

func (p *tField) String() string {
	if p == nil {
		return "<nil>"
	}
	return "<TField name:'" + p.name + "' type:" + string(p.typeId) + " field-id:" + string(p.id) + ">"
}

var ANONYMOUS_FIELD TField

type tFieldArray []TField

func (p tFieldArray) Len() int {
	return len(p)
}

func (p tFieldArray) Less(i, j int) bool {
	return p[i].Id() < p[j].Id()
}

func (p tFieldArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type TFieldContainer interface {
	FieldNameFromFieldId(id int) string
	FieldIdFromFieldName(name string) int
	FieldFromFieldId(id int) TField
	FieldFromFieldName(name string) TField
	At(i int) TField
}

type tFieldContainer struct {
	fields         []TField
	nameToFieldMap map[string]TField
	idToFieldMap   map[int]TField
}

func NewTFieldContainer(fields []TField) TFieldContainer {
	sortedFields := make([]TField, len(fields))
	nameToFieldMap := make(map[string]TField)
	idToFieldMap := make(map[int]TField)
	for i, field := range fields {
		sortedFields[i] = field
		idToFieldMap[field.Id()] = field
		if field.Name() != "" {
			nameToFieldMap[field.Name()] = field
		}
	}
	sort.Sort(tFieldArray(sortedFields))
	return &tFieldContainer{
		fields:         fields,
		nameToFieldMap: nameToFieldMap,
		idToFieldMap:   idToFieldMap,
	}
}

func (p *tFieldContainer) FieldNameFromFieldId(id int) string {
	if field, ok := p.idToFieldMap[id]; ok {
		return field.Name()
	}
	return ""
}

func (p *tFieldContainer) FieldIdFromFieldName(name string) int {
	if field, ok := p.nameToFieldMap[name]; ok {
		return field.Id()
	}
	return -1
}

func (p *tFieldContainer) FieldFromFieldId(id int) TField {
	if field, ok := p.idToFieldMap[id]; ok {
		return field
	}
	return ANONYMOUS_FIELD
}

func (p *tFieldContainer) FieldFromFieldName(name string) TField {
	if field, ok := p.nameToFieldMap[name]; ok {
		return field
	}
	return ANONYMOUS_FIELD
}

func (p *tFieldContainer) Len() int {
	return len(p.fields)
}

func (p *tFieldContainer) At(i int) TField {
	return p.FieldFromFieldId(i)
}

func (p *tFieldContainer) iterate(c chan<- TField) {
	for _, v := range p.fields {
		c <- v
	}
	close(c)
}

func init() {
	ANONYMOUS_FIELD = NewTField("", STOP, 0)
}
