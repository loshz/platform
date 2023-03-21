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

# Create 2 public subnets.
resource "aws_subnet" "public" {
  count = 2

  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 1)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = { "Name" = "public${count.index + 1}-${data.aws_availability_zones.available.names[count.index]}" }
}

# Create 2 private subnets.
resource "aws_subnet" "private" {
  count = 2

  vpc_id                  = aws_vpc.platform.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + 3)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = false

  tags = { "Name" = "private${count.index + 1}-${data.aws_availability_zones.available.names[count.index]}" }
}

# Create an IGW to enable internet-routable traffic.
resource "aws_internet_gateway" "public" {
  vpc_id = aws_vpc.platform.id
}

# Create a single public route table to direct network traffic.
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.platform.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.public.id
  }

  tags = { "Name" = "public" }
}

# Create 2 privte route tables to direct network traffic.
resource "aws_route_table" "private" {
  count = 2

  vpc_id = aws_vpc.platform.id

  tags = { "Name" = "private${count.index + 1}-${data.aws_availability_zones.available.names[count.index]}" }
}

# Associate the route table with the public subnets.
resource "aws_route_table_association" "public" {
  count = 2

  subnet_id      = element(aws_subnet.public.*.id, count.index)
  route_table_id = aws_route_table.public.id
}

# Associate the route table with the privte subnets.
resource "aws_route_table_association" "private" {
  count = 2

  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}

# Create a security group with public ingress/egress.
resource "aws_security_group" "load_balancer" {
  vpc_id = aws_vpc.platform.id

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
