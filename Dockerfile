FROM golang:alpine as builder
MAINTAINER Genie.C <ygenie.chae@gmail.com>
ENV TZ=Asia/Seoul
RUN apk update && apk add --no-cache git
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -ldflags '-w -s' -o main .

FROM golang:alpine
COPY --from=builder /build/main /app/

WORKDIR /app
EXPOSE 9801
ENV PATH /app
ENTRYPOINT ["/app/main"]
