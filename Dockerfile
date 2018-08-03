FROM centos:7

RUN yum update -y && \
  yum install -y centos-release-scl && \
  yum install -y docker sclo-git212-git unzip java-1.8.0-openjdk-devel java-1.8.0-openjdk-devel.i686 which zip wget make hg svn bzr gcc https://centos7.iuscommunity.org/ius-release.rpm && \
  yum clean all

ENV PATH=/opt/rh/sclo-git212/root/bin:$PATH

RUN curl --retry 999 --retry-max-time 0  -sSL https://bintray.com/artifact/download/fabric8io/helm-ci/helm-v0.1.0%2B825f5ef-linux-amd64.zip > helm.zip && \
  unzip helm.zip && \
  rm -f helm.zip && \
  mv helm /usr/bin/

# Maven
RUN curl -L http://mirrors.ukfast.co.uk/sites/ftp.apache.org/maven/maven-3/3.3.9/binaries/apache-maven-3.3.9-bin.tar.gz | tar -C /opt -xzv

ENV M2_HOME /opt/apache-maven-3.3.9
ENV maven.home $M2_HOME
ENV M2 $M2_HOME/bin
ENV PATH $M2:$PATH
RUN mkdir --parents --mode 777 /root/.mvnrepository

ENV GOLANG_VERSION 1.8.1
RUN wget https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
  rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV GLIDE_VERSION v0.12.3
ENV GO15VENDOREXPERIMENT 1
RUN wget https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.tar.gz && \
  tar -xzf glide-$GLIDE_VERSION-linux-amd64.tar.gz && \
  mv linux-amd64 /usr/local/glide && \
  rm glide-$GLIDE_VERSION-linux-amd64.tar.gz

ENV GH_RELEASE_VERSION 2.2.1
RUN wget https://github.com/progrium/gh-release/releases/download/v$GH_RELEASE_VERSION/gh-release_${GH_RELEASE_VERSION}_linux_x86_64.tgz && \
  tar -xzf gh-release_${GH_RELEASE_VERSION}_linux_x86_64.tgz && \
  mv gh-release /usr/local/gh-release && \
  rm gh-release_${GH_RELEASE_VERSION}_linux_x86_64.tgz


ENV PATH $PATH:/usr/local/go/bin
ENV PATH $PATH:/usr/local/glide
ENV PATH $PATH:/usr/local/
ENV GOROOT /usr/local/go
ENV PATH $PATH:/go/bin
ENV GOPATH=/go

RUN go get github.com/fabric8io/gobump

ENV FABRIC8_USER_NAME=fabric8

RUN useradd --user-group --create-home --shell /bin/false ${FABRIC8_USER_NAME}

ENV HOME=/home/${FABRIC8_USER_NAME}
ENV WORKSPACE=$GOPATH/src/github.com/fabric8io/kubernetes-model
RUN mkdir $WORKSPACE

COPY . $WORKSPACE

WORKDIR $WORKSPACE/