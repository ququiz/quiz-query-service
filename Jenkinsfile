pipeline {
    agent any
   
    stages {
        stage ('build and push') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[credentialsId: 'github', url: 'https://github.com/ququiz/quiz-query-service']])
                sh 'chmod 777 ./push.sh'
                sh './push.sh'
                    
    
            }
        }
    }

}