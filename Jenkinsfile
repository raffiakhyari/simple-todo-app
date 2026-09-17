pipeline {
    agent any

    environment {
        APP_NAME   = 'todo-api'

        REGISTRY   = 'docker.io'
        IMAGE_NAME = 'raffiakhyari/todo-api'

        IMAGE_TAG  = "${BUILD_NUMBER}"
        IMAGE      = "${IMAGE_NAME}:${IMAGE_TAG}"
    }

    options {
        timestamps()
        disableConcurrentBuilds()

        buildDiscarder(
            logRotator(
                numToKeepStr: '10'
            )
        )
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Lint') {
            steps {
                sh '''
                    set -e

                    echo "=== Go Version ==="
                    go version

                    echo "=== gofmt check ==="

                    if [ -n "$(gofmt -l .)" ]; then
                        echo "ERROR: Files are not formatted:"
                        gofmt -l .
                        exit 1
                    fi

                    echo "=== go vet ==="
                    go vet ./...

                    echo "Lint passed."
                '''
            }
        }

        stage('Unit Test') {
            steps {
                sh '''
                    set -e

                    echo "=== Running unit tests ==="

                    go test ./... -v

                    echo "Unit tests passed."
                '''
            }
        }

        stage('Build Image') {
            steps {
                sh '''
                    set -e

                    echo "=== Building Docker Image ==="
                    echo "Image: ${IMAGE}"

                    docker build \
                        -t ${IMAGE} \
                        .

                    echo "Docker image built successfully."
                '''
            }
        }

        stage('Trivy Scan') {
            steps {
                sh '''
                    set -e

                    echo "=== Trivy Vulnerability Scan ==="

                    trivy image \
                        --severity HIGH,CRITICAL \
                        --exit-code 1 \
                        --ignore-unfixed \
                        ${IMAGE}

                    echo "Trivy scan passed."
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
                        set -e

                        echo "=== Docker Login ==="

                        echo "${DOCKER_PASSWORD}" | docker login \
                            ${REGISTRY} \
                            --username "${DOCKER_USERNAME}" \
                            --password-stdin

                        echo "=== Push Image ==="

                        docker push ${IMAGE}

                        docker logout ${REGISTRY}

                        echo "Image pushed successfully."
                    '''
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

            Application : ${APP_NAME}
            Image       : ${IMAGE}
            Build       : ${BUILD_NUMBER}

            ==========================================
            """
        }

        failure {
            echo """
            ==========================================
            PIPELINE FAILED
            ==========================================

            Application : ${APP_NAME}
            Build       : ${BUILD_NUMBER}

            ==========================================
            """
        }
    }
}