# Create App: generate and install a minimal complete App scaffold.
resource "criblio_app" "create_app" {
  id           = "terraform-created-app"
  display_name = "Terraform Created App"
  version      = "1.0.0"
}

# Import from File: upload and install the archive included with this example.
resource "criblio_app" "import_from_file" {
  id       = "terraform-example-app"
  filename = abspath("${path.module}/terraform-example-app-1.0.0.tgz")

  depends_on = [criblio_app.create_app]
}

# Import from URL: install a complete hosted App archive.
resource "criblio_app" "import_from_url" {
  id     = "url-imported-app"
  source = "https://github.com/criblio/apm/releases/download/v0.10.0/apm-0.10.0.tgz"
  force  = true

  depends_on = [criblio_app.import_from_file]
}

# Import from Git: clone and install a complete App repository.
resource "criblio_app" "import_from_git" {
  id     = "git-repository-example"
  source = "git+https://github.com/criblapps/cribl-ai-o11y.git"

  depends_on = [criblio_app.import_from_url]
}
