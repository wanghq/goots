// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ots2
package goots

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	// . "github.com/GiterLab/goots/protobuf"
	"github.com/GiterLab/goots/urllib"
)

var OTSDebugEnable bool = false  // OTS调试默认关闭
var OTSLoggerEnable bool = false // OTS运行logger记录

const (
	DEFAULT_ENCODING       = "utf8"
	DEFAULT_SOCKET_TIMEOUT = 50
	DEFAULT_MAX_CONNECTION = 50
	DEFAULT_LOGGER_NAME    = "ots2-client"
)

var defaultOTSSetting = OTSClient{
	"http://127.0.0.1:8800",     // EndPoint
	"OTSMultiUser177_accessid",  // AccessId
	"OTSMultiUser177_accesskey", // AccessKey
	"TestInstance177",           // InstanceName
	DEFAULT_SOCKET_TIMEOUT,      // SocketTimeout
	DEFAULT_MAX_CONNECTION,      // MaxConnection
	DEFAULT_LOGGER_NAME,         // LoggerName
	DEFAULT_ENCODING,            // Encoding
	&defaultProtocol,
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
	if defaultOTSSetting.MaxConnection == 0 {
		defaultOTSSetting.MaxConnection = 50
	}
}

// 创建一个新的OTSClient实例
func New(end_point, accessid, accesskey, instance_name string, kwargs ...interface{}) (o *OTSClient, err error) {
	// init logger
	err = LoggerInit()
	if err != nil {
		return nil, err
	}

	o = &defaultOTSSetting
	o.EndPoint = end_point
	o.AccessId = accessid
	o.AccessKey = accesskey
	o.InstanceName = instance_name

	for i, v := range kwargs {
		switch i {
		case 0: // SocketTimeout --> int32
			if _, ok := v.(int); ok {
				o.SocketTimeout = v.(int)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.SocketTimeout should be int type, not %v", reflect.TypeOf(v)))
			}

		case 1: // MaxConnection --> int32
			if _, ok := v.(int); ok {
				o.MaxConnection = v.(int)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.MaxConnection should be int type, not %v", reflect.TypeOf(v)))
			}

		case 2: // LoggerName --> string
			if _, ok := v.(string); ok {
				o.LoggerName = v.(string)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.LoggerName should be string type, not %v", reflect.TypeOf(v)))
			}

		case 3: // Encoding --> string
			if _, ok := v.(string); ok {
				o.Encoding = v.(string)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.Encoding should be string type, not %v", reflect.TypeOf(v)))
			}
		}
	}

	// parse end point
	end_point_url, err := url.Parse(end_point)
	if err != nil {
		return nil, (OTSClientError{}.Set("url parse error", err))
	}
	if end_point_url.Scheme != "http" && end_point_url.Scheme != "https" {
		return nil, (OTSClientError{}.Set("protocol of end_point must be 'http' or 'https', e.g. http://ots.aliyuncs.com:80."))
	}

	if end_point_url.Host == "" {
		return nil, (OTSClientError{}.Set("host of end_point should be specified, e.g. http://ots.aliyuncs.com:80."))
	}

	// set default setting for urllib
	url_setting := urllib.HttpSettings{
		false,            // ShowDebug
		"GiterLab",       // UserAgent
		60 * time.Second, // ConnectTimeout
		60 * time.Second, // ReadWriteTimeout
		nil,              // TlsClientConfig
		nil,              // Proxy
		nil,              // Transport
		false,            // EnableCookie
	}
	if o.SocketTimeout != 0 {
		url_setting.ConnectTimeout = time.Duration(o.SocketTimeout) * time.Second
		url_setting.ReadWriteTimeout = time.Duration(o.SocketTimeout) * time.Second
	}
	if OTSDebugEnable {
		url_setting.ShowDebug = true
	}
	urllib.SetDefaultSetting(url_setting)

	o.protocol = newProtocol(nil)
	o.protocol.Set(o.AccessId, o.AccessKey, o.InstanceName, o.Encoding, o.LoggerName)

	return o, nil
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
	SocketTimeout int
	// 连接池的最大连接数。默认为50
	// golang http自带连接池管理，此参数留作以后扩展用
	MaxConnection int

	// 用来在请求中打DEBUG日志，或者在出错时打ERROR日志
	LoggerName string

	// 字符编码格式，此参数未用,兼容python
	Encoding string

	// protocol
	protocol *ots_protocol
}

