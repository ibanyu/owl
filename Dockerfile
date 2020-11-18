FROM hub.pri.ibanyu.com/devops/centos:7.5

LABEL MAINTAINER=devops@ipalfish.com

COPY dbinjection /data/home/serv/deploy/service/bin/

ENTRYPOINT ["/data/home/serv/deploy/service/bin/dbinjection"]