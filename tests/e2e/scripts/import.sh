#!/bin/bash

terraform import -no-color criblio_group.my_edge_fleet "my-edge-fleet" 
echo
terraform import -no-color criblio_group.syslog_worker_group "syslog-workers" 
echo
terraform import -no-color criblio_subscription.my_subscription '{"group_id": "default", "id": "my_subscription"}'
echo
terraform import -no-color criblio_subscription.my_subscription_with_enabled '{"group_id": "default", "id": "my_subscription_with_enabled"}'
echo
terraform import -no-color criblio_regex.my_regex '{"group_id": "default", "id": "test_regex_2"}'
echo
terraform import -no-color criblio_database_connection.my_databaseconnection '{"group_id": "default", "id": "my_databaseconnection"}'
echo
terraform import -no-color criblio_hmac_function.my_hmacfunction '{"group_id": "default", "id": "my_hmacfunction"}'
echo
terraform import -no-color criblio_parquet_schema.my_parquet_schema '{"group_id": "default", "id": "my_parquet_schema"}'
echo
terraform import -no-color criblio_project.my_project '{"group_id": "default", "id": "my_project"}'
echo
terraform import -no-color criblio_schema.my_schema '{"group_id": "default", "id": "my_schema"}'
echo
terraform import -no-color criblio_global_var.my_globalvar '{"group_id": "default", "id": "sample_globalvar"}'
echo
terraform import -no-color criblio_grok.my_grok '{"group_id": "default", "id": "test_grok"}'
echo
terraform import -no-color criblio_appscope_config.my_appscopeconfig '{"group_id": "default", "id": "sample_appscope_config"}'
echo
terraform import -no-color criblio_pack_vars.my_packvars '{"group_id": "default", "pack": "pack-with-vars"}'
echo
terraform import -no-color criblio_pack_lookups.my_packlookups '{"group_id": "default", "pack": "pack-with-lookups", "id": "my_id"}'
echo
terraform import -no-color criblio_pack.my_pack '{"group_id": "default", "id": "pack-from-source"}'
echo
terraform import -no-color criblio_pack.breakers_pack '{"group_id": "default", "id": "pack-breakers"}'
echo
terraform import -no-color criblio_pack.dest_pack '{"group_id": "default", "id": "pack-with-dest"}'
echo
terraform import -no-color criblio_pack.lookups_pack '{"group_id": "default", "id": "pack-with-lookups"}'
echo
terraform import -no-color criblio_pack.pipeline_pack '{"group_id": "default", "id": "pack-with-pipeline"}'
echo
terraform import -no-color criblio_pack.routes_pack '{"group_id": "default", "id": "pack-with-routes"}'
echo
terraform import -no-color criblio_pack.source_pack '{"group_id": "default", "id": "pack-with-source"}'
echo
terraform import -no-color criblio_pack.vars_pack '{"group_id": "default", "id": "pack-with-vars"}'
echo
terraform import -no-color criblio_pack.syslog_pack '{"group_id": "syslog-workers", "id": "syslog-processing"}'
echo
terraform import -no-color criblio_pack_breakers.my_packbreakers '{"group_id": "default", "pack": "pack-breakers", "id": "test_packbreakers"}'
echo
terraform import -no-color criblio_destination.cribl_lake '{"group_id": "syslog-workers", "id": "cribl-lake-2"}'
echo
terraform import -no-color criblio_source.syslog_source '{"group_id": "syslog-workers", "id": "syslog-input"}'
echo
terraform import -no-color criblio_event_breaker_ruleset.my_eventbreakerruleset '{"group_id": "default", "id": "test_eventbreakerruleset"}'
echo
terraform import -no-color criblio_pack_source.my_packsource '{"group_id": "default", "pack": "pack-with-source", "id": "my_id"}'
echo
terraform import -no-color criblio_pack_destination.my_packdest '{"group_id": "default", "pack": "pack-with-dest", "id": "test"}' 
echo
terraform import -no-color criblio_pack_routes.my_packroutes '{"group_id": "default", "pack": "pack-with-routes"}' 
echo
terraform import -no-color criblio_pack_pipeline.my_packpipeline '{"group_id": "default", "pack": "pack-with-pipeline"}'
echo
