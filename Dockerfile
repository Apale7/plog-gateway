FROM php:7.1-fpm

ENV TZ Asia/Shanghai

##添加镜像
RUN cp /etc/apt/sources.list /etc/apt/sources.list.back \
    && sed -i 's#http://deb.debian.org#https://mirrors.aliyun.com#g' /etc/apt/sources.list

##安装gd扩展
RUN apt-get update && apt-get install -y \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libpng-dev \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) gd
RUN docker-php-ext-install bcmath
#安装mysql,redis扩展
RUN docker-php-ext-install pdo_mysql \
    && cp /usr/local/etc/php/php.ini-development /usr/local/etc/php/php.ini \
    && curl -L -o /tmp/redis.tar.gz  http://pecl.php.net/get/redis-5.3.1.tgz \
    && tar zxvf /tmp/redis.tar.gz \
    && rm -r /tmp/redis.tar.gz \
    && mkdir -p /usr/src/php/ext \
    && mv ./redis-5.3.1 /usr/src/php/ext/redis \
    && docker-php-ext-install redis \
    #删除包缓存中的所有包
    && apt-get clean \
    && apt-get autoclean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*