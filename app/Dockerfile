FROM golang:1.14.3-alpine3.11

ENV APP_NAME foodsearching
ENV PORT 5555

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN apk add git
RUN go get github.com/gin-gonic/gin

RUN go get ./

EXPOSE ${PORT}

RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}