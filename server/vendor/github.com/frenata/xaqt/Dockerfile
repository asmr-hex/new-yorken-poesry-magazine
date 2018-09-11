############################################################
# Dockerfile to build sandbox for executing user code
# Based on Ubuntu
############################################################

FROM ubuntu:bionic
MAINTAINER Andrew Nichols

RUN apt-get update
RUN apt-get upgrade -y
#Install all the languages/compilers we are supporting.
RUN apt-get install -y gcc
RUN apt-get install -y g++
RUN apt-get install -y php-cli
RUN apt-get install -y ruby
RUN apt-get install -y python
RUN apt-get install -y mono-xsp4 mono-xsp4-base

RUN apt-get install -y mono-vbnc
RUN apt-get install -y npm
RUN apt-get install -y golang-go	
RUN apt-get install -y nodejs

#RUN npm cache clean -f
#RUN npm install -g n
#RUN n stable

RUN npm install -g underscore request express pug shelljs passport http sys jquery lodash async mocha moment connect validator restify ejs ws co when helmet fs-extra mustache should backbone forever  debug get-stdin

ENV NODE_PATH /usr/local/lib/node_modules/

RUN apt-get install -y clojure


#prepare for Java download
# commenting out according to: https://askubuntu.com/questions/422975/e-package-python-software-properties-has-no-installation-candidate
# RUN apt-get install -y python-software-properties
RUN apt-get install -y software-properties-common

#grab oracle java (auto accept licence)
RUN add-apt-repository -y ppa:webupd8team/java
RUN apt-get update
RUN echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections
RUN apt-get install -y oracle-java8-installer


RUN apt-get install -y gobjc
RUN apt-get install -y gnustep-devel &&  sed -i 's/#define BASE_NATIVE_OBJC_EXCEPTIONS     1/#define BASE_NATIVE_OBJC_EXCEPTIONS     0/g' /usr/include/GNUstep/GNUstepBase/GSConfig.h


#RUN apt-get install -y scala
RUN useradd -m mysql 
RUN apt-get install -y mysql-server
RUN apt-get install -y perl

RUN apt-get install -y curl
#RUN mkdir -p /opt/rust && \
#    curl https://sh.rustup.rs -sSf | HOME=/opt/rust sh -s -- --no-modify-path -y && \
#    chmod -R 777 /opt/rust

RUN apt-get install -y sudo
RUN apt-get install -y bc

RUN echo "mysql ALL = NOPASSWD: /usr/sbin/service mysql start" | cat >> /etc/sudoers

# install stack for haskell
RUN apt-get install -y haskell-platform

# copy entrypoint directory into image.
COPY ./entrypoint /entrypoint

RUN chmod 700 /entrypoint/usercode
