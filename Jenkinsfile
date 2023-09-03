pipeline {
    //install golang 1.19 on Jenkins node
    agent any
    tools {
      go 'go-1.21.0'
    }
   
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        POSTGRES_VOLUME = "postgres_data"
        DOCKER_NETWORK = "pika-network"
    }
   
    options {
        buildDiscarder(logRotator(daysToKeepStr: '90', numToKeepStr: '1', artifactDaysToKeepStr: '90', artifactNumToKeepStr: '20'))
        timeout(time: 120, unit: 'MINUTES')
    }
   
     stages {
        stage("initialize github") {
            steps {
                echo 'COMPILING PROJECT FROM GITHUB REPO'
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[credentialsId: '0e57d3d4-8232-46b1-956c-e44ce7aea3d0', url: 'https://github.com/CossyCossy/docker-deployment-test']])
            }
        }
        stage("initialize docker") {
            steps {
                echo 'LOGIN TO DOCKER'
                script {
                    withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'dockerhubPassword', usernameVariable: 'dockerhubUser')]) {
                    bat "docker login -u ${env.dockerhubUser} -p ${env.dockerhubPassword}"
                }
              }
            }
        }
        stage("create-volumes") {
            agent any
            steps {
                echo 'CREATE POSTGRESS VOLUME'
                script {
                    bat "docker volume create ${POSTGRES_VOLUME}"
                    bat 'docker volume ls'
              }
            }
        }
        // stage("build-db-image") {
        //     agent any
        //     steps {
        //         echo 'PULL AND PUSH POSGRESS IMAGE'
        //         script {
        //             bat 'docker pull postgres:15.2-alpine'
        //             bat 'docker tag postgres:15.2-alpine cossycossy/db:latest'
        //             bat 'docker image list'

        //       }
        //     }
        // }
        stage("build-bg-image") {
            steps {
                dir('bg') {
                    echo 'STARTED BACKEND BUILD'
                    script {
                        bat 'go version'
                        bat 'docker version'
                        bat 'go get ./...'
                        bat 'docker build . -t bg:latest'
                        bat 'docker tag bg:latest cossycossy/bg:latest'
                        bat 'docker image list'
                    }
                   
                }
            }
        }
        stage("build-front-image") {
            steps {
                dir('front') {
                    echo 'STARTED BACKEND BUILD'
                    script {
                        bat 'docker build . -t front'
                        bat 'docker tag front:latest cossycossy/front:latest'
                        bat 'docker image list'
                    }
                }
            }
        }
        stage("create-docker-network") {
            steps {
                    echo 'CREATING DOCKER NETWORK'
                    script {
                        bat "docker network ls | findstr ${DOCKER_NETWORK} || docker network create ${DOCKER_NETWORK}"
                    }
                   
            }
        }
        stage("move-images-hub and start containers") {
            steps {
                    echo 'PUSHING DB TO DOCKERHUB AND RUN CONTAINER'
                    // script {
                    //     bat 'docker push cossycossy/db:latest'
                    //     bat "docker run -d --name db --network ${DOCKER_NETWORK} -e POSTGRES_DB=bg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=admin -p 5432:5432 -v /${POSTGRES_VOLUME}:/var/lib/postgresql/data cossycossy/db"                   
                    //      }
                   
                    echo 'BG TO HUB'
                    script {
                        bat 'docker push cossycossy/bg:latest'
                        bat "docker run --name bg --network ${DOCKER_NETWORK} --restart on-failure:5 -p 8000:8000 -d cossycossy/bg"
                    }
                   
                    echo 'UI TO HUB'
                    script {
                        bat 'docker push cossycossy/front:latest'
                        bat "docker run --name front --network ${DOCKER_NETWORK} --restart always  -p 80:80 -d cossycossy/front" 
                    }
                   
            }
        }
     }
}