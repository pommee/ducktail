#!/bin/bash

# Function to generate a random log entry using jq
generate_log_entry() {
    jq -n \
    --arg timestamp "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
    --arg level "$(echo -e "INFO\nWARN\nERROR\nDEBUG" | shuf -n1)" \
    --arg action "$(echo -e "login\nlogout\npurchase\nview\nupdate" | shuf -n1)" \
    --arg user_id "$(uuidgen)" \
    --arg session_id "$(uuidgen)" \
    --arg ip_address "$(printf "%d.%d.%d.%d" $((RANDOM % 256)) $((RANDOM % 256)) $((RANDOM % 256)) $((RANDOM % 256)))" \
    --argjson request_time_ms "$((RANDOM % 1000))" \
    --argjson status_code "$((200 + RANDOM % 300))" \
    '{
        timestamp: $timestamp,
        level: $level,
        action: $action,
        user_id: $user_id,
        session_id: $session_id,
        ip_address: $ip_address,
        request_time_ms: $request_time_ms,
        status_code: $status_code
    }' | jq .
}

while true; do
    generate_log_entry
    sleep $(awk "BEGIN {srand(); print rand() * 1.9 + 0.1}")
done
