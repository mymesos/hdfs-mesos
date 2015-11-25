FROM mymesos/java7-mesos:0.25.0
MAINTAINER Dave Choi <goodoi09@gmail.com>

COPY . /usr/src/hdfs
WORKDIR /usr/src/hdfs

RUN tar zxvf hdfs-mesos-0.1.6.tgz && rm -rf hdfs-mesos-0.1.6.tgz
WORKDIR hdfs-mesos-0.1.6

CMD ./bin/hdfs-mesos