def checkoutConfigMap = [
    branchName: BRANCH_NAME,
    credentialsId: 'gitlab',
    repoUrl: 'https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev'
]

def dockerConfigMap = [
    imageName:"malinkamedok/currency_app",
    tag:"latest",
    pathToDockerfile:"./Dockerfile",
    registryUrl: "docker.io",
]

@Library("shared_library") _
node ('swarm') {
    try {
        stage("Checkout repo") {
        checkoutRepo(checkoutConfigMap)
        }
        stage('Build docker image') {
            buildDockerImage(dockerConfigMap)
        }
        stage('Push docker image to docker hub') {
            pushDockerImage(dockerConfigMap)
        }
    } finally {
        cleanWs(notFailBuild: true)
    }
}