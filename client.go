// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ots2
package goots

import (
	"errors"
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

var OTSDebugEnable bool = false     // OTS调试默认关闭
var OTSLoggerEnable bool = false    // OTS运行logger记录
var OTSHttpDebugEnable bool = false // OTS HTTP调试记录

const (
	DEFAULT_ENCODING       = "utf8"
	DEFAULT_SOCKET_TIMEOUT = 50
	DEFAULT_MAX_CONNECTION = 50
	DEFAULT_LOGGER_NAME    = "ots2-client"
)

var defaultOTSSetting = OTSClient{
	"your_instance_address", // EndPoint
	"your_accessid",         // AccessId
	"your_accesskey",        // AccessKey
	"your_instance_name",    // InstanceName
	DEFAULT_SOCKET_TIMEOUT,  // SocketTimeout
	DEFAULT_MAX_CONNECTION,  // MaxConnection
	DEFAULT_LOGGER_NAME,     // LoggerName
	DEFAULT_ENCODING,        // Encoding
	&defaultProtocol,        // default protocol
	&OTSDefaultRetryPolicy,  // default retry policy
}
var settingMutex sync.Mutex

//		Overwrite default settings
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

//		创建一个新的OTSClient实例
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

	// initialize the retry policy
	if o.RetryPolicy == nil {
		o.RetryPolicy = OTSDefaultRetryPolicy
	}

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
		true,             // Gzip
		true,             // DumpBody
	}
	if o.SocketTimeout != 0 {
		url_setting.ConnectTimeout = time.Duration(o.SocketTimeout) * time.Second
		url_setting.ReadWriteTimeout = time.Duration(o.SocketTimeout) * time.Second
	}
	if OTSHttpDebugEnable {
		url_setting.ShowDebug = true
	} else {
		url_setting.ShowDebug = false
	}
	urllib.SetDefaultSetting(url_setting)

	o.protocol = newProtocol(nil)
	o.protocol.Set(o.AccessId, o.AccessKey, o.InstanceName, o.Encoding, o.LoggerName)

	return o, nil
}

//		创建一个新的OTSClient实例
func NewWithRetryPolicy(end_point, accessid, accesskey, instance_name string, retry_policy RetryPolicyInterface, kwargs ...interface{}) (o *OTSClient, err error) {
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

	// initialize the retry policy
	if retry_policy != nil {
		o.RetryPolicy = retry_policy
	} else {
		o.RetryPolicy = OTSDefaultRetryPolicy
	}

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
		true,             // Gzip
		true,             // DumpBody
	}
	if o.SocketTimeout != 0 {
		url_setting.ConnectTimeout = time.Duration(o.SocketTimeout) * time.Second
		url_setting.ReadWriteTimeout = time.Duration(o.SocketTimeout) * time.Second
	}
	if OTSHttpDebugEnable {
		url_setting.ShowDebug = true
	} else {
		url_setting.ShowDebug = false
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

	// 定义了重试策略，默认的重试策略为 DefaultRetryPolicy。
	// 你可以继承 RetryPolicy 来实现自己的重试策略，请参考 DefaultRetryPolicy 的代码。
	RetryPolicy RetryPolicyInterface
}

func (o *OTSClient) String() string {
	r := ""
	r = r + fmt.Sprintln("#### OTSClinet Config ####")
	r = r + fmt.Sprintln("API_VERSION    :", API_VERSION)
	r = r + fmt.Sprintln("DebugEnable    :", OTSDebugEnable)
	r = r + fmt.Sprintln("EndPoint       :", o.EndPoint)
	r = r + fmt.Sprintln("AccessId       :", o.AccessId)
	r = r + fmt.Sprintln("AccessKey      :", o.AccessKey)
	r = r + fmt.Sprintln("InstanceName   :", o.InstanceName)
	r = r + fmt.Sprintln("SocketTimeout  :", o.SocketTimeout)
	r = r + fmt.Sprintln("MaxConnection  :", o.MaxConnection)
	r = r + fmt.Sprintln("OTSLoggerEnable:", OTSLoggerEnable)
	r = r + fmt.Sprintln("LoggerName     :", o.LoggerName)
	// r = r + fmt.Sprintln("Encoding:", o.Encoding)
	r = r + fmt.Sprintln("##########################")

	return r
}

