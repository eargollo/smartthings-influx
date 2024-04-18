FROM golang:1.22.2 as build
RUN apt-get update; apt-get install ca-certificates -y
WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD . /src
RUN CGO_ENABLED=0 go build -o app -ldflags "-extldflags=-static" .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /src/app /app
ENTRYPOINT [ "/app" ]