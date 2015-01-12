FROM openshift/jenkins-1-centos
MAINTAINER fabric8.io (http://fabric8.io/)

WORKDIR ~

USER root

RUN yum install -y bzr mercurial

RUN wget -q https://storage.googleapis.com/golang/go1.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.4.linux-amd64.tar.gz

RUN rm -rf go1.4.linux-amd64.tar.gz

RUN mkdir /go
RUN chown jenkins:jenkins /go

USER jenkins

ENV PATH $PATH:/usr/local/go/bin
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

RUN go get github.com/tools/godep
RUN cd $GOPATH/src/github.com/tools/godep && go install

# configure Jenkins
ADD http://updates.jenkins-ci.org/latest/ghprb.hpi $JENKINS_HOME/plugins/ghprb.hpi
ADD http://updates.jenkins-ci.org/latest/ssh-credentials.hpi $JENKINS_HOME/plugins/ssh-credentials.hpi
ADD http://updates.jenkins-ci.org/latest/git-client.hpi $JENKINS_HOME/plugins/git-client.hpi
ADD http://updates.jenkins-ci.org/latest/scm-api.hpi $JENKINS_HOME/plugins/scm-api.hpi
ADD http://updates.jenkins-ci.org/latest/credentials.hpi $JENKINS_HOME/plugins/credentials.hpi
ADD http://updates.jenkins-ci.org/latest/commit-message-trigger-plugin.hpi $JENKINS_HOME/plugins/commit-message-trigger-plugin.hpi
ADD http://updates.jenkins-ci.org/latest/ruby-runtime.hpi $JENKINS_HOME/plugins/
ADD http://updates.jenkins-ci.org/latest/git.hpi $JENKINS_HOME/plugins/git.hpi
ADD http://updates.jenkins-ci.org/latest/ssh-agent.hpi $JENKINS_HOME/plugins/ssh-agent.hpi
ADD http://updates.jenkins-ci.org/latest/github-api.hpi $JENKINS_HOME/plugins/github-api.hpi
ADD http://updates.jenkins-ci.org/latest/github.hpi $JENKINS_HOME/plugins/github.hpi

ADD config.xml $JENKINS_HOME/jobs/origin-schema-generator/config.xml

USER root
RUN chown -R jenkins:jenkins /var/jenkins_home
USER jenkins