FROM golang:1.18.3-alpine3.15 as builder
WORKDIR /build
COPY . .
RUN apk --update add make && make compile 

FROM alpine:3.15.4 
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/dist/cir-rotator /usr/bin/cir-rotator
ENTRYPOINT ["cir-rotator"]
CMD ["--help"]