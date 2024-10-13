FROM ubuntu:24.04

ADD . /var/www
WORKDIR /var/www

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV LC_CTYPE C.UTF-8

ARG DEBIAN_FRONTEND=noninteractive

RUN apt -y update

RUN apt -y install vim
RUN apt -y install screen
RUN apt -y install bash
RUN apt -y install git
RUN apt -y install htop
RUN apt -y install curl
RUN apt -y install wget
RUN apt -y install file
RUN apt -y install iputils-ping
RUN apt -y install psmisc
RUN apt -y install dnsutils
RUN apt -y install software-properties-common
RUN apt -y install sudo
RUN apt -y install telnet
RUN apt -y install libpq5
RUN apt -y install postgresql-client

# install latest Golang
RUN rm -rf update-golang
RUN git clone https://github.com/udhos/update-golang.git update-golang
WORKDIR update-golang
RUN sudo ./update-golang.sh
WORKDIR ..
RUN rm -rf update-golang

ENV GOPATH /root/go/
ENV PATH="${PATH}:/usr/local/go/bin/"

WORKDIR /var/www
CMD ["go", "run", "."]
