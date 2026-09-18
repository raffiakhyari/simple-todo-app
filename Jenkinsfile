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
                        script: "git log -1 --format=%aN ${env.GIT_COMMIT}",
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
        // 2. Prepare Docker Image
        // ==========================================

        stage('Prepare Image') {

            steps {

                script {

                    if (env.BRANCH_NAME == 'main') {

                        env.DOCKER_NAME = env.DOCKER_PROD

                        env.DEPLOY_ENV = 'production'

                    } else if (env.BRANCH_NAME == 'develop') {

                        env.DOCKER_NAME = env.DOCKER_DEV

                        env.DEPLOY_ENV = 'dev'

                    } else {

                        error(
                            "Unsupported branch for CI/CD: ${env.BRANCH_NAME}"
                        )
                    }

                    env.IMAGE = "${env.DOCKER_NAME}:${env.BUILD_NUMBER}"

                    echo """
                    ==========================================
                    DOCKER IMAGE
                    ==========================================

                    Branch      : ${env.BRANCH_NAME}
                    Environment : ${env.DEPLOY_ENV}
                    Image       : ${env.IMAGE}

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

                        echo ""

                        echo "Please run:"

                        echo "gofmt -w ."

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
        // 7. Push Docker Image
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
                    '''
                }
            }
        }
    }


    // ==========================================
    // Post Actions
    // ==========================================

    post {

        // ==========================================
        // CI SUCCESS → TRIGGER CD
        // ==========================================

        success {

            script {

                echo """
                ==========================================
                PIPELINE SUCCESS
                ==========================================

                Application : todo-api
                Branch      : ${env.BRANCH_NAME}
                Build       : ${env.BUILD_NUMBER}
                Image       : ${env.IMAGE}
                Environment : ${env.DEPLOY_ENV}
                Author      : ${env.AUTHOR_NAME}

                ==========================================
                """

                if (
                    env.BRANCH_NAME == 'develop' ||
                    env.BRANCH_NAME == 'main'
                ) {

                    echo """
                    ==========================================
                    TRIGGERING CD
                    ==========================================

                    Job         : todo-api-delivery
                    Image Repo  : ${env.DOCKER_NAME}
                    Image Tag   : ${env.BUILD_NUMBER}
                    Environment : ${env.DEPLOY_ENV}

                    ==========================================
                    """

                    build job: 'todo-api-delivery',
                        parameters: [

                            string(
                                name: 'IMAGE_REPO',
                                value: env.DOCKER_NAME
                            ),

                            string(
                                name: 'IMAGE_TAG',
                                value: env.BUILD_NUMBER
                            ),

                            string(
                                name: 'ENVIRONMENT',
                                value: env.DEPLOY_ENV
                            )

                        ],
                        wait: false

                } else {

                    echo """
                    ==========================================
                    CD SKIPPED
                    ==========================================

                    Branch ${env.BRANCH_NAME} is not configured
                    for automatic deployment.

                    ==========================================
                    """
                }
            }
        }

        // ==========================================
        // CI FAILURE
        // ==========================================

        failure {

            echo """
            ==========================================
            PIPELINE FAILED
            ==========================================

            Application : todo-api
            Branch      : ${env.BRANCH_NAME}
            Build       : ${env.BUILD_NUMBER}

            CD WILL NOT BE TRIGGERED

            ==========================================
            """
        }

        // ==========================================
        // ALWAYS
        // ==========================================

        always {

            script {

                if (env.IMAGE) {

                    sh """

                        echo "=========================================="
                        echo "Cleaning Local Docker Image"
                        echo "=========================================="

                        docker image rm \
                            "${env.IMAGE}" \
                            || true
                    """
                }
            }
        }
    }
}