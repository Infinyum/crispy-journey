pipeline {
    agent any

    stages {
        ws('${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/') {
            withEnv(['GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}']) {
                env.PATH="${GOPATH}/bin:$PATH"
                
                stage('Go test and build') {
                    steps {
                        sh 'go version'
                    }
                }
            }
        }
    }
}