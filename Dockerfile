# *************************************
#
# OpenGM
#
# *************************************

FROM alpine:3.14

MAINTAINER XTech Cloud "xtech.cloud"

ENV container docker
ENV MSA_MODE release

EXPOSE 18810

ADD bin/ogm-actor /usr/local/bin/
RUN chmod +x /usr/local/bin/ogm-actor

CMD ["/usr/local/bin/ogm-actor"]
