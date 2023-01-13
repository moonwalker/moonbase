FROM golang:alpine AS build
WORKDIR /work
COPY .git/ ./.git
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal
COPY pkg ./pkg
COPY vendor ./vendor
RUN apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go build -o moonbase ./cmd/moonbase

FROM cgr.dev/chainguard/alpine-base
RUN apk --no-cache add ca-certificates
COPY --from=build /work/moonbase /usr/bin
ENTRYPOINT ["moonbase"]
