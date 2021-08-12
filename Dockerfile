FROM hub.pri.ibanyu.com/devops/centos:7.5

LABEL MAINTAINER=infrastructure@ipalfish.com

COPY ./bin/dbinjection /service/bin/
COPY ./config /service/config
COPY ./static /service/static

WORKDIR /service

ENTRYPOINT ["bin/owl"]
