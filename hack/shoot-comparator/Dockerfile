FROM golang:1.23.0-alpine3.20 as build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workdir

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

ARG BIN
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o shoot-comparator cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /workdir/shoot-comparator .
USER 65532:65532

ENTRYPOINT ["/shoot-comparator"]
