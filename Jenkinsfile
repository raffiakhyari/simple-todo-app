pipeline {

    agent any

    environment {
        REGISTRY = 'docker.io'
        DOCKER_DEV = 'raffiakhyari/todo-api-dev'
        DOCKER_PROD = 'raffiakhyari/todo-api'
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

        // ==========================================
        // 1. Checkout
        // ==========================================
        stage('Git') {
            steps {
                step([$class: 'WsCleanup'])

                checkout scm

                script {
                    env.AUTHOR_NAME = sh(
                        script: "git log -n 1 ${env.GIT_COMMIT} --format=%aN",
                        returnStdout: true
                    ).trim()

                    env.COMMIT_MESSAGE = sh(
                        script: "git log -1 --format=%B ${env.GIT_COMMIT}",
                        returnStdout: true
                    ).trim()

                    echo """
                    ==========================================
                    BUILD INFORMATION
                    ==========================================
                    Branch  : ${env.BRANCH_NAME}
                    Commit  : ${env.GIT_COMMIT}
                    Author  : ${env.AUTHOR_NAME}
                    Message : ${env.COMMIT_MESSAGE}
                    Build   : ${env.BUILD_NUMBER}
                    ==========================================
                    """
                }
            }
        }

        // ==========================================
        // 2. Prepare Image Name
        // ==========================================
        stage('Prepare Image') {
            steps {
                script {

                    if (env.BRANCH_NAME == 'main') {

                        env.DOCKER_NAME = env.DOCKER_PROD

                    } else {

                        env.DOCKER_NAME = env.DOCKER_DEV
                    }

                    env.IMAGE = "${env.DOCKER_NAME}:${env.BUILD_NUMBER}"

                    echo """
                    ==========================================
                    DOCKER IMAGE
                    ==========================================
                    Branch : ${env.BRANCH_NAME}
                    Image  : ${env.IMAGE}
                    ==========================================
                    """
                }
            }
        }

        // ==========================================
        // 3. Lint
        // ==========================================
        stage('Lint') {
            steps {
                sh '''
                    set -e

                    echo "=========================================="
                    echo "Go Version"
                    echo "=========================================="

                    go version

                    echo "=========================================="
                    echo "Checking gofmt"
                    echo "=========================================="

                    if [ -n "$(gofmt -l .)" ]; then
                        echo "ERROR: The following files are not formatted:"

                        gofmt -l .

                        exit 1
                    fi

                    echo "gofmt passed."

                    echo "=========================================="
                    echo "Running go vet"
                    echo "=========================================="

                    go vet ./...

                    echo "go vet passed."

                    echo "=========================================="
                    echo "Lint SUCCESS"
                    echo "=========================================="
                '''
            }
        }

        // ==========================================
        // 4. Unit Test
        // ==========================================
        stage('Unit Test') {
            steps {
                sh '''
                    set -e

                    echo "=========================================="
                    echo "Running Unit Tests"
                    echo "=========================================="

                    go test ./... -v

                    echo "=========================================="
                    echo "Unit Test SUCCESS"
                    echo "=========================================="
                '''
            }
        }

        // ==========================================
        // 5. Build Docker Image
        // ==========================================
        stage('Build Image') {
            steps {
                sh '''
                    set -e

                    echo "=========================================="
                    echo "Building Docker Image"
                    echo "=========================================="

                    echo "Image: ${IMAGE}"

                    DOCKER_BUILDKIT=1 docker build \
                        --pull \
                        -t "${IMAGE}" \
                        .

                    echo "=========================================="
                    echo "Docker Build SUCCESS"
                    echo "=========================================="

                    docker images "${DOCKER_NAME}"
                '''
            }
        }

        // ==========================================
        // 6. Trivy Vulnerability Scan
        // ==========================================
        stage('Trivy Scan') {
            steps {
                sh '''
                    set -e

                    echo "=========================================="
                    echo "Trivy Vulnerability Scan"
                    echo "=========================================="

                    trivy image \
                        --severity HIGH,CRITICAL \
                        --exit-code 1 \
                        --ignore-unfixed \
                        "${IMAGE}"

                    echo "=========================================="
                    echo "Trivy Scan SUCCESS"
                    echo "=========================================="
                '''
            }
        }

        // ==========================================
        // 7. Push Image
        // ==========================================
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

                        echo "=========================================="
                        echo "Docker Login"
                        echo "=========================================="

                        echo "${DOCKER_PASSWORD}" | docker login \
                            "${REGISTRY}" \
                            --username "${DOCKER_USERNAME}" \
                            --password-stdin

                        echo "=========================================="
                        echo "Pushing Image"
                        echo "=========================================="

                        docker push "${IMAGE}"

                        echo "=========================================="
                        echo "Docker Logout"
                        echo "=========================================="

                        docker logout "${REGISTRY}"

                        echo "=========================================="
                        echo "Push SUCCESS"
                        echo "=========================================="
                    }
                }
            }
        }
    }

    // ==========================================
    // Post
    // ==========================================
    post {

        success {
            echo """
            ==========================================
            PIPELINE SUCCESS
            ==========================================

            Application : todo-api
            Branch      : ${env.BRANCH_NAME}
            Build       : ${env.BUILD_NUMBER}
            Image       : ${env.IMAGE}
            Author      : ${env.AUTHOR_NAME}

            ==========================================
            """
        }

        ailure {
            echo """
            ==========================================
            PIPELINE FAILED
            ==========================================

            Application : todo-api
            Branch      : ${env.BRANCH_NAME}
            Build       : ${env.BUILD_NUMBER}

            ==========================================
            """
        }

        always {
            script {
                if (env.IMAGE) {
                    sh """
                        echo "Cleaning local image..."

                        docker image rm \
                            "${env.IMAGE}" \
                            || true
                    """
                }
            }
        }
    }
}