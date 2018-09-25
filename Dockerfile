FROM golang:alpine
RUN apk update && apk add nano && apk add git

RUN rm -rf /var/lib/apt/lists/*

RUN mkdir -p /app/profiler

WORKDIR /app/

RUN go get "github.com/globalsign/mgo"
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/gorilla/handlers"
RUN go get "github.com/montanaflynn/stats"

COPY $PWD/ /app/

RUN GOOS=linux GOARCH=amd64 go build -o runner

CMD ["/bin/true"]
