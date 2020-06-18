FROM golang:1.14.4-alpine3.12

RUN apk add --no-cache make

ENV GO111MODULE=on
WORKDIR /src
COPY . .

RUN make all

FROM alpine:3.12

RUN apk update && apk upgrade && \
    apk add --no-cache bash curl

COPY --from=0 /src/bin/* /bin/

ENTRYPOINT ["/bin/matrix"]