func (o *OTSClient) String() string {
	r := ""
	r = r + fmt.Sprintln("#### OTSClinet Config ####")
	r = r + fmt.Sprintln("API_VERSION  :", API_VERSION)
	r = r + fmt.Sprintln("DebugEnable  :", OTSDebugEnable)
	r = r + fmt.Sprintln("EndPoint     :", o.EndPoint)
	r = r + fmt.Sprintln("AccessId     :", o.AccessId)
	r = r + fmt.Sprintln("AccessKey    :", o.AccessKey)
	r = r + fmt.Sprintln("InstanceName :", o.InstanceName)
	r = r + fmt.Sprintln("SocketTimeout:", o.SocketTimeout)
	r = r + fmt.Sprintln("MaxConnection:", o.MaxConnection)
	r = r + fmt.Sprintln("LoggerName   :", o.LoggerName)
	// r = r + fmt.Sprintln("Encoding:", o.Encoding)
	r = r + fmt.Sprintln("##########################")

	return r
}

// 在OTSClinet创建后（既调用了New函数），需要重新修改OTSClinet的参数时
// 可以调用此函数进行设置，参数使用字典方式，可以使用的字典如下：
// Debug --> bool
// EndPoint --> string
// AccessId --> string
// AccessKey --> string
// InstanceName --> string
// SocketTimeout --> int
// MaxConnection --> int
// LoggerName --> string
// Encoding --> string
// 注：具体参数意义请查看OTSClinet定义处的注释
func (o *OTSClient) Set(kwargs DictString) *OTSClient {
	if len(kwargs) != 0 {
		for k, v := range kwargs {
			switch k {
			case "Debug":
				if v1, ok := v.(bool); ok {
					setting := urllib.GetDefaultSetting()
					setting.ShowDebug = v1
				} else {
					panic(OTSClientError{}.Set("Debug should be bool, not %v", reflect.TypeOf(v)))
				}
			case "EndPoint":
				if v1, ok := v.(string); ok {
					o.EndPoint = v1
				} else {
					panic(OTSClientError{}.Set("EndPoint should be string, not %v", reflect.TypeOf(v)))
				}
				// parse end point
				end_point_url, err := url.Parse(v.(string))
				if err != nil {
					panic(OTSClientError{}.Set("url parse error", err))
				}
				if end_point_url.Scheme != "http" && end_point_url.Scheme != "https" {
					panic(OTSClientError{}.Set("protocol of end_point must be 'http' or 'https', e.g. http://ots.aliyuncs.com:80."))
				}

				if end_point_url.Host == "" {
					panic(OTSClientError{}.Set("host of end_point should be specified, e.g. http://ots.aliyuncs.com:80."))
				}

			case "AccessId":
				if v1, ok := v.(string); ok {
					o.AccessId = v1
				} else {
					panic(OTSClientError{}.Set("AccessId should be string, not %v", reflect.TypeOf(v)))
				}

			case "AccessKey":
				if v1, ok := v.(string); ok {
					o.AccessKey = v1
				} else {
					panic(OTSClientError{}.Set("AccessKey should be string, not %v", reflect.TypeOf(v)))
				}

			case "InstanceName":
				if v1, ok := v.(string); ok {
					o.InstanceName = v1
				} else {
					panic(OTSClientError{}.Set("InstanceName should be string, not %v", reflect.TypeOf(v)))
				}

			case "SocketTimeout":
				if v1, ok := v.(int); ok {
					o.SocketTimeout = v1
				} else {
					panic(OTSClientError{}.Set("SocketTimeout should be int, not %v", reflect.TypeOf(v)))
				}

			case "MaxConnection":
				if v1, ok := v.(int); ok {
					o.MaxConnection = v1
				} else {
					panic(OTSClientError{}.Set("MaxConnection should be int, not %v", reflect.TypeOf(v)))
				}

			case "LoggerName":
				if v1, ok := v.(string); ok {
					o.LoggerName = v1
				} else {
					panic(OTSClientError{}.Set("LoggerName should be string, not %v", reflect.TypeOf(v)))
				}

			case "Encoding":
				if v1, ok := v.(string); ok {
					o.Encoding = v1
				} else {
					panic(OTSClientError{}.Set("Encoding should be string, not %v", reflect.TypeOf(v)))
				}

			default:
				panic(OTSClientError{}.Set("Unknown param %s", k))
			}
		}
	}

	return o
}

