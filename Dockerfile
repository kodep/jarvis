FROM golang:1.20.4-alpine as builder

WORKDIR /build

RUN apk add --no-cache upx

COPY go.mod go.sum ./
RUN go mod download -x

COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN go build -o bin/ -ldflags "-s -w -X main.DefaultMode=production" ./cmd/jarvis && \
  upx bin/*


FROM cgr.dev/chainguard/wolfi-base

COPY --from=builder /build/bin /app

ENV PATH="/app:${PATH}"

CMD ["jarvis"]
