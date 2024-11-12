FROM goreleaser/goreleaser:v2.3.2 AS builder
WORKDIR /build
COPY . .
RUN goreleaser build --clean --single-target

FROM alpine:3.20.3
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/dist/cir-rotator_linux_amd64_v1/cir-rotator /usr/bin/cir-rotator
ENTRYPOINT ["gitlab-token-updater"]
CMD ["--help"]