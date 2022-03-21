FROM registry.ci.openshift.org/openshift/release:golang-1.17 as watcher-bin
ENV GOFLAGS "-mod=mod"
WORKDIR /go/src/github.com/jwmatthews/case_watcher
COPY . /go/src/github.com/jwmatthews/case_watcher
RUN go mod download
RUN go build -a -o /build/case_watcher

FROM registry.redhat.io/ubi8/ubi:latest
COPY --from=watcher-bin  /build/case_watcher /usr/local/bin/case_watcher
COPY --from=watcher-bin  /go/src/github.com/jwmatthews/case_watcher/.case_watcher.yml.example /case_watcher.yml.example

RUN dnf -y install sqlite
ENTRYPOINT ["/usr/local/bin/case_watcher"]
