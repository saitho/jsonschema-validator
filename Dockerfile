FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o ./bin/validator .

FROM alpine
COPY --from=builder /build/bin/validator /bin/
WORKDIR /data
CMD ["/bin/validator"]
