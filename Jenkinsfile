node {
    // Ensure the desired Go version is installed
    def root = tool type: 'go', name: 'go1.16'
            
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        dir('src/crispy') {
            // TODO: static code analysis
            sh 'go test -v'
        }
    }

    //env.PATH="${root}/bin:$PATH"
    
    def img = docker.build 'crispy'
}