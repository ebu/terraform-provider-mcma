provider "mcma" {
  url = "https://service-registry-example.mcma.io/services"
  auth {
    type = "AWS4"
    data = {
      accessKey = var.access_key
      secretKey = var.secret_key
    }
  }
}