# syntax = tonistiigi/dockerfile:runmount20180618
FROM golang:1.11-alpine AS compiler
RUN apk add --update gcc musl-dev linux-headers ca-certificates
ARG version
RUN --mount=src=/,dst=/go/src,ro=true \
 CGO_ENABLED=0 go build -o /go/bin/redirector -i --tags="debug verbose" -ldflags "-extldflags \"-static\"" daemon

FROM scratch
COPY --from=compiler /go/bin/redirector /var/lib/cycle/redirector
EXPOSE 80 443
CMD ["/var/lib/cycle/redirector"]