func (o *OTSClient) _request_helper(api_name string, args ...interface{}) (resp []reflect.Value, ots_service_error *OTSServiceError) {
	ots_service_error = new(OTSServiceError)

	// 1. make_request
	query, reqheaders, reqbody, err := o.protocol.make_request(api_name, args...)
	if err != nil {
		return nil, ots_service_error.SetErrorMessage("%s", err)
	}

	// 2. http send_receive
	req := urllib.Post(o.EndPoint + query)
	if OTSDebugEnable {
		req.Debug(true)
	}
	req.Body(reqbody)
	if reqheaders != nil {
		for k, v := range reqheaders {
			req.Header(k, v.(string))
		}
	}
	response, err := req.Response()
	if err != nil {
		return nil, ots_service_error.SetErrorMessage("%s", err)
	}
	status := response.StatusCode // e.g. 200
	reason := response.Status     // e.g. "200 OK"
	var resheaders = DictString{}
	if response.Header != nil {
		for k, v := range response.Header {
			resheaders[strings.ToLower(k)] = v[0] // map[string][]string
		}
	}
	if response.Body == nil {
		return nil, ots_service_error.SetErrorMessage("Http body is empty")
	}
	defer response.Body.Close()
	resbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, ots_service_error.SetErrorMessage("%s", err)
	}

	// for debug
	if OTSDebugEnable {
		fmt.Println("==== Aliyun OTS Response ====")
		fmt.Println("status:", status)
		fmt.Println("reason:", reason)
		fmt.Println("headers:", resheaders)
		if resbody != nil {
			if len(resbody) == 0 {
				fmt.Println("body-raw:", "None")
				fmt.Println("body-string:", "None")
			} else {
				fmt.Println("body-raw:", resbody)
				fmt.Println("body-string:", string(resbody))
			}

		} else {
			fmt.Println("body-raw:", "None")
			fmt.Println("body-string:", "None")
		}
		fmt.Println("-----------------------------")
	}

	// 3. handle_error
	ots_service_error = o.protocol.handle_error(api_name, query, reason, status, resheaders, resbody)
	if ots_service_error != nil {
		return nil, ots_service_error
	}

	// 4. parse_response
	resp, ots_service_error = o.protocol.parse_response(api_name, reason, status, resheaders, resbody)
	if ots_service_error != nil {
		return nil, ots_service_error
	}

	return resp, nil
}

// 说明：根据表信息创建表。
//
// ``table_meta``是``otstype.OTSTableMeta``类的实例，它包含表名和PrimaryKey的schema，
// 请参考``OTSTableMeta``类的文档。当创建了一个表之后，通常要等待1分钟时间使partition load
// 完成，才能进行各种操作。
// ``reserved_throughput``是``otstype.ReservedThroughput``类的实例，表示预留读写吞吐量。
//
// 返回：无。
//       错误信息。
//
// 示例：
//
// table_meta := &OTSTableMeta{
// 	TableName: "myTable",
// 	SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
// 		"gid": "INTEGER",
// 		"uid": "INTEGER",
// 	},
// }
//
// reserved_throughput := &OTSReservedThroughput{
// 	OTSCapacityUnit{100, 100},
// }
//
// ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
//
func (o *OTSClient) CreateTable(table_meta *OTSTableMeta, reserved_throughput *OTSReservedThroughput) (err *OTSError) {
	err = new(OTSError)
	r, service_err := o._request_helper("CreateTable", table_meta, reserved_throughput)
	if service_err != nil {
		return err.SetServiceError(service_err)
	}

	if r[0].Interface() != nil {
		return err.SetClientMessage("[CreateTable] %s", r[0].Interface().(error))
	}

	return nil
}

func (o *OTSClient) DeleteTable() {

}

// 说明：获取所有表名的列表。
//
// 返回：表名列表。
//       错误信息。
//
// ``table_list``表示获取的表名列表，类型为OTSListTableResponse。
//
// 示例：
//
//     table_list, ots_err := ots_client.ListTable()
//
func (o *OTSClient) ListTable() (table_list *OTSListTableResponse, err *OTSError) {
	err = new(OTSError)
	r, service_err := o._request_helper("ListTable")
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}
	if r == nil {
		return nil, err.SetClientMessage("[ListTable] Not expect error")
	}

	return r[0].Interface().(*OTSListTableResponse), nil
}

func (o *OTSClient) UpdateTable() {

}

func (o *OTSClient) DescribeTable() {

}

func (o *OTSClient) GetRow() {

}

func (o *OTSClient) PutRow() {

}

func (o *OTSClient) UpdateRow() {

}

func (o *OTSClient) DeleteRow() {

}

func (o *OTSClient) BatchGetRow() {

}

func (o *OTSClient) BatchWriteRow() {

}

func (o *OTSClient) GetRange() {

}

func (o *OTSClient) XGetRange() {

}
