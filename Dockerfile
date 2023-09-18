FROM golang:1.19-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

#dependencies
COPY ["app/go.mod","app/go.sum","./"]
RUN go mod download

#build
COPY app ./
RUN go build -o ./bin/app cmd/main/app.go


FROM alpine as runner

COPY --from=builder /usr/local/src/bin/app /

COPY configs/config.yml /configs/config.yml

COPY migrations /migrations

CMD ["/app"]
