terraform {
  required_providers {
    beget = {
      source  = "hashicorp.com/edu/beget"
    }
  }
  required_version = ">= 1.1.0"
}

provider "beget" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}

resource "beget_server" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "edu_server" {
  value = beget_server.edu
}
