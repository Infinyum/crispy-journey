pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    environment {
        DOCKER_CREDENTIALS_ID = 'dockerhub'
        //KUBECONFIG_CREDENTIALS_ID = 'minikube'
        DOCKER_REGISTRY = 'docker.io'
        DOCKER_REGISTRY_NAMESPACE = 'elabrom'
        APP_NAME = 'crispy'
    }

    stages {
        stage('Testing') {
            steps {
                // Temporarily add Golang to path 
                withEnv(["PATH+GO=/usr/local/go/bin/"]) {
                    dir('src/crispy') {
                        sh 'go test -v'
                    }
                }
            }
        }

        stage('Docker build and push') {
            steps {
                // Temporarily add Docker to path
                withEnv(["PATH+DOCKER=/usr/local/bin"]) {
                    sh 'docker build -t $DOCKER_REGISTRY/$DOCKER_REGISTRY_NAMESPACE/$APP_NAME .'

                    withCredentials([usernamePassword(credentialsId: DOCKER_CREDENTIALS_ID, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh 'echo $DOCKER_PASSWORD | docker login docker.io -u $DOCKER_USERNAME --password-stdin'
                        sh 'docker push $DOCKER_REGISTRY/$DOCKER_REGISTRY_NAMESPACE/$APP_NAME'
                    }
                }
            }
        }
    }
}