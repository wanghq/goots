// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Base type for golang
//
// int，rune
// int8 ,int16 ,int32 ,int64
// byte ,uint8 ,uint16 ,uint32 ,uint64
// float32 ，float64
// bool
// string
// complex128，complex64
package otstype

import (
	"math"
)

func NewInt(v int) *int {
	return &v
}

func NewRune(v rune) *rune {
	return &v
}

func NewInt8(v int8) *int8 {
	return &v
}

func NewInt16(v int16) *int16 {
	return &v
}

func NewInt32(v int32) *int32 {
	return &v
}

func NewInt64(v int64) *int64 {
	return &v
}

func Newbyte(v byte) *byte {
	return &v
}

func NewUint8(v uint8) *uint8 {
	return &v
}

func NewUint16(v uint16) *uint16 {
	return &v
}

func NewUint32(v uint32) *uint32 {
	return &v
}

func NewUint64(v uint64) *uint64 {
	return &v
}

func NewFloat32(v float32) *float32 {
	return &v
}

func NewFloat64(v float64) *float64 {
	return &v
}

func NewBool(v bool) *bool {
	return &v
}

func NewString(v string) *string {
	return &v
}

func NewComplex64(v complex64) *complex64 {
	return &v
}

func NewComplex128(v complex128) *complex128 {
	return &v
}

func GetInt8Min() int8 {
	return math.MinInt8
}

func GetInt8Max() int8 {
	return math.MaxInt8
}

func GetInt16Min() int16 {
	return math.MinInt16
}

func GetInt16Max() int16 {
	return math.MaxInt16
}

func GetInt32Min() int32 {
	return math.MinInt32
}

func GetInt32Max() int32 {
	return math.MaxInt32
}

func GetUint8Min() uint8 {
	return 0
}

func GetUint8Max() uint8 {
	return math.MaxUint8
}

func GetUint16Min() uint16 {
	return 0
}

func GetUint16Max() uint16 {
	return math.MaxUint16
}

func GetUint32Min() uint32 {
	return 0
}

func GetUint32Max() uint32 {
	return math.MaxUint32
}

func GetFloat32Mix() float32 {
	return math.SmallestNonzeroFloat32
}

func GetFloat32Max() float32 {
	return math.MaxFloat32
}

func GetFloat64Mix() float64 {
	return math.SmallestNonzeroFloat64
}

func GetFloat64Max() float64 {
	return math.MaxFloat64
}
