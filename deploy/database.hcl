job "dance-database" {
    datacenters = ["dc1"]
    type = "service"
    
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
                image = "eveld/da-dance-database:v0.1.0"
            }
        }
    }
}