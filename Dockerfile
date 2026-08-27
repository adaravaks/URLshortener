FROM golang:1.27 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/api

FROM gcr.io/distroless/static-debian12

COPY --from=build /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]