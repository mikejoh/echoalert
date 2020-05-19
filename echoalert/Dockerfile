# Build
FROM golang:1.14 AS build
ADD main.go /app/
ADD go.mod /app/
WORKDIR /app
RUN go build -o echoalert .

# Final
FROM gcr.io/distroless/base-debian10
COPY --from=build /app/echoalert /bin/
ENTRYPOINT ["echoalert"]