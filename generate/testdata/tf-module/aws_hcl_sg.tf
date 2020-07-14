# ALB

resource "aws_lb" "front" {
  name            = "some name"
  security_groups = ["${aws_security_group.lb-front.id}"]

  tags = {
    Name = "name"
    role = "front"
  }
}

resource "aws_security_group" "lb-front" {
  name        = "some name"
  description = "Front "

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "name"
    role = "front"
  }
}

# Launch Template

resource "aws_launch_template" "front" {
  name_prefix = "name"

  network_interfaces {
    security_groups = [
      "${aws_security_group.front.id}",
    ]
  }
  lifecycle {
    create_before_destroy = true
  }
  tags = {
    Name = "name"
    role = "fronttemplate"
  }
}

resource "aws_security_group" "front" {
  name        = "anem"
  description = "Front"

  # Allow to get myeasyapi nginx, openapi nginx, mypages nginx
  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = ["${aws_security_group.lb-front.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "name"
    role = "front"
  }
}

# RDS

resource "aws_security_group" "rds" {
  name        = "anme"
  description = "rds"

  ingress {
    from_port = 3306
    to_port   = 3306
    protocol  = "tcp"

    security_groups = ["${aws_security_group.front.id}"]
  }

  tags = {
    Name = "name"
    role = "rds"
  }
}

resource "aws_db_instance" "application" {
  identifier = "rds"

  vpc_security_group_ids = ["${aws_security_group.rds.id}"]

  tags = {
    Name = "name"
    type = "master"
    role = "rds"
  }
}

# Redis

resource "aws_security_group" "redis" {
  name        = "name"
  description = "desc"

  ingress {
    from_port = 3306
    to_port   = 3306
    protocol  = "tcp"

    security_groups = [
      "${aws_security_group.front.id}",
    ]
  }

  tags = {
    Name = "name"
    role = "redis"
  }
}

resource "aws_elasticache_cluster" "redis" {
  security_group_ids = ["${aws_security_group.redis.id}"]

  tags = {
    Name = "name"
    role = "redis"
  }
}
