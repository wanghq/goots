# -*- coding: utf8 -*-

from example_config import *
from ots2 import *
import time

table_name = 'TableOperationsExample'

def create_table(ots_client):
    schema_of_primary_key = [('gid', 'INTEGER'), ('uid', 'INTEGER')]
    table_meta = TableMeta(table_name, schema_of_primary_key)
    reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
    ots_client.create_table(table_meta, reserved_throughput)
    print 'Table has been created.'

def list_table(ots_client):
    tables = ots_client.list_table()
    print 'All the tables you have created:'
    for table in tables:
        print table

def describe_table(ots_client):
    describe_response = ots_client.describe_table(table_name)
    print u'TableName: %s' % describe_response.table_meta.table_name
    print u'PrimaryKey: %s' % describe_response.table_meta.schema_of_primary_key
    print u'Reserved read throughput: %s' % describe_response.reserved_throughput_details.capacity_unit.read
    print u'Reserved write throughput: %s' % describe_response.reserved_throughput_details.capacity_unit.write
    print u'Last increase throughput time: %s' % describe_response.reserved_throughput_details.last_increase_time
    print u'Last decrease throughput time: %s' % describe_response.reserved_throughput_details.last_decrease_time
    print u'Total decrease count in an UTC day: %s' % describe_response.reserved_throughput_details.number_of_decreases_today
        
def update_table(ots_client):
    time.sleep(60) # you should wait at least 1 minute after table created to do updating.
    reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
    update_response = ots_client.update_table(table_name, reserved_throughput)
    print u'Reserved read throughput: %s' % update_response.reserved_throughput_details.capacity_unit.read
    print u'Reserved write throughput: %s' % update_response.reserved_throughput_details.capacity_unit.write
    print u'Last increase throughput time: %s' % update_response.reserved_throughput_details.last_increase_time
    print u'Last decrease throughput time: %s' % update_response.reserved_throughput_details.last_decrease_time
    print u'Total decrease count in an UTC day: %s' % update_response.reserved_throughput_details.number_of_decreases_today

def delete_table(ots_client):
    ots_client.delete_table(table_name)
    print 'Table \'%s\' has been deleted.' % table_name

if __name__ == '__main__':
    ots_client = OTSClient(OTS_ENDPOINT, OTS_ID, OTS_SECRET, OTS_INSTANCE)
    create_table(ots_client)
    list_table(ots_client)
    describe_table(ots_client)
    update_table(ots_client)
    delete_table(ots_client)

