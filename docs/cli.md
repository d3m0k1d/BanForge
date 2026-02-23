# CLI commands BanForge

BanForge provides a command-line interface (CLI) to manage IP blocking,
configure detection rules, and control the daemon process.

## Commands

### init - Create configuration files

```shell
banforge init
```

**Description**  
This command creates the necessary directories and base configuration files
required for the daemon to operate:
- `/etc/banforge/config.toml` — main configuration
- `/etc/banforge/rules.toml` — default rules file
- `/etc/banforge/rules.d/` — directory for individual rule files

---

### version - Display BanForge version

```shell
banforge version
```

**Description**  
This command displays the current version of the BanForge software.

---

### daemon - Starts the BanForge daemon process

```shell
banforge daemon
```

**Description**  
This command starts the BanForge daemon process in the background.
The daemon continuously monitors incoming requests, detects anomalies,
and applies firewall rules in real-time.

---

### firewall - Manages firewall rules

```shell
banforge ban <ip>
banforge unban <ip>
```

**Description**  
These commands provide an abstraction over your firewall. If you want to simplify the interface to your firewall, you can use these commands.

| Flag        | Description                    |
| ----------- | ------------------------------ |
| `-t`, `-ttl` | Ban duration (default: 1 year) |

**Examples:**
```bash
# Ban IP for 1 hour
banforge ban 192.168.1.100 -t 1h

# Unban IP
banforge unban 192.168.1.100
```

---

### ports - Open and close ports on firewall

```shell
banforge open -port <port> -protocol <protocol>
banforge close -port <port> -protocol <protocol>
```

**Description**  
These commands provide an abstraction over your firewall. If you want to simplify the interface to your firewall, you can use these commands.

| Flag          | Required | Description              |
| ------------- | -------- | ------------------------ |
| `-port`       | +        | Port number (e.g., 80)   |
| `-protocol`   | +        | Protocol (tcp/udp)       |

**Examples:**
```bash
# Open port 80 for TCP
banforge open -port 80 -protocol tcp

# Close port 443
banforge close -port 443 -protocol tcp
```

---

### list - List blocked IP addresses

```shell
banforge list
```

**Description**  
This command outputs a table of IP addresses that are currently blocked.

---

### rule - Manage detection rules

Rules are stored in `/etc/banforge/rules.d/` as individual `.toml` files.

#### Add a new rule

```shell
banforge rule add -n <name> -s <service> [options]
```

**Flags:**

| Flag                | Required | Description                              |
| ------------------- | -------- | ---------------------------------------- |
| `-n`, `--name`      | +        | Rule name (used as filename)             |
| `-s`, `--service`   | +        | Service name (nginx, apache, ssh, etc.)  |
| `-p`, `--path`      | -        | Request path to match                    |
| `-m`, `--method`    | -        | HTTP method (GET, POST, etc.)            |
| `-c`, `--status`    | -        | HTTP status code (403, 404, etc.)        |
| `-t`, `--ttl`       | -        | Ban duration (default: 1y)               |
| `-r`, `--max_retry` | -        | Max retries before ban (default: 0)      |

**Note:** At least one of `-p`, `-m`, or `-c` must be specified.

**Examples:**
```bash
# Ban on 403 status
banforge rule add -n "Forbidden" -s nginx -c 403 -t 30m

# Ban on path pattern
banforge rule add -n "Admin Access" -s nginx -p "/admin/*" -t 2h -r 3

# SSH brute force protection
banforge rule add -n "SSH Bruteforce" -s ssh -c "Failed" -t 1h -r 5
```

---

#### List all rules

```shell
banforge rule list
```

**Description**  
Displays all configured rules in a table format.

**Example output:**
```
+------------------+---------+--------+--------+--------+----------+---------+
| NAME             | SERVICE | PATH   | STATUS | METHOD | MAXRETRY | BANTIME |
+------------------+---------+--------+--------+--------+----------+---------+
| SSH Bruteforce   | ssh     |        | Failed |        | 5        | 1h      |
| Nginx 404        | nginx   |        | 404    |        | 3        | 30m     |
| Admin Panel      | nginx   | /admin |        |        | 2        | 2h      |
+------------------+---------+--------+--------+--------+----------+---------+
```

---

#### Edit an existing rule

```shell
banforge rule edit -n <name> [options]
```

**Description**  
Edit fields of an existing rule. Only specified fields will be updated.

| Flag                | Required | Description                     |
| ------------------- | -------- | ------------------------------- |
| `-n`, `--name`      | +        | Rule name to edit               |
| `-s`, `--service`   | -        | New service name                |
| `-p`, `--path`      | -        | New path                        |
| `-m`, `--method`    | -        | New method                      |
| `-c`, `--status`    | -        | New status code                 |

**Examples:**
```bash
# Update ban time for existing rule
banforge rule edit -n "SSH Bruteforce" -t 2h

# Change status code
banforge rule edit -n "Forbidden" -c 403
```

---

#### Remove a rule

```shell
banforge rule remove <name>
```

**Description**  
Permanently delete a rule by name.

**Example:**
```bash
banforge rule remove "Old Rule"
```

---

## Ban time format

Use the following suffixes for ban duration:

| Suffix | Duration |
| ------ | -------- |
| `s`    | Seconds  |
| `m`    | Minutes  |
| `h`    | Hours    |
| `d`    | Days     |
| `M`    | Months (30 days) |
| `y`    | Years (365 days) |

**Examples:** `30s`, `5m`, `2h`, `1d`, `1M`, `1y`
