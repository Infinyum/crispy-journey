node {
    // Ensure the desired Go version is installed
    def root = tool type: 'go', name: 'go1.16'
            
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        // Move to directory where Go module is
        dir('src/crispy') {
            sh 'go version'
            sh 'go test -v'
        }
    }
}