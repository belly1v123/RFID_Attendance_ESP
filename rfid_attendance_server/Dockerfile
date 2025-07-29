# # Build stage
# FROM golang:1.22 AS build
# WORKDIR /app
# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . ./
# RUN go build -o main .

# # Run stage
# FROM alpine:latest  
# WORKDIR /root/
# COPY --from=build /app/main .
# EXPOSE 8000
# CMD ["./main"]
