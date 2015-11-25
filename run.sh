#!/bin/bash

./create_config -output_dir ${HDFS_MESOS_NAME}/etc/hadoop

pushd $HDFS_MESOS_NAME
./bin/hdfs-mesos