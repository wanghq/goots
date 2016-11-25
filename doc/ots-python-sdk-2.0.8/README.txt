Copyright 2015 Aliyun Inc.

Dependency
==========

ots_python_sdk依赖及版本要求：

python: 2.7
setuptools: 1.3.2
distribute: 0.7.3
protobuf: 2.5.0
urllib3: 1.14
certifi: 2016.2.28

上面依赖包可以从https://pypi.python.org/pypi网站下载，或者执行pymodules/install_modules_for_ots_python_sdk.sh安装。

Installation
============

执行以下命令安装OTS SDK：

sudo python2.7 setup.py install

安装成功后即可使用OTS SDK，使用方法请参考DOCUMENT和ots2/example。

Performance
===========

由于默认安装python-protobuf的Python实现性能比较差，如果要获得更好的SDK端性能，建议安装python-protobuf的C++实现，安装方法如下。

1. 删除python-protobuf的Python实现。

    sudo rm -rf $PYTHON_PACKAGE_PATH/protobuf-*

根据系统环境的不同，其中$PYTHON_PACKAGE_PATH可能为：
    /usr/local/lib/python2.7/site-packages
    /usr/lib/python2.7/site-packages
    /usr/local/lib/python2.7/dist-packages
    /usr/lib/python2.7/dist-packages

如果是64位机器，可能需要将lib换为lib64。

2. 下载Protobuf源码包，地址https://protobuf.googlecode.com/files/protobuf-2.5.0.tar.gz。解压并编译安装。

    tar zxf protobuf-2.5.0.tar.gz
    cd protobuf-2.5.0
    ./configure
    make
    sudo make install

3. 编译安装python-protobuf的C++实现（在protobuf-2.5.0/python目录里）。

    cd python
    sudo sh -c "export PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=cpp; python2.7 setup.py install"

4. 使用ots_python_sdk时需要设置下面的环境变量。

    export PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=cpp

Debug
=====

使用ots_python_sdk时如果需要调试，可以打开debug日志，方法为在创建OTSClient之前加入下面的代码。

    logger = logging.getLogger(OTSClient.DEFAULT_LOGGER_NAME)
    logger.setLevel(logging.DEBUG)
    handler = logging.FileHandler('test.log')
    logger.addHandler(handler)

其中test.log为输出日志的文件名，可自己定义。

