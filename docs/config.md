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

  # Actions are executed after successful ban
  [[rule.action]]
    type = "email"
    enabled = true
    email = "admin@example.com"
    email_sender = "banforge@example.com"
    email_subject = "BanForge Alert: IP Banned"
    smtp_host = "smtp.example.com"
    smtp_port = 587
    smtp_user = "user@example.com"
    smtp_password = "password"
    smtp_tls = true
    body = "IP {ip} has been banned for rule {rule}"

  [[rule.action]]
    type = "webhook"
    enabled = true
    url = "https://hooks.example.com/alert"
    method = "POST"
    headers = { "Content-Type" = "application/json", "Authorization" = "Bearer token" }
    body = "{\"ip\": \"{ip}\", \"rule\": \"{rule}\", \"service\": \"{service}\"}"

  [[rule.action]]
    type = "script"
    enabled = true
    script = "/usr/local/bin/notify.sh"
    interpretator = "bash"
```
**Description**
The [[rule]] section require name and one of the following parameters: service, path, status, method. To add a rule, create a [[rule]] block and specify the parameters.
ban_time require in format "1m", "1h", "1d", "1M", "1y".
If you want to ban all requests to PHP files (e.g., path = "*.php") or requests to the admin panel (e.g., path = "/admin/*").
If max_retry = 0 ban on first request.

## Actions

Actions are executed after a successful IP ban. You can configure multiple actions per rule.

### Action Types

#### 1. Email Notification

Send email alerts when an IP is banned.

```toml
[[rule.action]]
  type = "email"
  enabled = true
  email = "admin@example.com"
  email_sender = "banforge@example.com"
  email_subject = "BanForge Alert"
  smtp_host = "smtp.example.com"
  smtp_port = 587
  smtp_user = "user@example.com"
  smtp_password = "password"
  smtp_tls = true
  body = "IP {ip} has been banned"
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | + | Must be "email" |
| `enabled` | + | Enable/disable this action |
| `email` | + | Recipient email address |
| `email_sender` | + | Sender email address |
| `email_subject` | - | Email subject (default: "BanForge Alert") |
| `smtp_host` | + | SMTP server host |
| `smtp_port` | + | SMTP server port |
| `smtp_user` | + | SMTP username |
| `smtp_password` | + | SMTP password |
| `smtp_tls` | - | Use TLS connection (default: false) |
| `body` | - | Email body text |

#### 2. Webhook Notification

Send HTTP webhook requests when an IP is banned.

```toml
[[rule.action]]
  type = "webhook"
  enabled = true
  url = "https://hooks.example.com/alert"
  method = "POST"
  headers = { "Content-Type" = "application/json", "Authorization" = "Bearer token" }
  body = "{\"ip\": \"{ip}\", \"rule\": \"{rule}\"}"
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | + | Must be "webhook" |
| `enabled` | + | Enable/disable this action |
| `url` | + | Webhook URL |
| `method` | - | HTTP method (default: "POST") |
| `headers` | - | HTTP headers as key-value pairs |
| `body` | - | Request body (supports variables) |

#### 3. Script Execution

Execute a custom script when an IP is banned.

```toml
[[rule.action]]
  type = "script"
  enabled = true
  script = "/usr/local/bin/notify.sh"
  interpretator = "bash"
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | + | Must be "script" |
| `enabled` | + | Enable/disable this action |
| `script` | + | Path to script file |
| `interpretator` | - | Script interpretator (e.g., "bash", "python"). If empty, script runs directly |

### Variables

The following variables can be used in `body` fields (email, webhook):

| Variable | Description |
|----------|-------------|
| `{ip}` | Banned IP address |
| `{rule}` | Rule name that triggered the ban |
| `{service}` | Service name |
| `{ban_time}` | Ban duration |
