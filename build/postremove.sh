#!/bin/sh

if command -v systemctl >/dev/null 2>&1; then
  # for systemd based systems
  systemctl stop banforge 2>/dev/null || true
  systemctl disable banforge 2>/dev/null || true
  rm -f /etc/systemd/system/banforge.service
  systemctl daemon-reload
fi

if command -v rc-service >/dev/null 2>&1; then
  # for openrc based systems
  rc-service banforge stop 2>/dev/null || true
  rc-update del banforge 2>/dev/null || true
  rm -f /etc/init.d/banforge
fi

rm -rf /etc/banforge/
rm -rf /var/lib/banforge/
rm -rf /var/log/banforge/
