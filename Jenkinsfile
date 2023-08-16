pipeline {
    agent { node { label 'kubeagent' }}
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
                sh 'apt update && apt install -y build-essential'
                sh 'go version'
                sh 'go mod download'
            }
        }
        
        stage('Build') {
            steps {
                echo 'Building'
                sh '''
                /kaniko/executor --dockerfile `pwd`/Dockerfile \
                --context `pwd` --destination=adarshtw/passwordly_backend:${BUILD_NUMBER}
                '''
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running test'
                    sh 'make unit-tests'
                }
            }
        }
        
    }
}