FROM golang:1.22-alpine as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN mkdir build
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/ipg ./cmd/server
RUN echo "#!/bin/sh" > /build/reload && echo "kill -HUP 1" >> /build/reload && chmod a+x /build/reload


FROM alpine:edge as runner 
COPY --from=build /build /bin
ENV TZ=Asia/Tehran
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
EXPOSE 8080

# Run
CMD ["/bin/ipg"]

