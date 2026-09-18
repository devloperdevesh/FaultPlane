# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/faultplane-daemon ./cmd/daemon
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/healthcheck ./deployments/daemon/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/faultplane-daemon /faultplane-daemon
COPY --from=build /out/healthcheck /usr/local/bin/healthcheck
COPY --from=build /src/internal/kernel/bpf/sockmap.bpf.o /usr/local/lib/faultplane/bpf/sockmap.bpf.o

EXPOSE 8080

ENV FAULTPLANE_GATEWAY_PORT=8080
ENV FAULTPLANE_BPF_OBJECT=/usr/local/lib/faultplane/bpf/sockmap.bpf.o

ENTRYPOINT ["/faultplane-daemon"]
