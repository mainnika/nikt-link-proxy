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
                binary = docker.build(tag, "--progress=plain --target binary .")
            }
            stage ('Upload') {
                docker.withRegistry("https://" + REMOTE_REGISTRY, REMOTE_SECRET) {
                    binary.push()
                }
            }
            stage ('Deploy') {
                if (env.BRANCH_NAME == DEPLOY_BRANCH) {
                    build job: '../nikt-link-proxy-deploy', wait: false, parameters: [
                              string(name: 'VERSION', value: PROJECT_VERSION),
                              string(name: 'RELEASE_NAME', value: "nikt-link-proxy"),
                              string(name: 'RELEASE_NAMESPACE', value: "nikt-link-proxy"),
                          ]
                }
            }
        }
    } finally {
        sh "docker rmi -f ${tag}        || true"
        sh "docker rm -vf ${binary.id}  || true"
    }
}
