FROM alpine:3.11
ADD omo.msa.user /usr/bin/omo.msa.user
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.user" ]
