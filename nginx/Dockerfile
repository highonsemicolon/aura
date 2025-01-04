FROM nginx:1.18

RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    libpcre3-dev \
    zlib1g-dev \
    openssl libssl-dev \
    git \
    libjwt-dev \
    libjansson-dev

RUN git clone https://github.com/TeslaGov/ngx-http-auth-jwt-module.git /tmp/ngx-http-auth-jwt-module

RUN wget http://nginx.org/download/nginx-$(nginx -v 2>&1 | grep -oP '\d+\.\d+\.\d+').tar.gz -O /tmp/nginx.tar.gz && \
    tar -xvzf /tmp/nginx.tar.gz -C /tmp && \
    cd /tmp/nginx-* && \
    ./configure --with-compat --add-dynamic-module=/tmp/ngx-http-auth-jwt-module && \
    make modules && \
    cp objs/*.so /usr/lib/nginx/modules

COPY nginx.conf /etc/nginx/nginx.conf

RUN apt-get purge -y build-essential git wget && apt-get autoremove -y && apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/*
