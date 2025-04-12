FROM --platform=$BUILDPLATFORM busybox AS run

RUN mkdir /tmp/objcache
RUN mkdir /tmp/scene

COPY ./gifcreator/scene /tmp/scene
COPY ./frontend/templates /templates

COPY ./render/render /render
COPY ./gifcreator/gifcreator /gifcreator
COPY ./frontend/frontend /frontend

ENV FRONTEND_TEMPLATES_DIR=/templates
