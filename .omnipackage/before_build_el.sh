#!/usr/bin/env bash
# RHEL clones (AlmaLinux / RockyLinux): the X11 and pcsc-lite -devel packages
# live in CRB (CodeReady Builder / PowerTools), which is disabled by default —
# so they can't go in build_dependencies (that install runs before this script).
# Enable CRB + EPEL here, install the build libs, then ensure Go.
set -xEeuo pipefail

dnf -y install dnf-plugins-core || true
# CRB is 'crb' on EL9/EL10 and 'powertools' on EL8.
dnf config-manager --set-enabled crb 2>/dev/null \
  || dnf config-manager --set-enabled powertools 2>/dev/null || true
dnf -y install epel-release || true

dnf -y install \
  pkgconf-pkg-config pcsc-lite-devel mesa-libGL-devel \
  libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel \
  libXi-devel libXxf86vm-devel golang

# Reuse the shared Go bootstrap (no-op if the golang just installed is >= 1.21).
exec "$(dirname "$0")/install_go.sh"
