pipeline {
    agent any

    environment {
        REGISTRY    = 'docker.io'
        DOCKER_DEV  = 'raffiakhyari/todo-api-dev'
        DOCKER_PROD = 'raffiakhyari/todo-api'
    }

    options {
        timestamps()
        disableConcurrentBuilds()
        skipDefaultCheckout(true)
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {

        stage('Checkout') {
            steps {
                cleanWs()
                checkout scm

                script {
                    env.ACTUAL_COMMIT = sh(
                        script: 'git rev-parse HEAD',
                        returnStdout: true
                    ).trim()

                    env.AUTHOR_NAME = sh(
                        script: 'git log -1 --format=%aN HEAD',
                        returnStdout: true
                    ).trim()

                    env.COMMIT_MESSAGE = sh(
                        script: 'git log -1 --format=%s HEAD',
                        returnStdout: true
                    ).trim()

                    echo """
==========================================
BUILD INFORMATION
==========================================
Branch  : ${env.BRANCH_NAME}
Commit  : ${env.ACTUAL_COMMIT}
Author  : ${env.AUTHOR_NAME}
Message : ${env.COMMIT_MESSAGE}
Build   : ${env.BUILD_NUMBER}
==========================================
"""
                }
            }
        }

        stage('Prepare') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'develop') {
                        env.DOCKER_NAME = env.DOCKER_DEV
                        env.DEPLOY_ENV = 'dev'
                    } else if (env.BRANCH_NAME == 'main') {
                        env.DOCKER_NAME = env.DOCKER_PROD
                        env.DEPLOY_ENV = 'production'
                    } else {
                        error("Unsupported branch: ${env.BRANCH_NAME}")
                    }

                    env.IMAGE = "${env.DOCKER_NAME}:${env.BUILD_NUMBER}"

                    echo """
==========================================
IMAGE INFORMATION
==========================================
Branch      : ${env.BRANCH_NAME}
Environment : ${env.DEPLOY_ENV}
Image       : ${env.IMAGE}
==========================================
"""
                }
            }
        }

        stage('Lint & Test') {
            steps {
                sh '''
                    set -eu

                    echo "=== Go Version ==="
                    go version

                    echo "=== Check gofmt ==="
                    UNFORMATTED=$(gofmt -l .)

                    if [ -n "$UNFORMATTED" ]; then
                        echo "ERROR: Files need formatting:"
                        echo "$UNFORMATTED"
                        echo "Run: gofmt -w <files>"
                        exit 1
                    fi

                    echo "=== Go Vet ==="
                    go vet ./...

                    echo "=== Unit Tests ==="
                    go test ./... -v

                    echo "Lint and tests passed."
                '''
            }
        }

        stage('Build Image') {
            steps {
                sh '''
                    set -eu

                    DOCKER_BUILDKIT=1 docker build \
                        --pull \
                        -t "${IMAGE}" \
                        .

                    docker image inspect "${IMAGE}" >/dev/null
                '''
            }
        }

        stage('Trivy Scan') {
            steps {
                sh '''
                    set -eu

                    trivy image \
                        --severity HIGH,CRITICAL \
                        --exit-code 1 \
                        --ignore-unfixed \
                        "${IMAGE}"
                '''
            }
        }

        stage('Push Image') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'dockerhub-credentials',
                        usernameVariable: 'DOCKER_USERNAME',
                        passwordVariable: 'DOCKER_PASSWORD'
                    )
                ]) {
                    sh '''
                        set -eu

                        echo "$DOCKER_PASSWORD" | docker login \
                            "$REGISTRY" \
                            --username "$DOCKER_USERNAME" \
                            --password-stdin

                        docker push "$IMAGE"
                    '''
                }
            }
        }

        stage('Trigger CD') {
            steps {
                script {
                    def cdJob = "todo-api-delivery/${env.BRANCH_NAME}"

                    build job: cdJob,
                        parameters: [
                            string(
                                name: 'IMAGE_REPO',
                                value: env.DOCKER_NAME
                            ),
                            string(
                                name: 'IMAGE_TAG',
                                value: env.BUILD_NUMBER
                            )
                        ],
                        wait: false
                }
            }
        }
    }

    post {
        success {
            echo """
==========================================
PIPELINE SUCCESS
==========================================
Application : simple-todo-app
Branch      : ${env.BRANCH_NAME}
Environment : ${env.DEPLOY_ENV}
Build       : ${env.BUILD_NUMBER}
Image       : ${env.IMAGE}
Author      : ${env.AUTHOR_NAME}
==========================================
"""
        }

        failure {
            echo """
==========================================
PIPELINE FAILED
==========================================
Job    : ${env.JOB_NAME}
Branch : ${env.BRANCH_NAME}
Build  : ${env.BUILD_NUMBER}
CD     : Not triggered if CI failed
==========================================
"""
        }

        always {
            script {
                if (env.IMAGE) {
                    sh '''
                        docker image rm "${IMAGE}" || true
                    '''
                }

                sh 'docker logout "$REGISTRY" >/dev/null 2>&1 || true'
            }
        }
    }
}