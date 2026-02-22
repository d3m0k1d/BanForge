# Configs

## config.toml
Main configuration file for BanForge.

Example:
```toml
[firewall]
  name = "nftables"
  config = "/etc/nftables.conf"

[[service]]
  name = "nginx"
  logging = "file"
  log_path = "/home/d3m0k1d/test.log"
  enabled = true

[[service]]
  name = "nginx"
  logging = "journald"
  log_path = "nginx"
  enabled = false
```
**Description**
The [firewall] section defines firewall parameters. The banforge init command automatically detects your installed firewall (nftables, iptables, ufw, firewalld). For firewalls that require a configuration file, specify the path in the config parameter.

The [[service]] section is configured manually. Currently, only nginx is supported. To add a service, create a [[service]] block and specify the log_path to the nginx log file you want to monitor.
logging require in format "file" or "journald"
if you use journald logging, log_path require in format "service_name"

## rules.toml
Rules configuration file for BanForge.

If you wanna configure rules by cli command see [here](https://github.com/d3m0k1d/BanForge/blob/main/docs/cli.md)

Example:
```toml
[[rule]]
  name = "304 http"
  service = "nginx"
  path = ""
  status = "304"
  max_retry = 3
  method = ""
  ban_time = "1m"
```
**Description**
The [[rule]] section require name and one of the following parameters: service, path, status, method. To add a rule, create a [[rule]] block and specify the parameters.
ban_time require in format "1m", "1h", "1d", "1M", "1y".
If you want to ban all requests to PHP files (e.g., path = "*.php") or requests to the admin panel (e.g., path = "/admin/*")
