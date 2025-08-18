# Deployment Guide for Packulator Application

This guide provides step-by-step instructions for deploying the Packulator application to AWS using CodeBuild, ECS, and CloudFront.

## Architecture Overview

The deployment consists of:
- **Backend**: Go application running in ECS Fargate containers behind an Application Load Balancer
- **Frontend**: React application served from CloudFront + S3 (optional) or embedded in the Go app
- **Database**: PostgreSQL on RDS
- **CI/CD**: AWS CodeBuild for automated builds and deployments
- **Infrastructure**: Managed via CloudFormation templates

## Prerequisites

1. **AWS Account** with appropriate permissions
2. **AWS CLI** configured with your credentials
3. **Docker** installed locally (for testing)
4. **GitHub repository** with your code
5. **Domain name** (optional, for HTTPS and custom domain)

## Step 1: Set up AWS Infrastructure

### 1.1 Deploy Base Infrastructure

```bash
# Deploy the main infrastructure stack
aws cloudformation create-stack \
  --stack-name packulator-prod-infrastructure \
  --template-body file://cloudformation/infrastructure.yml \
  --parameters ParameterKey=ApplicationName,ParameterValue=packulator \
               ParameterKey=Environment,ParameterValue=prod \
               ParameterKey=DBPassword,ParameterValue=YourSecurePassword123! \
  --capabilities CAPABILITY_IAM
```

### 1.2 Wait for Stack Creation

```bash
# Monitor stack creation
aws cloudformation wait stack-create-complete \
  --stack-name packulator-prod-infrastructure

# Get stack outputs
aws cloudformation describe-stacks \
  --stack-name packulator-prod-infrastructure \
  --query 'Stacks[0].Outputs'
```

## Step 2: Set up CodeBuild Project

### 2.1 Create CodeBuild Service Role

```bash
# Create IAM role for CodeBuild
aws iam create-role \
  --role-name PackulatorCodeBuildRole \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "codebuild.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }'

# Attach policies
aws iam attach-role-policy \
  --role-name PackulatorCodeBuildRole \
  --policy-arn arn:aws:iam::aws:policy/CloudWatchLogsFullAccess

aws iam attach-role-policy \
  --role-name PackulatorCodeBuildRole \
  --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser

aws iam attach-role-policy \
  --role-name PackulatorCodeBuildRole \
  --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

aws iam attach-role-policy \
  --role-name PackulatorCodeBuildRole \
  --policy-arn arn:aws:iam::aws:policy/CloudFrontFullAccess
```

### 2.2 Create CodeBuild Project

```bash
# Get your AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)

# Create CodeBuild project
aws codebuild create-project \
  --name packulator-build \
  --source '{
    "type": "GITHUB",
    "location": "https://github.com/YOUR_USERNAME/packulator.git",
    "buildspec": "buildspec.yml"
  }' \
  --artifacts '{
    "type": "NO_ARTIFACTS"
  }' \
  --environment '{
    "type": "LINUX_CONTAINER",
    "image": "aws/codebuild/amazonlinux2-x86_64-standard:4.0",
    "computeType": "BUILD_GENERAL1_MEDIUM",
    "privilegedMode": true
  }' \
  --service-role arn:aws:iam::${AWS_ACCOUNT_ID}:role/PackulatorCodeBuildRole
```

### 2.3 Set Environment Variables

```bash
# Set environment variables for CodeBuild
aws codebuild update-project \
  --name packulator-build \
  --environment '{
    "type": "LINUX_CONTAINER",
    "image": "aws/codebuild/amazonlinux2-x86_64-standard:4.0",
    "computeType": "BUILD_GENERAL1_MEDIUM",
    "privilegedMode": true,
    "environmentVariables": [
      {
        "name": "AWS_DEFAULT_REGION",
        "value": "'${AWS_REGION}'"
      },
      {
        "name": "AWS_ACCOUNT_ID",
        "value": "'${AWS_ACCOUNT_ID}'"
      },
      {
        "name": "IMAGE_REPO_NAME",
        "value": "packulator"
      },
      {
        "name": "DEPLOY_FRONTEND_TO_S3",
        "value": "true"
      },
      {
        "name": "S3_FRONTEND_BUCKET",
        "value": "packulator-prod-frontend"
      }
    ]
  }'
```

## Step 3: Deploy ECS Service

### 3.1 Get Infrastructure Outputs

```bash
# Get database endpoint
DB_HOST=$(aws cloudformation describe-stacks \
  --stack-name packulator-prod-infrastructure \
  --query 'Stacks[0].Outputs[?OutputKey==`DatabaseEndpoint`].OutputValue' \
  --output text)

# Get ECR repository URI
ECR_URI=$(aws cloudformation describe-stacks \
  --stack-name packulator-prod-infrastructure \
  --query 'Stacks[0].Outputs[?OutputKey==`ECRRepository`].OutputValue' \
  --output text)
```

### 3.2 Build and Push Initial Image

```bash
# Build and push initial Docker image
docker build -f Dockerfile.production -t packulator:latest .

# Tag for ECR
docker tag packulator:latest ${ECR_URI}:latest

# Login to ECR
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_URI}

# Push image
docker push ${ECR_URI}:latest
```

