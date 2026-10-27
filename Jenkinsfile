pipeline {
agent any

stages {
    stage('Checkout SCM') {
        steps {
            checkout scm

            sh '''
                echo "=== Repository Information ==="
                git rev-parse --show-toplevel
                git branch --show-current
                git log -1 --format="Commit: %h%nAuthor: %aN%nMessage: %s" HEAD
            '''
        }
    }
}

post {
    success {
        echo 'Pipeline completed successfully.'
    }
    failure {
        echo 'Pipeline failed. Check the console output.'
    }
}

}