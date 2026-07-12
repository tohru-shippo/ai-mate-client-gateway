ARG GO_IMAGE=golang:1.26.3-bookworm
ARG GATEWAY_DIR=.

FROM ${GO_IMAGE} AS build
WORKDIR /src
ARG GATEWAY_DIR

COPY ${GATEWAY_DIR}/go.mod ${GATEWAY_DIR}/go.sum ./
COPY ai-mate-server /ai-mate-server
RUN go mod edit -replace github.com/tohru-shippo/ai-mate-server=/ai-mate-server && go mod download

COPY ${GATEWAY_DIR}/ .
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
