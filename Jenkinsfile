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
                    sh 'go test ./... -v'
            }
        }
    }

    stage('Build-Docker-Image') {
      steps {
        container('docker') {
          sh 'docker build -t adarshtw/passwordly_backend:${BUILD_NUMBER} .'
        }
      }
    }

    stage('Login-Into-Docker') {
        steps {
            container('docker') {
                script {
                    withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
                      sh "docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD"
                    }
                }
            }
        }
    }

    stage('Push-Images-Docker-to-DockerHub') {
      steps {
        container('docker') {
          sh 'docker push adarshtw/passwordly_backend:${BUILD_NUMBER}'
        }
      }
    }

    stage('Update K8s Deployment') {
      steps {
        container('docker') {
          sh 'apk add gettext'
          sh 'envsubst < deployment/k8s/manifest.yaml > manifest.yaml'
          withKubeConfig([namespace: "default"]) {
            sh 'kubectl apply -f manifest.yaml'
          }
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