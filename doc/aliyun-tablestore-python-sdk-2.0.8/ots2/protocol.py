# -*- coding: utf8 -*-

import hashlib
import urllib
import hmac
import base64
import time
import urlparse
import calendar
import logging
from email.utils import formatdate

import google.protobuf.text_format as text_format

from ots2.error import *
from ots2.protobuf.encoder import OTSProtoBufferEncoder
from ots2.protobuf.decoder import OTSProtoBufferDecoder
import ots2.protobuf.ots_protocol_2_pb2 as pb2


class OTSProtocol:

    api_version = '2014-08-08'

    encoder_class = OTSProtoBufferEncoder
    decoder_class = OTSProtoBufferDecoder

    api_list = {
        'CreateTable',
        'ListTable',
        'DeleteTable',
        'DescribeTable',
        'UpdateTable',
        'GetRow',
        'PutRow',
        'UpdateRow',
        'DeleteRow',
        'BatchGetRow',
        'BatchWriteRow',
        'GetRange'
    }

    def __init__(self, user_id, user_key, instance_name, encoding, logger):
        self.user_id = user_id
        self.user_key = user_key
        self.instance_name = instance_name
        self.encoder = self.encoder_class(encoding)
        self.decoder = self.decoder_class(encoding)
        self.logger = logger

    def _make_headers_string(self, headers):
        headers_item = ["%s:%s" % (k.lower(), v.strip()) for k, v in headers.iteritems() if k.startswith('x-ots-') and k != 'x-ots-signature']
        return "\n".join(sorted(headers_item))

    def _call_signature_method(self, signature_string):
        # The signature method is supposed to be HmacSHA1
        # A switch case is required if there is other methods available
        signature = base64.b64encode(hmac.new(
            self.user_key, signature_string, hashlib.sha1
        ).digest())
        return signature

    def _make_request_signature(self, query, headers):
        uri, param_string, query_string = urlparse.urlparse(query)[2:5]

        # TODO a special query should be input to test query sorting,
        # because none of the current APIs uses query map, but the sorting
        # is required in the protocol document.
        query_pairs = urlparse.parse_qsl(query_string)
        sorted_query = urllib.urlencode(sorted(query_pairs))
        signature_string = uri + '\n' + 'POST' + '\n' + sorted_query + '\n'

        headers_string = self._make_headers_string(headers)
        signature_string += headers_string + '\n'
        signature = self._call_signature_method(signature_string)
        return signature

    def _make_headers(self, body, query):
        # compose request headers and process request body if needed

        md5 = base64.b64encode(hashlib.md5(body).digest())

        date = formatdate(time.time(), usegmt=True)
        
        headers = {
            'x-ots-date' : date,
            'x-ots-apiversion' : self.api_version,
            'x-ots-accesskeyid' : self.user_id,
            'x-ots-instancename' : self.instance_name,
            'x-ots-contentmd5' : md5,
        }

        signature = self._make_request_signature(query, headers)
        headers['x-ots-signature'] = signature
        headers['User-Agent'] = "aliyun-sdk-python 2.0.6"
        return headers

    def _make_response_signature(self, query, headers):
        uri = urlparse.urlparse(query)[2]
        headers_string = self._make_headers_string(headers)

        signature_string = headers_string + '\n' + uri
        signature = self._call_signature_method(signature_string)
        return signature

    def _convert_urllib3_headers(self, headers):
        """
        old urllib3 headers: {'header1':'value1', 'header2':'value2'} 
        new urllib3 headers: {'header1':('header1', 'value1'), 'header2':('header2', 'value2')} 
        """
        std_headers = {}
        for k,v in headers.iteritems():
            if isinstance(v, tuple) and len(v) == 2:
                std_headers[k.lower()] = v[1]
            else:
                std_headers[k.lower()] = v

        return std_headers

    def _check_headers(self, headers, body, status=None):
        # check the response headers and process response body if needed.

        # 1, make sure we have all headers
        header_names = [
            'x-ots-contentmd5', 
            'x-ots-requestid', 
            'x-ots-date', 
            'x-ots-contenttype',
        ]

        if status >= 200 and status < 300:
            for name in header_names:
                if not name in headers:
                    raise OTSClientError('"%s" is missing in response header.' % name)

        # 2, check md5
        if 'x-ots-contentmd5' in headers:
            md5 = base64.b64encode(hashlib.md5(body).digest())
            if md5 != headers['x-ots-contentmd5']:
                raise OTSClientError('MD5 mismatch in response.')

        # 3, check date 
        if 'x-ots-date' in headers:
            try:
                server_time = time.strptime(headers['x-ots-date'], '%a, %d %b %Y %H:%M:%S %Z')
            except ValueError:
                raise OTSClientError('Invalid date format in response.')
        
            # 4, check date range
            server_unix_time = calendar.timegm(server_time)
            now_unix_time = time.time()
            if abs(server_unix_time - now_unix_time) > 15 * 60:
                raise OTSClientError('The difference between date in response and system time is more than 15 minutes.')

    def _check_authorization(self, query, headers, status=None):
        auth = headers.get('authorization')
        if auth is None:
            if status >= 200 and status < 300:
                raise OTSClientError('"Authorization" is missing in response header.')
            else:
                return

        # 1, check authorization
        if not auth.startswith('OTS '):
            raise OTSClientError('Invalid Authorization in response.')

        # 2, check accessid
        access_id, signature = auth[4:].split(':')
        if access_id != self.user_id:
            raise OTSClientError('Invalid accesskeyid in response.')

        # 3, check signature
        if signature != self._make_response_signature(query, headers):
            raise OTSClientError('Invalid signature in response.')

    def make_request(self, api_name, *args, **kwargs):
        
        if api_name not in self.api_list:
            raise OTSClientError('API %s is not supported.' % api_name)

        proto = self.encoder.encode_request(api_name, *args, **kwargs)
        body = proto.SerializeToString()
            
        query = '/' + api_name
        headers = self._make_headers(body, query)

        if self.logger.level <= logging.DEBUG:
            # prevent to generate formatted message which is time consuming 
            self.logger.debug("OTS request, API: %s, Headers: %s, Protobuf: %s" % (
                api_name, headers, 
                text_format.MessageToString(proto, as_utf8=True, as_one_line=True)
            ))

        return query, headers, body

    def _get_request_id_string(self, headers):
        request_id = headers.get('x-ots-requestid')
        if request_id is None:
            request_id = ""
        return request_id

    def parse_response(self, api_name, status, headers, body):
        if api_name not in self.api_list:
            raise OTSClientError("API %s is not supported." % api_name)

        headers = self._convert_urllib3_headers(headers)

        try:
            ret, proto = self.decoder.decode_response(api_name, body)
        except Exception, e:
            request_id = self._get_request_id_string(headers)
            error_message = 'Response format is invalid, %s, RequestID: %s, " \
                "HTTP status: %s, Body: %s.' % (str(e), request_id, status, body)
            self.logger.error(error_message)
            raise OTSClientError(error_message, status)

        if self.logger.level <= logging.DEBUG:
            # prevent to generate formatted message which is time consuming 
            request_id = self._get_request_id_string(headers)
            self.logger.debug("OTS response, API: %s, RequestID: %s, Protobuf: %s." % (
                api_name, request_id, 
                text_format.MessageToString(proto, as_utf8=True, as_one_line=True)
            ))
        return ret

    def handle_error(self, api_name, query, status, reason, headers, body):
        # convert headers according to different urllib3 versions.
        std_headers = self._convert_urllib3_headers(headers)

        if self.logger.level <= logging.DEBUG:
            # prevent to generate formatted message which is time consuming 
            self.logger.debug("OTS response, API: %s, Status: %s, Reason: %s, " \
                "Headers: %s" % (api_name, status, reason, std_headers))

        if api_name not in self.api_list:
            raise OTSClientError('API %s is not supported.' % api_name)


        try:
            self._check_headers(std_headers, body, status=status)
            if status != 403:
                self._check_authorization(query, std_headers, status=status)
        except OTSClientError, e:
            e.http_status = status
            e.message += " HTTP status: %s." % status
            raise e
        
        if status >= 200 and status < 300:
            return
        else:
            request_id = self._get_request_id_string(std_headers)

            try:
                error_proto = pb2.Error()
                error_proto.ParseFromString(body)
                error_code = error_proto.code
                error_message = error_proto.message
            except:
                error_message = "HTTP status: %s, reason: %s." % (status, reason)
                self.logger.error(error_message)
                raise OTSClientError(error_message, status)

            try:
                if status == 403 and error_proto.code != "OTSAuthFailed":
                    self._check_authorization(query, std_headers)
            except OTSClientError, e:
                e.http_status = status
                e.message += " HTTP status: %s." % status
                raise e

            self.logger.error("OTS request failed, API: %s, HTTPStatus: %s, " \
                "ErrorCode: %s, ErrorMessage: %s, RequestID: %s." % (
                api_name, status, error_proto.code, error_proto.message, request_id)
            )
            raise OTSServiceError(status, error_proto.code, error_proto.message, request_id)
