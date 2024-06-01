pipeline {
    agent any
   
    stages {
        stage ('build and push') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[credentialsId: 'github', url: 'https://github.com/ququiz/quiz-query-service']])
                sh 'chmod 777 ./push.sh'
                sh './push.sh'
                sh 'docker stop quiz-query-service && docker rm quiz-query-service'
                sh 'docker rmi lintangbirdas/quiz-query-service:v1'
            }
        }
        stage ('docker compose up') {
            steps {
                build job: "ququiz-compose", wait: true
            }
        }
    }

}