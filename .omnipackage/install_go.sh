#!/usr/bin/env bash
# Ensure a Go toolchain capable of building this module (go.mod requires 1.26).
# Runs inside the build container, after build_dependencies are installed.
#
# Strategy, in order:
#   1. If the distro already ships go >= 1.21, use it — GOTOOLCHAIN=auto will
#      transparently fetch the exact 1.26 toolchain at build time.
#   2. Otherwise download the official toolchain (works wherever curl/wget is
#      healthy).
#   3. If even that fails (e.g. the broken libcurl/ngtcp2 on some rolling
#      images), fall back to whatever go the distro provides.
set -xEeuo pipefail

GO_VERSION=1.26.2

have_usable_go() {
    command -v go >/dev/null 2>&1 || return 1
    # Accept go1.21+ (1.21..1.99) and any go2.x — enough for GOTOOLCHAIN=auto.
    go version | grep -qE 'go1\.(2[1-9]|[3-9][0-9])|go[2-9]\.'
}

if have_usable_go; then
    go version
    exit 0
fi

url="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
if curl -fLsS "$url" -o /tmp/go.tgz 2>/dev/null || wget -qO /tmp/go.tgz "$url" 2>/dev/null; then
    rm -rf /usr/local/go
    tar -zxf /tmp/go.tgz -C /usr/local
    ln -sf /usr/local/go/bin/* /usr/bin/
    go version
    exit 0
fi

# Last resort: rely on the distro toolchain (may be older; GOTOOLCHAIN=auto
# bootstraps the rest if it is >= 1.21).
go version
