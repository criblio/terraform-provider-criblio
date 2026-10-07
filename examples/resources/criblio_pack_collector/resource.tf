resource "criblio_pack_collector" "rest" {
  group_id = "default"
  pack     = criblio_pack.source_pack.id
  id       = "rest-api"

  input_collector_rest = {
    collector = {
      type = "rest"
      conf = {
        collect_url    = "'https://api.example.com/events'"
        collect_method = "get"
        authentication = "none"
      }
    }
    input = {
      type           = "collection"
      send_to_routes = true
    }
    ttl = "4h"
  }
}

resource "criblio_pack" "source_pack" {
  id           = "pack-with-collector"
  group_id     = "default"
  description  = "Pack with Collector"
  disabled     = true
  display_name = "Pack from Collector"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}
