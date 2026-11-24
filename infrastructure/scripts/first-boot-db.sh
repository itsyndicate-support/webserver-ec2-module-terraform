#!/bin/bash
set -Eeuo pipefail

exec > >(tee -a /var/log/week11-db-bootstrap.log /dev/console) 2>&1
trap 'status=$?; echo "MARIADB_BOOTSTRAP_FAILED status=$status line=$LINENO" >&2; exit "$status"' ERR

export DEBIAN_FRONTEND=noninteractive

for source_file in /etc/apt/sources.list /etc/apt/sources.list.d/*.list /etc/apt/sources.list.d/*.sources; do
  if [[ -f "$source_file" ]]; then
    sed -i -E \
      's#https?://security\.ubuntu\.com/ubuntu#http://us-east-1.ec2.archive.ubuntu.com/ubuntu#g' \
      "$source_file"
  fi
done

apt_options=(
  -o Acquire::ForceIPv4=true
  -o Acquire::Retries=3
  -o Acquire::http::Timeout=20
  -o Acquire::https::Timeout=20
  -o DPkg::Lock::Timeout=300
)

retry_apt() {
  local attempt
  for attempt in 1 2 3; do
    if apt-get "${apt_options[@]}" "$@"; then
      return 0
    fi
    if [[ "$attempt" -lt 3 ]]; then
      echo "APT attempt $attempt/3 failed; retrying in 10 seconds" >&2
      sleep 10
    fi
  done
  return 1
}

retry_apt -o APT::Update::Error-Mode=any update
retry_apt install -y mariadb-server

systemctl enable --now mariadb
systemctl is-active --quiet mariadb

result="$(mariadb --batch --skip-column-names -e 'SELECT 1;')"
[[ "$result" == "1" ]]

echo MARIADB_BOOTSTRAP_READY
