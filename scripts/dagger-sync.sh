#!/usr/bin/env bash

set -euo pipefail

# Function to log messages
log() {
  local type="$1"; shift
  printf "[%s] %s\n" "$type" "$*"
}

# Find all directories containing a 'dagger' subdirectory
modules=$(find . -type d -name 'dagger' -print)

if [ -z "$modules" ]; then
  log "WARNING" "No Dagger modules found."
  exit 1
fi

# Initialize counters
total_modules=0
successful_modules=0

# Count total modules
for dir in $modules; do
  ((total_modules++))
done

# Function to update progress bar
update_progress() {
  local progress=$((successful_modules * 100 / total_modules))
  printf "\rProgress: [%-50s] %d%%" $(printf "#%.0s" $(seq 1 $((progress / 2)))) $progress
}

log "INFO" "Syncing Dagger modules..."
echo

for dir in $modules; do
  module=$(dirname "$dir")
  echo -n "Syncing module: $module... "
  
  if (cd "$dir" && dagger mod sync) > /dev/null 2>&1; then
    echo "✅"
    ((successful_modules++))
  else
    echo "❌"
    log "ERROR" "Failed to sync module: $module"
  fi
  
  update_progress
done

echo
echo

if [ $successful_modules -eq $total_modules ]; then
  log "SUCCESS" "Dagger sync completed for all modules successfully! 🎉"
else
  log "WARNING" "Dagger sync completed with some failures. Please check the output above."
  exit 1
fi
