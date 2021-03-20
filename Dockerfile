# build stage
FROM golang:1.16 as build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV GO111MODULE on

WORKDIR /go/src/gifinator
COPY . .

WORKDIR /go/src/gifinator/render
RUN go get -v
#RUN go vet -v
RUN go install -v

WORKDIR /go/src/gifinator/gifcreator
RUN go get -v
#RUN go vet -v
RUN go install -v

WORKDIR /go/src/gifinator/frontend
RUN go get -v
#RUN go vet -v
RUN make installwithvars

# run stage
FROM busybox as run
COPY --from=build /go/bin/render /render
COPY --from=build /go/bin/gifcreator /gifcreator
COPY --from=build /go/bin/frontend /frontend
RUN mkdir /tmp/objcache
COPY ./gifcreator/scene /scene
COPY ./frontend/templates /templates
ENV FRONTEND_TEMPLATES_DIR=/templates
