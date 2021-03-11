# the entire file needs changes
FROM golang:1.16

WORKDIR /go/src

COPY ./render .
COPY ./gifcreator .
COPY ./frontend .

WORKDIR /go/src/render
RUN go install -v ./...
WORKDIR /go/src/gifcreator
RUN go install -v ./...
WORKDIR /go/src/frontend
RUN go install -v ./...

COPY ./frontend/static /static
COPY ./frontend/templates /templates

COPY ./gifcreator/scene /scene

ENV FRONTEND_TEMPLATES_DIR=/templates
ENV FRONTEND_STATIC_DIR=/static
ENV SCENE_PATH=/scene
