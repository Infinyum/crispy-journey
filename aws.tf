terraform {
  backend "remote" {
    organization = "elabrom"
    workspaces {
      name = "crispy"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_key_pair" "id_rsa" {
  key_name   = "id_rsa"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCcYxZpvNsBy5W9Hnrp8hT/RMDCVv5ibtIVnj4co8lUUB1/A+b4G/VLVPjVAy7SessBEE26X7uDg3VSCk0WHHYlo+ah6HYRhcdX1ysskVnqQn5mEKX0u5wz81YeFgTNb7U5HMqcSmNb2/RL3rLMVfRr0yE5sseMtX7NjTILydJPsCoVk0+4jafF4425dsGRwimRSnfyXXj2Acc8UemppcUG/egSHmTdN0kog/lKsfegkjQdN01XJYlT0Sfm9jGoLMEUKQEIYP8iMfXCcqi7/lkfe74rvKMJ8pAt4UHqku2dRpGUT3rsO6iSQKrXpV1JFe/450Ibhe0+HagdYQajMYaBQ/2RSwrs58RrMwq5HMBFQ9vavipVmPJce80HufoSwg0jchVRdMA+EY9yaVgFmFpszQmEf4ldbfqGRvO2fWxZPyzCF7DtT3OmAl3bR4u2b4LS7DIcUpXL01yhNcRXlRFUs8sT6Fwch/07he+pYN/x1qRv5f7ioo46CnX5OGioDnNbarsu04pconnKXk+n9+8mOz7HQZ6Ckc+VjAKB4NiZ7LnH+EcjRvRBdgOM1n9te5POj/kidnNpROfLLOyx8lPjNQqvdEve4QlJV2HAMt89c9kFRYi/xywAZXdOpdJVea9u/Y/ccp9GJ4x83/k1v5/fUZ4fuDbgzXlteeJN8M4kRQ== robin@macbook-air"
}

resource "aws_security_group" "allow_ssh" {
  name = "allow_ssh"

  # Allow ssh to the EC2 instance
  ingress {
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
  }

  # Allow all external traffic from the EC2 instance to the internet
  # Useful when trying to install packages...
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "linux" {
  ami           = "ami-047a51fa27710816e" // Amazon Linux 2 x86_64
  instance_type = "t2.micro"

  security_groups = [aws_security_group.allow_ssh.name]
  key_name        = aws_key_pair.id_rsa.key_name
}

output "ip" {
  value = aws_instance.linux.public_ip
}