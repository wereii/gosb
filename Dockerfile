FROM golang:alpine as build

RUN mkdir /opt/build
COPY gosb.go go.mod go.sum /opt/build/
WORKDIR /opt/build
RUN go mod download
RUN go build ./gosb.go

FROM alpine as app

COPY --from=build /opt/build/gosb /usr/bin/

CMD ["/usr/bin/gosb"]