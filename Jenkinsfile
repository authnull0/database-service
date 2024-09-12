pipeline {
    agent any

    environment {
    GITHUB_REPO = 'https://github.com/authnull0/database-service.git'
    GITHUB_BRANCH = 'production'
    DOCKER_REGISTRY = 'docker-repo.authnull.com'
    DOCKER_REGISTRY_CREDENTIALS = credentials('authnull-repo')
    DOCKER_IMAGE = 'docker-repo.authnull.com/database-service:production'
    }

    triggers {
        githubPush()
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '10', daysToKeepStr: '', numToKeepStr: '10')
    }

    stages {
        stage('Checkout') {
            steps {
            git credentialsId: '3d061b0d-704f-434d-8c70-dce67d4c42d1', url: "${GITHUB_REPO}", branch: "${GITHUB_BRANCH}"
            }
        }
        stage('Build Docker Image') {
            steps {
                sh 'docker build -t ${DOCKER_IMAGE} .'
            }
        }
        stage('Login to Docker Artifactory') {
            steps {
                sh 'echo ${DOCKER_REGISTRY_CREDENTIALS_PSW} | docker login ${DOCKER_REGISTRY} -u ${DOCKER_REGISTRY_CREDENTIALS_USR} --password-stdin'
            }
        }
        stage('Push Docker Image') {
            steps {
                sh 'docker push ${DOCKER_IMAGE}'
            }
        }
        stage('Logout from Docker Artifactory') {
            steps {
                sh 'docker logout ${DOCKER_REGISTRY}'
            }
        }
        stage('Remove Docker Image') {
            steps {
                sh 'docker rmi ${DOCKER_IMAGE}'
            }
        }
    }
}
