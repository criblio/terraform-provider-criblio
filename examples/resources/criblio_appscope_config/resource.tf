resource "criblio_appscope_config" "my_appscopeconfig" {
  config = {
    cribl = {
      authtoken = "...my_authtoken..."
      enable    = true
      transport = {
        buffer = "line"
        host   = "...my_host..."
        path   = "...my_path..."
        port   = 1.73
        tls = {
          cacertpath     = "...my_cacertpath..."
          enable         = false
          validateserver = true
        }
        type = "...my_type..."
      }
      use_scope_source_transport = true
    }
    custom = [
      {
        ancestor = "...my_ancestor..."
        arg      = "...my_arg..."
        config = {
          cribl = {
            authtoken = "...my_authtoken..."
            enable    = false
            transport = {
              buffer = "full"
              host   = "...my_host..."
              path   = "...my_path..."
              port   = 9.85
              tls = {
                cacertpath     = "...my_cacertpath..."
                enable         = false
                validateserver = false
              }
              type = "...my_type..."
            }
            use_scope_source_transport = true
          }
          event = {
            enable = true
            format = {
              enhancefs      = true
              maxeventpersec = 5.37
            }
            transport = {
              buffer = "line"
              host   = "...my_host..."
              path   = "...my_path..."
              port   = 4.5
              tls = {
                cacertpath     = "...my_cacertpath..."
                enable         = true
                validateserver = true
              }
              type = "...my_type..."
            }
            type = "ndjson"
            watch = [
              {
                allowbinary = true
                enabled     = true
                field       = "...my_field..."
                headers = [
                  "..."
                ]
                name  = "...my_name..."
                type  = "...my_type..."
                value = "...my_value..."
              }
            ]
          }
          libscope = {
            commanddir  = "...my_commanddir..."
            configevent = false
            log = {
              level = "info"
              transport = {
                buffer = "line"
                host   = "...my_host..."
                path   = "...my_path..."
                port   = 2.31
                tls = {
                  cacertpath     = "...my_cacertpath..."
                  enable         = false
                  validateserver = false
                }
                type = "...my_type..."
              }
            }
            summaryperiod = 9.88
          }
          metric = {
            enable = true
            format = {
              statsdmaxlen = 8.61
              statsdprefix = "...my_statsdprefix..."
              type         = "...my_type..."
              verbosity    = 4.06
            }
            transport = {
              buffer = "line"
              host   = "...my_host..."
              path   = "...my_path..."
              port   = 3.15
              tls = {
                cacertpath     = "...my_cacertpath..."
                enable         = false
                validateserver = true
              }
              type = "...my_type..."
            }
            watch = [
              {
                # ...
              }
            ]
          }
          payload = {
            dir    = "...my_dir..."
            enable = false
          }
          protocol = [
            {
              binary  = false
              detect  = false
              len     = 1.8
              name    = "...my_name..."
              payload = true
              regex   = "...my_regex..."
            }
          ]
          tags = [
            {
              key   = "...my_key..."
              value = "...my_value..."
            }
          ]
        }
        env      = "...my_env..."
        hostname = "...my_hostname..."
        procname = "...my_procname..."
        username = "...my_username..."
      }
    ]
    event = {
      enable = false
      format = {
        enhancefs      = true
        maxeventpersec = 3.37
      }
      transport = {
        buffer = "line"
        host   = "...my_host..."
        path   = "...my_path..."
        port   = 6.04
        tls = {
          cacertpath     = "...my_cacertpath..."
          enable         = true
          validateserver = true
        }
        type = "...my_type..."
      }
      type = "ndjson"
      watch = [
        {
          allowbinary = true
          enabled     = true
          field       = "...my_field..."
          headers = [
            "..."
          ]
          name  = "...my_name..."
          type  = "...my_type..."
          value = "...my_value..."
        }
      ]
    }
    libscope = {
      commanddir  = "...my_commanddir..."
      configevent = true
      log = {
        level = "none"
        transport = {
          buffer = "line"
          host   = "...my_host..."
          path   = "...my_path..."
          port   = 9.86
          tls = {
            cacertpath     = "...my_cacertpath..."
            enable         = true
            validateserver = true
          }
          type = "...my_type..."
        }
      }
      summaryperiod = 9.16
    }
    metric = {
      enable = false
      format = {
        statsdmaxlen = 5.96
        statsdprefix = "...my_statsdprefix..."
        type         = "...my_type..."
        verbosity    = 1.33
      }
      transport = {
        buffer = "full"
        host   = "...my_host..."
        path   = "...my_path..."
        port   = 5.7
        tls = {
          cacertpath     = "...my_cacertpath..."
          enable         = true
          validateserver = true
        }
        type = "...my_type..."
      }
      watch = [
        {
          # ...
        }
      ]
    }
    payload = {
      dir    = "...my_dir..."
      enable = true
    }
    protocol = [
      {
        binary  = true
        detect  = true
        len     = 9.26
        name    = "...my_name..."
        payload = true
        regex   = "...my_regex..."
      }
    ]
    tags = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
  }
  description = "My Appscope Lib Entry description"
  group_id    = "myExistingGroupId"
  id          = "myAppscopeId"
  lib         = "cribl"
  tags        = "myTags"
}