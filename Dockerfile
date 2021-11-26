FROM golang:1.17 as build
WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD . /src
RUN CGO_ENABLED=0 go build -o app -ldflags "-extldflags=-static" .

FROM scratch
COPY --from=build /src/app /app
ENTRYPOINT [ "/app" ]