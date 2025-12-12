FROM ubuntu:24.04

WORKDIR /opt/www/sayhi/backend/

# 设定时区
#ENV TZ=Asia/Shanghai
ENV TZ=Asia/Yangon
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y gcc
RUN apt-get install -y tzdata
RUN apt-get install -y openssl 
RUN apt-get install -y curl

CMD ["/opt/www/sayhi/backend/sayhi"]
