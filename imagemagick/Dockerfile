FROM alpine:3.13

RUN apk add --update imagemagick libheif \
	&& rm -rf /var/cache/apk/*

CMD ["/bin/sh"]