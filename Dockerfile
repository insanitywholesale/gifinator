FROM --platform=$BUILDPLATFORM gcr.io/distroless/static-debian12:nonroot

COPY ./gifcreator/scene /scene
COPY ./frontend/templates /templates

COPY ./render/render /render
COPY ./gifcreator/gifcreator /gifcreator
COPY ./frontend/frontend /frontend

ENV FRONTEND_TEMPLATES_DIR=/templates