### 3.3 Deploy ECS Service

```bash
# Deploy ECS service stack
aws cloudformation create-stack \
  --stack-name packulator-prod-ecs \
  --template-body file://cloudformation/ecs-service.yml \
  --parameters ParameterKey=ApplicationName,ParameterValue=packulator \
               ParameterKey=Environment,ParameterValue=prod \
               ParameterKey=ImageURI,ParameterValue=${ECR_URI}:latest \
               ParameterKey=DBHost,ParameterValue=${DB_HOST} \
               ParameterKey=DBPassword,ParameterValue=YourSecurePassword123! \
  --capabilities CAPABILITY_IAM
```

## Step 4: Set up CI/CD Pipeline

### 4.1 Create GitHub Webhook (Optional)

```bash
# Start a build manually
aws codebuild start-build --project-name packulator-build
```

### 4.2 Monitor Build

```bash
# Get build logs
aws codebuild batch-get-builds \
  --ids $(aws codebuild list-builds-for-project \
    --project-name packulator-build \
    --query 'ids[0]' --output text) \
  --query 'builds[0].logs.groupName'
```

## Step 5: Configure Domain and SSL (Optional)

### 5.1 Get CloudFront Distribution ID

```bash
CF_DISTRIBUTION_ID=$(aws cloudformation describe-stacks \
  --stack-name packulator-prod-infrastructure \
  --query 'Stacks[0].Outputs[?OutputKey==`CloudFrontDistributionId`].OutputValue' \
  --output text)
```

### 5.2 Get ALB DNS Name

```bash
ALB_DNS=$(aws cloudformation describe-stacks \
  --stack-name packulator-prod-infrastructure \
  --query 'Stacks[0].Outputs[?OutputKey==`ALBDNSName`].OutputValue' \
  --output text)

echo "Your application is accessible at: http://${ALB_DNS}"
```

## Step 6: Verification and Testing

### 6.1 Health Check

```bash
# Test API health
curl http://${ALB_DNS}/health/check

# Test API endpoints
curl http://${ALB_DNS}/packs/list
```

### 6.2 Frontend Access

If deployed to CloudFront:
```bash
# Get CloudFront domain
CF_DOMAIN=$(aws cloudfront get-distribution \
  --id ${CF_DISTRIBUTION_ID} \
  --query 'Distribution.DomainName' \
  --output text)

echo "Frontend accessible at: https://${CF_DOMAIN}"
```

## Deployment Options

### Option 1: All-in-One (Backend serves Frontend)
- Frontend is built and embedded in the Go application
- Single endpoint serves both API and static files
- Simpler deployment, single container

### Option 2: Separate Frontend (S3 + CloudFront)
- Frontend deployed to S3 and served via CloudFront
- Backend only serves API endpoints
- Better performance and caching for static assets

## Environment Variables

The following environment variables are used:

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment name | `production` |
| `DB_HOST` | Database hostname | `packulator-prod-db.xxx.rds.amazonaws.com` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `packulator` |
| `DB_USER` | Database username | `packulator_user` |
| `DB_PASSWORD` | Database password | `SecurePassword123!` |

## Monitoring and Logs

### CloudWatch Logs
```bash
# View ECS logs
aws logs describe-log-groups --log-group-name-prefix /ecs/packulator

# Stream logs
aws logs tail /ecs/packulator-prod --follow
```

### Application Metrics
- ECS Service metrics available in CloudWatch
- Auto-scaling configured based on CPU utilization
- Health checks monitor application status

## Troubleshooting

### Common Issues

1. **Build Failures**
   - Check CodeBuild logs
   - Verify environment variables
   - Ensure ECR permissions

2. **ECS Task Failures**
   - Check CloudWatch logs
   - Verify database connectivity
   - Check environment variables

3. **Load Balancer Health Checks**
   - Ensure `/health/check` endpoint returns 200
   - Check security group rules
   - Verify container port mapping

### Debug Commands

```bash
# Check ECS service status
aws ecs describe-services \
  --cluster packulator-prod-cluster \
  --services packulator-prod-service

# Check task logs
aws logs tail /ecs/packulator-prod --follow

# Test database connectivity
aws rds describe-db-instances \
  --db-instance-identifier packulator-prod-db
```

## Cost Optimization

- Use Fargate Spot for non-production environments
- Configure auto-scaling to scale down during low usage
- Use CloudFront caching for static assets
- Consider RDS reserved instances for production

## Security Best Practices

- Store secrets in AWS Secrets Manager
- Use IAM roles instead of access keys
- Enable VPC Flow Logs
- Configure WAF for the load balancer
- Enable GuardDuty for threat detection

## Cleanup

To delete all resources:

```bash
# Delete ECS service stack
aws cloudformation delete-stack --stack-name packulator-prod-ecs

# Delete infrastructure stack
aws cloudformation delete-stack --stack-name packulator-prod-infrastructure

# Delete CodeBuild project
aws codebuild delete-project --name packulator-build

# Delete ECR repository
aws ecr delete-repository --repository-name packulator --force
```