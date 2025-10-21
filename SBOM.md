# Software Bill of Materials (SBOM)

This document describes the Software Bill of Materials (SBOM) implementation for the terraform-provider-criblio project.

## Overview

A Software Bill of Materials (SBOM) is a comprehensive inventory of all components within our software, including open-source libraries and dependencies. This enhances transparency, aids in vulnerability management, and ensures compliance with licensing requirements.

## Generated SBOM Files

The following SBOM files are automatically generated:

- **`sbom-spdx.json`** - SPDX format (ISO/IEC 5962:2021 standard)
- **`sbom-cyclonedx.json`** - CycloneDX format (industry standard)

## Generation Process

### Automatic Generation

SBOMs are automatically generated during:

1. **Release Process** - When creating a new release (tagged with `v*`)
2. **Pull Requests** - On every PR to main branch
3. **Main Branch Updates** - On every push to main branch
4. **Weekly Schedule** - Every Sunday at midnight UTC

### Manual Generation

You can manually trigger SBOM generation by:

1. **GitHub Actions UI** - Go to Actions → "SBOM Generation" → Run workflow
2. **API Call** - Use GitHub's dependency graph API
3. **Local Generation** - Use the tools directly (see below)

## Tools Used

### Primary Tools

- **[Anchore Syft](https://github.com/anchore/syft)** - Comprehensive SBOM generation
- **[GitHub Dependency Graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph)** - Native GitHub integration

### Supported Formats

- **SPDX** (Software Package Data Exchange) - ISO/IEC 5962:2021 standard
- **CycloneDX** - OWASP standard for software supply chain security

## Integration Points

### GitHub Features

- **Security Tab** → Dependency graph
- **Insights Tab** → Dependency graph → Export SBOM
- **Dependabot Alerts** - Automatic vulnerability scanning
- **Code Scanning** - Security vulnerability detection

### Release Assets

SBOM files are automatically attached to GitHub releases as:
- `terraform-provider-criblio_{version}_sbom.json`

## Vulnerability Scanning

The SBOM generation process includes:

- **Trivy Scanner** - Filesystem vulnerability scanning
- **SARIF Upload** - Results uploaded to GitHub Security tab
- **Dependency Analysis** - Go module dependency scanning

## Compliance Benefits

This SBOM implementation provides:

- ✅ **Transparency** - Clear visibility into all dependencies
- ✅ **Security** - Proactive vulnerability management
- ✅ **Compliance** - Meet regulatory requirements (NTIA, CISA guidelines)
- ✅ **Supply Chain Security** - Enhanced security posture
- ✅ **License Management** - Track open-source license compliance

## Workflow Files

- **`.github/workflows/sbom-generation.yml`** - Dedicated SBOM generation
- **`.github/workflows/tf_provider_release.yaml`** - Release workflow with SBOM
- **`.goreleaser.yml`** - GoReleaser configuration for SBOM inclusion

## Local Development

To generate SBOMs locally:

```bash
# Install Syft
go install github.com/anchore/syft/cmd/syft@latest

# Generate SPDX SBOM
syft packages . -o spdx-json -f sbom-spdx.json

# Generate CycloneDX SBOM
syft packages . -o cyclonedx-json -f sbom-cyclonedx.json
```

## Monitoring

- **GitHub Security Tab** - View dependency vulnerabilities
- **Actions Artifacts** - Download SBOM files from workflow runs
- **Release Assets** - Access SBOM files from GitHub releases

## Maintenance

The SBOM generation is fully automated and requires no manual maintenance. The system will:

- Automatically detect new dependencies
- Update vulnerability information
- Generate fresh SBOMs on every build
- Upload results to GitHub's security features

## Support

For questions or issues related to SBOM generation:

1. Check the GitHub Actions logs
2. Review the Security tab for dependency alerts
3. Consult the [Anchore Syft documentation](https://github.com/anchore/syft)
4. Review [GitHub's dependency graph documentation](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph)
