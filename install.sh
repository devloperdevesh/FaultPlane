#!/bin/sh
# SPDX-License-Identifier: Apache-2.0

set -eu

APP_NAME="faultplane"
DAEMON_NAME="kmaskd"

PREFIX="${FAULTPLANE_PREFIX:-/usr/local}"
BIN_DIR="${PREFIX}/bin"
LIB_DIR="${PREFIX}/lib/${APP_NAME}"
CONFIG_DIR="${FAULTPLANE_CONFIG_DIR:-/etc/${APP_NAME}}"
BPF_DIR="${LIB_DIR}/bpf"

DAEMON="${BIN_DIR}/${DAEMON_NAME}"
BPF_OBJECT="${BPF_DIR}/sockmap.bpf.o"

log() {
    printf '[FaultPlane] %s\n' "$1"
}

fail() {
    printf '[FaultPlane] ERROR: %s\n' "$1" >&2
    exit 1
}

require_linux() {
    os="$(uname -s 2>/dev/null || true)"

    if [ "$os" != "Linux" ]; then
        fail "Linux is required; detected ${os:-unknown}."
    fi
}

require_root() {
    if [ "$(id -u)" -ne 0 ]; then
        fail "run the installer with sudo."
    fi
}

check_kernel_support() {
    kernel="$(uname -r)"

    log "Detected Linux kernel: ${kernel}"

    if [ -r /sys/kernel/btf/vmlinux ]; then
        log "Kernel BTF: available"
    else
        log "Kernel BTF: unavailable"
        log "BPF loading may require an environment-specific fallback."
    fi
}

install_binary() {
    if [ -x "./${DAEMON_NAME}" ]; then
        log "Installing local ${DAEMON_NAME} binary..."

        install -Dm0755 \
            "./${DAEMON_NAME}" \
            "${DAEMON}"
        return
    fi

    fail "${DAEMON_NAME} was not found in the current directory."
}

install_bpf_object() {
    source=""

    if [ -f "./sockmap.bpf.o" ]; then
        source="./sockmap.bpf.o"
    elif [ -f "./internal/kernel/bpf/sockmap.bpf.o" ]; then
        source="./internal/kernel/bpf/sockmap.bpf.o"
    fi

    if [ -z "$source" ]; then
        log "No sockmap.bpf.o artifact found."
        log "Daemon installed; BPF artifact was not installed."
        return
    fi

    install -Dm0644 \
        "$source" \
        "${BPF_OBJECT}"

    log "Installed BPF object: ${BPF_OBJECT}"
}

configure_capabilities() {
    if ! command -v setcap >/dev/null 2>&1; then
        log "setcap is unavailable; skipping Linux capability configuration."
        return
    fi

    setcap 'cap_net_admin,cap_sys_admin+ep' "${DAEMON}"

    log "Configured required network/BPF capabilities."
}

write_environment() {
    install -d -m0755 "${CONFIG_DIR}"

    cat > "${CONFIG_DIR}/environment" <<EOF
# FaultPlane runtime configuration
FAULTPLANE_BPF_OBJECT=${BPF_OBJECT}
FAULTPLANE_CONFIG_DIR=${CONFIG_DIR}
EOF

    chmod 0644 "${CONFIG_DIR}/environment"
}

verify_installation() {
    [ -x "${DAEMON}" ] ||
        fail "daemon installation verification failed."

    log "Daemon: ${DAEMON}"

    if command -v getcap >/dev/null 2>&1; then
        capabilities="$(getcap "${DAEMON}" 2>/dev/null || true)"

        if [ -n "$capabilities" ]; then
            log "Capabilities: ${capabilities}"
        else
            log "Capabilities: not configured"
        fi
    fi

    if [ -f "${BPF_OBJECT}" ]; then
        log "BPF object: ${BPF_OBJECT}"
    fi
}

main() {
    log "Starting FaultPlane installation."

    require_linux
    require_root
    check_kernel_support

    install -d -m0755 \
        "${BIN_DIR}" \
        "${LIB_DIR}" \
        "${BPF_DIR}" \
        "${CONFIG_DIR}"

    install_binary
    install_bpf_object
    configure_capabilities
    write_environment
    verify_installation

    log "Installation completed successfully."
    log "Daemon path: ${DAEMON}"
    log "Configuration: ${CONFIG_DIR}/environment"
}

main "$@"