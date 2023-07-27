FROM alpine
ADD pod_api /pod_api
#ADD filebeat.yml /filebeat.yml

ENTRYPOINT [ "/pod_api" ]