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
node {
    try {
        stage('Checkout repo') {
            checkoutRepo(checkoutConfigMap)
        }
        stage('Run Go tests') {
            def root = tool type: 'go', name: '1.22.2'
            withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                sh 'go install github.com/jstemmer/go-junit-report/v2@latest'
                sh 'go test -v 2>&1 ./repo/test | go-junit-report -set-exit-code > report.xml'
            }
        }
        stage('Build docker image') {
            buildDockerImage(dockerConfigMap)
        }
        stage('Push docker image to docker hub') {
            pushDockerImage(dockerConfigMap)
        }
    } finally {
        junit 'report.xml'
        cleanWs(notFailBuild: true)
    }
}