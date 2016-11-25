#!/bin/python
# -*- coding: utf8 -*-

import time
import logging
import unittest

from ots2 import *

ENDPOINT = 'your_instance_address'
ACCESSID = 'your_accessid'
ACCESSKEY = 'your_accesskey'
INSTANCENAME = 'your_instance_name'

def main():

    ots_client = OTSClient(ENDPOINT, ACCESSID, ACCESSKEY, INSTANCENAME)

    # delete_table
    try:
        ots_client.delete_table('myTable')
    except OTSServiceError,e:
        print e.get_error_code()
        print e.get_error_message()
    print u'表已删除。'
    
    # create_table
    # 注意：OTS是按设置的ReservedThroughput计量收费，即使没有读写也会产生费用。
    schema_of_primary_key = [('gid', 'INTEGER'), ('uid', 'INTEGER')]
    table_meta = TableMeta('myTable', schema_of_primary_key)
    reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
    ots_client.create_table(table_meta, reserved_throughput)
    print u'表已创建。'

    time.sleep(5)
    
    # list_table
    list_response = ots_client.list_table()
    print u'表的列表如下：'
    for table_name in list_response:
        print table_name
    
    # update_table
    try:
        # 由于刚创建表，需要2分钟之后才能调整表的预留读写吞吐量。
        # 注意：OTS是按设置的ReservedThroughput计量收费，即使没有读写也会产生费用。
        reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
        update_response = ots_client.update_table('myTable', reserved_throughput)
        print u'表的预留读吞吐量：%s' % update_response.reserved_throughput_details.capacity_unit.read
        print u'表的预留写吞吐量：%s' % update_response.reserved_throughput_details.capacity_unit.write
        print u'最后一次上调预留读写吞吐量时间：%s' % update_response.reserved_throughput_details.last_increase_time
        print u'最后一次下调预留读写吞吐量时间：%s' % update_response.reserved_throughput_details.last_decrease_time
        print u'UTC自然日内总的下调预留读写吞吐量次数：%s' % update_response.reserved_throughput_details.number_of_decreases_today
    except OTSServiceError,e:
        print e.get_error_code()
        print e.get_error_message()
    
    # describe_table
    describe_response = ots_client.describe_table('myTable')
    print u'表的名称: %s' % describe_response.table_meta.table_name
    print u'表的主键: %s' % describe_response.table_meta.schema_of_primary_key
    print u'表的预留读吞吐量：%s' % describe_response.reserved_throughput_details.capacity_unit.read
    print u'表的预留写吞吐量：%s' % describe_response.reserved_throughput_details.capacity_unit.write
    print u'最后一次上调预留读写吞吐量时间：%s' % describe_response.reserved_throughput_details.last_increase_time
    print u'最后一次下调预留读写吞吐量时间：%s' % describe_response.reserved_throughput_details.last_decrease_time
    print u'UTC自然日内总的下调预留读写吞吐量次数：%s' % describe_response.reserved_throughput_details.number_of_decreases_today
    
    # put_row
    primary_key = {'gid':1, 'uid':101}
    attribute_columns = {'name':'张三', 'mobile':111111111, 'address':'中国A地', 'age':20}
    condition = Condition('EXPECT_NOT_EXIST')
    consumed = ots_client.put_row('myTable', condition, primary_key, attribute_columns)
    print u'成功插入数据，消耗的写CapacityUnit为：%s' % consumed.write
    
    # get_row
    primary_key = {'gid':1, 'uid':101}
    columns_to_get = ['name', 'address', 'age']
    consumed, primary_key_columns, attribute_columns = ots_client.get_row('myTable', primary_key, columns_to_get)
    print u'成功读取数据，消耗的读CapacityUnit为：%s' % consumed.read
    print u'name信息：%s' % attribute_columns.get('name')
    print u'address信息：%s' % attribute_columns.get('address')
    print u'age信息：%s' % attribute_columns.get('age')

    # update_row
    primary_key = {'gid':1, 'uid':101}
    update_of_attribute_columns = {
        'put' : {'name':'张三丰', 'address':'中国B地'},
        'delete' : ['mobile', 'age'],
    }
    condition = Condition('EXPECT_EXIST')
    consumed = ots_client.update_row('myTable', condition, primary_key, update_of_attribute_columns) 
    print u'成功更新数据，消耗的写CapacityUnit为：%s' % consumed.write
    
    # delete_row
    primary_key = {'gid':1, 'uid':101}
    condition = Condition('IGNORE')
    consumed = ots_client.delete_row('myTable', condition, primary_key) 
    print u'成功删除数据，消耗的写CapacityUnit为：%s' % consumed.write
    
    # batch_write_row
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

    # 每一行操作都是独立的，需要分别判断是否成功。对于失败子操作进行重试。
    retry_count = 0
    operation_list = ['put', 'update', 'delete']
    while retry_count < 3:
        failed_batch_list = []
        for i in range(len(batch_write_response)):
            table_item = batch_write_response[i]
            for operation in operation_list:
                operation_item = table_item.get(operation)
                if not operation_item:
                    continue
                print u'操作：%s' % operation
                for j in range(len(operation_item)):
                    row_item = operation_item[j]
                    print u'操作是否成功：%s' % row_item.is_ok
                    if not row_item.is_ok:
                        print u'错误码：%s' % row_item.error_code
                        print u'错误信息：%s' % row_item.error_message
                        add_batch_write_item(failed_batch_list, batch_list[i]['table_name'], operation, batch_list[i][operation][j])
                    else:
                        print u'本次操作消耗的写CapacityUnit为：%s' % row_item.consumed.write

        if not failed_batch_list:
            break
        retry_count += 1
        batch_list = failed_batch_list
        batch_write_response = ots_client.batch_write_row(batch_list)
    
    # batch_get_row
    row1_primary_key = {'gid':1, 'uid':101}
    row2_primary_key = {'gid':2, 'uid':202}
    row3_primary_key = {'gid':3, 'uid':303}
    columns_to_get = ['name', 'address', 'mobile', 'age']
    batch_list = [('myTable', [row1_primary_key, row2_primary_key, row3_primary_key], columns_to_get)]
    batch_list = [('notExistTable', [row1_primary_key, row2_primary_key, row3_primary_key], columns_to_get)]
    batch_get_response = ots_client.batch_get_row(batch_list) 

    # 每一行操作都是独立的，需要分别判断是否成功。对于失败子操作进行重试。
    retry_count = 0
    while retry_count < 3:
        failed_batch_list = []
        for i in range(len(batch_get_response)):
            table_item = batch_get_response[i]
            for j in range(len(table_item)):
                row_item = table_item[j]
                print u'操作是否成功：%s' % row_item.is_ok
                if not row_item.is_ok:
                    print u'错误码：%s' % row_item.error_code
                    print u'错误信息：%s' % row_item.error_message
                    add_batch_get_item(failed_batch_list, batch_list[i][0], batch_list[i][1][j], batch_list[i][2])
                else:
                    print u'name信息：%s' % row_item.attribute_columns.get('name')
                    print u'address信息：%s' % row_item.attribute_columns.get('address')
                    print u'mobile信息：%s' % row_item.attribute_columns.get('mobile')
                    print u'age信息：%s' % row_item.attribute_columns.get('age')
                    print u'本次操作消耗的读CapacityUnit为：%s' % row_item.consumed.read

        if not failed_batch_list:
            break
        retry_count += 1
        batch_list = failed_batch_list
        batch_get_response = ots_client.batch_get_row(batch_list)
    
    # get_range
    # 查询区间：[(1, INF_MIN), (4, INF_MAX))，左闭右开。
    inclusive_start_primary_key = {'gid':1, 'uid':INF_MIN} 
    exclusive_end_primary_key = {'gid':4, 'uid':INF_MAX} 
    columns_to_get = ['name', 'address', 'mobile', 'age']
    consumed, next_start_primary_key, row_list = ots_client.get_range(
                'myTable', 'FORWARD', 
                inclusive_start_primary_key, exclusive_end_primary_key,
                columns_to_get, 100
    )
    for row in row_list:
        attribute_columns = row[1]
        print u'name信息为：%s' % attribute_columns.get('name')
        print u'address信息为：%s' % attribute_columns.get('address')
        print u'mobile信息为：%s' % attribute_columns.get('mobile')
        print u'age信息为：%s' % attribute_columns.get('age')
    print u'本次操作消耗的读CapacityUnit为：%s' % consumed.read
    print u'下次开始的主键：%s' % next_start_primary_key

    # xget_range
    consumed_counter = CapacityUnit(0, 0)
    inclusive_start_primary_key = {'gid':1, 'uid':INF_MIN} 
    exclusive_end_primary_key = {'gid':4, 'uid':INF_MAX} 
    columns_to_get = ['gid', 'uid', 'name', 'address', 'mobile', 'age']
    range_iter = ots_client.xget_range(
                'myTable', 'FORWARD', 
                inclusive_start_primary_key, exclusive_end_primary_key,
                consumed_counter, columns_to_get, 100
    )
    for (primary_key_columns, attribute_columns) in range_iter:
        print u'gid信息为：%s' % primary_key_columns.get('gid')
        print u'uid信息为：%s' % primary_key_columns.get('uid')
        print u'name信息为：%s' % attribute_columns.get('name')
        print u'address信息为：%s' % attribute_columns.get('address')
        print u'mobile信息为：%s' % attribute_columns.get('mobile')
        print u'age信息为：%s' % attribute_columns.get('age')

def add_batch_write_item(batch_list, table_name, operation, item):
    for table_item in batch_list:
        if table_item.get('table_name') == table_name:
            operation_item = table_item.get(operation)
            if not operation_item:
                table_item[operation] = [item]
            else:
                operation_item.append(item)
            return
    # not found
    table_item = {'table_name':table_name, operation:[item]}
    batch_list.append(table_item)

def add_batch_get_item(batch_list, table_name, item, columns_to_get):
    for table_item in batch_list:
        if table_item[0] == table_name:
            row_item_list = table_item[1]
            row_item_list.append(item)
            return
    # not found
    table_item = (table_name, [item], columns_to_get)
    batch_list.append(table_item)
                
if __name__ == '__main__':
    main()
