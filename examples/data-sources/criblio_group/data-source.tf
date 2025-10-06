data "criblio_group" "my_group" {
  fields = "git.commit,git.localChanges"
  id     = "workers-west"
}