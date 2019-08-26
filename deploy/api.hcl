job "dance" {
    datacenters = ["dc1"]
    type = "service"

    group "api" {
        count = 1

        network {
            mode  = "bridge"
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
                LISTEN_ADDR = "localhost:9090"
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

    group "database" {
        count = 1

        network {
            mode  = "bridge"
        }

        service {
            name = "dance-database"
            port = 5432

            connect {
                sidecar_service {}
            }
        }

        task "postgres" {
            driver = "docker"

            env {
                POSTGRES_DB = "dda"
                POSTGRES_USER = "secret_user"
                POSTGRES_PASSWORD = "secret_password"
            }

            config {
                image = "eveld/da-dance-database"
            }
        }
    }
}