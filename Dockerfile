FROM golang:latest
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's/security.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
RUN apt update && apt install python-pip -y
RUN pip2 install fabric==1.14 -i https://pypi.douban.com/simple

WORKDIR /code
ADD . /code
