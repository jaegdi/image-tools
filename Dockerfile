# Multistage Stage 0 as base
FROM default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-baseimages/fedora:40

LABEL maintainer="Dirk Jäger <dirk.jaeger@schufa.de>"

ENV LANG=en_US.UTF-8

USER root

# get the credentials from the secret, that is mounted in the build config
COPY build/scripts/fix-permissions.sh /usr/local/bin/fix-permissions.sh
COPY build/scripts/create-netrc.sh /usr/local/bin/create-netrc.sh
COPY build/pipeline-secrets/password /password
COPY build/pipeline-secrets/username /username
COPY build/certs/*.crt /etc/pki/ca-trust/source/anchors/
COPY image-tool /usr/bin/
COPY clusterconfig.json /usr/bin/

RUN chmod a+rx /usr/local/bin/fix-permissions.sh /usr/local/bin/create-netrc.sh \
    && /usr/local/bin/create-netrc.sh 'my-password-file' ${ART_HOSTNAME} \
    && chgrp root /usr/bin/image-tool \
    && chgrp root /usr/bin/clusterconfig.json \
    && chmod a+rx /usr/bin/image-tool \
    && chmod a+rw /usr/bin/clusterconfig.json \
    && update-ca-trust extract \
    && rm my-password-file


# COPY scptools-bitbucket/password /password
# COPY scptools-bitbucket/username /username

# RUN /usr/local/bin/create-netrc.sh 'art-password-file' ${ART_HOSTNAME} \
#     && curl -k --netrc-file art-password-file \
#         https://artifactory-pro.sf-rz.de:8443/artifactory/scptools-bin-develop/tools/kustomize/kustomize \
#         -o /usr/local/bin/kustomize       \
#     && curl -k --netrc-file art-password-file \
#         https://artifactory-pro.sf-rz.de:8443/artifactory/scptools-bin-develop/tools/yq/yq \
#         -o /usr/local/bin/yq              \
#     && curl -k --netrc-file art-password-file \
#         https://artifactory-pro.sf-rz.de:8443/artifactory/scptools-bin-develop/tools/pc/pc \
#         -o /usr/local/bin/pc              \
#     && curl -k --netrc-file art-password-file \
#         https://artifactory-pro.sf-rz.de:8443/artifactory/scptools-bin-develop/tools/cosign/cosign-linux-amd64 \
#         -o /usr/local/bin/cosign          \
#     && chmod a+x /usr/local/bin/kustomize  \
#     && chmod a+x /usr/local/bin/yq         \
#     && chmod a+x /usr/local/bin/pc         \
#     && chmod a+x /usr/local/bin/cosign     \
#     && ls -l /usr/local/bin                \
#     && rm art-password-file

# RUN export CURL_CA_BUNDLE=/etc/pki/tls/certs/ca-bundle.crt \
#     && INSTALL_PKGS="gettext rsync skopeo jq git" \
#     && yum -y --setopt=tsflags=nodocs install $INSTALL_PKGS \
#     && yum update -a \
#     && yum clean all

CMD ["image-tools", "-socks5=no", "-server", "-statcfg"]
