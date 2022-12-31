terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name        = "main-vpc"
    Environment = "production"
  }
}

resource "aws_cloudwatch_log_group" "main" {
  name = "main-service-logs"

  tags = {
    Environment = "production"
    Application = "main"
  }
}

# Internet Gateway for connecting to internet
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "main-igw"
    Environment = "production"
  }
}

# Public subnet 
resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id

  cidr_block              = "10.0.0.0/20"
  map_public_ip_on_launch = true

  tags = {
    Name        = "main-public-subnet-1"
    Environment = "production"
  }
}

# Routing table for public subnet
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "public-route-table-1"
    Environment = "production"
  }
}

resource "aws_route" "public_internet_gateway" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw.id
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

#  Private subnet
resource "aws_subnet" "private" {
  vpc_id = aws_vpc.main.id

  cidr_block              = "10.0.128.0/20"
  map_public_ip_on_launch = false

  tags = {
    Name        = "private-subnet-1"
    Environment = "production"
  }
}

# Routing table for private subnet
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "private-route-table-1"
    Environment = "production"
  }
}

resource "aws_route_table_association" "private" {
  subnet_id      = aws_subnet.private.id
  route_table_id = aws_route_table.private.id
}

# default security group
resource "aws_security_group" "ecs_sg" {
  name        = "ecs-sg"
  description = "ecs security group to allow inbound/outbound from the VPC"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol         = "tcp"
    from_port        = 8080
    to_port          = 8080
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name        = "ecs-sg"
    Environment = "production"
  }
}

# ECR container registry
resource "aws_ecr_repository" "main" {
  name = "main"
}

# IAM task role
resource "aws_iam_role" "ecs_task_role" {
  name = "ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        },
        Effect = "Allow",
        Sid    = ""
      }
    ]
  })

  tags = {
    name        = "ecs-task-role"
    environment = "production"
  }
}

# DynamoDB table for authentication
resource "aws_dynamodb_table" "authentication_dynamodb_table" {
  name           = "Authentication"
  hash_key       = "PK"
  read_capacity  = 10
  write_capacity = 10

  attribute {
    name = "PK"
    type = "S"
  }

  tags = {
    Name        = "Authentication"
    Environment = "production"
  }
}

# DynamoDB task policy for table access
resource "aws_iam_policy" "dynamodb_task_policy" {
  name        = "dynamodb-task-policy"
  description = "Policy that allows access to DynamoDB"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "dynamodb:CreateTable",
          "dynamodb:UpdateTimeToLive",
          "dynamodb:PutItem",
          "dynamodb:DescribeTable",
          "dynamodb:ListTables",
          "dynamodb:DeleteItem",
          "dynamodb:GetItem",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:UpdateItem",
          "dynamodb:UpdateTable"
        ],
        Resource = aws_dynamodb_table.authentication_dynamodb_table.arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs-task-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.dynamodb_task_policy.arn
}

# IAM task execution role
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        },
        Effect = "Allow",
        Sid    = ""
      }
    ]
  })

  tags = {
    name        = "ecs-task-execution-role"
    environment = "production"
  }
}

resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# AWS ECS task definition
resource "aws_ecs_task_definition" "task" {
  family = "golangsocial"

  network_mode = "awsvpc"
  cpu          = 256
  memory       = 512

  requires_compatibilities = ["FARGATE"]

  container_definitions = jsonencode([
    {
      name  = "golangsocial-container-production"
      image = "${aws_ecr_repository.main.repository_url}:latest"

      environment = [
        {
          name  = "AWS_REGION"
          value = "us-east-1"
        },
        {
          name  = "DYNAMODB_ENDPOINT"
          value = "https://dynamodb.us-east-1.amazonaws.com"
        }
      ]

      essential = true
      portMappings = [
        {
          protocol      = "tcp"
          containerPort = 8080
          hostPort      = 8080
        }
      ]

      logConfiguration = {
        logDriver = "awslogs",
        options = {
          awslogs-group         = "${aws_cloudwatch_log_group.main.id}",
          awslogs-region        = "us-east-1",
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  execution_role_arn = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn      = aws_iam_role.ecs_task_role.arn
}

resource "aws_ecs_cluster" "main" {
  name = "golangsocial-cluster-production"
}

resource "aws_ecs_service" "service" {
  name            = "golangsocial-service-production"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.task.id
  desired_count   = 1

  launch_type         = "FARGATE"
  scheduling_strategy = "REPLICA"

  network_configuration {
    security_groups  = [aws_security_group.ecs_sg.id]
    subnets          = [aws_subnet.public.id]
    assign_public_ip = true
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}
