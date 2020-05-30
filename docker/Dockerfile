FROM alpine:3.8
MAINTAINER Jeff Plourde <jeff@jeffplourde.com>

ARG build_tag

ADD soba /soba

RUN apk add --update --no-cache ca-certificates bash curl git \
    && rm -f "/var/cache/apk/*"
#    && curl -L https://dl.bintray.com/jeffplourde/soba/soba/$build_tag/soba_linux_amd64 -o /soba \
#    && chmod 755 /soba
ENTRYPOINT ["/bin/bash", "-c", "/soba \"$@\"", "--"]
