// Copyright 2016 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// retry for ots2
package goots

import (
	"math"
	"math/rand"
	"net"
	"time"
)

// RetryPolicy 是重试策略的接口，包含2个未实现的方法和它们的参数列表。要实现一个重试策略，
// 继承这个类并实现它的2个方法。
type RetryPolicyInterface interface {
	GetRetryDelay(retry_times int, exception *OTSServiceError, api_name string) float64
	ShouldRetry(retry_times int, exception *OTSServiceError, api_name string) bool
}

// 默认重试策略
// 最大重试次数为3，最大重试间隔为2000毫秒，对流控类错误以及读操作相关的服务端内部错误进行了重试。
//
// Methods defined here:
//
// GetRetryDelay(retry_times int, exception *OTSServiceError, api_name string) float64
//
// ShouldRetry(retry_times int, exception *OTSServiceError, api_name string) bool
//
// ----------------------------------------------------------------------
// Data and other attributes defined here:
//
// MaxRetryDelay = 2
//
// MaxRetryTimes = 3
//
// ScaleFactor = 2
//
// ServerThrottlingExceptionDelayFactor = 0.5
//
// StabilityExceptionDelayFactor = 0.2
var OTSDefaultRetryPolicy DefaultRetryPolicy

// 不进行任何重试的重试策略
var OTSNoRetryPolicy NoRetryPolicy

// 没有延时的重试策略
// 最大重试次数为3
var OTSNoDelayRetryPolicy NoDelayRetryPolicy

func init() {
	OTSDefaultRetryPolicy.MaxRetryTimes = 6
	OTSDefaultRetryPolicy.MaxRetryDelay = 2
	OTSDefaultRetryPolicy.ScaleFactor = 2
	OTSDefaultRetryPolicy.ServerThrottlingExceptionDelayFactor = 0.5
	OTSDefaultRetryPolicy.StabilityExceptionDelayFactor = 0.2

	OTSNoDelayRetryPolicy.MaxRetryTimes = 3
}

func should_retry_no_matter_which_api(exception *OTSServiceError) bool {
	if exception != nil {
		error_code := exception.Code
		error_message := exception.Message

		if error_code == "OTSRowOperationConflict" ||
			error_code == "OTSNotEnoughCapacityUnit" ||
			error_code == "OTSTableNotReady" ||
			error_code == "OTSPartitionUnavailable" ||
			error_code == "OTSServerBusy" ||
			error_code == "OTSOperationThrottled" {
			return true
		}

		if error_code == "OTSQuotaExhausted" && error_message == "Too frequent table operations." {
			return true
		}
	}

	return false
}

func is_repeatable_api(api_name string) bool {
	if api_name == "ListTable" ||
		api_name == "DescribeTable" ||
		api_name == "GetRow" ||
		api_name == "BatchGetRow" ||
		api_name == "GetRange" ||
		api_name == "DescrieStream" ||
		api_name == "GetShardIterator" ||
		api_name == "GetStreamRecord" ||
		api_name == "ListStream" {
		return true
	}
	return false
}

func should_retry_when_api_repeatable(retry_times int, exception *OTSServiceError, api_name string) bool {
	if exception != nil {
		if exception.Err != nil {
			if _, ok := exception.Err.(net.Error); ok {
				return true
			} else if exception.Err == ErrNonResponseBody || exception.Err == ErrReadResponse {
				return true
			}
		}

		error_code := exception.Code
		//error_message := exception.Message
		http_status := exception.HttpStatus

		if error_code == "OTSTimeout" ||
			error_code == "OTSInternalServerError" ||
			error_code == "OTSServerUnavailable" {
			return true
		}

		if http_status == 500 || http_status == 502 || http_status == 503 {
			return true
		}
	}

	// TODO handle network error & timeout
	return false
}

