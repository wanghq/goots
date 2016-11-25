# -*- coding: utf8 -*-

__version__ = '2.0.8'
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


from ots2.client import OTSClient

from ots2.metadata import *
from ots2.error import *
from ots2.retry import *

