#!/bin/bash

set -e -x

pip install --user awscli
export PATH=$PATH:$HOME/.local/bin

add-apt-repository ppa:eugenesan/ppa
apt-get update
apt-get install jq -y

curl https://raw.githubusercontent.com/silinternational/ecs-deploy/master/ecs-deploy | \
  sudo tee -a /usr/bin/ecs-deploy
sudo chmod +x /usr/bin/ecs-deploy

eval $(aws ecr get-login --no-include-email --region $AWS_REGION)

aws ecs update-service --cluster ${CLUSTER} --service ${API_SERVICE} --region ${AWS_REGION} --force-new-deployment
aws ecs update-service --cluster ${CLUSTER} --service ${CONSUMER_SERVICE} --region ${AWS_REGION} --force-new-deployment
