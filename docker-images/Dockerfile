FROM alpine:3.5

# update apk repository mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/' /etc/apk/repositories

RUN apk update && apk upgrade && \
  apk add ruby-json ruby-irb && \
  apk add build-base ruby-dev && \
  apk add --no-cache python && \
  apk add ca-certificates wget && \
  apk add tzdata && \
  cat /usr/share/zoneinfo/Asia/Shanghai > /etc/localtime && \
#  gem sources --add https://gems.ruby-china.org/ --remove https://rubygems.org/ && \
#  gem sources -l  && \
  gem install fluentd -v "~> 0.12.0" --no-ri --no-rdoc && \
  gem install fluent-plugin-elasticsearch --no-ri --no-rdoc && \
  gem install gelf -v "~> 3.0.0" --no-ri --no-rdoc && \
  gem install aliyun_sls_sdk -v ">=0.0.9" --no-ri --no-rdoc && \
  apk del build-base ruby-dev && \
  rm -rf /root/.gem

COPY plugins/ /etc/fluentd/plugins/

VOLUME /etc/fluentd/conf.d

COPY pilot fluentd.tpl entrypoint config.default healthcheck.sh /pilot/
VOLUME /pilot/pos

EXPOSE 24224
WORKDIR /pilot/
#CMD exec fluentd -c /fluentd/etc/$FLUENTD_CONF -p /fluentd/plugins $FLUENTD_OPT
CMD /pilot/entrypoint

#HEALTHCHECK --interval=1s --timeout=1s --retries=1 \
#	CMD ret=`pidof fluentd |wc -w`;if [ $ret -ne 2 ] ;then exit 1; fi