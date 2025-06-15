#!/bin/bash
# runme_local.sh

# This runs a locally built docker instance of the WhatsNewService and allows your to test your built package.
# This needs to have APP_ID, INSTALLATION_ID and the path to the GITHUB PEMFILE available in ENV or will fail

# This runs very privileged (host network mode, privileged, etc) and its not meant for production. Use the K8s Envelope for more 'production' ready packaging and runtime goodness


app=whatsnew-service
version=latest

docker stop ${app}-instance
docker rm  ${app}-instance

APP_ID=1408585
INSTALLATION_ID=71368278
PEM_CONTENTS="$(cat whatsnew-github.pem)"


if [[ -z "$app" || -z "$version" ]]; then
  echo "ERROR: app or version variable is not set!"
  exit 1
fi

echo "Running: docker run ... $app:$version"
docker run -d -p 10020:10020 -p 9200:9200 --name ${app}-instance \
  --env APP_ID=${APP_ID} \
  --env INSTALLATION_ID=${INSTALLATION_ID} \
  --env PEM_CONTENTS="${PEM_CONTENTS}" \
  --restart=always \
  ${app}:${version}

