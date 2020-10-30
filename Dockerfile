FROM golang:1.15.2 as build

ENV DEBIAN_FRONTEND noninteractive
ENV TZ Asia/Tokyo
ENV WORK_DIR /opt/tiny-cluster

WORKDIR ${WORK_DIR}

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        make \
 && apt-get clean \
 && apt-get autoremove --purge \
 && rm -rf /var/lib/apt/list/*

COPY . . 

RUN go mod download

RUN make

FROM scratch as prod

COPY --from=build /opt/tiny-cluster/bin/ /usr/local/bin/

CMD ["bootserver",  "start",  "-p", "8080"]