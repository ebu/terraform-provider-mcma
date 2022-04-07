# No auth
provider "mcma" {
  services_url = "https://service-registry-example.mcma.io/api/services"
}

# AWS auth with profile
provider "mcma" {
  services_url = "https://service-registry-example.mcma.io/api/services"
  aws4_auth {
    region  = "us-east-1"
    profile = "myprofile"
  }
}

# AWS auth with keys
provider "mcma" {
  services_url = "https://service-registry-example.mcma.io/api/services"
  aws4_auth {
    region     = "us-east-1"
    access_key = "accesskey"
    secret_key = "secretkey"
  }
}