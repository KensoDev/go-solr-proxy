FROM sameersbn/ubuntu:14.04.20150613

RUN \
    apt-get update && \
    apt-get install -y wget && \
    cd /tmp && \
    wget https://github.com/KensoDev/go-solr-proxy/releases/download/v0.4.0/proxy_linux_amd64 && \
    mv proxy_linux_amd64 /usr/local/bin/solrproxy

ADD run.sh /run.sh

RUN chmod 755 /usr/local/bin/solrproxy
RUN chmod 755 /run.sh

EXPOSE 8982

CMD /run.sh