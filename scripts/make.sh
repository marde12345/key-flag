#!/usr/bin/env bash
set -e

# Action List
ACTIONS=(
   build-binary 
)

SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd -P)"

runaction() {
    local action="$1"; shift
    echo "Running Action: $(basename "$action")(in $SCRIPT_DIR)" 
    source "${SCRIPT_DIR}/$action"
}

if [ $# -lt 1 ]; then
    action=${ACTIONS[*]}
else
    action=${*}
fi
for action in ${action[*]}; do
    runaction "$action"
    echo
done