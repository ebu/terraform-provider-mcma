# No auth
provider "mcma" {
  service_registry_url = "https://service-registry-example.mcma.io/api/"
}

# AWS auth with profile
provider "mcma" {
  service_registry_url = "https://service-registry-example.mcma.io/api/"
  aws4_auth {
    region  = "us-east-1"
    profile = "myprofile"
  }
}

# AWS auth with keys
provider "mcma" {
  service_registry_url = "https://service-registry-example.mcma.io/api/"
  aws4_auth {
    region     = "us-east-1"
    access_key = "accesskey"
    secret_key = "secretkey"
  }
}

# MCMA API Key auth
provider "mcma" {
  service_registry_url = "https://service-registry-example.mcma.io/api/"
  mcma_api_key_auth {
    api_key = "abcd1234efgh5678"
  }
}