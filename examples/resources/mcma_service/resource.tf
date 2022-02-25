resource "mcma_service" "example" {
  name      = "example"
  auth_type = "AWS4"
  job_type  = "AmeJob"

  resource {
    resource_type = "JobAssignment"
    http_endpoint = "https://some.endpoint.com/api/job-assignments"
  }
  resource {
    resource_type = "JobAssignment"
    http_endpoint = "https://some.endpoint.com/api/job-assignments"
    auth_type     = "JWT"
  }

  job_profile_ids = [
    "https://service.registry.com/api/job-profiles/12345",
    "https://service.registry.com/api/job-profiles/67890"
  ]
}