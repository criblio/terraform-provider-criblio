resource "criblio_notification" "my_notification" {
  condition = "...my_condition..."
  conf = {
    # ...
  }
  disabled = false
  group    = "...my_group..."
  id       = "...my_id..."
  metadata = [
    {
      name  = "...my_name..."
      value = "...my_value..."
    }
  ]
  target_configs = [
    {
      conf = {
        body = "...my_body..."
        email_recipient = {
          bcc = "Marquise55@gmail.com"
          cc  = "Adonis56@gmail.com"
          to  = "Tate.Flatley@yahoo.com"
        }
        subject = "...my_subject..."
      }
      id = "...my_id..."
    }
  ]
  targets = [
    "..."
  ]
}