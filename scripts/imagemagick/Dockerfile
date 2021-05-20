FROM alpine:3.13

RUN apk add --update imagemagick libheif \
	&& rm -rf /var/cache/apk/*

COPY script /script

RUN chmod +x /script

CMD ["/bin/sh"]