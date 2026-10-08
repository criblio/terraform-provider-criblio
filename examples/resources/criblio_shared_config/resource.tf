resource "criblio_shared_config" "example" {
  id   = "shared-firing-defaults"
  name = "Shared Firing Defaults"

  # Partial monitor configuration inherited by monitors that reference it.
  # Opaque JSON (Partial<AetosMonitorConf>); fingerprint is computed by the server.
  config = jsonencode({
    firingCondition = {
      fireDelay  = 300
      clearDelay = 60
    }
  })
}
