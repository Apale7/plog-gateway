FROM golang:latest as builder

ENV TZ Asia/Shanghai

## 添加镜像
RUN cp /etc/apt/source.list /etc/apt/sources.list.back \
    && sed -i '-s#http://deb.debian.org#https://mirrors.aliyun.com#g' /etc/apt/sources.list

##添加mongodb扩展

##添加redis扩展

##添加kafka扩展