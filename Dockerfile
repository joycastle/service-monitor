FROM golang:1.18 AS base
MAINTAINER Levin

RUN apt update && apt upgrade -y && \
    apt install -y git make openssh-client

ADD ./ /app/
ADD /home/ec2-user/cert /app/cert

WORKDIR /app
RUN cd /app && go build

ENTRYPOINT ./service-monitor