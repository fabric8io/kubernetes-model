#!/usr/bin/groovy
/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def dummy
clientsTemplate{
  mavenNode {
    ws{
      checkout scm
      readTrusted 'release.groovy'
      if (env.BRANCH_NAME.startsWith('PR-')){
        echo 'Running CI pipeline'
        container(name: 'maven') {
          sh 'mvn clean install'
        }
      } else if (env.BRANCH_NAME.equals('master')){
        echo 'Running CD pipeline'
        sh "git remote set-url origin git@github.com:fabric8io/kubernetes-model.git"

        def pipeline = load 'release.groovy'

        stage 'Stage'
        def stagedProject = pipeline.stage()

        stage 'Promote'
        pipeline.release(stagedProject)

        stage 'Update downstream dependencies'
        pipeline.updateDownstreamDependencies(stagedProject)
      }
    }
  }
}
