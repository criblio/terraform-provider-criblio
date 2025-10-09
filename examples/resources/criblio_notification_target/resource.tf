resource "criblio_notification_target" "my_notificationtarget" {
  slack_target = {
    id = "mySlackTarget"
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    type = "slack"
    url  = "https://hooks.slack.com/services/THISisanExample"
  }
}