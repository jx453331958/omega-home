FROM golang:1.25-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o omega-home .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=build /app/omega-home /usr/local/bin/
COPY --from=build /app/static /app/static
WORKDIR /app
RUN mkdir -p /app/data /app/static/uploads
EXPOSE 3000
CMD ["omega-home"]
