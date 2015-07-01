./proxy --master=http://192.168.99.98:8983 \
  --slaves=http://192.168.99.98:8983,http://192.168.99.98:8983,http://192.168.99.98:8983 \
  --log-location=proxy.log \
  --bucket-name=gogobot-solr-docs \
  --aws-region=us-west-2 \
  --aws-endpoint=https://s3-us-west-2.amazonaws.com \
  --bucket-prefix=kensodev-dev \
  --listen-port=8982