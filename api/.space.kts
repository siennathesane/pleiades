// only run when the image has changed or on Sunday to ensure the image gets updates
job("Build CI Image") {
    startOn {
        gitPush {
            pathFilter {
                +"Dockerfile"
            }
        }
        schedule { cron("0 0 * * 0") }
    }
    kaniko {
        build {
            context = "docker"
            dockerfile = "Dockerfile"
        }
        push("anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
            tags {
                +"latest"
            }
        }
    }
}

job("Mirror") {
    git {
        refSpec {
            +"refs/*:refs/*"
        }
        depth = UNLIMITED_DEPTH
    }
    container(displayName = "mirror", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
        env["DEPLOY_PUBLIC_KEY"] = Secrets("pleiades-deploy-public-key")
        env["DEPLOY_PRIVATE_KEY"] = Secrets("pleiades-deploy-private-key")
        shellScript {
            location = "./ci/mirror.sh"
        }
    }
}

// run on any push
job("Lint & Build") {
    container(displayName = "buf lint", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
        shellScript {
            location = "./ci/lint.sh"
        }
    }
}

// only run on mainline
job("Push to Registry") {
    startOn {
        codeReviewClosed{}
    }
    container(displayName = "buf lint", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
        env["BUF_USER"] = "sienna-al"
        env["BUF_API_TOKEN"] = Secrets("buf-api-token")
        shellScript {
            location = "./ci/push.sh"
        }
    }
}
