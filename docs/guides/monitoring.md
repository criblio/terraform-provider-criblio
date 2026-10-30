---
page_title: "Monitoring resources"
subcategory: ""
description: |-
  Manage customer-metrics monitors, muting rules, notifications, and presets.
---

# Monitoring resources

| UI configuration | Terraform resource | REST collection |
| --- | --- | --- |
| Monitors | `criblio_monitor` | `/products/lakehouse_engine_metrics/monitors` |
| Muting Rules | `criblio_monitor_silence` | `/a/__cribl_metrics/products/lakehouse_engine_metrics/silences` |
| Notification Targets | `criblio_notification_target` | `/notification-targets` |
| Notification Templates | `criblio_notification_template` | `/notification-templates?engine=handlebars` |
| Notification Policies | `criblio_notification_policy` | `/notification-policies` |
| Shared Configurations | `criblio_monitor_shared_config` | `/products/lakehouse_engine_metrics/shared-configs` |
| Profiles | `criblio_monitor_profile` | `/products/lakehouse_engine_metrics/config-profiles` |

Each configuration resource supports import by ID. The four new resource types also have single-object and list data sources. Their documentation is generated from the provider schemas.

Silence CRUD, data sources, and import/export use the metrics app prefix `/a/__cribl_metrics`, matching the Muting Rules UI request. Cribl retains the app authorization context when rewriting this URL to its customer-metrics handler.

Active Alerts and preset Activity are backend-generated state and audit history. They do not have a configuration lifecycle. `criblio_monitor_active_alerts` reads current monitor instance states from `/products/lakehouse_engine_metrics/engines/{engineId}/datasets/{datasetId}/alert_activity/current`, including confirmed and pending statuses. Audit history remains available through each preset's `/history` endpoint in Cribl.

Policy conditions, shared configurations, and muting-rule tags use `jsonencode(...)` to preserve their complete upstream shapes. Muting-rule times are Unix milliseconds; creation time and creator are computed. Shared-configuration fingerprints are also computed.

Profile changes use replacement: the current backend PATCH implementation does not rebuild its hierarchy metadata. Unlink referencing monitors and child profiles before replacing or deleting a profile. Cribl also prevents deleting shared configurations still attached to profiles. Express dependencies through resource references to preserve creation and destruction order.

The muting-rule and notification-policy collection routes exist in the Cribl repository but are absent from the published upstream spec. The Terraform overlay adds their route definitions and references existing upstream schemas. Other resource routes are annotated directly from upstream.

Live acceptance tests require customer metrics and V2 Notifications to be enabled. The existing on-prem fixture does not enable these features, so monitoring tests skip there. Current-state acceptance tests additionally require `MONITOR_ENGINE_ID` and `MONITOR_DATASET_ID`. Offline tests validate all new examples and lifecycle configurations without credentials.
