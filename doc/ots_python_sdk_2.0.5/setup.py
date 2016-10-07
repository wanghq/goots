try:
    import setuptools
except ImportError:
    print "Missing module setuptools."
    exit(1)

try:
    import distutils.core
except ImportError:
    print "Missing module distutils.core."
    exit(1)

try:
    import google.protobuf
except ImportError:
    print "Missing module google.protobuf."
    exit(1)

try:
    import urllib3
except ImportError:
    print "Missing module urllib3."
    exit(1)

from distutils.core import setup
setup(
    name='ots2_python_sdk',
    description='SDK of Open Table Service',
    author='Haowei YAO, Kunpeng HAN',
    author_email='haowei.yao@aliyun-inc.com, kunpeng.hkp@aliyun-inc.com',
    url='http://ots.aliyun.com',
    version='2.0.5',
    packages=['ots2', 'ots2.protobuf', 'ots2.example'],
    requires=['protobuf'],
)
