// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Implement call func by name of function
package coder

import (
	"fmt"
	"reflect"
)

type Funcmap map[string]reflect.Value

// Create a new function map
func NewFuncmap() Funcmap {
	return make(Funcmap, 2)
}

// Bind functon by name
func (f Funcmap) Bind(name string, fn interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s - bind function failed", name)
			fmt.Println(err) // tell developer error message
		}
	}()

	v := reflect.ValueOf(fn)
	// panic if the type's Kind is not Func.
	v.Type().NumIn()
	v.Type().NumOut()
	f[name] = v
	return nil
}

// Check function map whether has name method specified
func (f Funcmap) Has(name string) bool {
	_, ok := f[name]
	return ok
}

// Call function by name
func (f Funcmap) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := f[name]; !ok {
		err = fmt.Errorf("%s - fucntion can not found", name)
		return nil, err
	}

	if len(params) != f[name].Type().NumIn() {
		err = fmt.Errorf("%s - input params can not adapt", name)
		return nil, err
	}

	in := make([]reflect.Value, len(params))
	for k, v := range params {
		in[k] = reflect.ValueOf(v)
	}
	result = f[name].Call(in)
	return result, nil
}
