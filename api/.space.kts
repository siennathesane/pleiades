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
                buf breaking --against "ssh://git@git.jetbrains.space/anthroposlabs/pleiades/Pleiades.git#branch=mainline"
            """
        }
    }
}
