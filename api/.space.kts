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
    container(displayName = "buf lint", image = "anthroposlabs.registry.jetbrains.space/p/pleiades/containers/api-ci") {
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
