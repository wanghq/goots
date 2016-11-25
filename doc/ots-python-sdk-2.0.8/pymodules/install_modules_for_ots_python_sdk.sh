#!/bin/sh -xe

sudo rm -rf distribute-0.7.3 setuptools-1.3.2 protobuf-2.5.0 urllib3-1.14 certifi-2016.2.28

unzip distribute-0.7.3.zip
cd distribute-0.7.3
sudo python2.7 setup.py install
cd ..

tar xvf setuptools-1.3.2.tar.gz
cd setuptools-1.3.2
sudo python2.7 setup.py install
cd ..

tar xvf protobuf-2.5.0.tar.gz
cd protobuf-2.5.0
sudo python2.7 setup.py install
cd ..

tar xvf urllib3-1.14.tar.gz
cd urllib3-1.14
sudo python2.7 setup.py install
cd ..

tar xvf certifi-2016.2.28.tar.gz
cd certifi-2016.2.28
sudo python2.7 setup.py install
cd ..

sudo rm -rf distribute-0.7.3 setuptools-1.3.2 protobuf-2.5.0 urllib3-1.14 certifi-2016.2.28
