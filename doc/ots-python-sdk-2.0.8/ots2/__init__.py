# -*- coding: utf8 -*-

r"""
开放结构化数据服务（Open Table Service，OTS，http://www.aliyun.com/product/ots）是构建在飞天大规模分布式计算系统之上的海量结构化和半结构化数据存储与实时查询的服务。它通过RESTful API来提供服务，并且也有方便的WEB GUI。

This document is encoded in UTF-8.

OTS Python SDK提供了通过Python API访问OTS的方式。它实现了连接池管理，链接超时，日志输出等功能。

这个SDK的文档包含了对每个接口的详细说明和使用样例。而作为开头，这里有一个建表，插入一行，并读取的例子。

    from ots2 import *

    client = OTSClient(
        'http://your_ots_address/',
        'your_access_id',
        'your_access_key',
        'your_instance_name',
    )

    table_meta = TableMeta(
        'sample_table',                                             # table name
        [('sample_pk1', 'STRING'), ('sample_pk2', 'INTEGER')],      # primary key schema
    )
    capacity_unit = CapacityUnit(
        100,                                                        # read capacity
        100                                                         # write capacity
    )
    reserved_throughput = ReservedThroughput(capacity_unit)
    client.create_table(table_meta, reserved_throughput)

    # you probably should wait for a while until the table partition is loaded

    condition = Condition('IGNORE')
    consumed = client.put_row(
        'sample_table',                                             # table name
        condition,                                                  # condition
        {'sample_pk1':'Hangzhou', 'sample_pk2':123456},             # primary key
        {'sample_col1':True, 'sample_col2':3.14}                    # attribute_columns
    )

    #return:
    #   consumed                                                    # consumed capacity unit

    consumed, primary_key_columns, attribute_columns = client.get_row(
        'sample_table',                                             # table name
        {'sample_pk1':'Hangzhou', 'sample_pk2':123456},             # primary key
        ['sample_pk1', 'sample_col1', 'sample_col2']                # columns to get
    )

    #return: 
    #   consumed                                                    # consumed capacity unit
    #   {'sample_pk1':'Hangzhou'},                                  # primary key columns
    #   {'sample_col1':True, 'sample_col2':3.14}                    # attribute_columns

请注意：

    (1) 用python doc的console模式去查看本文档时，比较长的行会出现显示问题。遇到这样的情况，请用PageUp/PageDown键进行翻页。

"""
__version__ = '2.0.5'
__all__ = [
    'OTSClient',

    # Data Types
    'INF_MIN',
    'INF_MAX',
    'TableMeta',
    'CapacityUnit',
    'ReservedThroughput',
    'ReservedThroughputDetails',
    'UpdateTableResponse',
    'DescribeTableResponse',
    'RowDataItem',
    'Condition',
    'PutRowItem',
    'UpdateRowItem',
    'DeleteRowItem',
    'BatchWriteRowResponseItem',
    'OTSClientError',
    'OTSServiceError',
    'DefaultRetryPolicy',
]

__author__ = 'Haowei YAO <haowei.yao@aliyun-inc.com>; Kunpeng HAN <kunpeng.hkp@aliyun-inc.com>'

from ots2.client import OTSClient

from ots2.metadata import *
from ots2.error import *
from ots2.retry import *

