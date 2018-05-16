#  # build stage
#  FROM golang:alpine AS build-env
#  COPY main.go /src/main.go
#  COPY vendor /go/src/
#  RUN cd /src && go build -o app && ls -la
#
#  # final stage
#  FROM alpine
#  COPY --from=build-env /src/app /bin/app
#  COPY check.sh /opt/resource/check
#  COPY in.sh /opt/resource/in
#  COPY out.sh /opt/resource/check

FROM ubuntu
COPY kube /kube/
COPY main /bin/k8s-app
COPY check.sh /opt/resource/check
COPY in.sh /opt/resource/in
COPY out.sh /opt/resource/out
RUN chmod +x /bin/k8s-app /opt/resource/check /opt/resource/in /opt/resource/check