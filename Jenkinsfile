pipeline {
    agent any
    tools {
        go 'go1.19'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'apt install -y build-essential'
                sh 'go version'
                sh 'go mod download'
            }
        }
        
        // stage('Build') {
        //     steps {
        //         echo 'Compiling and building'
        //         sh 'go build'
        //     }
        // }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running test'
                    sh 'make unit-test'
                }
            }
        }
        
    }
}