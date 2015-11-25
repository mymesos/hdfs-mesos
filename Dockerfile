FROM mymesos/java7-mesos:0.24.1
MAINTAINER Dave Choi <goodoi09@gmail.com>

COPY . /usr/src/hdfs
WORKDIR /usr/src/hdfs

ENV HDFS_MESOS_NAME hdfs-mesos-0.1.5
RUN tar zxvf $HDFS_MESOS_NAME.tgz && rm -rf $HDFS_MESOS_NAME.tgz

CMD ./run.sh