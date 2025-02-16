FROM cgr.dev/chainguard/go as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /liftplan-server serve/main.go

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder /liftplan-server /usr/local/bin/
EXPOSE 9000
CMD ["liftplan-server"]
