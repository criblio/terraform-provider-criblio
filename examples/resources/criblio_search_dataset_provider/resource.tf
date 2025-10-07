resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  api_okta_provider = {
    account_configs = [
      {
        api_token       = "00aBcDefGhIjKlMnOpQrStUvWxYz123456"
        domain_endpoint = "https://dev-123456.okta.com"
        name            = "okta-account-1"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
}