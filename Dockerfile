FROM google/debian:wheezy
MAINTAINER Luis Mora Medina <luismoramedina@gmail.com>

ADD goauth goauth
ADD goauth.yaml goauth.yaml
ENV PORT 14000
EXPOSE 14000
ENTRYPOINT ["/goauth"]