# build stage
FROM golang:1.21 as build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV GO111MODULE on

WORKDIR /go/src/gifinator
COPY . .

WORKDIR /go/src/gifinator/render
RUN go get
RUN go vet
RUN go install

WORKDIR /go/src/gifinator/gifcreator
RUN go get
RUN go vet
RUN go install

WORKDIR /go/src/gifinator/frontend
RUN go get
RUN go vet
RUN make installwithvars

RUN ls /go/bin

# run stage
FROM busybox as run
RUN mkdir /tmp/objcache
RUN mkdir /tmp/scene
COPY ./gifcreator/scene /tmp/scene
COPY ./frontend/templates /templates
COPY --from=build /go/bin/render /render
COPY --from=build /go/bin/gifcreator /gifcreator
COPY --from=build /go/bin/frontend /frontend
ENV FRONTEND_TEMPLATES_DIR=/templates