func is_server_throttling_exception(exception *OTSServiceError) bool {
	if exception != nil {
		error_code := exception.Code
		error_message := exception.Message

		if error_code == "OTSServerBusy" ||
			error_code == "OTSNotEnoughCapacityUnit" ||
			error_code == "OTSOperationThrottled" {
			return true
		}

		if error_code == "OTSQuotaExhausted" && error_message == "Too frequent table operations." {
			return true
		}
	}

	return false
}

// 默认重试策略
// 最大重试次数为3，最大重试间隔为2000毫秒，对流控类错误以及读操作相关的服务端内部错误进行了重试。
type DefaultRetryPolicy struct {
	// 最大重试次数
	MaxRetryTimes int

	// 最大重试间隔，单位为秒
	MaxRetryDelay float64

	// 每次重试间隔的递增倍数
	ScaleFactor float64

	// 两种错误的起始重试间隔，单位为秒
	ServerThrottlingExceptionDelayFactor float64
	StabilityExceptionDelayFactor        float64
}

func (self DefaultRetryPolicy) _max_retry_time_reached(retry_times int, exception *OTSServiceError, api_name string) bool {
	return retry_times >= self.MaxRetryTimes
}

func (self DefaultRetryPolicy) _can_retry(retry_times int, exception *OTSServiceError, api_name string) bool {
	if should_retry_no_matter_which_api(exception) {
		return true
	}

	if is_repeatable_api(api_name) && should_retry_when_api_repeatable(retry_times, exception, api_name) {
		return true
	}

	return false
}

func (self DefaultRetryPolicy) GetRetryDelay(retry_times int, exception *OTSServiceError, api_name string) float64 {
	var delay_factor float64
	if is_server_throttling_exception(exception) {
		delay_factor = self.ServerThrottlingExceptionDelayFactor
	} else {
		delay_factor = self.StabilityExceptionDelayFactor
	}

	delay_limit := delay_factor * math.Pow(self.ScaleFactor, float64(retry_times))

	if delay_limit >= self.MaxRetryDelay {
		delay_limit = self.MaxRetryDelay
	}

	real_delay := delay_limit*0.5 + delay_limit*0.5*rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
	return real_delay
}

func (self DefaultRetryPolicy) ShouldRetry(retry_times int, exception *OTSServiceError, api_name string) bool {
	if self._max_retry_time_reached(retry_times, exception, api_name) {
		return false
	}

	if self._can_retry(retry_times, exception, api_name) {
		return true
	}

	return false
}

// 不进行任何重试的重试策略
type NoRetryPolicy struct {
}

func (self NoRetryPolicy) GetRetryDelay(retry_times int, exception *OTSServiceError, api_name string) float64 {
	return 0
}

func (self NoRetryPolicy) ShouldRetry(retry_times int, exception *OTSServiceError, api_name string) bool {
	return false
}

// 没有延时的重试策略
type NoDelayRetryPolicy struct {
	// 最大重试次数
	MaxRetryTimes int
}

func (self NoDelayRetryPolicy) _max_retry_time_reached(retry_times int, exception *OTSServiceError, api_name string) bool {
	return retry_times >= self.MaxRetryTimes
}

func (self NoDelayRetryPolicy) _can_retry(retry_times int, exception *OTSServiceError, api_name string) bool {
	if should_retry_no_matter_which_api(exception) {
		return true
	}

	if is_repeatable_api(api_name) && should_retry_when_api_repeatable(retry_times, exception, api_name) {
		return true
	}

	return false
}

func (sels NoDelayRetryPolicy) GetRetryDelay(retry_times int, exception *OTSServiceError, api_name string) float64 {
	return 0
}

func (self NoDelayRetryPolicy) ShouldRetry(retry_times int, exception *OTSServiceError, api_name string) bool {
	if self._max_retry_time_reached(retry_times, exception, api_name) {
		return false
	}

	if self._can_retry(retry_times, exception, api_name) {
		return true
	}

	return false
}
