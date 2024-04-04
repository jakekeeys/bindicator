FROM golang as build
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -v -o app ./cmd/bindicator

FROM chromedp/headless-shell:latest
RUN apt-get update && apt-get install -y dumb-init ca-certificates

WORKDIR /service
COPY --from=build /build/app .
ENTRYPOINT ["dumb-init", "--"]
CMD ./app