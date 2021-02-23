node {
    // Temporarily add Golang to path 
    withEnv(["PATH+GO=/usr/local/go/bin/"]) {
        dir('src/crispy') {
            sh 'go test -v'
        }
    }

    // Temporarily add Docker to path
    withEnv(["PATH+DOCKER=/usr/local/bin"]) {
        withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
            sh 'echo $DOCKER_PASSWORD | docker login docker.io -u $DOCKER_USERNAME --password-stdin'
            docker.build('elabrom/crispy').push('latest')
        }
    }
}