FROM golang:1.12-alpine
WORKDIR /lambda-control-plane
RUN apk add --no-cache git
ADD . .
RUN go build -o landad cmd/landa/main.go

FROM alpine:3.9
COPY --from=0 /lambda-control-plane/landad /bin/landad
CMD ["landad"]
