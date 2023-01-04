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
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-east-1f"

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

# Public subnet 2
resource "aws_subnet" "public_2" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-east-1b"

  cidr_block              = "10.0.16.0/20"
  map_public_ip_on_launch = true

  tags = {
    Name        = "main-public-subnet-2"
    Environment = "production"
  }
}

# Routing table for public subnet
resource "aws_route_table" "public_2" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "public-route-table-2"
    Environment = "production"
  }
}

resource "aws_route" "public_internet_gateway_2" {
  route_table_id         = aws_route_table.public_2.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw.id
}

resource "aws_route_table_association" "public_2" {
  subnet_id      = aws_subnet.public_2.id
  route_table_id = aws_route_table.public_2.id
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
resource "aws_security_group" "alb_sg" {
  name        = "alb-sg"
  description = "alb security group to allow inbound/outbound from the VPC"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol         = "tcp"
    from_port        = 80
    to_port          = 80
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    protocol         = "tcp"
    from_port        = 443
    to_port          = 443
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
}

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

# application load balancer
resource "aws_lb" "main" {
  name               = "main-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public.id, aws_subnet.public_2.id]

  enable_deletion_protection = false
}

resource "aws_alb_target_group" "main" {
  name        = "main-alb-tg"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id
}

resource "aws_alb_listener" "http" {
  load_balancer_arn = aws_lb.main.id
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = 443
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_alb_listener" "https" {
  load_balancer_arn = aws_lb.main.id
  port              = 443
  protocol          = "HTTPS"

  ssl_policy      = "ELBSecurityPolicy-2016-08"
  certificate_arn = "arn:aws:acm:us-east-1:582250362323:certificate/121dbb1a-53c6-41f4-8325-d7b3ee58dca3"

  default_action {
    target_group_arn = aws_alb_target_group.main.id
    type             = "forward"
  }
}

# ECR container registry
resource "aws_ecr_repository" "main" {
  name = "main"
}

resource "aws_ecr_lifecycle_policy" "main" {
  repository = aws_ecr_repository.main.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "keep last 10 images"
      action = {
        type = "expire"
      }
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 5
      }
    }]
  })
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
  name           = "AUTHENTICATION"
  hash_key       = "PK"
  read_capacity  = 10
  write_capacity = 10

  attribute {
    name = "PK"
    type = "S"
  }

  tags = {
    Name        = "AUTHENTICATION"
    Environment = "production"
  }
}

# DynamoDB table for authentication
resource "aws_dynamodb_table" "application_dynamodb_table" {
  name           = "APPLICATION"
  hash_key       = "PK"
  range_key      = "SK"
  read_capacity  = 10
  write_capacity = 10

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }
  
  attribute {
    name = "GSI1"
    type = "S"
  }

  global_secondary_index {
    name               = "APPLICATION_GSI_1"
    hash_key           = "GSI1"
    range_key          = "SK"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "ALL"
  }

  tags = {
    Name        = "APPLICATION"
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
        Resource = [
            aws_dynamodb_table.authentication_dynamodb_table.arn,
            aws_dynamodb_table.application_dynamodb_table.arn,
        ]
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
        },
        {
          name  = "GIN_MODE"
          value = "release"
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
    subnets          = [aws_subnet.public.id, aws_subnet.public_2.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.main.arn
    container_name   = "golangsocial-container-production"
    container_port   = 8080
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}

# DNS domain 'terraform import aws_route53_zone.main <ZONE_ID>'
resource "aws_route53_zone" "main" {
  name = "amuel.org"
}

# DNS routing from domain to alb
resource "aws_route53_record" "main" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "api.amuel.org"
  type    = "A"

  alias {
    name                   = aws_lb.main.dns_name
    zone_id                = aws_lb.main.zone_id
    evaluate_target_health = true
  }
}
