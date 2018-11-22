FROM golang:alpine as builder
RUN mkdir -p /go/src/github.com/writeameer/aci
ADD . /go/src/github.com/writeameer/aci
WORKDIR /go/src/github.com/writeameer/aci
RUN apk add git
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /go/src/github.com/writeameer/aci/main /app/
COPY --from=builder /go/src/github.com/writeameer/aci/template /app/template
WORKDIR /app
CMD ["./main"]
