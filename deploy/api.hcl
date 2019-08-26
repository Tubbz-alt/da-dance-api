job "dance-api" {
    datacenters = ["dc1"]
    type = "service"

    group "api" {
        count = 1

        network {
            mode  = "bridge"
            port "http" {
                to = 9090
            }
        }

        service {
            name = "dance-api"
            port = 9090

            connect {
                sidecar_service {
                    proxy {
                        upstreams {
                            destination_name = "dance-database"
                            local_bind_port = 5432
                        }
                    }
                }
            }
        }

        task "server" {
            driver = "docker"

            env {
                LISTEN_ADDR = "0.0.0.0:9090"
                POSTGRES_HOST = "localhost"
                POSTGRES_PORT = 5432
                POSTGRES_USER = "secret_user"
                POSTGRES_PASSWORD = "secret_password"
                POSTGRES_DATABASE = "dda"
            }

            config {
                image = "eveld/da-dance-api"
            }
        }
    }
}