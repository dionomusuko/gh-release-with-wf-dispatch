FROM golang:1.20-buster AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY *.go ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -a -o /bin/app

FROM gcr.io/distroless/static
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
