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
	"sort"
)

/**
 * Helper class that encapsulates field metadata.
 *
 */
type TField interface {
	Name() string
	TypeId() TType
	Id() int16
	String() string
	CompareTo(other interface{}) (int, bool)
	Equals(other interface{}) bool
}

type tField struct {
	name   string
	typeId TType
	id     int16
}

func NewTFieldDefault() TField {
	return ANONYMOUS_FIELD
}

func NewTField(n string, t TType, i int16) TField {
	return &tField{name: n, typeId: t, id: i}
}

func (p *tField) Name() string {
	return p.name
}

func (p *tField) TypeId() TType {
	return p.typeId
}

func (p *tField) Id() int16 {
	return p.id
}

func (p *tField) String() string {
	return "<TField name:'" + p.name + "' type:" + string(p.typeId) + " field-id:" + string(p.id) + ">"
}

func (p *tField) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	if data, ok := other.(TField); ok {
		if p.Id() != data.Id() {
			return CompareInt16(p.Id(), data.Id()), true
		}
		if p.TypeId() != data.TypeId() {
			return CompareByte(byte(p.TypeId()), byte(data.TypeId())), true
		}
		return CompareString(p.Name(), data.Name()), true
	}
	return 0, false
}

func (p *tField) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	if data, ok := other.(TField); ok {
		return p.TypeId() == data.TypeId() && p.Id() == data.Id()
	}
	return false
}

var ANONYMOUS_FIELD TField

type TFieldContainer interface {
	TContainer
	Less(i, j int) bool
	Swap(i, j int)
	At(i int16) TField
	Iter() <-chan TField
}

type tFieldContainer struct {
	fields   []TField
	isSorted bool
}

func NewTFieldContainer(fields []TField) TFieldContainer {
	return &tFieldContainer{fields: fields, isSorted: false}
}

func (p *tFieldContainer) Len() int {
	return len(p.fields)
}

func (p *tFieldContainer) Sort() {
	if !p.isSorted {
		sort.Sort(p)
		p.isSorted = true
	}
}

func (p *tFieldContainer) At(i int16) TField {
	p.Sort()
	for _, field := range p.fields {
		if i == field.Id() {
			return field
		}
		if i > field.Id() {
			return nil
		}
	}
	return nil
}

func (p *tFieldContainer) Contains(data interface{}) bool {
	if i, ok := data.(int); ok {
		for _, field := range p.fields {
			if field.Id() == int16(i) {
				return true
			}
		}
	} else if i, ok := data.(int16); ok {
		for _, field := range p.fields {
			if field.Id() == i {
				return true
			}
		}
	} else if s, ok := data.(string); ok {
		for _, field := range p.fields {
			if field.Name() == s {
				return true
			}
		}
	} else if f, ok := data.(TField); ok {
		for _, field := range p.fields {
			if field.Equals(f) {
				return true
			}
		}
	}
	return false
}

func (p *tFieldContainer) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	if data, ok := other.(TFieldContainer); ok {
		if p.Len() != data.Len() {
			return false
		}
		for _, field := range p.fields {
			if !data.Contains(field) {
				return false
			}
		}
		return true
	}
	return false
}

func (p *tFieldContainer) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	if data, ok := other.(TFieldContainer); ok {
		cont, ok2 := data.(*tFieldContainer)
		if ok2 && p == cont {
			return 0, true
		}
		if cmp := CompareInt(p.Len(), data.Len()); cmp != 0 {
			return cmp, true
		}
		for _, field := range p.fields {
			if cmp, ok3 := field.CompareTo(data.At(field.Id())); !ok3 || cmp != 0 {
				return cmp, ok3
			}
		}
		return 0, true
	}
	return 0, false
}

func (p *tFieldContainer) Less(i, j int) bool {
	return p.fields[i].Id() < p.fields[j].Id()
}

func (p *tFieldContainer) Swap(i, j int) {
	p.fields[i], p.fields[j] = p.fields[j], p.fields[i]
}

func (p *tFieldContainer) Iter() <-chan TField {
	c := make(chan TField)
	go p.iterate(c)
	return c
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
