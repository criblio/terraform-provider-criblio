data "criblio_config_version" "my_configversion" {
  id = "syslog-workers"
  depends_on = [criblio_commit.my_commit]
}
