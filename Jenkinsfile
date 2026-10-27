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
                    env.AUTHOR_NAME = sh(
                        script: 'git log -1 --format=%aN',
                        returnStdout: true
                    ).trim()

                    env.COMMIT_MESSAGE = sh(
                        script: 'git log -1 --format=%s',
                        returnStdout: true
                    ).trim()

                    echo "Branch: ${env.BRANCH_NAME}"
                    echo "Commit: ${env.GIT_COMMIT}"
                    echo "Author: ${env.AUTHOR_NAME}"
                    echo "Message: ${env.COMMIT_MESSAGE}"
                }
            }
        }

        stage('Prepare') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'develop') {
                        env.DOCKER_NAME = env.DOCKER_DEV
                    } else if (env.BRANCH_NAME == 'main') {
                        env.DOCKER_NAME = env.DOCKER_PROD
                    } else {
                        error("Unsupported branch: ${env.BRANCH_NAME}")
                    }

                    env.IMAGE = "${env.DOCKER_NAME}:${env.BUILD_NUMBER}"

                    echo "Docker Image: ${env.IMAGE}"
                }
            }
        }

        stage('Lint & Test') {
            steps {
                sh '''
                    set -eu

                    docker run --rm \
                        -v "$PWD:/app" \
                        -w /app \
                        golang:1.25 \
                        sh -c '
                            go version

                            UNFORMATTED=$(gofmt -l .)
                            if [ -n "$UNFORMATTED" ]; then
                                echo "ERROR: Files need formatting:"
                                echo "$UNFORMATTED"
                                exit 1
                            fi

                            go vet ./...
                            go test ./... -v
                        '
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

                        echo "${DOCKER_PASSWORD}" | docker login \
                            "${REGISTRY}" \
                            --username "${DOCKER_USERNAME}" \
                            --password-stdin

                        docker push "${IMAGE}"

                        docker logout "${REGISTRY}"
                    '''
                }
            }
        }
    }

    post {
        success {
            script {
                def cdJob = "todo-api-delivery/${env.BRANCH_NAME}"

                echo "CI SUCCESS: ${env.IMAGE}"
                echo "Triggering CD job: ${cdJob}"

                build job: cdJob,
                    parameters: [
                        string(name: 'IMAGE_REPO', value: env.DOCKER_NAME),
                        string(name: 'IMAGE_TAG', value: env.BUILD_NUMBER)
                    ],
                    wait: false
            }
        }

        failure {
            echo "CI FAILED: ${env.JOB_NAME} #${env.BUILD_NUMBER}. CD will not be triggered."
        }

        always {
            script {
                if (env.IMAGE) {
                    sh 'docker image rm "${IMAGE}" || true'
                }
            }
        }
    }
}