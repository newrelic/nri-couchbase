FROM golang:1.10 as builder
COPY . /go/src/github.com/newrelic/nri-couchbase/
RUN cd /go/src/github.com/newrelic/nri-couchbase && \
    make && \
    strip ./bin/nri-couchbase

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-couchbase/bin/nri-couchbase /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-couchbase
COPY --from=builder /go/src/github.com/newrelic/nri-couchbase/couchbase-definition.yml /nri-sidecar/newrelic-infra/newrelic-integrations/definition.yml
USER 1000
