#!/bin/bash

# デプロイの動作確認に用いるスクリプト

set -e

if [ $# != 1 ]; then
    echo 引数エラー: AWS_PROFILEを指定してください
    exit 1
fi

AWS_PROFILE=$1
REPOSITORY_URI="471112567496.dkr.ecr.ap-northeast-1.amazonaws.com"
IMAGE_TAG="debug"
APP_REPO_NAME="e-privado-stg-repo"

# push to ecr
aws ecr get-login-password --region ap-northeast-1 --profile $AWS_PROFILE | docker login --username AWS --password-stdin $REPOSITORY_URI
docker build ./api -t $REPOSITORY_URI/$APP_REPO_NAME:$IMAGE_TAG --build-arg ENTRYPOINT=./run
docker push $REPOSITORY_URI/$APP_REPO_NAME:$IMAGE_TAG

pushd cdk
npm i
npm install -g aws-cdk
npm run build
cdk --version
cdk synth -c tag=$IMAGE_TAG -c stage=stg --profile $AWS_PROFILE
cdk deploy --all -c tag=$IMAGE_TAG -c stage=stg --profile $AWS_PROFILE

popd

set +e
