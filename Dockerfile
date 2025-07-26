ARG GO_VERSION=1.15.6
 
FROM golang:${GO_VERSION}-alpine AS build
 
RUN apk add --no-cache git
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app
 
FROM scratch AS final
 
COPY --from=build . /app
 
EXPOSE 3000

# run binary; use vector form
ENTRYPOINT ["/app"]