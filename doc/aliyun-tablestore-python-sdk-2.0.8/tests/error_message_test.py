# -*- coding: utf8 -*-
import unittest
import time
import urllib
import urlparse
import base64
import hashlib

from lib.ots2_api_test_base import OTS2APITestBase
import lib.restriction
import lib.test_config as test_config
from ots2 import *
from ots2.error import *
from ots2.protocol import OTSProtocol
from ots2.connection import ConnectionPool
from ots2.protobuf.encoder import OTSProtoBufferEncoder
import ots2.protobuf.ots_protocol_2_pb2 as pb2

class ByPassHeaderCheckProtocol(OTSProtocol):

    def _check_headers(self, headers, body, status=None):
        pass

    def _check_authorization(self, query, headers, status=None):
        pass


class ErrorMessageTest(OTS2APITestBase):

    def _get_missed_header_client(self, missed_header):
    
        class MissingHeaderProtocol(ByPassHeaderCheckProtocol):
            def _make_headers(self, body, query):
                headers = OTSProtocol._make_headers(self, body, query)
                del headers[missed_header]
                if missed_header != 'x-ots-signature':
                    signature = self._make_request_signature(query, headers)
                    headers['x-ots-signature'] = signature
                return headers

        class MissingHeaderClient(OTSClient):
            protocol_class = MissingHeaderProtocol

        client = MissingHeaderClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        return client

    def test_missing_header(self):
        """请求中缺少某个头，期望返回OTSParameterInvalid"""

        headers = [
            'x-ots-date', 
            'x-ots-contentmd5', 
            'x-ots-signature',
            'x-ots-accesskeyid',
            'x-ots-instancename', 
        ]
        for missed_header in headers:

            client = self._get_missed_header_client(missed_header)

            try:
                client.list_table()
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, 'OTSMissingHeader', "Missing header: '%s'." % missed_header)

    def test_missing_apiversion_in_header(self):

        client = self._get_missed_header_client('x-ots-apiversion')
        try:
            client.list_table()
            self.assert_false()
        except OTSClientError as e:
            self.assertEqual(e.http_status, 403)


    def test_invalid_http_method(self):
        for method in ['PUT', 'DELETE', 'CONNECT', 'HEAD', 'TRACE']:
            class WrongHTTPMethodConnection(ConnectionPool):

                def send_receive(self, url, request_headers, request_body):
                    response = self.pool.urlopen(
                        method, self.host + self.path + url, 
                        body=request_body,
                        headers=request_headers,
                        redirect=True,
                        assert_same_host=False,
                    )
             
                    response_headers = dict(response.getheaders())
                    response_body = response.data
             
                    return response.status, response.reason, response_headers, response_body

            class WrongHTTPMethodProtocol(ByPassHeaderCheckProtocol):

                def _make_request_signature(self, query, headers):
                    uri, param_string, query_string = urlparse.urlparse(query)[2:5]
                    query_pairs = urlparse.parse_qsl(query_string)
                    sorted_query = urllib.urlencode(sorted(query_pairs))
                    # the original use 'POST' method in hard code
                    signature_string = uri + '\n' + method + '\n' + sorted_query + '\n'
             
                    headers_string = self._make_headers_string(headers)
                    signature_string += headers_string + '\n'
                    signature = self._call_signature_method(signature_string)
                    return signature 

            class WrongHTTPMethodClient(OTSClient):

                connection_pool_class = WrongHTTPMethodConnection
                protocol_class = WrongHTTPMethodProtocol

            client = WrongHTTPMethodClient(
                test_config.OTS_ENDPOINT,
                test_config.OTS_ID,
                test_config.OTS_SECRET,
                test_config.OTS_INSTANCE
            )

            try:
                client.list_table()
                self.assert_false()
            except OTSClientError as e:
                self.assertEqual(method, 'TRACE')
                self.assertEqual(e.http_status, 405)
            except OTSServiceError as e:
                if method == 'HEAD':
                    self.assertEqual(e.http_status, 405)
                else:
                    self.assert_error(e, 405, 'OTSMethodNotAllowed', 'Only POST method for requests is supported.')

    def test_access_id_not_exist(self):
        client = OTSClient(
            test_config.OTS_ENDPOINT,
            'blahblahblah',
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )
        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 403, 'OTSAuthFailed', 'The AccessKeyID does not exist.')

    def _test_access_id_disabled(self):
        raise NotImplementedError

    def test_instance_not_found(self):
        client = OTSClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            'blahblahblah'
        )

        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 403, 'OTSAuthFailed', 'The instance is not found.')

    def test_operation_not_supported(self):

        class BadOpConnectionPool(ConnectionPool):
            def send_receive(self, url, request_headers, request_body):
                return ConnectionPool.send_receive(self, '/HelloWorld', request_headers, request_body)

        class BadOpClient(OTSClient):

            connection_pool_class = BadOpConnectionPool

        client = BadOpClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )
        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSUnsupportOperation', "Unsupported operation: 'HelloWorld'.")

    def test_invalid_date_format(self):

        class InvalidDateFormatProtocol(ByPassHeaderCheckProtocol):

            def _make_headers(self, body, query):
                # compose request headers and process request body if needed
         
                md5 = base64.b64encode(hashlib.md5(body).digest())
         
                date = time.strftime('%a, %d %b %Y %H:%M:%S GMT', time.gmtime())
                 
                headers = {
                    'x-ots-date' : "blahblah",
                    'x-ots-apiversion' : self.api_version,
                    'x-ots-accesskeyid' : self.user_id,
                    'x-ots-instancename' : self.instance_name,
                    'x-ots-contentmd5' : md5,
                }
         
                signature = self._make_request_signature(query, headers)
                headers['x-ots-signature'] = signature
         
                return headers

        class InvalidDateFormatClient(OTSClient):

            protocol_class = InvalidDateFormatProtocol

        client = InvalidDateFormatClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'Invalid date format: blahblah.')

    def test_date_expired(self):

        class InvalidDateFormatProtocol(OTSProtocol):

            def _make_headers(self, body, query):
                # compose request headers and process request body if needed
         
                md5 = base64.b64encode(hashlib.md5(body).digest())
         
                date = time.strftime('%a, %d %b %Y %H:%M:%S GMT', time.gmtime())
                headers = {
                    'x-ots-date' : "Sat, 07 Jun 2014 14:25:40 GMT",
                    'x-ots-apiversion' : self.api_version,
                    'x-ots-accesskeyid' : self.user_id,
                    'x-ots-instancename' : self.instance_name,
                    'x-ots-contentmd5' : md5,
                }
         
                signature = self._make_request_signature(query, headers)
                headers['x-ots-signature'] = signature
         
                return headers

        class InvalidDateFormatClient(OTSClient):

            protocol_class = InvalidDateFormatProtocol

        client = InvalidDateFormatClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 403, 'OTSAuthFailed', 'Mismatch between system time and x-ots-date: Sat, 07 Jun 2014 14:25:40 GMT.')

    def test_md5_mismatch(self):

        class InvalidDateFormatProtocol(OTSProtocol):

            def _make_headers(self, body, query):
                # compose request headers and process request body if needed
         
                md5 = base64.b64encode(hashlib.md5(body).digest())
         
                date = time.strftime('%a, %d %b %Y %H:%M:%S GMT', time.gmtime())
                headers = {
                    'x-ots-date' : date,
                    'x-ots-apiversion' : self.api_version,
                    'x-ots-accesskeyid' : self.user_id,
                    'x-ots-instancename' : self.instance_name,
                    'x-ots-contentmd5' : 'blahblah',
                }
         
                signature = self._make_request_signature(query, headers)
                headers['x-ots-signature'] = signature
         
                return headers

        class InvalidDateFormatClient(OTSClient):

            protocol_class = InvalidDateFormatProtocol

        client = InvalidDateFormatClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 403, 'OTSAuthFailed', 'Mismatch between MD5 value of request body and x-ots-contentmd5 in header.')

    def test_failed_parse_pb(self):
    
        class BadBodyConnectionPool(ConnectionPool):
            def send_receive(self, url, request_headers, request_body):
                return ConnectionPool.send_receive(self, url, request_headers, request_body[:100])

           
        class BadBodyProtocol(OTSProtocol):

            def _make_headers(self, body, query):
                body = body[:100]
                return OTSProtocol._make_headers(self, body, query)

        class BadBodyClient(OTSClient):

            connection_pool_class = BadBodyConnectionPool
            protocol_class = BadBodyProtocol
 
        
        client = BadBodyClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        try:
            client.delete_table('X' * 200)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'Failed to parse the ProtoBuf message.')


    def test_signature_mismatch(self):
    
        class BadSignatureProtocol(OTSProtocol):

            def _make_headers(self, body, query):
                # compose request headers and process request body if needed
                headers = OTSProtocol._make_headers(self, body, query)
                headers['x-ots-signature'] = 'blahblah'
                return headers

        class BadSignatureClient(OTSClient):

            protocol_class = BadSignatureProtocol

        client = BadSignatureClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 403, 'OTSAuthFailed', 'Signature mismatch.')

    def test_both_read_write_CU_should_set_when_create_table(self):

        class NoReadCUEncoder(OTSProtoBufferEncoder):
            
            def _make_capacity_unit(self, proto, capacity_unit):
                proto.write = self._get_int32(capacity_unit.write)

        class NoWriteCUEncoder(OTSProtoBufferEncoder):
            
            def _make_capacity_unit(self, proto, capacity_unit):
                proto.read = self._get_int32(capacity_unit.read)

        class NoReadCUProtocol(OTSProtocol):
            encoder_class = NoReadCUEncoder

        class NoWriteCUProtocol(OTSProtocol):
            encoder_class = NoWriteCUEncoder

        class NoReadCUClient(OTSClient):
            protocol_class = NoReadCUProtocol

        class NoWriteCUClient(OTSClient):
            protocol_class = NoWriteCUProtocol

        no_read_cu_client = NoReadCUClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )
        
        no_write_cu_client = NoWriteCUClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )


        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('NCVonline', [('PK', 'STRING')])

        try:
            no_read_cu_client.create_table(table_meta, reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'Both read and write capacity unit are required to create table.')

        try:
            no_write_cu_client.create_table(table_meta, reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'Both read and write capacity unit are required to create table.')

    def _test_neither_read_nor_write_is_set_when_update_table(self):
        # N/A
    
        class NoReadCUEncoder(OTSProtoBufferEncoder):
            def _make_update_capacity_unit(self, proto, capacity_unit):
                proto.write = self._get_int32(capacity_unit.write)

        class NoWriteCUEncoder(OTSProtoBufferEncoder):
            def _make_update_capacity_unit(self, proto, capacity_unit):
                proto.read = self._get_int32(capacity_unit.read)
                
        class NoCUEncoder(OTSProtoBufferEncoder):
            def _make_update_capacity_unit(self, proto, capacity_unit):
                pass

        class NoReadCUProtocol(OTSProtocol):
            encoder_class = NoReadCUEncoder

        class NoWriteCUProtocol(OTSProtocol):
            encoder_class = NoWriteCUEncoder
            
        class NoCUProtocol(OTSProtocol):
            encoder_class = NoCUEncoder

        class NoReadCUClient(OTSClient):
            protocol_class = NoReadCUProtocol

        class NoWriteCUClient(OTSClient):
            protocol_class = NoWriteCUProtocol

        class NoCUClient(OTSClient):
            protocol_class = NoCUProtocol

        no_read_cu_client = NoReadCUClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )
        
        no_write_cu_client = NoWriteCUClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        no_cu_client = NoCUClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )

        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('NCVonline', [('PK', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)

        time.sleep(restriction.AdjustCapacityUnitIntervalForTest)
        no_read_cu_client.update_table('NCVonline', reserved_throughput)
        time.sleep(restriction.AdjustCapacityUnitIntervalForTest)
            
        no_write_cu_client.update_table('NCVonline', reserved_throughput)
        time.sleep(restriction.AdjustCapacityUnitIntervalForTest)
        
        try:
            no_cu_client.update_table('NCVonline', reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'Neither read nor writes capacity unit is set.')

    def _test_column_utf8_encoding(self):
        # N/A

        class NoUtf8Encoder(OTSProtoBufferEncoder):

            def _get_unicode(self, value):
                return value

        class NoUtf8Protocol(OTSProtocol):
            encoder_class = NoUtf8Encoder

        class NoUtf8Client(OTSClient):
            protocol_class = NoUtf8Protocol

        client = NoUtf8Client(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )


        client.put_row('T0', Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col' : '中文'.decode('utf8').encode('gb2312')})

    def _test_length_of_column(self):
        # N/A
        raise NotImplementedError

    def _test_length_of_primay_key_column_name(self):
        # N/A
        raise NotImplementedError

    def test_duplicated_primary_key_when_batch_get_row(self):
        """OTSParameterInvalid  Duplicated primary key:  '{PKName}' of getting row #{RowIndex} in table '{TableName}'."""
        try:
            self.client_test.batch_get_row([
                ('T0', [{'PK0' : 'XXXX'}, {'PK0' : '---'}], []),
                ('T1', [{'PK0' : 'XXXX'}, {'PK0' : 'XXXX'}], [])
            ])
            self.assert_false()
        except OTSServiceError as e:
            # self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated primary key: 'PK0' of getting row #1 in table 'T1'.")
            self.assert_error(e, 400, 'OTSParameterInvalid', "The input parameter is invalid.")


    def test_duplicated_primary_key_of_put_when_batch_write_row(self):
        """OTSParameterInvalid  Duplicated primary key:  '{PKName}' of writing row #{RowIndex} in table: '{TableName}'."""
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'put' : [
                        PutRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col' : 'XXXX'}),
                        PutRowItem(Condition('IGNORE'), {'PK0' : '----'}, {'Col' : 'XXXX'}),
                    ],
                },
                {
                    'table_name': 'T1', 
                    'put' : [
                        PutRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col' : 'XXXX'}),
                        PutRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col' : 'XXXX'}),
                    ],
                }
            ])
            self.assert_false()
        except OTSServiceError as e:
            # self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated primary key: 'PK0' of writing row #1 in table 'T1'.")
            self.assert_error(e, 400, 'OTSParameterInvalid', "The input parameter is invalid.")

    def test_duplicated_column_name_of_put_when_batch_write_row(self):
        """OTSParameterInvalid  Duplicated column name with primary key column: '{PKName}' while putting row # {RowIndex} in table: '{TableName}'."""

        class MyDict:
            
            def iteritems(self):
                return [('Col0', 'XXXX'), ('Col0', 'XXXX')]
                
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('T0', [('PK0', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)

        time.sleep(10)

        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'put' : [
                        PutRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col0' : 'XXXX', 'Col1' : 'XXXX'}),
                        PutRowItem(Condition('IGNORE'), {'PK0' : '----'}, MyDict()),
                    ],
                },
            ])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated attribute column name: 'Col0' while putting row #1 in table: 'T0'.")

    
    def test_duplicated_primay_key_of_update_when_batch_write_row(self):
        """OTSParameterInvalid  Duplicate primary key '{PKName}' of updating row #{RowIndex} in table: '{TableName}'."""
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'update' : [
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'put':{'Col' : 'XXXX'}}),
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : '----'}, {'put':{'Col' : 'XXXX'}}),
                    ],
                },
                {
                    'table_name': 'T1', 
                    'update' : [
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'put':{'Col' : 'XXXX'}}),
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'put':{'Col' : 'XXXX'}}),
                    ],
                }
            ])
            self.assert_false()
        except OTSServiceError as e:
            # self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated primary key: 'PK0' of updating row #1 in table 'T1'.")
            self.assert_error(e, 400, 'OTSParameterInvalid', "The input parameter is invalid.")

    def test_duplicated_column_name_of_update_when_batch_write_row(self):
        """OTSParameterInvalid  Duplicated column name with primary key column: '{PKName}' while updating row # {RowIndex} in table: '{TableName}'."""
        
        class MyDict(dict):
            
            def iteritems(self):
                return [('Col0', 'XXXX'), ('Col0', 'XXXX')]
                
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('T0', [('PK0', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)

        time.sleep(10)

        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'update' : [
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'put':{'Col0' : 'XXXX', 'Col1' : 'XXXX'}}),
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : '----'}, {'put': MyDict()}),
                    ],
                },
            ])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated attribute column name: 'Col0' while updating row #1 in table: 'T0'.")

        
    def test_duplicated_primay_key_of_delete_when_batch_write_row(self):
        """OTSParameterInvalid  Duplicated primary key '{PKName}' of deleting row #{RowIndex} in table: '{TableName}'."""
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'delete' : [
                        DeleteRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}),
                        DeleteRowItem(Condition('IGNORE'), {'PK0' : '----'}),
                    ],
                },
                {
                    'table_name': 'T1', 
                    'delete' : [
                        DeleteRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}),
                        DeleteRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}),
                    ],
                }
            ])
            self.assert_false()
        except OTSServiceError as e:
            # self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated primary key: 'PK0' of deleting row #1 in table 'T1'.")
            self.assert_error(e, 400, 'OTSParameterInvalid', "The input parameter is invalid.")

    def test_invalid_condition_of_update_when_batch_write_row(self):
        """OTSParameterInvalid  Invalid condition: {RowExistence} while updating row #{RowIndex} in table : '{TableName}'."""
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'update' : [
                        UpdateRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'put':{'Col' : 'XXXX'}}),
                        UpdateRowItem(Condition('EXPECT_NOT_EXIST'), {'PK0' : '----'}, {'put':{'Col' : 'XXXX'}}),
                    ],
                },
            ])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Invalid condition: EXPECT_NOT_EXIST while updating row #1 in table: 'T0'.")

    def test_invalid_condition_of_delete_when_batch_write_row(self):
        """OTSParameterInvalid  Invalid condition: {RowExistence} while deleting row #{RowIndex} in table : '{TableName}'."""
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'delete' : [
                        DeleteRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}),
                        DeleteRowItem(Condition('EXPECT_NOT_EXIST'), {'PK0' : '----'}),
                    ],
                },
            ])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Invalid condition: EXPECT_NOT_EXIST while deleting row #1 in table: 'T0'.")

    def _test_duplicated_primary_key(self):
        """OTSParameterInvalid  Duplicated primary key column: '{PKName}'."""
        # 所有相关API都要测到
        raise NotImplementedError

    def test_duplicated_column_name_when_put_row(self):
        """OTSParameterInvalid  Duplicated attribute column name with primary key column: '{ColumnName}' while putting row."""
        
        class MyDict:
            def iteritems(self):
                return [('Col0', 'XXXX'), ('Col0', 'XXXX')]
                
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('T0', [('PK0', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)

        time.sleep(10)

        
        try:
            self.client_test.put_row('T0', Condition('IGNORE'), {'PK0' : '----'}, MyDict())
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated attribute column name: 'Col0' while putting row.")

    def test_duplicated_column_name_when_update_row(self):
        """OTSParameterInvalid  Duplicated attribute column name with primary key column: '{ColumnName}' while updating row."""
        
        class MyDict(dict):
            def iteritems(self):
                return [('Col0', 'XXXX'), ('Col0', '----')]
                
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('T0', [('PK0', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)
        time.sleep(10)
        
        try:
            self.client_test.update_row('T0', Condition('IGNORE'), {'PK0' : '----'}, {'put':MyDict()})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Duplicated attribute column name: 'Col0' while updating row.")

    def test_invalid_column_type(self):
        """OTSParameterInvalid  Invalid column type, only STRING|INTEGER|BOOLEAN|DOUBLE|BINARY is allowed."""
        """OTSParameterInvalid  Invalid column type: {ColumnType}."""
       
        try:
            self.client_test.update_row('T0', Condition('IGNORE'), {'PK0' : INF_MIN}, {'put':{'Col0' : INF_MIN}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'INF_MIN is an invalid type for the primary key.')
            
        try:
            self.client_test.update_row('T0', Condition('IGNORE'), {'PK0' : '----'}, {'put':{'Col0' : INF_MIN}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'INF_MIN is an invalid type for the attribute column.')
            
        try:
            self.client_test.update_row('T0', Condition('IGNORE'), {'PK0' : '----'}, {'put':{'Col0' : INF_MAX}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', 'INF_MAX is an invalid type for the attribute column.')


    def test_option_field_not_set_in_column_value(self):

        class MissedFieldEncoder(OTSProtoBufferEncoder):

            def _make_column_value(self, proto, value):
                if isinstance(value, str) or isinstance(value, unicode):
                    string = self._get_unicode(value)
                    proto.type = pb2.STRING
                elif isinstance(value, bool):
                    proto.type = pb2.BOOLEAN
                elif isinstance(value, int) or isinstance(value, long):
                    proto.type = pb2.INTEGER
                elif isinstance(value, float):
                    proto.type = pb2.DOUBLE
                elif isinstance(value, bytearray):
                    proto.type = pb2.BINARY

        class MissedFieldProtocol(OTSProtocol):
            encoder_class = MissedFieldEncoder

        class MissedFieldClient(OTSClient):
            protocol_class = MissedFieldProtocol

        client = MissedFieldClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )
        
        
        class MissedFieldEncoder2(OTSProtoBufferEncoder):

            def _make_column_value(self, proto, value):
                if isinstance(value, str) or isinstance(value, unicode):
                    string = self._get_unicode(value)
                    proto.type = pb2.STRING
                    proto.v_string = string
                elif isinstance(value, bool):
                    proto.type = pb2.BOOLEAN
                elif isinstance(value, int) or isinstance(value, long):
                    proto.type = pb2.INTEGER
                elif isinstance(value, float):
                    proto.type = pb2.DOUBLE
                elif isinstance(value, bytearray):
                    proto.type = pb2.BINARY

        class MissedFieldProtocol2(OTSProtocol):
            encoder_class = MissedFieldEncoder2

        class MissedFieldClient2(OTSClient):
            protocol_class = MissedFieldProtocol2

        client2 = MissedFieldClient2(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            test_config.OTS_INSTANCE
        )


        try:
            client.put_row('T0', Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col0' : 'XXXX'})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Optional field 'v_string' must be set as ColumnType is STRING.")
            
        try:
            client.put_row('T0', Condition('IGNORE'), {'PK0' : 123}, {'Col0' : 'XXXX'})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Optional field 'v_int' must be set as ColumnType is INTEGER.")

        try:
            client2.put_row('T0', Condition('IGNORE'), {'PK0' : 'XXX'}, {'Col0' : 3.14})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Optional field 'v_double' must be set as ColumnType is DOUBLE.")

        try:
            client2.put_row('T0', Condition('IGNORE'), {'PK0' : 'XXXXX'}, {'Col0' : True})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Optional field 'v_bool' must be set as ColumnType is BOOLEAN.")
            
        try:
            client2.put_row('T0', Condition('IGNORE'), {'PK0' : 'XXXXX'}, {'Col0' : bytearray('341324213')})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "Optional field 'v_binary' must be set as ColumnType is BINARY.")


    def test_table_is_not_ready(self):
        """404 OTSTableNotReady The table is not ready."""
        
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        table_meta = TableMeta('T0', [('PK0', 'STRING')])
        self.client_test.create_table(table_meta, reserved_throughput)
        
        try:
            self.client_test.update_row('T0', Condition('IGNORE'), {'PK0' : '----'}, {'put':{'Col0' :  1}})
        except OTSServiceError as e:
            self.assert_error(e, 404, 'OTSTableNotReady', "The table is not ready.")

    def _test_partition_unavailabe(self):
        """503 OTSPartitionUnavailable The partition is not available."""
        raise NotImplementedError

    def test_batch_write_row_data_size_exceeded(self):
        """OTSParameterInvalid The total data size of single BatchWriteRow request exceeded the limit."""
        cell_num = 2 * 1024 / 64
        string = 'X' * 64 * 1024
        
        try:
            self.client_test.batch_write_row([
                {
                    'table_name': 'T0', 
                    'put' : [
                        PutRowItem(Condition('IGNORE'), {'PK0' : 'XXXX'}, {'Col' : string}),
                    ] * cell_num,
                },
            ])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, 'OTSParameterInvalid', "The total data size of single BatchWriteRow request exceeded the limit.")
        except OTSClientError as e:
            self.assertEqual(e.http_status, 413)


    def _test_storage_conflict(self):
        """409 OTSConflict Data is being modified by the other request."""
        raise NotImplementedError

    def _test_server_is_busy(self):
        """503 OTSServerBusy Server is busy."""
        raise NotImplementedError

if __name__ == '__main__':
    unittest.main()
