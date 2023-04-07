terraform {
  required_providers {
    beget = {
      source = "hashicorp.com/edu/beget"
    }
  }
}

provider "beget" {
  username = "education"
  password = "test123"
}

data "beget_software" "example" {}
