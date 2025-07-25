variable "public_subnets_cidrs" {
  type = list(string)
  default = [ "10.0.1.0/24", "10.0.2.0/24" ]
}

variable "private_subnets_cidrs" {
  type    = list(string)
  default = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "azs" {
  type = list(string)
  default = [ "us-east-1a", "us-east-1b" ]
}

variable "app_name" {
  default = "sample-django"
}

variable "project_name" {
  default = "django"
  type = string
}
