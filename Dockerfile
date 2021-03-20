# build stage
FROM golang:1.16 as build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

COPY ./render /go/src/render
WORKDIR /go/src/render
RUN go get -v
RUN go vet -v
RUN go install -v

COPY ./gifcreator /go/src/gifcreator
WORKDIR /go/src/gifcreator
RUN go get -v
RUN go vet -v
RUN go install -v

COPY ./frontend /go/src/frontend
WORKDIR /go/src/frontend
RUN go get -v
RUN go vet -v
RUN go install -v

# run stage
FROM busybox as run
COPY --from=build /go/bin/render /render
COPY --from=build /go/bin/gifcreator /gifcreator
COPY --from=build /go/bin/frontend /frontend
RUN mkdir /tmp/objcache
COPY ./gifcreator/scene /scene
ENV FRONTEND_TEMPLATES_DIR=/templates
