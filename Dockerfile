FROM alpine:3.13

LABEL description="Image for 32 degrees server"

RUN addgroup -S app && adduser -S -G app app
RUN apk --update upgrade && \
    apk add ca-certificates && \
    update-ca-certificates && \
    apk add su-exec && \
    apk add curl && \
    rm -rf /var/cache/apk/*

RUN apk add --no-cache tzdata

RUN mkdir /home/app/bin
COPY ./main /home/app/bin/
COPY ./entrypoint.sh /home/app/bin/entrypoint.sh

WORKDIR /home/app/bin

CMD ["./entrypoint.sh"]
