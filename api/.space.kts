job("build-ci-image") {
    startOn {
        gitPush {
            pathFilter {
                +"Dockerfile"
            }
        }
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

job("lint") {
    git {
        refSpec {
            +"refs/heads/mainline"
        }
    }
    container(displayName = "buf lint", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
        shellScript {
            interpreter = "/bin/bash"
            content = """
                buf lint
                export BUF_INPUT_HTTPS_USERNAME=${'$'}JB_SPACE_CLIENT_ID
                export BUF_INPUT_HTTPS_PASSWORD=${'$'}JB_SPACE_CLIENT_SECRET
                buf breaking --against "https://git.jetbrains.space/anthroposlabs/pleiades/Pleiades.git#branch=mainline"
            """
        }
    }
}