// 		在OTSClinet创建后（既调用了New函数），需要重新修改OTSClinet的参数时
// 		可以调用此函数进行设置，参数使用字典方式，可以使用的字典如下：
// 		Debug --> bool
// 		EndPoint --> string
// 		AccessId --> string
// 		AccessKey --> string
// 		InstanceName --> string
// 		SocketTimeout --> int
// 		MaxConnection --> int
// 		LoggerName --> string
// 		Encoding --> string
// 		注：具体参数意义请查看OTSClinet定义处的注释
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
	var reason string
	var status int
	var resheaders = DictString{}
	var resbody []byte

	ots_service_error = new(OTSServiceError)

	// 1. make_request
	query, reqheaders, reqbody, err := o.protocol.make_request(api_name, args...)
	if err != nil {
		return nil, ots_service_error.SetErrorMessage("%s", err)
	}

	retry_times := 0

	for {
		// 2. http send_receive
		req := urllib.Post(o.EndPoint + query)
		if OTSHttpDebugEnable {
			req.Debug(true)
		} else {
			req.Debug(false)
		}
		req.Body(reqbody)
		if reqheaders != nil {
			for k, v := range reqheaders {
				req.Header(k, v.(string))
			}
		}
		response, err := req.Response()
		if err != nil {
			ots_service_error.SetErrorMessage("%s", err)
			if o.RetryPolicy.ShouldRetry(retry_times, ots_service_error, api_name) {
				retry_delay := o.RetryPolicy.GetRetryDelay(retry_times, ots_service_error, api_name)
				time.Sleep(time.Duration(retry_delay*1000) * time.Millisecond)
				retry_times += 1
			} else {
				return nil, ots_service_error
			}
		}
		status = response.StatusCode // e.g. 200
		reason = response.Status     // e.g. "200 OK"
		if response.Header != nil {
			for k, v := range response.Header {
				resheaders[strings.ToLower(k)] = v[0] // map[string][]string
			}
		}
		if response.Body == nil {
			ots_service_error.SetErrorMessage("Http body is empty")
			if o.RetryPolicy.ShouldRetry(retry_times, ots_service_error, api_name) {
				retry_delay := o.RetryPolicy.GetRetryDelay(retry_times, ots_service_error, api_name)
				time.Sleep(time.Duration(retry_delay*1000) * time.Millisecond)
				retry_times += 1
			} else {
				return nil, ots_service_error
			}
		}
		defer response.Body.Close()
		resbody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			ots_service_error.SetErrorMessage("%s", err)
			if o.RetryPolicy.ShouldRetry(retry_times, ots_service_error, api_name) {
				retry_delay := o.RetryPolicy.GetRetryDelay(retry_times, ots_service_error, api_name)
				time.Sleep(time.Duration(retry_delay*1000) * time.Millisecond)
				retry_times += 1
			} else {
				return nil, ots_service_error
			}
		}

		// for debug
		if OTSDebugEnable {
			fmt.Println("==== Aliyun OTS Response ====")
			fmt.Println("status:", status)
			fmt.Println("reason:", reason)
			fmt.Println("headers:", resheaders)
			// if resbody != nil {
			// 	if len(resbody) == 0 {
			// 		fmt.Println("body-raw:", "None")
			// 		fmt.Println("body-string:", "None")
			// 	} else {
			// 		fmt.Println("body-raw:", resbody)
			// 		fmt.Println("body-string:", string(resbody))
			// 	}
			//
			// } else {
			// 	fmt.Println("body-raw:", "None")
			// 	fmt.Println("body-string:", "None")
			// }
			// fmt.Println("-----------------------------")
		}

		// 3. handle_error
		ots_service_error = o.protocol.handle_error(api_name, query, reason, status, resheaders, resbody)
		if ots_service_error != nil {
			if o.RetryPolicy.ShouldRetry(retry_times, ots_service_error, api_name) {
				retry_delay := o.RetryPolicy.GetRetryDelay(retry_times, ots_service_error, api_name)
				fmt.Println(retry_delay, time.Duration(retry_delay*1000)*time.Millisecond)
				time.Sleep(time.Duration(retry_delay*1000) * time.Millisecond)
				retry_times += 1
			} else {
				return nil, ots_service_error
			}
		} else {
			break
		}
	} // end for

	// 4. parse_response
	resp, ots_service_error = o.protocol.parse_response(api_name, reason, status, resheaders, resbody)
	if ots_service_error != nil {
		return nil, ots_service_error
	}

	return resp, nil
}

