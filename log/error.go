// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// errors for ots2
package log

import (
	"fmt"
)

type OTSError struct {
	ClientError  OTSClientError
	ServiceError OTSServiceError
}

type OTSClientError struct {
	Message    string
	HttpStatus string
}

func (o OTSClientError) Set(format string, a ...interface{}) (e error) {
	defer func() {
		panic(e)
	}()
	e = fmt.Errorf(format, a)
	return e
}

func (o *OTSClientError) Error() string {
	return o.Message
}

func (o *OTSClientError) String() string {
	return o.Message
}

func (o *OTSClientError) GetHttpStatus() string {
	return o.HttpStatus
}

func (o *OTSClientError) GetErrorMessage() string {
	return o.Message
}

type OTSServiceError struct {
	HttpStatus string
	Code       string
	Message    string
	RequestId  string
}

func (o OTSServiceError) Set(format string, a ...interface{}) (e error) {
	defer func() {
		panic(e)
	}()
	e = fmt.Errorf(format, a)
	return e
}

func (o *OTSServiceError) Error() string {
	return o.Message
}

func (o *OTSServiceError) String() string {
	return fmt.Sprintf("ErrorCode: %s, ErrorMessage: %s, RequestID: %s",
		o.Code, o.Message, o.RequestId)
}

func (o *OTSServiceError) GetHttpStatus() string {
	return o.HttpStatus
}

func (o *OTSServiceError) GetErrorCode() string {
	return o.Code
}

func (o *OTSServiceError) GetErrorMessage() string {
	return o.Message
}

func (o *OTSServiceError) GetRequestId() string {
	return o.RequestId
}
