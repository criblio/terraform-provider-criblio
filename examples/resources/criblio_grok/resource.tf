resource "criblio_grok" "my_grok" {
  content  = "%%{IPV4:client} %%{WORD:method} %%{URIPATHPARAM:request}"
  group_id = "Cribl"
  id       = "myGrokFile"
  tags     = "logs,patterns"
}