func (o *OTSClient) _check_request_helper_error(resp []reflect.Value) (r interface{}, e error) {
	// parse the following two cases
	// 1. (err error)
	// 2. (x *xxx, err error)
	switch len(resp) {
	case 1: // (err error)
		if resp[0].Interface() != nil {
			if err, ok := resp[0].Interface().(error); ok {
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("Illegal data parameters, parse err failed")
			}
		}
		return nil, nil

	case 2: // (x *xxx, err error)
		if resp[1].Interface() != nil {
			if err, ok := resp[1].Interface().(error); ok {
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("Illegal data parameters, parse err failed")
			}
		}
		return resp[0].Interface(), nil

	default:
		return nil, errors.New("Illegal data parameters")
	}

	return nil, errors.New("The program will not perform here")
}

// 说明：根据表信息创建表。
//
// 		``table_meta``是``otstype.OTSTableMeta``类的实例，它包含表名和PrimaryKey的schema，
// 		请参考``OTSTableMeta``类的文档。当创建了一个表之后，通常要等待1分钟时间使partition load
// 		完成，才能进行各种操作。
// 		``reserved_throughput``是``otstype.ReservedThroughput``类的实例，表示预留读写吞吐量。
//
// 		返回：无。
// 		      错误信息。
//
// 		示例：
//
// 		table_meta := &OTSTableMeta{
// 			TableName: "myTable",
// 			SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
// 				"gid": "INTEGER",
// 				"uid": "INTEGER",
// 			},
// 		}
//
// 		reserved_throughput := &OTSReservedThroughput{
// 			OTSCapacityUnit{100, 100},
// 		}
//
// 		ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
//
func (o *OTSClient) CreateTable(table_meta *OTSTableMeta, reserved_throughput *OTSReservedThroughput) (err *OTSError) {
	err = new(OTSError)
	if table_meta == nil {
		return err.SetClientMessage("[CreateTable] table_meta should not be nil")
	}
	if reserved_throughput == nil {
		return err.SetClientMessage("[CreateTable] reserved_throughput should not be nil")
	}

	resp, service_err := o._request_helper("CreateTable", table_meta, reserved_throughput)
	if service_err != nil {
		return err.SetServiceError(service_err)
	}

	_, e := o._check_request_helper_error(resp)
	if e != nil {
		return err.SetClientMessage("[CreateTable] %s", e)
	}

	return nil
}

// 说明：根据表名删除表。
//
// 		``table_name``是对应的表名。
//
// 		返回：无。
// 		      错误信息。
//
// 		示例：
//
// 		ots_client.DeleteTable("myTable")
//
func (o *OTSClient) DeleteTable(table_name string) (err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return err.SetClientMessage("[DeleteTable] table_name should not be empty")
	}

	resp, service_err := o._request_helper("DeleteTable", table_name)
	if service_err != nil {
		return err.SetServiceError(service_err)
	}

	_, e := o._check_request_helper_error(resp)
	if e != nil {
		return err.SetClientMessage("[DeleteTable] %s", e)
	}

	return nil
}

// 说明：获取所有表名的列表。
//
// 		返回：表名列表。
// 		      错误信息。
//
// 		``table_list``表示获取的表名列表，类型为``otstype.OTSListTableResponse``。
//
// 		示例：
//
// 		table_list, ots_err := ots_client.ListTable()
//
func (o *OTSClient) ListTable() (table_list *OTSListTableResponse, err *OTSError) {
	err = new(OTSError)
	resp, service_err := o._request_helper("ListTable")
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)

	if e != nil {
		return nil, err.SetClientMessage("[ListTable] %s", e)
	}

	return r.(*OTSListTableResponse), nil
}

