FROM golang:alpine as builder
MAINTAINER Genie.C <ygenie.chae@gmail.com>
ENV TZ=Asia/Seoul
RUN apk update && apk add --no-cache git gcc libc-dev make alsa-lib-dev
RUN mkdir /build
ADD . /build/
RUN ls -al /build/
WORKDIR /build
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/faiface/beep/speaker
RUN go get -u github.com/faiface/beep/wav
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -ldflags '-w -s' -o main .

FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/l2ping /app/
COPY --from=builder /build/static/* /app/static/
ADD lib /lib
ADD /usr/bin/aplay /usr/bin/aplay
ADD /bin/bash /bin/bash
WORKDIR /app
EXPOSE 9801
ENV PATH="/app:$PATH"
ENTRYPOINT ["/app/main"]
