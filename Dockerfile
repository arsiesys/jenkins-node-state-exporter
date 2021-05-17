############################
# STEP 1 build executable binary
############################
FROM golang as builder

WORKDIR $GOPATH/srv/aresiesys/jenkins-node-state-exporter/
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/jenkins-node-state-exporter -mod vendor main.go

############################
# STEP 2 build a small image
############################
FROM scratch

ENV GIN_MODE=release
WORKDIR /app/
# Import from builder.
COPY --from=builder /go/bin/jenkins-node-state-exporter /app/jenkins-node-state-exporter
ENTRYPOINT ["/app/jenkins-node-state-exporter"]
EXPOSE 9723
