ARG NGINX_VERSION=1.25.5

FROM nginx:${NGINX_VERSION}-alpine

COPY docker/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ca/taanni.online.crt /etc/ssl/certs/taanni.online.crt
COPY ca/taanni.online.key /etc/ssl/private/taanni.online.key

WORKDIR /var/www/public