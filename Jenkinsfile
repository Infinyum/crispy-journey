node {
    // Temporarily add Golang to path 
    withEnv(["PATH+GO=/usr/local/go/bin/"]) {
        dir('src/crispy') {
            sh 'go test -v'
        }
    }

    // Temporarily add Docker to path
    withEnv(["PATH+DOCKER=/usr/local/bin"]) {
        withDockerRegistry(credentialsId: 'dockerhub') {
            docker.build('elabrom/crispy').push('latest')
        }
    }
}