#!/usr/bin/env bash
set -euo pipefail

TAG="${1:-v0.0.0}"
cat <<NOTES
## Release ${TAG}

### New commands
- 

### Breaking changes
- None

### Known issues
- Windows token retrieval backend may require org-specific credential integration.

### Install and upgrade
- Download binary assets from this release.
NOTES
