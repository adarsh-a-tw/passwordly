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
                sh 'go version'
                sh 'go mod download'
            }
        }
        
        stage('Build') {
            steps {
                script {
                  // DOCKER HUB
                  
                  /* Build the container image */            
                  def dockerImage = docker.build("adarshtw/passwordly_backend:${BUILD_NUMBER}")
                        
                  /* Push the container to the docker Hub */
                  dockerImage.push()

                  /* Remove docker image*/
                  sh 'docker rmi -f my-image:${env.BUILD_ID}'

                } 
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