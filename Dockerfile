FROM registry-quay-quay.apps.pro-scp1.sf-rz.de/scp/image-tool:xtagx
LABEL maintainer="Dirk Jäger <dirk.jaeger@schufa.de>"

ENV LANG=en_US.UTF-8
ENV TZ=Europe/Berlin
ENV https_proxy=http://webproxy.sf-bk.de:8181/
ENV HTTPS_PROXY=$https_proxy

RUN rm -f /usr/bin/image-tool

COPY image-tool /usr/bin/image-tool

RUN chgrp root /usr/bin/image-tool

CMD ["/usr/bin/image-tool", "-socks5=no", "-server", "-statcfg"]
