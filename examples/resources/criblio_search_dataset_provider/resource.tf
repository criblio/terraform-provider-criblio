resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  prometheus = {
    auth_type       = "basic"
    description     = "my generic provider description"
    endpoint        = "https://prometheus.example.com"
    id              = "myUniqueGenericProviderId"
    max_concurrency = 10
    password        = "prom_pass"
    token           = "prometheusBearerToken123"
    type            = "generic"
    username        = "prom_user"
  }
}