// 说明：更新表属性，目前只支持修改预留读写吞吐量。
//
// 		``table_name``是对应的表名。
// 		``reserved_throughput``是``otstype.ReservedThroughput``类的实例，表示预留读写吞吐量。
//
// 		返回：针对该表的预留读写吞吐量的最近上调时间、最近下调时间和当天下调次数。
// 		      错误信息。
//
// 		``update_table_response``表示更新的结果，是``otstype.OTSUpdateTableResponse``类的实例。
//
// 		示例：
// 		reserved_throughput := &OTSReservedThroughput{
// 		 OTSCapacityUnit{5000, 5000},
// 		}
//
// 		// 每次调整操作的间隔应大于10分钟
// 		// 如果是刚创建表，需要10分钟之后才能调整表的预留读写吞吐量。
// 		update_response, ots_err := ots_client.UpdateTable("myTable", reserved_throughput)
//
func (o *OTSClient) UpdateTable(table_name string, reserved_throughput *OTSReservedThroughput) (update_table_response *OTSUpdateTableResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[UpdateTable] table_name should not be empty")
	}
	if reserved_throughput == nil {
		return nil, err.SetClientMessage("[UpdateTable] reserved_throughput should not be nil")
	}

	resp, service_err := o._request_helper("UpdateTable", table_name, reserved_throughput)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[UpdateTable] %s ", e)
	}

	return r.(*OTSUpdateTableResponse), nil
}

// 说明：获取表的描述信息。
//
// 		``table_name``是对应的表名。
//
// 		返回：表的描述信息。
// 		      错误信息。
//
// 		``describe_table_response``表示表的描述信息，是``otstype.OTSDescribeTableResponse``类的实例。
//
// 		示例：
//
// 		describe_response, ots_err := ots_client.DescribeTable("myTable")
//
func (o *OTSClient) DescribeTable(table_name string) (describe_table_response *OTSDescribeTableResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[DescribeTable] table_name should not be empty")
	}

	resp, service_err := o._request_helper("DescribeTable", table_name)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[DescribeTable] %s", e)
	}

	return r.(*OTSDescribeTableResponse), nil
}

// 说明：获取一行数据。
//
// 		``table_name``是对应的表名。
// 		``primary_key``是主键，类型为``otstype.OTSPrimaryKey``。
// 		``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为``otstype.OTSColumnsToGet``；如果填nil，表示获取所有列。
//
// 		返回：本次操作消耗的CapacityUnit、行数据（包含主键列和属性列）。
// 		      错误信息。
//
// 		``get_row_response``为``otstype.OTSGetRowResponse``类的实例包含了：
// 		``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
// 		``Row``表示一行的数据，是``otstype.OTSRow``的实例,也包含了:
// 		``PrimaryKeyColumns``表示主键列，类型为``otstype.OTSPrimaryKey``，如：{"PK0":value0, "PK1":value1}。
// 		``AttributeColumns``表示属性列，类型为``otstype.OTSAttribute``，如：{"COL0":value0, "COL1":value1}。
//
// 		示例：
//
// 		primary_key := &OTSPrimaryKey{
// 			"gid": 1,
// 			"uid": 101,
// 		}
// 		columns_to_get := &OTSColumnsToGet{
// 			"name", "address", "age",
// 		}
// 		// columns_to_get = nil // read all
// 		get_row_response, ots_err := ots_client.GetRow("myTable", primary_key, columns_to_get)
//
func (o *OTSClient) GetRow(table_name string, primary_key *OTSPrimaryKey, columns_to_get *OTSColumnsToGet) (get_row_response *OTSGetRowResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[GetRow] table_name should not be empty")
	}
	if primary_key == nil {
		return nil, err.SetClientMessage("[GetRow] primary_key should not be nil")
	}

	resp, service_err := o._request_helper("GetRow", table_name, primary_key, columns_to_get)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[GetRow] %s", e)
	}

	return r.(*OTSGetRowResponse), nil
}

