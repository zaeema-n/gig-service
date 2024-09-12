#!/bin/bash

CONFIG_FILE="/app/conf/app.conf"

# Function to update or add a configuration
update_config() {
    key=$1
    value=$2
    if grep -q "^$key\s*=" "$CONFIG_FILE"; then
        sed -i "s|^$key\s*=.*|$key = $value|" "$CONFIG_FILE"
    else
        echo "$key = $value" >> "$CONFIG_FILE"
    fi
}

# Update configurations based on environment variables
[ ! -z "$MONGO_PATH" ] && update_config "mongo.path" "$MONGO_PATH"

# Print the updated configuration (optional, for debugging)
echo "Updated configuration:"
cat "$CONFIG_FILE"

# Run the application
exec /app/run.sh