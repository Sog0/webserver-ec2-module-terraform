provider "aws" {
  region = "us-east-1"
}

terraform {
  backend "s3" {
    bucket         = "bucketformyneeds1"
    key            = "webserver/terraform.tfstate"
    region         = "us-east-1"
    use_lockfile = true
  }
}