// 说明：写入一行数据。返回本次操作消耗的CapacityUnit。
//
// 		``table_name``是对应的表名。
// 		``condition``表示执行操作前做条件检查，满足条件才执行，是string的实例。
// 		目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
// 		``primary_key``表示主键，类型为``otstype.OTSPrimaryKey``的实例。
// 		``attribute_columns``表示属性列，类型为``otstype.OTSAttribute``的实例。
//
// 		返回：本次操作消耗的CapacityUnit。
// 		      错误信息。
//
// 		``put_row_response``为``otstype.OTSGetRowResponse``类的实例包含了：
// 		``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
//
// 		示例：
//
// 		primary_key := &OTSPrimaryKey{
// 			"gid": 1,
// 			"uid": 101,
// 		}
// 		attribute_columns := &OTSAttribute{
// 			"name":    "张三",
// 			"mobile":  111111111,
// 			"address": "中国A地",
// 			"age":     20,
// 		}
// 		condition := OTSCondition_EXPECT_NOT_EXIST
// 		put_row_response, ots_err := ots_client.PutRow("myTable", condition, primary_key, attribute_columns)
//
func (o *OTSClient) PutRow(table_name string, condition string, primary_key *OTSPrimaryKey, attribute_columns *OTSAttribute) (put_row_response *OTSPutRowResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[PutRow] table_name should not be empty")
	}
	if condition == "" {
		return nil, err.SetClientMessage("[PutRow] condition should not be empty")
	}
	if primary_key == nil {
		return nil, err.SetClientMessage("[PutRow] primary_key should not be nil")
	}
	if attribute_columns == nil {
		return nil, err.SetClientMessage("[PutRow] attribute_columns should not be nil")
	}

	resp, service_err := o._request_helper("PutRow", table_name, condition, primary_key, attribute_columns)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[PutRow] %s", e)
	}

	return r.(*OTSPutRowResponse), nil
}

// 说明：更新一行数据。
//
// 		``table_name``是对应的表名。
// 		``condition``表示执行操作前做条件检查，满足条件才执行，是string的实例。
// 		目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
// 		``primary_key``表示主键，类型为``otstype.OTSPrimaryKey``的实例。
// 		``update_of_attribute_columns``表示属性列，类型为``otstype.OTSUpdateOfAttribute``的实例，可以包含put和delete操作。其中put是
// 		``otstype.OTSColumnsToPut`` 表示属性列的写入；delete是``otstype.OTSColumnsToDelete``，表示要删除的属性列的列名，
// 		见示例。
//
// 		返回：本次操作消耗的CapacityUnit。
// 		      错误信息。
//
// 		``update_row_response``为``otstype.OTSUpdateRowResponse``类的实例包含了：
// 		``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
//
// 		示例：
//
// 		primary_key := &OTSPrimaryKey{
// 			"gid": 1,
// 			"uid": 101,
// 		}
// 		update_of_attribute_columns := &OTSUpdateOfAttribute{
// 			OTSOperationType_PUT: OTSColumnsToPut{
// 				"name":    "张三丰",
// 				"address": "中国B地",
// 			},
//
// 			OTSOperationType_DELETE: OTSColumnsToDelete{
// 				"mobile", "age",
// 			},
// 		}
// 		condition := OTSCondition_EXPECT_EXIST
// 		update_row_response, ots_err := ots_client.UpdateRow("myTable", condition, primary_key, update_of_attribute_columns)
//
func (o *OTSClient) UpdateRow(table_name string, condition string, primary_key *OTSPrimaryKey, update_of_attribute_columns *OTSUpdateOfAttribute) (update_row_response *OTSUpdateRowResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[UpdateRow] table_name should not be empty")
	}
	if condition == "" {
		return nil, err.SetClientMessage("[UpdateRow] condition should not be empty")
	}
	if primary_key == nil {
		return nil, err.SetClientMessage("[UpdateRow] primary_key should not be nil")
	}
	if update_of_attribute_columns == nil {
		return nil, err.SetClientMessage("[UpdateRow] update_of_attribute_columns should not be nil")
	}

	resp, service_err := o._request_helper("UpdateRow", table_name, condition, primary_key, update_of_attribute_columns)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[UpdateRow] %s", e)
	}

	return r.(*OTSUpdateRowResponse), nil
}

