// pipeline {
//     agent { node { label 'kubeagent' }}
//     tools {
//         go 'go1.19'
//     }
//     environment {
//         GO114MODULE = 'on'
//         CGO_ENABLED = 0 
//         GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
//     }
//     stages {        
//         stage('Pre Test') {
//             steps {
//                 echo 'Installing dependencies'
//                 sh 'go version'
//                 sh 'go mod download'
//             }
//         }
        
//         stage('Build') {
//             steps {
//                 script {
//                   // DOCKER HUB
                  
//                   /* Build the container image */            
//                   def dockerImage = docker.build("adarshtw/passwordly_backend:${BUILD_NUMBER}")
                        
//                   /* Push the container to the docker Hub */
//                   dockerImage.push()

//                   /* Remove docker image*/
//                   sh 'docker rmi -f my-image:${env.BUILD_ID}'

//                 } 
//             }
//         }

//         stage('Test') {
//             steps {
//                 withEnv(["PATH+GO=${GOPATH}/bin"]){
//                     echo 'Running vetting'
//                     sh 'go vet .'
//                     echo 'Running test'
//                     sh 'make unit-tests'
//                 }
//             }
//         }
        
//     }
// }


pipeline {
  agent {
    kubernetes {
      yaml '''
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: go
            image: golang:1.19
            command:
            - cat
            tty: true
          - name: docker
            image: docker:latest
            command:
            - cat
            tty: true
            volumeMounts:
             - mountPath: /var/run/docker.sock
               name: docker-sock
          volumes:
          - name: docker-sock
            hostPath:
              path: /var/run/docker.sock    
        '''
    }
  }

 environment {
        DOCKER_CREDENTIAL = credentials('dockerhub')
 }

  stages {

    stage('Pre-Tests') {
        steps {
            container('go'){
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go mod download'
            }
        }
    }

    stage('Tests') {
        steps {
             container('go'){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running test'
                    sh 'make unit-tests'
            }
        }
    }


    stage('Build-Docker-Image') {
      steps {
        container('docker') {
          sh 'docker build -t adarshtw/passwordly_backend:latest .'
        }
      }
    }
    stage('Login-Into-Docker') {
      steps {
        container('docker') {
          script {
            withCredentials([string(credentialsId: 'dockerhub', variable: 'DOCKER_CREDENTIAL')]) {
                def creds = DOCKER_CREDENTIAL.tokenize(':')
                def dockerUsername = creds[0]
                def dockerPassword = creds[1]

                sh "docker login -u $dockerUsername -p $dockerPassword"
            }
        }
      }
    }
    }
     stage('Push-Images-Docker-to-DockerHub') {
      steps {
        container('docker') {
          sh 'docker push adarshtw/passwordly_backend:latest'
      }
    }
     }
  }
    post {
      always {
        container('docker') {
          sh 'docker logout'
      }
      }
    }
}