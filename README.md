# BanForge

Log-based IPS system written on Go for Linux based system.

# Table of contents
1. [Overview](#overview)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Usage](#usage)
5. [License](#license)

# Overview
BanForge is a simple IPS for replacement fail2ban in Linux system.
The project is currently in its early stages of development.
All release are available on my self-hosted [Gitea](https://gitea.d3m0k1d.ru/d3m0k1d/BanForge) because Github have limit for Actions.
If you have any questions or suggestions, create issue on [Github](https://github.com/d3m0k1d/BanForge/issues).

## Roadmap
- [ ] Real-time Nginx log monitoring
- [ ] Add support for other service
- [ ] Add support for user service with regular expressions
- [ ] TUI interface

# Requirements

- Go 1.21+
- ufw/iptables/nftables/firewalld

# Installation
currently no binary file if you wanna build the project yourself, you can use [Makefile](https://github.com/d3m0k1d/BanForge/blob/master/Makefile)

# Usage

# License
The project is licensed under the [GPL-3.0](https://github.com/d3m0k1d/BanForge/blob/master/LICENSE)
