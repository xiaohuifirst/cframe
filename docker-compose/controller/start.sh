#!/usr/bin/env bash
echo "$settings" > /config.toml
echo "setting:" $settings

/controller -c config.toml
