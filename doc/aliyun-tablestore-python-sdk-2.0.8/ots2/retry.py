# -*- coding: utf8 -*-
import random
import math

class RetryPolicy:
    """
    ```RetryPolicy```是重试策略的接口，包含2个未实现的方法和它们的参数列表。要实现一个重试策略，
    继承这个类并实现它的2个方法。
    """

    def should_retry(self, retry_times, exception, api_name):
        raise NotImplementedError()

    def get_retry_delay(self, retry_times, exception, api_name):
        raise NotImplementedError()


class RetryUtil:

    @classmethod
    def should_retry_no_matter_which_api(cls, exception):
        error_code = exception.code
        error_message = exception.message

        if (error_code == "OTSRowOperationConflict" or 
            error_code == "OTSNotEnoughCapacityUnit" or
            error_code == "OTSTableNotReady" or
            error_code == "OTSPartitionUnavailable" or
            error_code == "OTSServerBusy" or
            error_code == "OTSOperationThrottled"):
            return True

        if error_code == "OTSQuotaExhausted" and error_message == "Too frequent table operations.":
            return True

        return False

    @classmethod
    def is_repeatable_api(cls, api_name):
        if (api_name == "/ListTable" or
            api_name == "/DescribeTable" or
            api_name == "/GetRow" or
            api_name == "/BatchGetRow" or
            api_name == "/GetRange" or
            api_name == "/DescrieStream" or
            api_name == "/GetShardIterator" or
            api_name == "/GetStreamRecord" or
            api_name == "/ListStream"):
            return True

        return False

    @classmethod
    def should_retry_when_api_repeatable(cls, retry_times, exception, api_name):
        error_code = exception.code
        error_message = exception.message
        http_status = exception.http_status

        if (error_code == "OTSTimeout" or 
            error_code == "OTSInternalServerError" or 
            error_code == "OTSServerUnavailable"):
            return True

        if (http_status == 500 or http_status == 502 or http_status == 503):
            return True

        # TODO handle network error & timeout
        return False

    @classmethod
    def is_server_throttling_exception(cls, exception):
        error_code = exception.code
        error_message = exception.message

        if (error_code == "OTSServerBusy" or 
            error_code == "OTSNotEnoughCapacityUnit" or
            error_code == "OTSOperationThrottled"): 
            return True

        if error_code == "OTSQuotaExhausted" and error_message == "Too frequent table operations.":
            return True

        return False


class DefaultRetryPolicy(RetryPolicy):
    """
    默认重试策略
    最大重试次数为3，最大重试间隔为2000毫秒，对流控类错误以及读操作相关的服务端内部错误进行了重试。
    """

    # 最大重试次数
    max_retry_times = 3

    # 最大重试间隔，单位为秒
    max_retry_delay = 2   

    # 每次重试间隔的递增倍数
    scale_factor = 2

    # 两种错误的起始重试间隔，单位为秒
    server_throttling_exception_delay_factor = 0.5
    stability_exception_delay_factor = 0.2

    def _max_retry_time_reached(self, retry_times, exception, api_name):
        return retry_times >= self.max_retry_times

    def _can_retry(self, retry_times, exception, api_name):

        if RetryUtil.should_retry_no_matter_which_api(exception):
            return True

        if RetryUtil.is_repeatable_api(api_name) and RetryUtil.should_retry_when_api_repeatable(retry_times, exception, api_name):
            return True

        return False

    def get_retry_delay(self, retry_times, exception, api_name):

        if RetryUtil.is_server_throttling_exception(exception):
            delay_factor = self.server_throttling_exception_delay_factor
        else:
            delay_factor = self.stability_exception_delay_factor

        delay_limit = delay_factor * math.pow(self.scale_factor, retry_times)

        if delay_limit >= self.max_retry_delay:
            delay_limit = self.max_retry_delay

        real_delay = delay_limit * 0.5 + delay_limit * 0.5 * random.random()
        return real_delay

    def should_retry(self, retry_times, exception, api_name):
        
        if self._max_retry_time_reached(retry_times, exception, api_name):
            return False

        if self._can_retry(retry_times, exception, api_name):
            return True

        return False


class NoRetryPolicy(RetryPolicy):
    """
    不进行任何重试的重试策略
    """

    def get_retry_delay(self, retry_times, exception, api_name):
        return 0

    def should_retry(self, retry_times, exception, api_name):
        return False


class NoDelayRetryPolicy(DefaultRetryPolicy):
    """
    没有延时的重试策略
    """

    def get_retry_delay(self, retry_times, exception, api_name):
        return 0

