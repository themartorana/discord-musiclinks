#!/usr/bin/env sh
podman build --no-cache --tag 'discord-musiclinks' .

podman run -d -t -i \
        --name discord-musiclinks \
        --restart=always \
        --replace \
    'discord-musiclinks:latest'
