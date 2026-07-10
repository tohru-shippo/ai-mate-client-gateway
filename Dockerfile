ARG GO_IMAGE=golang:1.26.3-bookworm

FROM ${GO_IMAGE} AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/server
RUN mkdir -p /out/app && cp -r configs /out/app/configs && if [ -d docs ]; then cp -r docs /out/app/docs; fi

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/gateway /app/gateway
COPY --from=build /out/app/ /app/

CMD ["/app/gateway"]