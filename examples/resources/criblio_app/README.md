# App resource examples

This directory demonstrates the four **Add App** workflows without variables:

- **Create App** generates and uploads a minimal complete scaffold.
- **Import from File** uploads the included `.tgz` archive.
- **Import from URL** downloads a hardcoded public App archive.
- **Import from Git** clones a hardcoded complete App repository.

Cribl.Cloud supports all four workflows. On-prem deployments must use URL or
Git because the App upload endpoint is unavailable there.

The resources are serialized because concurrent App installations can fail on
some deployments. A valid App archive must contain `package.json` and
`static/index.html` at its root.
