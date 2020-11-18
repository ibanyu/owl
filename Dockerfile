FROM hub.pri.ibanyu.com/devops/centos:7.5

LABEL MAINTAINER=devops@ipalfish.com

COPY dbinjectionopensource /data/home/serv/deploy/service/bin/

ENTRYPOINT ["/data/home/serv/deploy/service/bin/dbinjectionopensource"]
