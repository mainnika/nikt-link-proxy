node('k8s') {

    def HELM_DOWNLOAD_URL = "https://get.helm.sh/helm-v3.7.2-linux-amd64.tar.gz"
    def WS_DIST = ".dist"
    def WS_BIN = ".bin"
    def WS_SCM = ".scm"

    parameters {
        string(
            name: "VERSION",
            defaultValue: "main",
            description: "version to release",
        )
        string(
            name: "RELEASE_NAME",
            defaultValue: "nikt-link-proxy",
            description: "release name",
        )
        string(
            name: "RELEASE_NAMESPACE",
            defaultValue: "nikt-link-proxy",
            description: "release namespace",
        )
        password(
            name: "CREDENTIALS_ID",
            defaultValue: "nikt-link-proxy-config",
            description: "release configuration",
        )
    }

    try {
        stage ('Download Helm') {
            def helmExists = fileExists("${WS_BIN}/linux-amd64/helm")
            if (!helmExists) {
                dir("${WS_DIST}") {
                    sh """
                        curl -fsSL -O "${HELM_DOWNLOAD_URL}";
                        mkdir -p "${env.WORKSPACE}/${WS_BIN}"
                        tar -C "${env.WORKSPACE}/${WS_BIN}" -xzf *;
                        ls -lah "${env.WORKSPACE}/${WS_BIN}/linux-amd64/helm";
                    """
                    deleteDir()
                }
            }
        }
        stage ('Checkout Git Repo') {
            dir("${WS_SCM}") {
                deleteDir()
                checkout scm
            }
        }
        stage("Make chart package") {
            dir("${WS_SCM}/chart.tgz") {
                sh """
                    "${env.WORKSPACE}/${WS_BIN}/linux-amd64/helm" package ../chart \
                      --app-version ${params.VERSION} \
                      ;
                """
            }
        }
        stage("Release packaged chart") {
            withCredentials([file(credentialsId: "${params.CREDENTIALS_ID}", variable: "RELEASE_CONFIG")]) {
                dir("${WS_SCM}/chart.tgz") {
                    sh """
                        "${env.WORKSPACE}/${WS_BIN}/linux-amd64/helm" upgrade ${params.RELEASE_NAME} * \
                          --namespace ${params.RELEASE_NAMESPACE} \
                          --create-namespace \
                          --dependency-update \
                          --reset-values \
                          -f "\$RELEASE_CONFIG" \
                          ;
                    """
                }
            }
        }
    } finally {
    }
}
