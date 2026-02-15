#!/bin/sh
# Convenience entrypoint: install.sh at repo root fetches and runs scripts/install.sh
# so both URLs work:
#   .../main/install.sh
#   .../main/scripts/install.sh
REPO="${CODEMINT_REPO:-neghani/code-mint-cli}"
curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/scripts/install.sh" | sh
