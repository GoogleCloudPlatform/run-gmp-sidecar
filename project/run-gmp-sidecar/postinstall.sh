if command -v systemctl >/dev/null 2>&1; then
    if [ -d /run/systemd/system ]; then
        systemctl daemon-reload
    fi
    systemctl enable rungmpcol.service
    if [ -f /etc/rungmpcol/config.yaml ]; then
        if [ -d /run/systemd/system ]; then
            systemctl restart rungmpcol.service
        fi
    fi
fi
