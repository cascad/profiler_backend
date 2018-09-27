FROM golang:alpine as builder

RUN apk add --no-cache git gcc libc-dev

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY $PWD/ /src/profiler/
WORKDIR /src/profiler/

RUN go get "github.com/globalsign/mgo"
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/gorilla/handlers"
RUN go get "github.com/montanaflynn/stats"

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/runner

FROM scratch

COPY --from=builder /src/profiler/profiler_frontend/build /app/static
COPY --from=builder /go/bin/runner /app/runner

WORKDIR /app

CMD ["/app/runner"]
