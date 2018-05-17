FROM ubuntu
COPY kube /kube/
COPY main /bin/k8s-app
COPY scripts /opt/resource/
RUN chmod +x /bin/k8s-app /opt/resource/check /opt/resource/in /opt/resource/check