// 说明：删除一行数据。
//
// 		``table_name``是对应的表名。
// 		``condition``表示执行操作前做条件检查，满足条件才执行，是string的实例。
// 		目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
// 		``primary_key``表示主键，类型为``otstype.OTSPrimaryKey``的实例。
//
// 		返回：本次操作消耗的CapacityUnit。
// 		      错误信息。
//
// 		``delete_row_response``为``otstype.OTSDeleteRowResponse``类的实例包含了：
// 		``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
//
// 		示例：
//
// 		primary_key := &OTSPrimaryKey{
// 			"gid": 1,
// 			"uid": 101,
// 		}
// 		condition := OTSCondition_IGNORE
// 		delete_row_response, ots_err := ots_client.DeleteRow("myTable", condition, primary_key)
//
func (o *OTSClient) DeleteRow(table_name string, condition string, primary_key *OTSPrimaryKey) (delete_row_response *OTSDeleteRowResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[DeleteRow] table_name should not be empty")
	}
	if condition == "" {
		return nil, err.SetClientMessage("[DeleteRow] condition should not be empty")
	}
	if primary_key == nil {
		return nil, err.SetClientMessage("[DeleteRow] primary_key should not be nil")
	}

	resp, service_err := o._request_helper("DeleteRow", table_name, condition, primary_key)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[DeleteRow] %s", e)
	}

	return r.(*OTSDeleteRowResponse), nil
}

// 说明：批量获取多行数据。
//
// 		``batch_list``表示获取多行的条件列表，格式如下：
//
// 		batch_list := &OTSBatchGetRowRequest{
// 			{
// 				// TableName
// 				TableName: "table_name0",
// 				// PrimaryKey
// 				Rows: OTSPrimaryKeyRows{
// 					{"gid": 1, "uid": 101},
// 					{"gid": 2, "uid": 202},
// 					{"gid": 3, "uid": 303},
// 				},
// 				// ColumnsToGet
// 				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
// 			},
// 			{
// 				// TableName
// 				TableName: "table_name1",
// 				// PrimaryKey
// 				Rows: OTSPrimaryKeyRows{
// 					{"gid": 1, "uid": 101},
// 					{"gid": 2, "uid": 202},
// 					{"gid": 3, "uid": 303},
// 				},
// 				// ColumnsToGet
// 				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
// 			},
// 			...
// 		}
//
// 		其中，Rows 为主键，类型为``otstype.OTSPrimaryKeyRows``。
//
// 		返回：对应行的结果列表。
// 		      错误信息
//
// 		``response_rows_list``为``otstype.OTSBatchGetRowResponse``的实例
// 		``response_rows_list.Tables``为返回的结果列表，与请求的顺序一一对应，格式如下：
// 		response_rows_list.Tables --> []*OTSTableInBatchGetRowResponseItem{
// 			{
// 				TableName: "table_name0",
// 				Rows : []*OTSRowInBatchGetRowResponseItem{
// 					row_data_item0, row_data_item1, ...
// 				},
// 			},
// 			{
// 				TableName: "table_name1",
// 				Rows : []*OTSRowInBatchGetRowResponseItem{
// 					row_data_item0, row_data_item1, ...
// 				},
// 			},
// 			...
// 		}
//
// 		其中，row_data_item0, row_data_item1为``otstype.OTSRowInBatchGetRowResponseItem``的实例。
//
// 		示例：
//
// 		batch_list_get := &OTSBatchGetRowRequest{
// 			{
// 				// TableName
// 				TableName: "myTable",
// 				// PrimaryKey
// 				Rows: OTSPrimaryKeyRows{
// 					{"gid": 1, "uid": 101},
// 					{"gid": 2, "uid": 202},
// 					{"gid": 3, "uid": 303},
// 				},
// 				// ColumnsToGet
// 				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
// 			},
// 			{
// 				// TableName
// 				TableName: "notExistTable",
// 				// PrimaryKey
// 				Rows: OTSPrimaryKeyRows{
// 					{"gid": 1, "uid": 101},
// 					{"gid": 2, "uid": 202},
// 					{"gid": 3, "uid": 303},
// 				},
// 				// ColumnsToGet
// 				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
// 			},
// 		}
// 		batch_get_response, ots_err := ots_client.BatchGetRow(batch_list_get)
//
func (o *OTSClient) BatchGetRow(batch_list *OTSBatchGetRowRequest) (response_rows_list *OTSBatchGetRowResponse, err *OTSError) {
	err = new(OTSError)
	if batch_list == nil {
		return nil, err.SetClientMessage("[BatchGetRow] primary_key should not be nil")
	}

	resp, service_err := o._request_helper("BatchGetRow", batch_list)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[BatchGetRow] %s", e)
	}

	return r.(*OTSBatchGetRowResponse), nil
}

