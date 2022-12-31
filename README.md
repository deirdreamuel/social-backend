# Go Project

### Running project locally using Docker Compose
1. `docker compose build`
2. `docker compose up`

### Uploading Docker image to AWS ECR
Visit: https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html
1. run `docker images` to list Docker images and copy Docker Image ID
2. aws ecr get-login-password --region `<AWS_REGION>` | docker login --username AWS --password-stdin `<AWS_ACCOUNT_ID>`.dkr.ecr.`<AWS_REGION>`.amazonaws.com
3. docker tag `<IMAGE_ID>` `<AWS_ACCOUNT_ID>`.dkr.ecr.`<AWS_REGION>`.amazonaws.com/main:latest
4. docker push `<AWS_ACCOUNT_ID>`.dkr.ecr.`<AWS_REGION>`.amazonaws.com/main:latest
