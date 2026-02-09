# CLI commands BanForge
BanForge provides a command-line interface (CLI) to manage IP blocking, 
configure detection rules, and control the daemon process.
## Commands
### init - create a deps file

```shell
banforge init
```

**Description**
This command creates the necessary directories and base configuration files 
required for the daemon to operate.

### version - Display BanForge version

```shell
banforge version
```

**Description**
This command displays the current version of the BanForge software.

### daemon - Starts the BanForge daemon process

```shell
banforge daemon
```

**Description**
This command starts the BanForge daemon process in the background. 
The daemon continuously monitors incoming requests, detects anomalies, 
and applies firewall rules in real-time.

### firewall - Manages firewall rules
```shell
banforge ban <ip>
banforge unban <ip>
```

**Description**
These commands provide an abstraction over your firewall. If you want to simplify the interface to your firewall, you can use these commands.

Flag -t or -ttl add bantime if not used default ban 1 year

### ports - Open and Close ports on firewall

```shell
banforge open -port <port> -protocol <protocol>
banforge close -port <port> -protocol <protocol>
```

**Description**
These commands provide an abstraction over your firewall. If you want to simplify the interface to your firewall, you can use these commands.

### list - Lists the IP addresses that are currently blocked
```shell
banforge list
```

**Description**
This command output table of IP addresses that are currently blocked

### rule - Manages detection rules

```shell
banforge rule add -n rule.name -c 403
banforge rule list 
```

**Description**
These command help you to create and manage detection rules in CLI interface.

| Flag        | Required |
| ----------- | -------- |
| -n -name    | +        |
| -s -service | +        |
| -p -path    | -        |
| -m -method  | -        |
| -c -status  | -        |
| -t -ttl     | -(if not used default ban 1 year)        |

You must specify at least 1 of the optional flags to create a rule.
