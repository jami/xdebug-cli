FROM ubuntu:17.10

RUN apt-get -y update && apt-get -y upgrade
RUN apt-get -y install php php-xdebug

COPY ./bin/xdbg_linux_amd64 /usr/bin/xdbg
COPY ./_example/test.php /opt/example/test.php

ENTRYPOINT ["/usr/bin/xdbg", "--host", "127.0.0.1", "--port", "9000", "listen"]