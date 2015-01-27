FROM busybox:latest

MAINTAINER tim@magnetic.io

ADD ./target/linux_i386/monarch /monarch

ENTRYPOINT ["/monarch"]
