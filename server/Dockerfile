FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./picourl-backend .

FROM scratch

WORKDIR /

COPY --from=build /app/picourl-backend /picourl-backend

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./picourl-backend"]