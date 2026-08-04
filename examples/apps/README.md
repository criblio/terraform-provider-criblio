# App example

The configuration maps the four **Add App** UI options to Terraform:

- **Create App**: omit both `source` and `filename`; the provider generates and
  uploads a minimal scaffold containing `static/index.html`.
- **Import from File**: set `filename` to a local `.crbl` or `.tgz` archive.
- **Import from URL**: set `source` to a direct artifact URL.
- **Import from Git**: set `source` to a `git+https://...` repository URL.

`fixture` is a minimal complete Cribl App. Its archive must contain
`package.json` and `static/index.html` at the archive root.

Create an uploadable archive:

```shell
cd fixture
tar -czf ../terraform-example-app-1.0.0.tgz package.json static
```

The example includes the resulting `terraform-example-app-1.0.0.tgz` archive.
The provider uploads it with `PUT /apps` and installs the returned `source` in
the same apply:

```hcl
resource "criblio_app" "local_archive" {
  id       = "terraform-example-app"
  filename = abspath("${path.module}/terraform-example-app-1.0.0.tgz")
}
```
