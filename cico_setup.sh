#!/bin/bash

# Output command before executing
set -x

# Exit on error
set -e

# Source environment variables of the jenkins slave
# that might interest this worker.
function load_jenkins_vars() {
  if [ -e "jenkins-env.json" ]; then
    eval "$(./env-toolkit load -f jenkins-env.json JENKINS_URL GIT_BRANCH GIT_COMMIT BUILD_NUMBER ghprbSourceBranch ghprbActualCommit BUILD_URL ghprbPullId)"
  fi
}

function install_deps() {
  # We need to disable selinux for now, XXX
  /usr/sbin/setenforce 0 || :

  # Get all the deps in
  yum -y install docker

  service docker start

  echo 'CICO: Dependencies installed'
}

function setup() {
  load_jenkins_vars
  install_deps
}

function build_run_tests() {
  docker exec $CID make gobuild
  docker exec $CID mvn clean install
}