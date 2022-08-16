#!/bin/bash

sudo yum install -y centos-release-SCL
yum repolist


sudo yum install -y rh-git227.x86_64 devtoolset-10.x86_64

echo 'you can use it like: yum --disablerepo="*"  --enablerepo="centos-sclo-rh" search git'
echo 'scl enable devtoolset-10 bash'
echo 'scl enable rh-git227 bash'
scl -l
