#!/usr/bin/env bash
#
# This script will build and tag a docker image, push it to ECR and update ECS.
#
# If your AWS CLI configuration has a named profile other than 'default' for AWS account
# then pass the profile name as the first argument.  If a profile is not passed as an argument
# the default profile is used.
#
# Usage:
#
# ./deploy-to-aws-ecs.sh <aws account id> <profile>
#
awsAccountId=$1
profile=${2:-default}
region=us-west-2
cluster=star-wars
ecrRepo=star-wars-server
serviceName=star-wars-server
taskFamily=star-wars-server
taskDefinitionFile=ecs-task-definition.json
ecrArn=${awsAccountId}.dkr.ecr.${region}.amazonaws.com/${ecrRepo}
version=$(git rev-parse --short HEAD)

check_version() {
	local imageId=$(\
		aws --profile $profile ecr list-images --repository-name $ecrRepo \
		| jq --arg version $version '.imageIds[] | select(.imageTag == $version)' \
	)
	if [ -n "$imageId" ]; then
		echo "${ecrArn}:${version} already exists in ECR"
		echo " - Add commit "
		exit 1
	fi
}

build_docker_image() {
	docker build -t ${ecrRepo} .
	docker tag ${ecrRepo}:latest ${ecrRepo}:${version}
	docker tag ${ecrRepo}:latest ${ecrArn}:latest
	docker tag ${ecrRepo}:latest ${ecrArn}:${version}
}

push_to_ecr() {
    $(aws --profile $profile ecr get-login --no-include-email)
	docker push ${ecrArn}:latest
	docker push ${ecrArn}:${version}
}

update_ecs() {
	echo 'Creating new task definition...'
	local taskDefinition=$(\
		aws --profile $profile ecs describe-task-definition --task-definition $taskFamily \
		| jq --arg image ${ecrArn}:${version} '.taskDefinition | .containerDefinitions[].image = $image' \
		| jq 'def maybe(k): if has(k) then {(k): .[k]} else null end;
			  maybe("taskRoleArn") +
			  {family, networkMode, containerDefinitions, volumes, placementConstraints}' \
	)
	echo $taskDefinition | jq '.' > $taskDefinitionFile
	local taskDefinitionArn=$(\
		aws --profile $profile ecs register-task-definition --cli-input-json file://${taskDefinitionFile} \
		| jq -r '.taskDefinition.taskDefinitionArn' \
	)
	rm $taskDefinitionFile
	echo 'Updating ECS service...'
	local service=$(\
		aws --profile $profile ecs update-service \
		--cluster $cluster \
		--service $serviceName \
		--desired-count 1 \
		--task-definition $taskDefinitionArn \
	)
}

#check_version
build_docker_image
push_to_ecr
#update_ecs
