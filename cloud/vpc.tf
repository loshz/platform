data "aws_availability_zones" "available" {
  state = "available"
}

# Create a single VPC with DNS support.
resource "aws_vpc" "platform" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = { "Name" = local.name }
}

# Create a public subnet.
resource "aws_subnet" "platform_public" {
  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, 0)
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = { "Name" = "public1-${data.aws_availability_zones.available.names[0]}" }
}

# Create a private subnet.
resource "aws_subnet" "platform_private" {
  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, 128)
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = { "Name" = "private1-${data.aws_availability_zones.available.names[0]}" }
}

# Create an IGW to enable internet-routable traffic.
resource "aws_internet_gateway" "platform" {
  vpc_id = aws_vpc.platform.id
}

# Create a single public route table to direct network traffic.
resource "aws_route_table" "platform_public" {
  vpc_id = aws_vpc.platform.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.platform.id
  }
}

# Create 2 privte route tables to direct network traffic.
resource "aws_route_table" "platform_private" {
  vpc_id = aws_vpc.platform.id
}

# Associate the route table with the public subnets.
resource "aws_route_table_association" "platform_public" {
  subnet_id      = aws_subnet.platform_public.id
  route_table_id = aws_route_table.platform_public.id
}

# Associate the route table with the privte subnets.
resource "aws_route_table_association" "platform_private" {
  subnet_id      = aws_subnet.platform_private.id
  route_table_id = aws_route_table.platform_private.id
}

# Create a security group with public egress.
resource "aws_security_group" "platform" {
  vpc_id = aws_vpc.platform.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}
