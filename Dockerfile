# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -mod=readonly -o /code/build/wasmd ./cmd/wasmd

# Final image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates jq bash libc6-compat

COPY --from=builder /code/build/wasmd /usr/bin/wasmd

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# grpc
EXPOSE 9090

ENTRYPOINT ["/usr/bin/wasmd"]
CMD ["version"]
