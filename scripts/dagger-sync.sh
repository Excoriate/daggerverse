#!/bin/bash

# Strict mode
set -euo pipefail
IFS=$'\n\t'

# Function to log messages
log() {
  local type="$1"; shift
  printf "[%s] %s\n" "$type" "$*"
}

# Function to process directories containing dagger.json
process_dagger_modules() {
  local dir=$1
  log "INFO" "Processing module in ${dir}..."
  if dagger mod sync -m "${dir}"; then
    log "SUCCESS" "Module sync successful in ${dir}"
  else
    log "ERROR" "Module sync failed in ${dir}"
  fi
}

# Main function
main() {
  log "INFO" "Scanning for Dagger modules..."
  local modules_found=0

  # Search one level below the root directory
  local search_path='*/dagger.json'

  # Find directories containing dagger.json and process them
  for dir in ${search_path}; do
    if [ -f "${dir}" ]; then
      process_dagger_modules "$(dirname "${dir}")"
      ((modules_found++))
      log "INFO" "Module found: $(dirname "${dir}")"
    fi
  done

  if ((modules_found == 0)); then
    log "WARNING" "No Dagger modules found."
  else
    log "INFO" "Processed ${modules_found} Dagger modules."
  fi
}

# Run the main function
main "$@"
