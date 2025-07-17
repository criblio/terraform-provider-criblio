import {
    to = criblio_appscope_config.my_appscopeconfig
    id = "sample_appscope_config"
}

import {
    to = criblio_global_var.my_globalvar
    id = "sample_globalvar"
}

import {
    to = criblio_pack.my_pack
    id = "pack-from-source"
}

import {
    to = criblio_schema.my_schema
    id = "my_schema"
}

import {
    to = criblio_grok.my_grok
    id = "my_grok"
}

import {
    to = criblio_parquet_schema.my_parquet_schema
    id = "my_parquet_schema"
}

import {
    to = criblio_source.syslog_source
    id = "syslog-input"
}

import {
    to = criblio_destination.cribl_lake
    id = "cribl-lake-2"
}

import {
    to = criblio_pack.syslog_pack
    id = "syslog-processing"
}

import {
    to = criblio_database_connection.my_databaseconnection
    id = "my_databaseconnection"
}

import {
    to = criblio_hmac_function.my_hmacfunction
    id = "my_hmacfunction"
}

import {
    to = criblio_project.my_project
    id = "my_project"
}

import {
    to = criblio_subscription.my_subscription
    id = "my_subscription"
}

import {
    to = criblio_subscription.my_subscription_with_enabled
    id = "my_subscription_with_enabled"
}

import {
    to = criblio_event_breaker_ruleset.my_eventbreakerruleset
    id = "test_eventbreakerruleset"
}

data "criblio_appscope_config" "my_appscopeconfig" {
  group_id   = "default"
}

data "criblio_global_var" "my_globalvar" {
  group_id   = "default"
}

data "criblio_pack" "my_pack" {
  group_id   = "default"
}

data "criblio_schema" "my_schema" {
  group_id   = "default"
}

data "criblio_grok" "my_grok" {
  group_id   = "default"
}

data "criblio_parquet_schema" "criblio_parquet_schemas" {
  group_id   = "default"
}

data "criblio_source" "syslog_source" {
  group_id   = "syslog-workers"
}

data "criblio_destination" "cribl_lake" {
  group_id = "default"
}

data "criblio_pack" "syslog_pack" {
  group_id = "syslog_workers"
}

data "criblio_database_connection" "my_databaseconnection" {
  group_id = "default"
}

data "criblio_hmac_function" "my_hmacfunction" {
  group_id = "default"
}

data "criblio_project" "my_project" {
  group_id = "default"
}

data "criblio_subscription" "my_subscription" {
  group_id = "default"
}

data "criblio_subscription" "my_subscription_with_enabled" {
  group_id = "default"
}

data "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  group_id = "default"
}
