FROM ubuntu:16.04

ENV \
 RUBY_INSTALL_VERSION="0.7.0" \
 RUBY_VERSION="2.3.1" \
 PATH="$PATH:/go/bin:/root/google-cloud-sdk/bin"

SHELL ["/bin/bash", "-c"]

RUN apt-get update \
 && apt-get install -y software-properties-common \
 && apt-get install -y \
       curl \
       libcurl4-gnutls-dev \
       libxml2-dev \
       python-pip \
       # python3-pip \
       zlib1g-dev \
 && curl -Lo /tmp/ruby-install-${RUBY_INSTALL_VERSION}.tar.gz https://github.com/postmodern/ruby-install/archive/v${RUBY_INSTALL_VERSION}.tar.gz \
 && tar zxvf /tmp/ruby-install-${RUBY_INSTALL_VERSION}.tar.gz --directory /root \
 && cd /root/ruby-install-${RUBY_INSTALL_VERSION} \
 && make install \
 && ruby-install -c -j4 --system ruby-${RUBY_VERSION} -- --disable-install-rdoc CC=cc \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/usr/local/go/bin" \
 GO111MODULE=on

RUN curl -Lo /tmp/go1.13.1.linux-amd64.tar.gz https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf /tmp/go1.13.1.linux-amd64.tar.gz \
 && rm /tmp/go1.13.1.linux-amd64.tar.gz

WORKDIR /app

COPY . .

CMD make build
