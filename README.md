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
- [x] Prometheus metrics
- [x] Actions (email, webhook, script)
- [ ] Add support for most popular web-service
- [ ] User regexp for custom services
- [ ] TUI interface

# Requirements

- Go 1.25+
- ufw/iptables/nftables/firewalld

# Installation
Search for a release on the [Gitea](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases) releases page and download it.
In release page you can find rpm, deb, apk packages, for amd or arm architecture.

## Installation guide for packages

### Debian/Ubuntu(.deb)
```bash
# Download the latest DEB package
wget https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases/download/v0.6.0/banforge_0.6.0_linux_amd64.deb

# Install
sudo dpkg -i banforge_0.6.0_linux_amd64.deb

# Verify installation
sudo systemctl status banforge
```

### RHEL-based(.rpm)
```bash

# Download
wget https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases/download/v0.6.0/banforge_0.6.0_linux_amd64.rpm

# Install
sudo rpm -i banforge_0.6.0_linux_amd64.rpm

# Or with dnf (CentOS 8+, AlmaLinux)
sudo dnf install banforge_0.6.0_linux_amd64.rpm

# Verify
sudo systemctl status banforge
```

### Alpine(.apk)
```bash

# Download
wget https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases/download/v0.6.0/banforge_0.6.0_linux_amd64.apk

# Install
sudo apk add --allow-untrusted banforge_0.6.0_linux_amd64.apk

# Verify
sudo rc-service banforge status
```

### Arch Linux(.pkg.tar.zst)
```bash

# Download
wget https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases/download/v0.6.0/banforge_0.6.0_linux_amd64.pkg.tar.zst

# Install
sudo pacman -U banforge_0.6.0_linux_amd64.pkg.tar.zst

# Verify
sudo systemctl status banforge
```
This is examples for other versions with different architecture or new versions check release page on [Gitea](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge/releases).

## Installation guide for source code
```bash
# Download
git clone https://github.com/d3m0k1d/BanForge.git
cd BanForge
make build-daemon
cd bin
mv banforge /usr/bin/banforge
cd ..
# Add init script and uses banforge init
cd build
./postinstall.sh
```
# Usage
For first steps use this commands
```bash
banforge init   # Create config files and database
banforge daemon # Start BanForge daemon (use systemd or another init system to create a service)
```
You can edit the config file with examples in
- `/etc/banforge/config.toml` main config file
- `/etc/banforge/rules.d/*.toml` individual rule files with actions support

For more information see the [docs](https://github.com/d3m0k1d/BanForge/docs).

## Actions

BanForge supports actions that are executed after a successful IP ban:
- **Email** - Send email notifications via SMTP
- **Webhook** - Send HTTP requests to external services (Slack, Telegram, etc.)
- **Script** - Execute custom scripts

See [configuration docs](https://github.com/d3m0k1d/BanForge/blob/main/docs/config.md#actions) for detailed setup instructions.

# License
The project is licensed under the [GPL-3.0](https://github.com/d3m0k1d/BanForge/blob/master/LICENSE)
