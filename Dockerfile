FROM cgr.dev/chainguard/go:1.19 as build
WORKDIR /work
COPY go.mod go.sum ./
COPY app ./app
COPY cmd ./cmd
COPY docs ./docs
COPY pkg ./pkg
COPY vendor ./vendor
RUN CGO_ENABLED=0 go build -o moonbase ./cmd/moonbase

FROM cgr.dev/chainguard/alpine-base
RUN apk --no-cache add ca-certificates
COPY --from=build /work/moonbase /usr/bin
CMD ["moonbase"]
