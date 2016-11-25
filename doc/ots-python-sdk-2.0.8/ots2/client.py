# -*- coding: utf8 -*-
# Implementation of OTSClient

__all__ = ['OTSClient']
__author__ = 'Haowei YAO<haowei.yao@aliyun-inc.com>, Kunpeng HAN<kunpeng.hkp@aliyun-inc.com>'

import sys
import logging
import urlparse
import time
import _strptime

from ots2.error import *
from ots2.protocol import OTSProtocol
from ots2.connection import ConnectionPool
from ots2.metadata import *
from ots2.retry import DefaultRetryPolicy


class OTSClient(object):
    """
    ``OTSClient``实现了OTS服务的所有接口。用户可以通过创建``OTSClient``的实例，并调用它的
    方法来访问OTS服务的所有功能。用户可以在初始化方法``__init__()``中设置各种权限、连接等参数。

    除非另外说明，``OTSClient``的所有接口都以抛异常的方式处理错误(请参考模块``ots.error``
    )，即如果某个函数有返回值，则会在描述中说明；否则返回None。
    """


    DEFAULT_ENCODING = 'utf8'
    DEFAULT_SOCKET_TIMEOUT = 50
    DEFAULT_MAX_CONNECTION = 50
    DEFAULT_LOGGER_NAME = 'ots2-client'

    protocol_class = OTSProtocol
    connection_pool_class = ConnectionPool 

    def __init__(self, end_point, accessid, accesskey, instance_name, **kwargs):
        """
        初始化``OTSClient``实例。

        ``end_point``是OTS服务的地址（例如 'http://instance.cn-hangzhou.ots.aliyun.com:80'），必须以'http://'开头。

        ``accessid``是访问OTS服务的accessid，通过官方网站申请或通过管理员获取。

        ``accesskey``是访问OTS服务的accesskey，通过官方网站申请或通过管理员获取。

        ``instance_name``是要访问的实例名，通过官方网站控制台创建或通过管理员获取。

        ``encoding``请求参数的字符串编码类型，默认是utf8。

        ``socket_timeout``是连接池中每个连接的Socket超时，单位为秒，可以为int或float。默认值为50。

        ``max_connection``是连接池的最大连接数。默认为50，

        ``logger_name``用来在请求中打DEBUG日志，或者在出错时打ERROR日志。

        ``retry_policy``定义了重试策略，默认的重试策略为 DefaultRetryPolicy。你可以继承 RetryPolicy 来实现自己的重试策略，请参考 DefaultRetryPolicy 的代码。


        示例：创建一个OTSClient实例

            from ots2.client import OTSClient

            ots_client = OTSClient('your_instance_endpoint', 'your_user_id', 'your_user_key', 'your_instance_name')
        """

        self.encoding = kwargs.get('encoding')
        if self.encoding is None:
            self.encoding = OTSClient.DEFAULT_ENCODING

        self.socket_timeout = kwargs.get('socket_timeout')
        if self.socket_timeout is None:
            self.socket_timeout = OTSClient.DEFAULT_SOCKET_TIMEOUT

        self.max_connection = kwargs.get('max_connection')
        if self.max_connection is None:
            self.max_connection = OTSClient.DEFAULT_MAX_CONNECTION

        # initialize logger
        logger_name = kwargs.get('logger_name')
        if logger_name is None:
            self.logger = logging.getLogger(OTSClient.DEFAULT_LOGGER_NAME)
            nullHandler = logging.NullHandler()
            self.logger.addHandler(nullHandler)
        else:
            self.logger = logging.getLogger(logger_name)

        # parse end point
        scheme, netloc, path = urlparse.urlparse(end_point)[:3]
        host = scheme + "://" + netloc

        if scheme != 'http' and scheme != 'https':
            raise OTSClientError(
                "protocol of end_point must be 'http' or 'https', e.g. http://ots.aliyuncs.com:80."
            )
        if host == '':
            raise OTSClientError(
                "host of end_point should be specified, e.g. http://ots.aliyuncs.com:80."
            )

        # intialize protocol instance via user configuration
        self.protocol = self.protocol_class(
            accessid, accesskey, instance_name, self.encoding, self.logger
        )
        
        # initialize connection via user configuration
        self.connection = self.connection_pool_class(
            host, path, timeout=self.socket_timeout, maxsize=self.max_connection,
        )

        # initialize the retry policy
        retry_policy = kwargs.get('retry_policy')
        if retry_policy is None:
            retry_policy = DefaultRetryPolicy()
        self.retry_policy = retry_policy

    def _request_helper(self, api_name, *args, **kwargs):

        query, reqheaders, reqbody = self.protocol.make_request(
            api_name, *args, **kwargs
        )

        retry_times = 0

        while True:

            try:
                status, reason, resheaders, resbody = self.connection.send_receive(
                    query, reqheaders, reqbody
                )
                self.protocol.handle_error(api_name, query, status, reason, resheaders, resbody)
                break

            except OTSServiceError as e:

                if self.retry_policy.should_retry(retry_times, e, api_name):
                    retry_delay = self.retry_policy.get_retry_delay(retry_times, e, api_name)
                    time.sleep(retry_delay)
                    retry_times += 1
                else:
                    raise e

        ret = self.protocol.parse_response(api_name, status, resheaders, resbody)

        return ret

    def create_table(self, table_meta, reserved_throughput):
        """
        说明：根据表信息创建表。

        ``table_meta``是``ots.metadata.TableMeta``类的实例，它包含表名和PrimaryKey的schema，
        请参考``TableMeta``类的文档。当创建了一个表之后，通常要等待1分钟时间使partition load
        完成，才能进行各种操作。
        ``reserved_throughput``是``ots.metadata.ReservedThroughput``类的实例，表示预留读写吞吐量。

        返回：无。

        示例：

            schema_of_primary_key = [('gid', 'INTEGER'), ('uid', 'INTEGER')]
            table_meta = TableMeta('myTable', schema_of_primary_key)
            reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
            ots_client.create_table(table_meta, reserved_throughput)
        """

        self._request_helper('CreateTable', table_meta, reserved_throughput)

    def delete_table(self, table_name):
        """
        说明：根据表名删除表。

        ``table_name``是对应的表名。

        返回：无。

        示例：

            ots_client.delete_table('myTable')
        """

        self._request_helper('DeleteTable', table_name)

    def list_table(self):
        """
        说明：获取所有表名的列表。
        
        返回：表名列表。

        ``table_list``表示获取的表名列表，类型为tuple，如：('MyTable1', 'MyTable2')。

        示例：

            table_list = ots_client.list_table()
        """

        table_names = self._request_helper('ListTable')
        return table_names

    def update_table(self, table_name, reserved_throughput):
        """ 
        说明：更新表属性，目前只支持修改预留读写吞吐量。
        
        ``table_name``是对应的表名。
        ``reserved_throughput``是``ots2.metadata.ReservedThroughput``类的实例，表示预留读写吞吐量。

        返回：针对该表的预留读写吞吐量的最近上调时间、最近下调时间和当天下调次数。

        ``update_table_response``表示更新的结果，是ots2.metadata.UpdateTableResponse类的实例。

        示例：

            reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
            update_response = ots_client.update_table('myTable', reserved_throughput)
        """

        update_table_response = self._request_helper(
                    'UpdateTable', table_name, reserved_throughput
        )
        return update_table_response

    def describe_table(self, table_name):
        """
        说明：获取表的描述信息。

        ``table_name``是对应的表名。

        返回：表的描述信息。

        ``describe_table_response``表示表的描述信息，是ots2.metadata.DescribeTableResponse类的实例。

        示例：

            describe_table_response = ots_client.describe_table('myTable')
        """

        describe_table_response = self._request_helper('DescribeTable', table_name)
        return describe_table_response

    def get_row(self, table_name, primary_key, columns_to_get=None):
        """
        说明：获取一行数据。

        ``table_name``是对应的表名。
        ``primary_key``是主键，类型为dict。
        ``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为list；如果不填，表示获取所有列。

        返回：本次操作消耗的CapacityUnit、主键列和属性列。

        ``consumed``表示消耗的CapacityUnit，是ots2.metadata.CapacityUnit类的实例。
        ``primary_key_columns``表示主键列，类型为dict，如：{'PK0':value0, 'PK1':value1}。
        ``attribute_columns``表示属性列，类型为dict，如：{'COL0':value0, 'COL1':value1}。

        示例：

            primary_key = {'gid':1, 'uid':101}
            columns_to_get = ['name', 'address', 'age']
            consumed, primary_key_columns, attribute_columns = ots_client.get_row('myTable', primary_key, columns_to_get)
        """

        (consumed, primary_key_columns, attribute_columns) = self._request_helper(
                    'GetRow', table_name, primary_key, columns_to_get
        )
        return consumed, primary_key_columns, attribute_columns

    def put_row(self, table_name, condition, primary_key, attribute_columns):
        """
        说明：写入一行数据。返回本次操作消耗的CapacityUnit。

        ``table_name``是对应的表名。
        ``condition``表示执行操作前做条件检查，满足条件才执行，是ots2.metadata.Condition类的实例。
        目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
        ``primary_key``表示主键，类型为dict。
        ``attribute_columns``表示属性列，类型为dict。

        返回：本次操作消耗的CapacityUnit。

        consumed表示消耗的CapacityUnit，是ots2.metadata.CapacityUnit类的实例。

        示例：

            primary_key = {'gid':1, 'uid':101}
            attribute_columns = {'name':'张三', 'mobile':111111111, 'address':'中国A地', 'age':20}
            condition = Condition('EXPECT_NOT_EXIST')
            consumed = ots_client.put_row('myTable', condition, primary_key, attribute_columns)
        """

        consumed = self._request_helper(
                    'PutRow', table_name, condition, primary_key, attribute_columns
        )
        return consumed
    
    def update_row(self, table_name, condition, primary_key, update_of_attribute_columns):
        """
        说明：更新一行数据。

        ``table_name``是对应的表名。
        ``condition``表示执行操作前做条件检查，满足条件才执行，是ots2.metadata.Condition类的实例。
        目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
        ``primary_key``表示主键，类型为dict。
        ``update_of_attribute_columns``表示属性列，类型为dict，可以包含put和delete操作。其中put是dict
        表示属性列的写入；delete是list，表示要删除的属性列的列名，见示例。

        返回：本次操作消耗的CapacityUnit。

        consumed表示消耗的CapacityUnit，是ots2.metadata.CapacityUnit类的实例。

        示例：

            primary_key = {'gid':1, 'uid':101}
            update_of_attribute_columns = {
                'put' : {'name':'张三丰', 'address':'中国B地'},
                'delete' : ['mobile', 'age'],
            }
            condition = Condition('EXPECT_EXIST')
            consumed = ots_client.update_row('myTable', condition, primary_key, update_of_attribute_columns) 
        """

        consumed = self._request_helper(
                    'UpdateRow', table_name, condition, primary_key, update_of_attribute_columns 
        )
        return consumed

    def delete_row(self, table_name, condition, primary_key):
        """
        说明：删除一行数据。

        ``table_name``是对应的表名。
        ``condition``表示执行操作前做条件检查，满足条件才执行，是ots2.metadata.Condition类的实例。
        目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
        ``primary_key``表示主键，类型为dict。

        返回：本次操作消耗的CapacityUnit。

        consumed表示消耗的CapacityUnit，是ots2.metadata.CapacityUnit类的实例。

        示例：

            primary_key = {'gid':1, 'uid':101}
            condition = Condition('IGNORE')
            consumed = ots_client.delete_row('myTable', condition, primary_key) 
        """

        consumed = self._request_helper(
                    'DeleteRow', table_name, condition, primary_key 
        )
        return consumed

    def batch_get_row(self, batch_list):
        """
        说明：批量获取多行数据。

        ``batch_list``表示获取多行的条件列表，格式如下：
        [
            (table_name0, [row0_primary_key, row1_primary_key, ...], [column_name0, column_name1, ...]),
            (table_name1, [row0_primary_key, row1_primary_key, ...], [column_name0, column_name1, ...])
            ...
        ]
        其中，row0_primary_key, row1_primary_key为主键，类型为dict。

        返回：对应行的结果列表。

        ``response_rows_list``为返回的结果列表，与请求的顺序一一对应，格式如下：
        [
            [row_data_item0, row_data_item1, ...],      # for table_name0
            [row_data_item0, row_data_item1, ...],      # for table_name1
            ...
        ]
        其中，row_data_item0, row_data_item1为ots2.metadata.RowDataItem的实例。

        示例：

            row1_primary_key = {'gid':1, 'uid':101}
            row2_primary_key = {'gid':2, 'uid':202}
            row3_primary_key = {'gid':3, 'uid':303}
            columns_to_get = ['name', 'address', 'mobile', 'age']
            batch_list = [('myTable', [row1_primary_key, row2_primary_key, row3_primary_key], columns_to_get)]
            batch_list = [('notExistTable', [row1_primary_key, row2_primary_key, row3_primary_key], columns_to_get)]
            batch_get_response = ots_client.batch_get_row(batch_list) 
        """

        response_rows_list = self._request_helper('BatchGetRow', batch_list)
        return response_rows_list

    def batch_write_row(self, batch_list):
        """
        说明：批量修改多行数据。

        ``batch_list``表示获取多行的条件列表，格式如下：
        [
            {
                'table_name':table_name0,
                'put':[put_row_item, ...],
                'update':[update_row_item, ...],
                'delete':[delete_row_item, ..]
            },
            {
                'table_name':table_name1,
                'put':[put_row_item, ...],
                'update':[update_row_item, ...],
                'delete':[delete_row_item, ..]
            },
            ...
        ]
        其中，put_row_item, 是ots2.metadata.PutRowItem类的实例；
              update_row_item, 是ots2.metadata.UpdateRowItem类的实例；
              delete_row_item, 是ots2.metadata.DeleteRowItem类的实例。

        返回：对应行的修改结果列表。

        ``response_items_list``为返回的结果列表，与请求的顺序一一对应，格式如下：
        [
            {                                       # for table_name0
                'put':[put_row_resp, ...],
                'update':[update_row_resp, ...],
                'delete':[delete_row_resp, ..])
            },
            {                                       # for table_name1
                'put':[put_row_resp, ...],
                'update':[update_row_resp, ...],
                'delete':[delete_row_resp, ..]
            },
            ...
        ]
        其中put_row_resp，update_row_resp和delete_row_resp都是ots2.metadata.BatchWriteRowResponseItem类的实例。

        示例：

            primary_key = {'gid':2, 'uid':202}
            attribute_columns = {'name':'李四', 'address':'中国某地', 'age':20}
            condition = Condition('EXPECT_NOT_EXIST')
            put_row_item = PutRowItem(condition, primary_key, attribute_columns)
            
            primary_key = {'gid':3, 'uid':303}
            condition = Condition('IGNORE')
            update_of_attribute_columns = {
                'put' : {'name':'张三', 'address':'中国某地'},
                'delete' : ['mobile', 'age'],
            }
            update_row_item = UpdateRowItem(condition, primary_key, update_of_attribute_columns)
            
            primary_key = {'gid':4, 'uid':404}
            condition = Condition('IGNORE')
            delete_row_item = DeleteRowItem(condition, primary_key)
            
            table_item1  = {'table_name':'myTable', 'put':[put_row_item], 'update':[update_row_item], 'delete':[delete_row_item]}
            table_item2  = {'table_name':'notExistTable', 'put':[put_row_item], 'update':[update_row_item], 'delete':[delete_row_item]}
            batch_list = [table_item1, table_item2]
            batch_write_response = ots_client.batch_write_row(batch_list) 
        """

        response_item_list = self._request_helper('BatchWriteRow', batch_list)
        return response_item_list

    def get_range(self, table_name, direction, 
                  inclusive_start_primary_key, 
                  exclusive_end_primary_key, 
                  columns_to_get=None, limit=None):
        """
        说明：根据范围条件获取多行数据。

        ``table_name``是对应的表名。
        ``direction``表示范围的方向，字符串格式，取值包括'FORWARD'和'BACKWARD'。
        ``inclusive_start_primary_key``表示范围的起始主键（在范围内）。
        ``exclusive_end_primary_key``表示范围的结束主键（不在范围内）。
        ``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为list；如果不填，表示获取所有列。
        ``limit``是可选参数，表示最多读取多少行；如果不填，则没有限制。

        返回：符合条件的结果列表。

        ``consumed``表示本次操作消耗的CapacityUnit，是ots2.metadata.CapacityUnit类的实例。
        ``next_start_primary_key``表示下次get_range操作的起始点的主健列，类型为dict。
        ``row_list``表示本次操作返回的行数据列表，格式为：[(primary_key_columns，attribute_columns), ...]。

        示例：

            inclusive_start_primary_key = {'gid':1, 'uid':INF_MIN} 
            exclusive_end_primary_key = {'gid':4, 'uid':INF_MAX} 
            columns_to_get = ['name', 'address', 'mobile', 'age']
            consumed, next_start_primary_key, row_list = ots_client.get_range(
                        'myTable', 'FORWARD', 
                        inclusive_start_primary_key, exclusive_end_primary_key,
                        columns_to_get, 100
            )
        """

        (consumed, next_start_primary_key, row_list) = self._request_helper(
                    'GetRange', table_name, direction, 
                    inclusive_start_primary_key, exclusive_end_primary_key,
                    columns_to_get, limit
        )
        return consumed, next_start_primary_key, row_list

    def xget_range(self, table_name, direction,
                   inclusive_start_primary_key,
                   exclusive_end_primary_key, consumed_counter,
                   columns_to_get=None, count=None):
        """
        说明：根据范围条件获取多行数据，iterator版本。

        ``table_name``是对应的表名。
        ``direction``表示范围的方向，字符串格式，取值包括'FORWARD'和'BACKWARD'。
        ``inclusive_start_primary_key``表示范围的起始主键（在范围内）。
        ``exclusive_end_primary_key``表示范围的结束主键（不在范围内）。
        ``consumed_counter``用于消耗的CapacityUnit统计，是ots2.metadata.CapacityUnit类的实例。
        ``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为list；如果不填，表示获取所有列。
        ``count``是可选参数，表示最多读取多少行；如果不填，则尽量读取整个范围内的所有行。

        返回：符合条件的结果列表。

        ``range_iterator``用于获取符合范围条件的行数据的iterator，每次取出的元素格式为：
        (primary_key_columns，attribute_columns)。其中，primary_key_columns为主键列，dict类型，
        attribute_columns为属性列，dict类型。其它用法见iter类型说明。

        示例：

            consumed_counter = CapacityUnit(0, 0)
            inclusive_start_primary_key = {'gid':1, 'uid':INF_MIN} 
            exclusive_end_primary_key = {'gid':4, 'uid':INF_MAX} 
            columns_to_get = ['name', 'address', 'mobile', 'age']
            range_iterator = client.xget_range(
                        'myTable', 'FORWARD', 
                        inclusive_start_primary_key, exclusive_end_primary_key,
                        consumed_counter, columns_to_get, 100
            )
            for row in range_iterator:
               pass 
        """

        if not isinstance(consumed_counter, CapacityUnit):
            raise OTSClientError(
                "consumed_counter should be an instance of CapacityUnit, not %s" % (
                    consumed_counter.__class__.__name__)
            )
        left_count = None
        if count is not None:
            if count <= 0:
                raise OTSClientError("the value of count must be larger than 0")
            left_count = count

        consumed_counter.read = 0
        consumed_counter.write = 0
        next_start_pk = inclusive_start_primary_key
        while next_start_pk:
            consumed, next_start_pk, row_list = self.get_range(
                table_name, direction,
                next_start_pk, exclusive_end_primary_key, 
                columns_to_get, left_count
            )
            consumed_counter.read += consumed.read
            for row in row_list:
                yield row
                if left_count is not None:
                    left_count -= 1
                    if left_count <= 0:
                        return 

