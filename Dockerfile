FROM golang:1.10-alpine


RUN apk --update add git


ENV GOBIN /go/bin

ENV mongoDB_address 127.0.0.1:27017

ENV mongoDB_database article_xtend

ENV mongoDB_collection article

ENV XS_ADMIN_USERID doc@xenonstack.com

ENV XS_ADMIN_PASSWORD xenonstack#admin


ADD . /go/src/l2-golang-auth
WORKDIR /go/src/l2-golang-auth

RUN go get -u github.com/golang/dep/...
RUN dep init
RUN dep ensure

RUN go install l2-golang-auth

EXPOSE 8080

ENTRYPOINT "/go/bin/Xenonstack-documentation-portal"
