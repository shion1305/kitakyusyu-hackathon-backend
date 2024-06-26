remote_state {
  backend = "gcs"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
  config = {
    project = "kyusyu-hackathon"
    bucket = "kitakyushu-hackathon-terraform"
    prefix = "${path_relative_to_include()}/terraform.tfstate"
    location = "asia-northeast1"
  }
}

terraform {
  extra_arguments "identity" {
    commands = [
      "apply",
      "plan",
      "import",
      "destroy",
      "init"
    ]

    env_vars = {
      GOOGLE_IMPERSONATE_SERVICE_ACCOUNT = "gh-terraform@kyusyu-hackathon.iam.gserviceaccount.com"
    }
  }
}

inputs = {
  environment = "dev"
}

generate "provider" {
  path      = "_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = file("../../shared/provider.tf")
}

generate "version" {
  path      = "_version.tf"
  if_exists = "overwrite_terragrunt"
  contents  = file("../../shared/version.tf")
}

generate "variables" {
  path      = "_variables.tf"
  if_exists = "overwrite_terragrunt"
  contents  = file("../../shared/variables.tf")
}
