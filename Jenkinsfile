pipeline {
  agent any
  stages {
    stage('Build') {
      parallel {
        stage('Build 1') {
          steps {
            sh 'echo Build 1'
          }
        }
        stage('Build 2') {
          steps {
            sh 'echo Build 2'
          }
        }
      }
    }
    stage('Deploy') {
      steps {
        sh 'echo Deploy'
      }
    }
  }
}