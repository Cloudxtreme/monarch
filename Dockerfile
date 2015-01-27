FROM busybox:latest

MAINTAINER tim@magnetic.io

ADD ./target/linux_i386/mesos-tester /mesos-tester

ENTRYPOINT ["/mesos-tester"]
