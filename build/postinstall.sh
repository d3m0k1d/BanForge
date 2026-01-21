#!/bin/sh

if command -v systemctl >/dev/null 2>&1; then
  # for systemd based systems
  banforge init
  cat > /etc/systemd/system/banforge.service << 'EOF'
[Unit]
Description=BanForge - IPS log based system
After=network-online.target
Wants=network-online.target
Documentation=https://github.com/d3m0k1d/BanForge

[Service]
Type=simple
ExecStart=/usr/local/bin/banforge daemon
User=root
Group=root
Restart=always
StandardOutput=journal
StandardError=journal
SyslogIdentifier=banforge
TimeoutStopSec=90
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
EOF
  chmod 644 /etc/systemd/system/banforge.service
  systemctl daemon-reload
  systemctl enable banforge
fi

if command -v rc-service >/dev/null 2>&1; then
  # for openrc based systems
  banforge init
  cat > /etc/init.d/banforge << 'EOF'
#!/sbin/openrc-run

description="BanForge - IPS log based system"
command="/usr/bin/banforge"
command_args="daemon"

pidfile="/run/${RC_SVCNAME}.pid"
command_background="yes"

depend() {
  need net
  after network
}

start_post() {
  einfo "BanForge is now running"
}

stop_post() {
  einfo "BanForge is now stopped"
}
EOF
  chmod 755 /etc/init.d/banforge
  rc-update add banforge
fi
