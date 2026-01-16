FROM golang:1.25-alpine AS apiBuilder
WORKDIR /opt
RUN apk update && apk add --no-cache make
COPY . /opt/
RUN go mod download
ARG GIT_TAG
ARG GIT_COMMIT
ARG GIT_COMMIT_DATE
RUN make build GIT_TAG=${GIT_TAG} GIT_COMMIT=${GIT_COMMIT} GIT_COMMIT_DATE=${GIT_COMMIT_DATE}

FROM alpine
ENV ENVIRONMENT=production
ENV CGO_ENABLED=0
WORKDIR /opt
RUN apk update && \
    apk add --no-cache curl ca-certificates && \
    update-ca-certificates
COPY --from=apiBuilder /opt/exitgatebot /opt/exitgatebot
CMD [ "./exitgatebot", "run" ]
