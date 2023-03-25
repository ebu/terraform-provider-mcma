resource "mcma_job_profile" "example" {
  name = "example"

  input_parameter {
    name = "param1"
    type = "string"
  }
  input_parameter {
    name     = "param2"
    type     = "number"
    optional = true
  }

  output_parameter {
    name = "outparam1"
    type = "string"
  }
  output_parameter {
    name = "outparam2"
    type = "number"
  }

  custom = {
    customprop1 = "customprop1val"
    customprop2 = "customprop2val"
  }
}