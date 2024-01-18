FROM golang:1.21 AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /app/docker-app

FROM alpine
LABEL maintainer="Robert Jabłecki <robert.jablecki@gmail.com>"

RUN mkdir /app

# Copy any other files required in the final image here
COPY --from=build /app/docker-app /app/docker-app

EXPOSE 8081

CMD ["/app/docker-app"]
