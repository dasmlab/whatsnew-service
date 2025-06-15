#!/bin/bash
# runme_local.sh

# This runs a locally built docker instance of the WhatsNewService and allows your to test your built package.
# This needs to have APP_ID, INSTALLATION_ID and the path to the GITHUB PEMFILE available in ENV or will fail

# This runs very privileged (host network mode, privileged, etc) and its not meant for production. Use the K8s Envelope for more 'production' ready packaging and runtime goodness

app=mcp-go-server-manager
version=latest

docker stop mcp-go-server-manager
docker rm  mcp-go-server-manager
docker run -d --name mcp-go-server-manager \
        --privileged \
        --restart=always \
        --network=host \
        ${app}:${version}

