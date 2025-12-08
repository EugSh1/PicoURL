FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o ./picourl-backend .

FROM alpine

WORKDIR /

COPY --from=build /app/picourl-backend /picourl-backend
COPY --from=build /app/db/migrations /db/migrations

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./picourl-backend"]