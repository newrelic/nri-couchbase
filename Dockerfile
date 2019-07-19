FROM golang:1.10 as builder
RUN go get -d github.com/newrelic/nri-couchbase/... && \
    cd /go/src/github.com/newrelic/nri-couchbase && \
    make && \
    strip ./bin/nr-couchbase

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-couchbase/bin/nr-couchbase /var/db/newrelic-infra/newrelic-integrations/bin/nr-couchbase
COPY --from=builder /go/src/github.com/newrelic/nri-couchbase/couchbase-definition.yml /var/db/newrelic-infra/newrelic-integrations/definition.yml
