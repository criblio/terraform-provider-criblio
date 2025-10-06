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
            config = {
              enable = true
              format = {
                level   = "info"
                maxline = 1024
              }
              log = {
                level = "debug"
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
          }
          metric = {
            enable       = true
            format       = "statsd"
            statsdmaxlen = 512
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
          payload = {
            dir    = "/var/log/appscope/payloads"
            enable = false
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
      commanddir = "path/to/dir"
      config = {
        enable = true
        format = {
          level   = "info"
          maxline = 1024
        }
        log = {
          level = "debug"
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
      configevent = true
    }
    metric = {
      enable       = true
      format       = "statsd"
      statsdmaxlen = 512
      statsdprefix = "myPrefix-"
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
      verbosity = 0
    }
    payload = {
      dir    = "/var/log/appscope/payloads"
      enable = false
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