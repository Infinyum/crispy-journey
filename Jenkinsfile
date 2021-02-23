node {
    // Ensure the desired Go version is installed
    def root = tool type: 'go', name: 'go1.16'
            
    // Temporarily add Golang to path 
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        dir('src/crispy') {
            sh 'go test -v'
        }
    }

    // Temporarily add Docker to path
    withEnv(["PATH+DOCKER=/usr/local/bin"]) {
        docker.withRegistry('dockerhub') {
            docker.build('elabrom/crispy').push('latest')
        }
    }
}