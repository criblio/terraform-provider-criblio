resource "criblio_commit" "my_commit" {
  effective = true
  group     = "default"
  message   = "test"
}

resource "criblio_deploy" "my_deploy" {
  id      = "default"
  version = criblio_commit.my_commit.items[0].commit
}

output "deploy" {
  value = criblio_deploy.my_deploy
}

output "commit" {
  value = criblio_commit.my_commit
}
