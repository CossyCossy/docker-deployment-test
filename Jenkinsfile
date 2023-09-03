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
        DOCKER_NETWORK = "docker-deployment-test_pika-network"
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
                    echo "${env.BUILD_ID} ${BUILD_ID}"
                    withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'dockerhubPassword', usernameVariable: 'dockerhubUser')]) {
                    bat "docker login -u ${env.dockerhubUser} -p ${env.dockerhubPassword}"
                }
              }
            }
        }
        stage("check volumes") {
            agent any
            steps {
                echo 'CHECK POSGRESS VOLUME'
                script {
                  //  bat "docker volume create ${POSTGRES_VOLUME}"
                    bat 'docker volume ls'
              }
            }
        }
        stage("build-bg-image") {
            steps {
                dir('bg') {
                    echo 'STARTED BACKEND BUILD'
                    script {
                        bat 'go version'
                        bat 'docker version'
                        bat 'go get ./...'
                        // bat 'docker build . -t bg:latest'
                        bat "docker build . -t bg:${BUILD_ID}"
                        bat "docker tag bg:${BUILD_ID} bg:latest"
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
                        bat "docker build . -t front:${BUILD_ID}"
                        bat "docker tag front:${BUILD_ID} front:latest"
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
                
                    echo 'PUSHING BG TO DOCKERHUB AND RUN CONTAINER'
                    script {

                        //stop the container
                        bat 'docker stop bg'

                        //remove the container
                        bat 'docker rm bg'

                        bat "docker run --name bg --network ${DOCKER_NETWORK} --restart on-failure:5 -p 8000:8000 -d bg:latest"

                        // create image for docker hub
                        bat 'docker tag bg:latest cossycossy/bg:latest'
                        
                        //push image to docker hub
                        bat 'docker push cossycossy/bg:latest'
                    }
            
                    echo 'PUSHING FRONT TO DOCKERHUB AND RUN CONTAINER'
                    script {
                          //stop the container
                        bat 'docker stop front'

                        //remove the container
                        bat 'docker rm front'

                        bat "docker run --name front --network ${DOCKER_NETWORK} --restart always -p 80:80 -d front:latest"

                        // create image for docker hub
                        bat 'docker tag front:latest cossycossy/front:latest'
                        
                        //push image to docker hub
                        bat 'docker push cossycossy/front:latest'
                    }
                   
            }
        }
         stage("remover local image") {
             steps {
                echo 'REMOVING LOCAL IMAGES'
                script {
                    echo 'REMOVING BG IMAGES'
                    bat "docker rmi bg:${BUILD_ID}"
                    bat 'docker rmi cossycossy/bg:latest'

                    echo 'REMOVING FRONT IMAGES'
                    bat "docker rmi front:${BUILD_ID}"
                    bat 'docker rmi cossycossy/front:latest'

                    bat 'docker image list'
                }
             }
         }
     }
}