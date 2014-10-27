// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ots2
package goots

import (
	"reflect"
	"sync"
	"time"

	_ "github.com/GiterLab/goots/protobuf"
	// _ "github.com/GiterLab/goots/protobuf/decoder"
	_ "github.com/GiterLab/goots/protobuf/coder"
)

const (
	DEFAULT_ENCODING       = "utf8"
	DEFAULT_SOCKET_TIMEOUT = 50
	DEFAULT_MAX_CONNECTION = 50
	DEFAULT_LOGGER_NAME    = "ots2-client"
)

var defaultOTSSetting = OTSClient{
	"http://127.0.0.1:8800",              // EndPoint
	"OTSMultiUser177_accessid",           // AccessId
	"OTSMultiUser177_accesskey",          // AccessKey
	"TestInstance177",                    // InstanceName
	DEFAULT_SOCKET_TIMEOUT * time.Second, // SocketTimeout
	DEFAULT_MAX_CONNECTION,               // MaxConnection
	DEFAULT_LOGGER_NAME,                  // LoggerName
	DEFAULT_ENCODING,                     // Encoding
}
var settingMutex sync.Mutex

// Overwrite default settings
func SetDefaultSetting(setting OTSClient) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultOTSSetting = setting
	if defaultOTSSetting.SocketTimeout == 0 {
		defaultOTSSetting.SocketTimeout = 50
	}
	if defaultSetting.ReadWriteTimeout == 0 {
		defaultSetting.ReadWriteTimeout = 60
	}
}

func New(end_point, accessid, accesskey, instance_name string, kwargs ...interface{}) *OTSClient {
	o := &defaultOTSSetting
	o.EndPoint = end_point
	o.AccessId = accessid
	o.AccessKey = accesskey
	o.InstanceName = instance_name

	for _, v := range kwargs {
		switch i {
		case 0: // SocketTimeout --> int32
			if _, ok := v.(int32); ok {
				o.SocketTimeout = v.(int32)
			} else {
				panic("OTSClient.SocketTimeout should be int32 type, not %v", reflect.TypeOf(v))
			}

		case 1: // MaxConnection --> int32
			if _, ok := v.(int32); ok {
				o.MaxConnection = v.(int32)
			} else {
				panic("OTSClient.MaxConnection should be int32 type, not %v", reflect.TypeOf(v))
			}

		case 2: // LoggerName --> string
			if _, ok := v.(string); ok {
				o.LoggerName = v.(string)
			} else {
				panic("OTSClient.LoggerName should be int32 type, not %v", reflect.TypeOf(v))
			}

		case 3: // Encoding --> string
			if _, ok := v.(string); ok {
				o.Encoding = v.(string)
			} else {
				panic("OTSClient.Encoding should be int32 type, not %v", reflect.TypeOf(v))
			}
		}
	}
	return
}

// OTSClient实例
type OTSClient struct {
	// OTS服务的地址（例如 'http://instance.cn-hangzhou.ots.aliyun.com:80'），必须以'http://'开头
	EndPoint string
	// 访问OTS服务的accessid，通过官方网站申请或通过管理员获取
	AccessId string
	// 访问OTS服务的accesskey，通过官方网站申请或通过管理员获取
	AccessKey string
	// 访问的实例名，通过官方网站控制台创建或通过管理员获取
	InstanceName string

	// 连接池中每个连接的Socket超时，单位为秒，可以为int或float。默认值为50
	SocketTimeout int32
	// 连接池的最大连接数。默认为50
	// golang http自带连接池管理，此参数留作以后扩展用
	MaxConnection int32

	// 用来在请求中打DEBUG日志，或者在出错时打ERROR日志
	LoggerName string

	// 字符编码格式，此参数未用,兼容python
	Encoding string
}

func (o *OTSClient) Set(end_point, accessid, accesskey, instance_name string, kwargs ...interface{}) *OTSClient {
	o.EndPoint = end_point
	o.AccessId = accessid
	o.AccessKey = accesskey
	o.InstanceName = instance_name

	return o
}
