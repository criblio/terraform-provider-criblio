resource "criblio_appscope_config" "my_appscopeconfig" {
  config = {
    cribl = {
      authtoken = "myAuthToken123"
      enable    = true
      transport = {
        buffer = "line"
        host   = "localhost"
        path   = "/var/run/appscope.sock"
        port   = 8080
        tls = {
          cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
          enable         = true
          validateserver = false
        }
        type = "tcp"
      }
      use_scope_source_transport = false
    }
    custom = [
      {
        ancestor = "parentProcess"
        arg      = "--debug"
        config = {
          cribl = {
            authtoken = "myAuthToken123"
            enable    = true
            transport = {
              buffer = "line"
              host   = "localhost"
              path   = "/var/run/appscope.sock"
              port   = 8080
              tls = {
                cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
                enable         = true
                validateserver = false
              }
              type = "tcp"
            }
            use_scope_source_transport = false
          }
          event = {
            enable = true
            format = {
              enhancefs      = false
              maxeventpersec = 100
            }
            transport = {
              buffer = "line"
              host   = "localhost"
              path   = "/var/run/appscope.sock"
              port   = 8080
              tls = {
                cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
                enable         = true
                validateserver = false
              }
              type = "tcp"
            }
            type = "ndjson"
            watch = [
              {
                allowbinary = false
                enabled     = true
                field       = "http.method"
                headers = [
                  "Content-Type"
                ]
                name  = "RequestEvents"
                type  = "match"
                value = "GET"
              }
            ]
          }
          libscope = {
            commanddir  = "/opt/appscope"
            configevent = false
            log = {
              level = "info"
              transport = {
                buffer = "line"
                host   = "localhost"
                path   = "/var/run/appscope.sock"
                port   = 8080
                tls = {
                  cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
                  enable         = true
                  validateserver = false
                }
                type = "tcp"
              }
            }
            summaryperiod = 60
          }
          metric = {
            enable = true
            format = {
              statsdmaxlen = 512
              statsdprefix = "webshop.prod."
              type         = "statsd"
              verbosity    = 3
            }
            transport = {
              buffer = "line"
              host   = "localhost"
              path   = "/var/run/appscope.sock"
              port   = 8080
              tls = {
                cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
                enable         = true
                validateserver = false
              }
              type = "tcp"
            }
            watch = [
              {
                # ...
              }
            ]
          }
          payload = {
            dir    = "/var/lib/appscope/payloads"
            enable = true
          }
          protocol = [
            {
              binary  = false
              detect  = true
              len     = 128
              name    = "http"
              payload = true
              regex   = ".*"
            }
          ]
          tags = [
            {
              key   = "env"
              value = "prod"
            }
          ]
        }
        env      = "production"
        hostname = "host123.example.com"
        procname = "myprocess"
        username = "appuser"
      }
    ]
    event = {
      enable = true
      format = {
        enhancefs      = false
        maxeventpersec = 100
      }
      transport = {
        buffer = "line"
        host   = "localhost"
        path   = "/var/run/appscope.sock"
        port   = 8080
        tls = {
          cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
          enable         = true
          validateserver = false
        }
        type = "tcp"
      }
      type = "ndjson"
      watch = [
        {
          allowbinary = false
          enabled     = true
          field       = "http.method"
          headers = [
            "Content-Type"
          ]
          name  = "RequestEvents"
          type  = "match"
          value = "GET"
        }
      ]
    }
    libscope = {
      commanddir  = "/var/run/appscope/commands"
      configevent = true
      log = {
        level = "info"
        transport = {
          buffer = "line"
          host   = "localhost"
          path   = "/var/run/appscope.sock"
          port   = 8080
          tls = {
            cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
            enable         = true
            validateserver = false
          }
          type = "tcp"
        }
      }
      summaryperiod = 60
    }
    metric = {
      enable = true
      format = {
        statsdmaxlen = 512
        statsdprefix = "webshop.prod."
        type         = "statsd"
        verbosity    = 3
      }
      transport = {
        buffer = "line"
        host   = "localhost"
        path   = "/var/run/appscope.sock"
        port   = 8080
        tls = {
          cacertpath     = "/etc/ssl/certs/ca-certificates.crt"
          enable         = true
          validateserver = false
        }
        type = "tcp"
      }
      watch = [
        {
          # ...
        }
      ]
    }
    payload = {
      dir    = "/var/lib/appscope/payloads"
      enable = true
    }
    protocol = [
      {
        binary  = false
        detect  = true
        len     = 128
        name    = "http"
        payload = true
        regex   = ".*"
      }
    ]
    tags = [
      {
        key   = "env"
        value = "prod"
      }
    ]
  }
  description = "Custom Appscope configuration for nginx"
  group_id    = "Cribl"
  id          = "appscopeConfig1"
  lib         = "cribl"
  tags        = "scope,nginx"
}