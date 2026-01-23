# BanForge

Log-based IPS system written in Go for Linux-based system.

[![Go Reference](https://pkg.go.dev/badge/github.com/d3m0k1d/BanForge/cmd/banforge.svg)](https://pkg.go.dev/github.com/d3m0k1d/BanForge)
[![License](https://img.shields.io/badge/license-%20%20GNU%20GPLv3%20-green?style=plastic)](https://github.com/d3m0k1d/BanForge/blob/master/LICENSE)
[![Build Status](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/actions/workflows/CI.yml/badge.svg?branch=master)](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/actions)
![GitHub Tag](https://img.shields.io/github/v/tag/d3m0k1d/BanForge)
# Table of contents
1. [Overview](#overview)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Usage](#usage)
5. [License](#license)

# Overview
BanForge is a simple IPS for replacement fail2ban in Linux system.
All release are available on my self-hosted [Gitea](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge) after release v1.0.0 are available on Github release page.
If you have any questions or suggestions, create issue on [Github](https://github.com/d3m0k1d/BanForge/issues).

## Roadmap
- [x] Rule system
- [x] Nginx and Sshd support
- [x] Working with ufw/iptables/nftables/firewalld
- [] Add support for most popular web-service
- [] User regexp for custom services
- [] TUI interface

# Requirements

- Go 1.25+
- ufw/iptables/nftables/firewalld

# Installation
Search for a release on the [Gitea](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases) releases page and download it. Then create or copy(/build dir) a systemd unit(openrc script0) file.
Or clone the repo and use the Makefile.
```
git clone https://gitea.d3m0k1d.ru/d3m0k1d/BanForge.git
cd BanForge
sudo make build-daemon
cd bin
```

# Usage
For first steps use this commands
```bash
banforge init   # Create config files and database
banforge daemon # Start BanForge daemon (use systemd or another init system to create a service)
```
You can edit the config file with examples in 
- `/etc/banforge/config.toml` main config file
- `/etc/banforge/rules.toml` ban rules
For more information see the [docs](https://github.com/d3m0k1d/BanForge/docs).

# License
The project is licensed under the [GPL-3.0](https://github.com/d3m0k1d/BanForge/blob/master/LICENSE)
