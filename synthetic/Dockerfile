FROM alpine:latest

RUN apk add --no-cache curl

COPY setup.sh /setup.sh
RUN chmod +x /setup.sh
RUN /setup.sh

ENTRYPOINT [ "hey", "-z", "10000m", "-c", "10", "--host", "foo.localhost", "http://ingress-nginx-controller.ingress-nginx.svc.cluster.local" ]
