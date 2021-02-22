pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'ls src/crispy'

                // Ensure the desired Go version is installed
                tool name: 'go1.16', type: 'go'

                // Export environment variables pointing to the directory where Go was installed
                withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                    echo "$GOROOT"
                    sh 'go version'
                }
            }
        }
    }
}