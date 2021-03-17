# build stage
FROM golang:1.16 as build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

COPY ./render /go/src/render
COPY ./gifcreator /go/src/gifcreator
COPY ./frontend /go/src/frontend

WORKDIR /go/src/render
RUN go install -v

WORKDIR /go/src/gifcreator
RUN go install -v

WORKDIR /go/src/frontend
RUN go install -v

# run stage
FROM busybox as run
COPY --from=build /go/bin/render /render
COPY --from=build /go/bin/gifcreator /gifcreator
COPY --from=build /go/bin/frontend /frontend
RUN mkdir /tmp/objcache
COPY ./gifcreator/scene /scene
ENV FRONTEND_TEMPLATES_DIR=/templates
