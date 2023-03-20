data "aws_availability_zones" "available" {
  state = "available"
}

resource "random_shuffle" "az" {
  input = tolist(data.aws_availability_zones.available.names)
}

# Create a single VPC with DNS support.
resource "aws_vpc" "platform" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = merge(
    local.tags,
    { "Name" = local.name }
  )
}

# Create 2 public subnets.
resource "aws_subnet" "platform_public" {
  count = 2

  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 1)
  availability_zone       = element(random_shuffle.az.result, count.index)
  map_public_ip_on_launch = true

  tags = local.tags
}

# Create 2 private subnets.
resource "aws_subnet" "platform_private" {
  count = 2

  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 3)
  availability_zone       = element(random_shuffle.az.result, count.index)
  map_public_ip_on_launch = true

  tags = local.tags
}

# Create an IGW to enable internet-routable traffic.
resource "aws_internet_gateway" "platform" {
  vpc_id = aws_vpc.platform.id

  tags = local.tags
}

# Create a single public route table to direct network traffic.
resource "aws_route_table" "platform_public" {
  vpc_id = aws_vpc.platform.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.platform.id
  }

  tags = local.tags
}

# Create 2 privte route tables to direct network traffic.
resource "aws_route_table" "platform_private" {
  vpc_id = aws_vpc.platform.id

  tags = local.tags
}

# Associate the route table with the public subnets.
resource "aws_route_table_association" "platform_public" {
  count = 2

  subnet_id      = element(aws_subnet.platform_public.*.id, count.index)
  route_table_id = aws_route_table.platform_public.id
}

# Associate the route table with the privte subnets.
resource "aws_route_table_association" "platform_private" {
  count = 2

  subnet_id      = element(aws_subnet.platform_private.*.id, count.index)
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

  tags = local.tags
}
