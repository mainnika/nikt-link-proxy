node('docker') {
    def scmVars = checkout scm

    def REMOTE_REGISTRY = "do.cr.tokarch.uk"
    def REMOTE_SECRET = "docr-mainnika"

    def DEPLOY_BRANCH = "main"

    def PROJECT_NAME = env.JOB_NAME.minus("/" + env.BRANCH_NAME)
    def PROJECT_VERSION = scmVars.GIT_COMMIT + ".${env.BUILD_NUMBER}"

    def tag = "${REMOTE_REGISTRY}/${PROJECT_NAME}:${PROJECT_VERSION}"
    def binary

    try {
        withEnv([
            "DOCKER_BUILDKIT=1",
        ]) {
            stage ('Build') {
                binary = docker.build(tag, "--target binary .")
            }
            stage ('Upload') {
                docker.withRegistry("https://" + REMOTE_REGISTRY, REMOTE_SECRET) {
                    binary.push()
                }
            }
        }
    } finally {
        sh "docker rmi -f ${tag}        || true"
        sh "docker rm -vf ${binary.id}  || true"
    }
}
