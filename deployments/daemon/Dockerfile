# syntax=docker/dockerfile:1.7

FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/faultplane-daemon ./cmd/daemon

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/faultplane-daemon /faultplane-daemon

EXPOSE 8080

ENV FAULTPLANE_GATEWAY_PORT=8080

ENTRYPOINT ["/faultplane-daemon"]
