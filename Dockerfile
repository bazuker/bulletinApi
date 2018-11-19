FROM golang:latest

RUN mkdir -p /go/src/bulletinApi

WORKDIR /go/src/bulletinApi

COPY . /go/src/bulletinApi

RUN go install bulletinApi

CMD /go/bin/bulletinApi

EXPOSE 8080