// 说明：批量修改多行数据。
//
// 		``batch_list``表示获取多行的条件列表，格式如下：
//
// 		batch_list := &OTSBatchWriteRowRequest{
// 			{
// 				TableName: "table_name0",
// 				PutRows: OTSPutRows{
// 					put_row_item, ...
// 				},
// 				UpdateRows: OTSUpdateRows{
// 					update_row_item, ...
// 				},
// 				DeleteRows: OTSDeleteRows{
// 					delete_row_item, ...
// 				},
// 			},
// 			{
// 				TableName: "table_name1",
// 				PutRows: OTSPutRows{
// 					put_row_item, ...
// 				},
// 				UpdateRows: OTSUpdateRows{
// 					update_row_item, ...
// 				},
// 				DeleteRows: OTSDeleteRows{
// 					delete_row_item, ...
// 				},
// 			},
// 			...
// 		}
//
// 		其中，put_row_item, 是``otstype.OTSPutRows``类的实例；
// 		      update_row_item, 是``otstype.OTSUpdateRows``类的实例；
// 		      delete_row_item, 是``otstype.OTSDeleteRows``类的实例。
//
// 		返回：对应行的修改结果列表。
// 		      错误信息。
//
// 		``response_items_list``为``otstype.OTSBatchWriteRowResponse``的实例
// 		``response_items_list.Tables``为返回的结果列表，与请求的顺序一一对应，格式如下：
// 		response_items_list.Tables --> []*OTSTableInBatchWriteRowResponseItem{
// 			{
// 				TableName: "table_name0", // for table_name0
// 				PutRows: []*OTSRowInBatchWriteRowResponseItem{
// 					put_row_resp, ...
// 				},
// 				UpdateRows: []*OTSRowInBatchWriteRowResponseItem{
// 					update_row_resp, ...
// 				},
// 				DeleteRows: []*OTSRowInBatchWriteRowResponseItem{
// 					delete_row_resp, ...
// 				}
// 			},
// 			{
// 				TableName: "table_name1", // for table_name1
// 				PutRows: []*OTSRowInBatchWriteRowResponseItem{
// 					put_row_resp, ...
// 				},
// 				UpdateRows: []*OTSRowInBatchWriteRowResponseItem{
// 					update_row_resp, ...
// 				},
// 				DeleteRows: []*OTSRowInBatchWriteRowResponseItem{
// 					delete_row_resp, ...
// 				}
// 			},
// 			...
// 		}
//
// 		其中put_row_resp，update_row_resp和delete_row_resp都是``*otstype.OTSRowInBatchWriteRowResponseItem``类的实例。
//
// 		示例：
//
// 		put_row_item := OTSPutRowItem{
// 			Condition: OTSCondition_EXPECT_NOT_EXIST, // OTSCondition_IGNORE
// 			PrimaryKey: OTSPrimaryKey{
// 				"gid": 2,
// 				"uid": 202,
// 			},
// 			AttributeColumns: OTSAttribute{
// 				"name":    "李四",
// 				"address": "中国某地",
// 				"age":     20,
// 			},
// 		}
// 		// [2] update_row
// 		update_row_item := OTSUpdateRowItem{
// 			Condition: OTSCondition_IGNORE,
// 			PrimaryKey: OTSPrimaryKey{
// 				"gid": 3,
// 				"uid": 303,
// 			},
// 			UpdateOfAttributeColumns: OTSUpdateOfAttribute{
// 				OTSOperationType_PUT: OTSColumnsToPut{
// 					"name":    "李三",
// 					"address": "中国某地",
// 				},
// 				OTSOperationType_DELETE: OTSColumnsToDelete{
// 					"mobile", "age",
// 				},
// 			},
// 		}
// 		// [3] delete_row
// 		delete_row_item := OTSDeleteRowItem{
// 			Condition: OTSCondition_IGNORE,
// 			PrimaryKey: OTSPrimaryKey{
// 				"gid": 4,
// 				"uid": 404,
// 			},
// 		}
// 		batch_list := &OTSBatchWriteRowRequest{
// 			{
// 				TableName: "myTable",
// 				PutRows: OTSPutRows{
// 					put_row_item,
// 				},
// 				UpdateRows: OTSUpdateRows{
// 					update_row_item,
// 				},
// 				DeleteRows: OTSDeleteRows{
// 					delete_row_item,
// 				},
// 			},
// 			{
// 				TableName: "notExistTable",
// 				PutRows: OTSPutRows{
// 					put_row_item,
// 				},
// 				UpdateRows: OTSUpdateRows{
// 					update_row_item,
// 				},
// 				DeleteRows: OTSDeleteRows{
// 					delete_row_item,
// 				},
// 			},
// 		}
// 		batch_write_response, ots_err := ots_client.BatchWriteRow(batch_list)
//
func (o *OTSClient) BatchWriteRow(batch_list *OTSBatchWriteRowRequest) (response_item_list *OTSBatchWriteRowResponse, err *OTSError) {
	err = new(OTSError)
	if batch_list == nil {
		return nil, err.SetClientMessage("[BatchWriteRow] primary_key should not be nil")
	}

	resp, service_err := o._request_helper("BatchWriteRow", batch_list)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[BatchWriteRow] %s", e)
	}

	return r.(*OTSBatchWriteRowResponse), nil
}

