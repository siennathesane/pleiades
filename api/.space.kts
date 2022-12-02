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
    git("API"){}

    container(displayName = "buf lint", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
        env["BUF_INPUT_HTTPS_USERNAME"] = Secrets("git_clone_https_user")
        env["BUF_INPUT_HTTPS_PASSWORD"] = Secrets("git_clone_https_user")
        shellScript {
            interpreter = "/bin/bash"
            content = """
                buf lint

                # kvstore
                pushd databaseapi
                buf breaking --against buf.build/anthropos-labs/kvstore
                popd

                # errors
                pushd errorapi
                buf breaking --against buf.build/anthropos-labs/errors
                popd

                # raft
                pushd raftapi
                buf breaking --against buf.build/anthropos-labs/raft
                popd
            """
        }
    }
}
