FROM golang:1.21 AS builder

WORKDIR /
# Uncomment next line when adding modules
# RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./main.go

FROM scratch
WORKDIR /bin/
COPY --from=builder /bin/server ./
CMD ["./server"]