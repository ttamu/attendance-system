resource "aws_security_group" "nat_sg" {
  name        = "${var.project_name}-nat-sg"
  description = "Security group for NAT instance"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [var.private_subnet_cidr_a, var.private_subnet_cidr_c]
  }
  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
    cidr_blocks = [var.admin_ip]
  }
  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "nat" {
  most_recent = true
  owners = ["amazon"]

  filter {
    name = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

resource "aws_instance" "nat" {
  ami                         = data.aws_ami.nat.id
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.public_a.id
  associate_public_ip_address = true
  security_groups = [aws_security_group.nat_sg.id]
  key_name                    = var.key_pair_name
  source_dest_check           = false

  user_data = <<-EOF
    #!/bin/bash
    echo 1 > /proc/sys/net/ipv4/ip_forward
    echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
    sysctl -p
    iptables -t nat -A POSTROUTING -o eth0 -s ${var.private_subnet_cidr_a} -j MASQUERADE
    iptables -t nat -A POSTROUTING -o eth0 -s ${var.private_subnet_cidr_c} -j MASQUERADE
    yum install -y iptables-services
    service iptables save
    service iptables restart
  EOF

  tags = {
    Name    = "${var.project_name}-nat"
    Project = var.project_name
  }
}
