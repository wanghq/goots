// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// like python's type, but it's simple
package otstype

import (
	"errors"
	"fmt"
)

// tuple
type Tuple []struct {
	K interface{}
	V interface{}
}

type TupleString struct {
	K string
	V interface{}
}

func (t *TupleString) GetKey() string {
	if t == nil {
		return ""
	}
	return t.K
}

func (t *TupleString) SetKey(k string) {
	if t != nil {
		(*t).K = k
	}
}

func (t *TupleString) GetValue() interface{} {
	if t == nil {
		return nil
	}
	return t.V
}

func (t *TupleString) SetValue(v interface{}) {
	if t != nil {
		(*t).V = v
	}
}

// for OTS protobuf adapter
func (t TupleString) GetName() string {
	return t.K
}

func (t TupleString) GetType() interface{} {
	return t.V
}

// dict
type Dict map[interface{}]interface{}

func (d *Dict) Add(k, v interface{}) {
	(*d)[k] = v
}

func (d *Dict) Get(k interface{}) (v interface{}, err error) {
	v, ok := (*d)[k]
	if ok {
		return v, nil
	}

	return v, errors.New("key not found")
}

type DictString map[string]interface{}

func (d DictString) String() string {
	result := ""
	for k, v := range d {
		result = result + fmt.Sprintf("%s:%v\n", k, v)
	}

	return result
}

// delete key
func (d DictString) Del(key string) {
	if d != nil {
		delete(d, key)
	}
}

// get value by key
func (d DictString) Get(key string) interface{} {
	if d == nil || key == "" {
		return nil
	}

	return d[key]
}

// set key and value to map
func (d DictString) Set(key string, value interface{}) {
	if d != nil {
		d[key] = value
	}
}
