#!/bin/bash

allEnvs=(
	'criblio_regex.my_regex/{"group_id": "default", "id": "test_regex_2"}'
	'criblio_database_connection.my_databaseconnection/{"group_id": "default", "id": "my_databaseconnection"}'
	'criblio_hmac_function.my_hmacfunction/{"group_id": "default", "id": "my_hmacfunction"}'
	'criblio_parquet_schema.my_parquet_schema/{"group_id": "default", "id": "my_parquet_schema"}'
	'criblio_schema.my_schema/{"group_id": "default", "id": "my_schema"}'
	'criblio_global_var.my_globalvar/{"group_id": "default", "id": "sample_globalvar"}'
	'criblio_appscope_config.my_appscopeconfig/{"group_id": "default", "id": "sample_appscope_config"}'
	'criblio_pack_vars.my_packvars/{"group_id": "default", "pack": "pack-with-vars"}'
	'criblio_pack_lookups.my_packlookups/{"group_id": "default", "pack": "pack-with-lookups", "id": "my_id"}'
	'criblio_pack.my_pack/{"group_id": "default", "id": "pack-from-source"}'
	'criblio_pack.breakers_pack/{"group_id": "default", "id": "pack-breakers"}'
	'criblio_pack.dest_pack/{"group_id": "default", "id": "pack-with-dest"}'
	'criblio_pack.lookups_pack/{"group_id": "default", "id": "pack-with-lookups"}'
	'criblio_pack.pipeline_pack/{"group_id": "default", "id": "pack-with-pipeline"}'
	'criblio_pack.routes_pack/{"group_id": "default", "id": "pack-with-routes"}'
	'criblio_pack.source_pack/{"group_id": "default", "id": "pack-with-source"}'
	'criblio_pack.vars_pack/{"group_id": "default", "id": "pack-with-vars"}'
	'criblio_pack_breakers.my_packbreakers/{"group_id": "default", "pack": "pack-breakers", "id": "test_packbreakers"}'
	'criblio_pack_source.my_packsource/{"group_id": "default", "pack": "pack-with-source", "id": "my_id"}'
	'criblio_pack_destination.my_packdest/{"group_id": "default", "pack": "pack-with-dest", "id": "test"}'
	'criblio_pack_routes.my_packroutes/{"group_id": "default", "pack": "pack-with-routes"}'
	'criblio_pack_pipeline.my_packpipeline/{"group_id": "default", "pack": "pack-with-pipeline"}'
	'criblio_event_breaker_ruleset.my_eventbreakerruleset/{"group_id": "default", "id": "test_eventbreakerruleset"}'
)

cloud=(
	'criblio_subscription.my_subscription[0]/{"group_id": "default", "id": "my_subscription"}'
	'criblio_subscription.my_subscription_with_enabled[0]/{"group_id": "default", "id": "my_subscription_with_enabled"}'
	'criblio_project.my_project[0]/{"group_id": "default", "id": "my_project"}'
	'criblio_grok.my_grok[0]/{"group_id": "default", "id": "test_grok"}'
	'criblio_group.syslog_worker_group[0]/"syslog-workers'
	'criblio_pack.syslog_pack[0]/{"group_id": "syslog-workers", "id": "syslog-processing"}'
	'criblio_destination.cribl_lake[0]/{"group_id": "syslog-workers", "id": "cribl-lake-2"}'
	'criblio_source.syslog_source[0]/{"group_id": "syslog-workers", "id": "syslog-input"}'
)

shipIt() {
	for key in "$@"; do
		resource=${key%%/*}
		id=${key#*/}
		terraform import -no-color "$resource" "$id"
		tC=$?
		if [[ $tC -ne 0 ]]; then
			rc=$tC
		fi
		sleep $sleepTime
		echo
	done
}

rc=0
onprem=false
sleepTime=0
if [[ -z $CRIBL_CLOUD_DOMAIN ]]; then
	onprem=true
	sleepTime=1
fi

shipIt "${allEnvs[@]}"

if [[ $onprem ]]; then
	echo "WARNING: Skipping resources since we are testing OnPrem"
else
	shipIt "${cloud[@]}"
fi

exit $rc
