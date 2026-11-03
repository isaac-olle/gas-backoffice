provider "postgresql" {
  host            = data.digitalocean_database_cluster.postgres.host
  port            = data.digitalocean_database_cluster.postgres.port
  database        = "defaultdb"
  username        = data.digitalocean_database_cluster.postgres.user
  password        = data.digitalocean_database_cluster.postgres.password
  sslmode         = "require"
  connect_timeout = 15
}

