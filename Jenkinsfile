node {
    // Ensure the desired Go version is installed
    def root = tool type: 'go', name: 'go1.16'
            
    // Temporarily add Golang to path 
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        dir('src/crispy') {
            // TODO: static code analysis
            sh 'go test -v'
        }
    }

    // Temporarily add Docker to path
    withEnv(["PATH+DOCKER=/usr/local/bin"]) {
        docker.build('crispy').push('latest')
    }
}