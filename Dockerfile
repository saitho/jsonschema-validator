FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o validator .

FROM alpine
COPY --from=builder /build/validator /bin/
WORKDIR /data
CMD ["/bin/validator"]
