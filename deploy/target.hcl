job "dance-target" {
    datacenters = ["dc1"]
    type = "service"

    group "sitting" {
        count = 40

        network {
            mode = "bridge"
        }

        service {
            name = "dance-target"
            port = 6379
        }

        task "duck" {
            driver = "docker"

            config {
                image = "redis"
            }
        }
    }
}