FROM golang:1.16 as build
WORKDIR /gateway
COPY go.mod .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o gateway .

FROM alpine:3.14
WORKDIR /
COPY --from=build /gateway/gateway ./
ENTRYPOINT ["./gateway"]