FROM golang:1.14

WORKDIR /go/src/
COPY . .

RUN go get github.com/tarantool/go-tarantool \
           github.com/spf13/afero            \
           github.com/yandex/pandora

RUN GOOS=linux GOARCH=amd64 go build tnt_queue_gun.go