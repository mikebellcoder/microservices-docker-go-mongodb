# base image
FROM golang:1.22.5
RUN adduser -D -g '' elf
WORKDIR /opt/app/
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-users ./cmd/app

# final image (smaller)
FROM alpine:3.17.3
LABEL language="golang"
LABEL org.opencontainer.image.source="https://github.com/mikebellcoder/microservices-docker-go-mongodb"
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder --chown=elf:1000 /go/bin/cinema-users /cinema-users
USER elf
ENTRYPOINT [ "./cinema-users" ]