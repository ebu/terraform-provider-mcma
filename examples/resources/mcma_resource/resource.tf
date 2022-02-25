resource "mcma_resource" "bm_content" {
  type = "BMContent"
  resource_json = jsonencode({
    metadata = {
      name        = "Terraform provider test BMContent"
      description = "Test asset"
    }
  })
}