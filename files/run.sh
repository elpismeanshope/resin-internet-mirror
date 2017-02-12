#!/bin/bash
set -euo pipefail

# Networking config
## FIXME: Move dhcp range here
WLAN_ADDRESS=10.10.0.1
WLAN_NETMASK=22

if [[ -z "${WLAN_SKIP:-}" ]]; then
  if ! ip link show wlan0; then
    echo "Interface wlan0 not found. Set WLAN_SKIP=1 to disable wlan configuration" >&2
    exit 1
  fi
  ip addr show wlan0 | grep "$WLAN_ADDRESS" \
    || ip addr add $WLAN_ADDRESS/$WLAN_NETMASK dev wlan0
fi

# SSH Config
ssh-keygen -A
if [[ -n "${SSH_PASS:-}" ]]; then
  echo "root:${SSH_PASS}" | chpasswd
fi

# Logfiles
mkdir -p /data/log

# Configure ka-lite
if [[ ! -d "/data/kalite" ]]; then
  install -o kalite -g kalite -d /data/kalite
  su -c 'kalite manage initialize_kalite' kalite
fi

if [[ -n "${KA_PASSWORD:-}" && -n "${KA_EMAIL:-}" ]]; then
  KA_USER=admin su -c 'kalite shell' kalite < /opt/ka-lite-password
fi

# Configuring mobile data uplink
## We don't run udev, so we need devtmpfs and devpts
mount -t devtmpfs none /dev
mount -t devpts none /dev/pts
args=$(lsusb.py \
  | awk '/HUAWEI/ { split($2,a,":"); print "-v " a[1] " -p " a[2]}')

if [[ -z "${args:-}" ]]; then
  echo "No modem found, skipping"
else
  if ! usb_modeswitch $args -J; then
    echo "Couldn't switch modem, attempting reboot"
    curl -X POST --header "Content-Type:application/json" \
      "$RESIN_SUPERVISOR_ADDRESS/v1/reboot?apikey=$RESIN_SUPERVISOR_API_KEY"
    exit 1
  fi
fi

# Exec s6
exec s6-svscan /etc/service