// 说明：根据范围条件获取多行数据。
//
// 		``table_name``是对应的表名。
// 		``direction``表示范围的方向，字符串格式，取值包括'FORWARD'和'BACKWARD'。
// 		``inclusive_start_primary_key``表示范围的起始主键（在范围内）。
// 		``exclusive_end_primary_key``表示范围的结束主键（不在范围内）。
// 		``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为``otstype.OTSColumnsToGet``；如果为nil，表示获取所有列。
// 		``limit``是可选参数，表示最多读取多少行；如果为0，则没有限制。
//
// 		返回：符合条件的结果列表。
// 		      错误信息。
//
// 		``response_row_list``为``otstype.OTSGetRangeResponse``类的实例包含了：
// 		``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
// 		``NextStartPrimaryKey``表示下次get_range操作的起始点的主健列，类型为``otstype.OTSPrimaryKey``。
// 		``Rows``表示本次操作返回的行数据列表，是``otstype.OTSRows``类的实例。
//
// 		示例：
//
// 		// get_range
// 		// 查询区间：[(1, INF_MIN), (4, INF_MAX))，左闭右开。
// 		inclusive_start_primary_key := &OTSPrimaryKey{
// 			"gid": 1,
// 			"uid": OTSColumnType_INF_MIN,
// 		}
// 		exclusive_end_primary_key := &OTSPrimaryKey{
// 			"gid": 4,
// 			"uid": OTSColumnType_INF_MAX,
// 		}
// 		columns_to_get := &OTSColumnsToGet{
// 			"gid", "uid", "name", "address", "mobile", "age",
// 		}
//
// 		// 选择方向
// 		// OTSDirection_FORWARD
// 		// OTSDirection_BACKWARD
// 		response_row_list, ots_err := ots_client.GetRange("myTable", OTSDirection_FORWARD,
// 			inclusive_start_primary_key, exclusive_end_primary_key, columns_to_get, 100)
//
func (o *OTSClient) GetRange(table_name string, direction string,
	inclusive_start_primary_key *OTSPrimaryKey,
	exclusive_end_primary_key *OTSPrimaryKey,
	columns_to_get *OTSColumnsToGet,
	limit int32) (response_row_list *OTSGetRangeResponse, err *OTSError) {
	err = new(OTSError)
	if table_name == "" {
		return nil, err.SetClientMessage("[GetRange] table_name should not be empty")
	}
	if direction != OTSDirection_FORWARD && direction != OTSDirection_BACKWARD {
		return nil, err.SetClientMessage("[GetRange] direction should be FORWARD or BACKWARD")
	}
	if exclusive_end_primary_key == nil {
		return nil, err.SetClientMessage("[GetRange] exclusive_end_primary_key should not be nil")
	}

	resp, service_err := o._request_helper("GetRange", table_name, direction, inclusive_start_primary_key, exclusive_end_primary_key, columns_to_get, limit)
	if service_err != nil {
		return nil, err.SetServiceError(service_err)
	}

	r, e := o._check_request_helper_error(resp)
	if e != nil {
		return nil, err.SetClientMessage("[GetRange] %s", e)
	}

	return r.(*OTSGetRangeResponse), nil
}

// func (o *OTSClient) XGetRange() {
//
// }

func (o *OTSClient) Version() string {
	return "ots_golang_sdk_" + VERSION
}
