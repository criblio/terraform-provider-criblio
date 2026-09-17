data "criblio_group" "existing" {
  id     = "default"
  fields = "git.commit,git.localChanges,git.log"
}

output "group_git" {
  value = data.criblio_group.existing.git
}

output "current_commit" {
  value = try(data.criblio_group.existing.git.commit, null)
}

output "local_change_count" {
  value = try(data.criblio_group.existing.git.local_changes, null)
}

output "recent_commits" {
  value = try(data.criblio_group.existing.git.log, [])
}