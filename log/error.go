// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// errors for ots2
package log

import (
	"fmt"
)

// 如果用户不希望使用panic模式，则设置此为false
var OTSErrorPanicMode bool = true // 默认开启panic模式

type OTSError struct {
	ClientError  *OTSClientError
	ServiceError *OTSServiceError
}

func LoggerInit() error {
	// TODO:
	// open log file

	return nil
}

func (o OTSError) Set(format string, a ...interface{}) (e error) {
	defer func() {
		if OTSErrorPanicMode {
			panic(e)
		}
	}()
	e = fmt.Errorf(format, a...)
	return e
}

func (o OTSError) Log(enable bool, format string, a ...interface{}) (e error) {
	e = fmt.Errorf(format, a...)

	// TODO:
	// log to file
	if enable {

	}

	return e
}

func (o *OTSError) Error() string {
	var client_message string
	var service_message string

	if o.ClientError == nil {
		client_message = "[C]-None"
	} else {
		client_message = "[C]-" + o.ClientError.Message
	}

	if o.ServiceError == nil {
		service_message = "[S]-None"
	} else {
		service_message = "[S]-" + o.ServiceError.Code + " @ " + o.ServiceError.Message
	}

	return client_message + " <--> " + service_message
}

func (o *OTSError) String() string {
	var client_message string
	var service_message string

	if o.ClientError == nil {
		client_message = "[C]-None"
	} else {
		client_message = o.ClientError.Message
	}

	if o.ServiceError == nil {
		service_message = "[S]-None"
	} else {
		service_message = "[S]-" + o.ServiceError.Code + " @ " + o.ServiceError.Message
	}

	return client_message + " <--> " + service_message
}

func (o *OTSError) SetClientError(client_err *OTSClientError) *OTSError {
	o.ClientError = client_err

	return o
}

func (o *OTSError) SetClientMessage(format string, a ...interface{}) *OTSError {
	o.ClientError = new(OTSClientError)
	o.ClientError.SetErrorMessage(format, a...)

	return o
}

func (o *OTSError) SetServiceError(service_err *OTSServiceError) *OTSError {
	o.ServiceError = service_err

	return o
}

func (o *OTSError) SetServiceMessage(format string, a ...interface{}) *OTSError {
	o.ServiceError = new(OTSServiceError)
	o.ServiceError.SetErrorMessage(format, a...)

	return o
}

type OTSClientError struct {
	Message    string
	HttpStatus string
}

func (o OTSClientError) Set(format string, a ...interface{}) (e error) {
	defer func() {
		if OTSErrorPanicMode {
			panic(e)
		}
	}()
	e = fmt.Errorf(format, a...)
	return e
}

func (o OTSClientError) Log(enable bool, format string, a ...interface{}) (e error) {
	e = fmt.Errorf(format, a...)

	// TODO:
	// log to file
	if enable {

	}

	return e
}

func (o *OTSClientError) Error() string {
	return "[C]-" + o.Message
}

func (o *OTSClientError) String() string {
	return "[C]-" + o.Message
}

func (o *OTSClientError) SetHttpStatus(status string) *OTSClientError {
	o.HttpStatus = status

	return o
}

func (o *OTSClientError) GetHttpStatus() string {
	return o.HttpStatus
}

func (o *OTSClientError) SetErrorMessage(format string, a ...interface{}) *OTSClientError {
	o.Message = fmt.Sprintf(format, a...)

	return o
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
		if OTSErrorPanicMode {
			panic(e)
		}
	}()
	e = fmt.Errorf(format, a...)
	return e
}

func (o OTSServiceError) Log(enable bool, format string, a ...interface{}) (e error) {
	e = fmt.Errorf(format, a...)

	// TODO:
	// log to file
	if enable {

	}

	return e
}

func (o *OTSServiceError) Error() string {
	return fmt.Sprintf("[S]-ErrorCode: %s, ErrorMessage: %s, RequestID: %s",
		o.Code, o.Message, o.RequestId)
}

func (o *OTSServiceError) String() string {
	return fmt.Sprintf("[S]-ErrorCode: %s, ErrorMessage: %s, RequestID: %s",
		o.Code, o.Message, o.RequestId)
}

func (o *OTSServiceError) SetHttpStatus(status string) *OTSServiceError {
	o.HttpStatus = status

	return o
}
func (o *OTSServiceError) GetHttpStatus() string {
	return o.HttpStatus
}

func (o *OTSServiceError) SetErrorCode(code string) *OTSServiceError {
	o.Code = code

	return o
}

func (o *OTSServiceError) GetErrorCode() string {
	return o.Code
}

func (o *OTSServiceError) SetErrorMessage(format string, a ...interface{}) *OTSServiceError {
	o.Message = fmt.Sprintf(format, a...)

	return o
}

func (o *OTSServiceError) GetErrorMessage() string {
	return o.Message
}

func (o *OTSServiceError) SetRequestId(request_id string) *OTSServiceError {
	o.RequestId = request_id

	return o
}

func (o *OTSServiceError) GetRequestId() string {
	return o.RequestId
}
