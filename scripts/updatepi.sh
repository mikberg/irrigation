#!/bin/bash

# --- begin runfiles.bash initialization v2 ---
# Copy-pasted from the Bazel Bash runfiles library v2.
set -uo pipefail; f=bazel_tools/tools/bash/runfiles/runfiles.bash
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \
source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null || \
source "$0.runfiles/$f" 2>/dev/null || \
source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
{ echo>&2 "ERROR: cannot find $f"; exit 1; }; f=; set -e
# --- end runfiles.bash initialization v2 ---

USER="mikael"
PI="irr.mikbe.com"
BINARY=$(rlocation irrigation/irrigation/irrigation-raspberrypi_/irrigation-raspberrypi)

SSH_COMBO="$USER@$PI"

# remove the file
ssh "$SSH_COMBO" 'rm -f /tmp/irrigation'

# copy the file
scp "$BINARY" "$SSH_COMBO:/tmp/irrigation"

# restart server
ssh "$SSH_COMBO" 'sudo supervisorctl stop irrigation && sudo mv /tmp/irrigation /bin/irrigation && sudo supervisorctl start irrigation'
