#!/bin/bash
# runme_local.sh

# This runs a locally built docker instance of the WhatsNewService and allows your to test your built package.
# This needs to have APP_ID, INSTALLATION_ID and the path to the GITHUB PEMFILE available in ENV or will fail

# This runs very privileged (host network mode, privileged, etc) and its not meant for production. Use the K8s Envelope for more 'production' ready packaging and runtime goodness

app=whatsnew-service
version=latest

### SET THESE TO YOU ENV
APP_ID=1408585
INSTALLATION_ID=71368278
PEMFILE=/app/whatsnew-github.pem

docker stop ${app}-instance
docker rm  ${app}-instance
docker run -d -p 10020:10020 -p 9200:9200 --name whatsnew-service-instance \
	--env APP_ID=${APP_ID} \
	--env INSTALLATION_ID=${INSTALLATION_ID} \
	--env PEMFILE=${PEMFILE} \
	-v $(pwd)/whatsnew-github.pem:/app/whatsnew-github.pem \
        --privileged \
        --restart=always \
        --network=host \
        ${app}:${version}

