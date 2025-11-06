# CLI Tool Design Specification

**Document Version:** 1.0  
**Date:** 2025-01-XX  
**Author:** Cribl Infrastructure Team  
**Related Issues:** INFRA-7290, INFRA-7432

## Executive Summary

This document outlines the design specification for a CLI tool that imports existing Cribl configuration from the `/system/diag/download` endpoint and converts it into Terraform modules. The tool will enable users to migrate their existing Cribl configurations to Infrastructure as Code (IaC) using Terraform.

## 1. End-User Usage Perspective

### 1.1 Tool Purpose

The `criblio-import` CLI tool allows users to:
- Download existing Cribl configuration from a running instance
- Automatically convert YAML/JSON configurations to Terraform HCL format
- Generate organized Terraform modules ready for version control
- Migrate from manual configuration management to Infrastructure as Code

### 1.2 Target Users

**Primary Users:**
- DevOps Engineers managing Cribl infrastructure
- Platform Engineers implementing IaC practices
- Infrastructure teams standardizing configuration management

**Use Cases:**
1. **Initial Migration**: Import existing Cribl configuration to Terraform for the first time
2. **Workspace Sync**: Import configuration from one workspace to another
3. **Backup & Recovery**: Generate Terraform from production for disaster recovery
4. **Environment Parity**: Ensure staging/test environments match production

### 1.3 User Workflow

#### Step 1: Authentication Setup
Users authenticate using one of the supported methods:

```bash
# Option 1: Environment Variables (Recommended)
export CRIBL_BEARER_TOKEN="your-token"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_ORGANIZATION_ID="org"

# Option 2: OAuth Credentials
export CRIBL_CLIENT_ID="client-id"
export CRIBL_CLIENT_SECRET="client-secret"
export CRIBL_ORGANIZATION_ID="org"
export CRIBL_WORKSPACE_ID="main"

# Option 3: Credentials File (~/.cribl/credentials)
# Configured once, reused across sessions
```

#### Step 2: Run Import Command
```bash
# Basic usage - imports all resources
criblio-import --output ./terraform-configs

# Selective import - only specific resource types
criblio-import --output ./terraform-configs \
  --include sources,destinations

# Preview mode - see what would be generated
criblio-import --output ./preview --dry-run

# With sensitive data handling
criblio-import --output ./terraform-configs \
  --var-file variables.tf
```

#### Step 3: Review Generated Files
```bash
cd terraform-configs
tree
# .
# ├── main.tf              # Provider configuration
# ├── variables.tf          # Sensitive variables (optional)
# └── modules/
#     ├── sources/
#     │   └── main.tf      # All source resources
#     ├── destinations/
#     │   └── main.tf      # All destination resources
#     └── ...
```

#### Step 4: Initialize and Apply
```bash
terraform init
terraform plan    # Review changes
terraform apply    # Apply if satisfied
```

### 1.4 Command-Line Interface

#### Command Structure
```bash
criblio-import [command] [flags]
```

#### Commands
- `import` - Import configuration (default command)
- `version` - Show version information
- `help` - Show help message

#### Flags

**Output Options:**
- `--output, -o` (required): Output directory for Terraform files
- `--format, -f`: Output format - `resources` or `modules` (default: `modules`)
- `--var-file, -v`: Generate variables.tf for sensitive fields

**Filtering Options:**
- `--include, -i`: Only include specified resource types (comma-separated)
- `--exclude, -e`: Exclude specified resource types (comma-separated)

**Execution Options:**
- `--dry-run, -d`: Preview what would be generated without writing files
- `--verbose`: Enable verbose logging

**Authentication Options:**
- `--bearer-token`: Bearer token for authentication
- `--client-id`: OAuth client ID
- `--client-secret`: OAuth client secret
- `--workspace-id`: Cribl workspace ID
- `--organization-id`: Cribl organization ID
- `--onprem-server-url`: On-prem server URL
- `--onprem-username`: On-prem username
- `--onprem-password`: On-prem password

### 1.5 Example Usage Scenarios

#### Scenario 1: First-Time Import
```bash
# User wants to import all configuration from production
criblio-import --output ./prod-terraform

# Review generated files
cd prod-terraform
terraform init
terraform plan

# Apply to create state
terraform apply
```

#### Scenario 2: Selective Import
```bash
# User only wants to import sources and destinations
criblio-import --output ./configs \
  --include sources,destinations

# This generates only source and destination resources
```

#### Scenario 3: On-Prem Deployment
```bash
# Import from on-prem Cribl instance
criblio-import --output ./configs \
  --onprem-server-url http://cribl-server:9000 \
  --onprem-username admin \
  --onprem-password secure-password
```

#### Scenario 4: Preview Before Import
```bash
# User wants to see what would be generated
criblio-import --output ./preview --dry-run

# Output shows:
# - Resource count
# - File structure
# - Sample resource definitions
# - Warnings or issues
```

### 1.6 Output Structure

The tool generates a well-organized Terraform structure:

```
output-directory/
├── main.tf                 # Provider configuration
├── variables.tf            # Input variables (if --var-file specified)
├── README.md               # Generated documentation
└── modules/
    ├── sources/
    │   ├── main.tf         # All source resources
    │   ├── variables.tf    # Module variables
    │   └── outputs.tf       # Module outputs
    ├── destinations/
    │   └── ...
    ├── pipelines/
    │   └── ...
    ├── routes/
    │   └── ...
    └── packs/
        └── ...
```

### 1.7 Error Handling and User Feedback

**Success Indicators:**
- Progress bar showing download/extraction progress
- Summary of resources imported
- File count and size information
- Next steps guidance

**Error Messages:**
- Clear authentication errors with resolution steps
- API endpoint errors with troubleshooting tips
- File system errors with permission guidance
- Validation errors showing which resources failed

**Verbose Mode:**
- Detailed logging of each step
- Resource-by-resource conversion status
- Field mapping warnings
- Skipped resources with reasons

## 2. Programming Language Selection and Justification

### 2.1 Selected Language: **Go (Golang)**

### 2.2 Justification

#### 2.2.1 Code Reusability ⭐ (Critical)
**Why it matters:** The Terraform provider is already written in Go with a complete SDK.

**Benefits:**
- **Direct SDK Import**: Can import `github.com/criblio/terraform-provider-criblio/internal/sdk` directly
- **Shared Types**: Reuse provider types, models, and validation logic
- **Authentication Logic**: **Automatic authentication via `CriblTerraformHook`** - no reimplementation needed
- **No Translation Layer**: Avoid converting between languages or reimplementing API clients

**Evidence:**
- Provider SDK located at `internal/sdk/` with complete API coverage
- **Authentication automatically handled by `CriblTerraformHook`** (registered in `internal/sdk/internal/hooks/registration.go`)
- Hook handles: Bearer token, OAuth, on-prem auth, credentials file, URL construction
- Resource schemas in `internal/provider/*_resource.go` provide mapping reference

#### 2.2.2 Deployment Simplicity ⭐ (Critical)
**Why it matters:** End users need a simple, reliable installation process.

**Benefits:**
- **Single Static Binary**: No runtime dependencies, no package managers needed
- **Cross-Platform**: One build per platform (Linux, macOS, Windows)
- **No Version Conflicts**: Self-contained executable
- **Easy Distribution**: Direct download from GitHub Releases

**Alternative Issues:**
- Python: Requires Python runtime, virtual environments, dependency management
- Node.js: Requires Node.js runtime, npm/yarn, potential version conflicts
- Java: Requires JVM, larger download size

#### 2.2.3 Performance ⭐ (Important)
**Why it matters:** Users may have large configurations with hundreds of resources.

**Benefits:**
- **Fast Compilation**: Quick build times for development
- **Fast Execution**: Efficient parsing and file I/O
- **Low Memory Footprint**: Suitable for CI/CD environments
- **Concurrent Processing**: Native goroutines for parallel resource conversion

**Performance Comparison:**
- Go: ~10-50MB binary, fast startup, efficient memory usage
- Python: Runtime overhead, slower startup, higher memory usage
- Node.js: Runtime overhead, slower for CPU-intensive tasks

#### 2.2.4 CLI Tool Ecosystem ⭐ (Important)
**Why it matters:** Rich ecosystem for building professional CLI tools.

**Benefits:**
- **Cobra**: Industry-standard CLI framework (used by Kubernetes, Docker, etc.)
- **Viper**: Powerful configuration management (env vars, files, flags)
- **HCL Support**: First-class HCL generation with `github.com/hashicorp/hcl/v2`
- **Color Output**: Rich terminal output with `github.com/fatih/color`

**Libraries Available:**
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `github.com/hashicorp/hcl/v2/hclwrite` - HCL generation
- `gopkg.in/yaml.v3` - YAML parsing
- `archive/tar`, `compress/gzip` - Archive handling (stdlib)

#### 2.2.5 Team Expertise ⭐ (Important)
**Why it matters:** Faster development and maintenance.

**Benefits:**
- **Existing Codebase**: Team already knows Go from provider development
- **Consistent Stack**: Same language across provider and tools
- **Code Review**: Easier for team to review Go code
- **Maintenance**: Easier long-term maintenance with consistent stack

#### 2.2.6 Terraform Integration ⭐ (Important)
**Why it matters:** Tight integration with Terraform ecosystem.

**Benefits:**
- **HCL Generation**: Native HCL library support
- **Terraform Formatting**: Can invoke `terraform fmt` or use HCL libraries
- **Schema Understanding**: Can reference provider schemas directly
- **Validation**: Can validate generated HCL against provider schemas

### 2.3 Alternative Languages Considered

#### Python
**Pros:**
- Rapid prototyping
- Rich ecosystem
- Good for data processing

**Cons:**
- ❌ Runtime dependency management
- ❌ No direct SDK reuse (would need API translation)
- ❌ Slower execution
- ❌ Deployment complexity (virtual environments, dependencies)

**Verdict:** Not suitable - deployment complexity and lack of code reuse outweigh benefits.

#### Node.js/TypeScript
**Pros:**
- Good for web-based tools
- Rich package ecosystem

**Cons:**
- ❌ Runtime overhead
- ❌ No SDK reuse
- ❌ Deployment complexity
- ❌ Not ideal for system tools

**Verdict:** Not suitable - wrong tool for the job.

#### Rust
**Pros:**
- Excellent performance
- Single binary deployment
- Memory safety

**Cons:**
- ❌ Steeper learning curve
- ❌ No SDK reuse (would need FFI or reimplementation)
- ❌ Longer development time
- ❌ Smaller ecosystem for Terraform/HCL

**Verdict:** Overkill - Go provides sufficient performance with better code reuse.

### 2.4 Conclusion on Language Selection

**Go is the optimal choice** because:
1. ✅ **Maximum code reuse** - Direct SDK import, shared authentication
2. ✅ **Simplest deployment** - Single static binary
3. ✅ **Team familiarity** - Same language as provider
4. ✅ **Performance** - Fast enough for all use cases
5. ✅ **Ecosystem** - Excellent CLI and HCL libraries

**Trade-offs:**
- Slightly longer development time than Python (but faster overall due to code reuse)
- Static typing requires more upfront work (but catches errors early)

## 3. Deployment Process and End-User Presentation

### 3.1 Deployment Methods

#### Method 1: GitHub Releases (Primary) ⭐

**How it works:**
- Automated builds via GitHub Actions
- Pre-built binaries for multiple platforms
- Versioned releases with release notes

**User Experience:**
```bash
# Download for macOS
curl -LO https://github.com/criblio/terraform-provider-criblio/releases/download/v1.0.0/criblio-import-darwin-amd64
chmod +x criblio-import-darwin-amd64
sudo mv criblio-import-darwin-amd64 /usr/local/bin/criblio-import

# Download for Linux
curl -LO https://github.com/criblio/terraform-provider-criblio/releases/download/v1.0.0/criblio-import-linux-amd64
chmod +x criblio-import-linux-amd64
sudo mv criblio-import-linux-amd64 /usr/local/bin/criblio-import

# Download for Windows
# Download criblio-import-windows-amd64.exe and add to PATH
```

**Advantages:**
- ✅ No dependencies
- ✅ Works on any system
- ✅ Version control
- ✅ Easy CI/CD integration

#### Method 2: Go Install (Developer-Friendly)

**How it works:**
- Users with Go installed can install directly
- Always gets latest version

**User Experience:**
```bash
go install github.com/criblio/terraform-provider-criblio/tools/criblio-import/cmd/criblio-import@latest
```

**Advantages:**
- ✅ Simple for developers
- ✅ Always up-to-date

**Disadvantages:**
- ❌ Requires Go installation
- ❌ Not suitable for all users

#### Method 3: Package Managers

**Homebrew (macOS):**
```bash
brew tap criblio/criblio
brew install criblio-import
```

**Chocolatey (Windows):**
```powershell
choco install criblio-import
```

**Advantages:**
- ✅ Familiar package manager experience
- ✅ Easy updates

**Disadvantages:**
- ❌ Requires package manager setup
- ❌ Platform-specific

#### Method 4: Docker Container

**How it works:**
- Containerized version for CI/CD
- No local installation needed

**User Experience:**
```bash
docker run --rm \
  -v $(pwd):/output \
  -e CRIBL_BEARER_TOKEN="$TOKEN" \
  criblio/criblio-import:latest \
  criblio-import --output /output/terraform-configs
```

**Advantages:**
- ✅ No local installation
- ✅ Consistent environment
- ✅ CI/CD friendly

**Disadvantages:**
- ❌ Requires Docker
- ❌ More complex for simple use cases

### 3.2 Distribution Strategy

**Primary Distribution:** GitHub Releases
- Binary for Linux (amd64, arm64)
- Binary for macOS (amd64, arm64 - Apple Silicon)
- Binary for Windows (amd64)

**Secondary Distribution:**
- Homebrew tap for macOS users
- Go install for developers
- Docker Hub for containerized usage

### 3.3 Installation Documentation

**For End Users:**

#### Quick Start Guide
```markdown
# Installation

## Option 1: Download Binary (Recommended)

### macOS / Linux
```bash
# Download latest release
curl -LO https://github.com/criblio/terraform-provider-criblio/releases/latest/download/criblio-import-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x criblio-import-*
sudo mv criblio-import-* /usr/local/bin/criblio-import
```

### Windows
1. Download `criblio-import-windows-amd64.exe` from GitHub Releases
2. Rename to `criblio-import.exe`
3. Add to PATH or use from current directory

## Option 2: Homebrew (macOS)
```bash
brew tap criblio/criblio
brew install criblio-import
```

## Option 3: Go Install
```bash
go install github.com/criblio/terraform-provider-criblio/tools/criblio-import/cmd/criblio-import@latest
```

## Verify Installation
```bash
criblio-import version
```

### 3.4 User Onboarding Flow

**Step 1: Installation**
- User visits GitHub Releases page
- Downloads appropriate binary for their platform
- Adds to PATH or uses directly

**Step 2: Authentication Setup**
- User reads authentication documentation
- Chooses authentication method (Bearer token, OAuth, credentials file)
- Sets up credentials

**Step 3: First Import**
- User runs `criblio-import --output ./my-terraform`
- Tool provides progress feedback
- Generated files are ready for review

**Step 4: Integration**
- User reviews generated Terraform
- Runs `terraform init` and `terraform plan`
- Applies configuration if satisfied

### 3.5 Documentation Structure

**For End Users:**
1. **README.md** - Quick start, installation, basic usage
2. **USAGE.md** - Detailed usage examples, scenarios
3. **AUTHENTICATION.md** - Authentication methods, troubleshooting
4. **EXAMPLES/** - Sample configurations and use cases

**For Developers:**
1. **CONTRIBUTING.md** - Development setup, contribution guidelines
2. **ARCHITECTURE.md** - Internal design, component structure
3. **TESTING.md** - Testing strategy, test examples

### 3.6 Versioning Strategy

**Semantic Versioning:**
- Major: Breaking changes to CLI or output format
- Minor: New features, resource types
- Patch: Bug fixes, improvements

**Release Process:**
1. Automated builds on tag creation
2. GitHub Actions builds for all platforms
3. Release notes auto-generated from commits
4. Binaries attached to GitHub Release

### 3.7 Update Mechanism

**Manual Updates:**
- Users download new version from GitHub Releases
- Replace existing binary

**Automated Updates (Future):**
- `criblio-import self-update` command
- Checks GitHub Releases for newer version
- Downloads and replaces binary

**Notification:**
- Version check on startup (optional)
- Warns if newer version available
- Non-blocking, user can dismiss

### 3.8 CI/CD Integration

**GitHub Actions Example:**
```yaml
- name: Import Cribl Config
  uses: criblio/criblio-import-action@v1
  with:
    bearer-token: ${{ secrets.CRIBL_TOKEN }}
    workspace-id: main
    organization-id: org
    output-path: ./terraform-configs
```

**Advantages:**
- ✅ Automated infrastructure updates
- ✅ Version control for configuration
- ✅ Pull request reviews for changes

## 4. Acceptance Criteria Verification

### ✅ AC1: Design Specification Outlining Tool Usage

**Status:** Complete  
**Deliverable:** This document (Section 1)

**Coverage:**
- ✅ User workflow and scenarios
- ✅ Command-line interface design
- ✅ Example usage scenarios
- ✅ Output structure
- ✅ Error handling and feedback

### ✅ AC2: Programming Language Selection with Justification

**Status:** Complete  
**Deliverable:** This document (Section 2)

**Coverage:**
- ✅ Language selected: Go
- ✅ Detailed justification with 6 key factors
- ✅ Alternative languages considered and evaluated
- ✅ Trade-offs documented

### ✅ AC3: Deployment Process and End-User Presentation

**Status:** Complete  
**Deliverable:** This document (Section 3)

**Coverage:**
- ✅ Multiple deployment methods
- ✅ Installation instructions
- ✅ User onboarding flow
- ✅ Documentation structure
- ✅ Versioning and update strategy

## 5. Next Steps

### Implementation Phases

1. **Phase 1: Foundation** (Week 1-2)
   - Authentication module
   - API client wrapper
   - Archive extraction

2. **Phase 2: Parsing & Mapping** (Week 3)
   - YAML parser
   - Resource type detection
   - Field mapping

3. **Phase 3: Converters** (Week 4-5)
   - Resource converters (sources, destinations, etc.)
   - HCL generation

4. **Phase 4: Integration** (Week 6)
   - CLI integration
   - Testing
   - Documentation

### Dependencies

- ✅ Cribl API endpoint `/system/diag/download` confirmed
- ✅ Terraform provider SDK available
- ✅ Sample configuration structure needed

### Open Questions

1. Exact structure of diag bundle archive
2. Resource dependency handling
3. Sensitive field detection strategy
4. Terraform import command generation

## 6. Appendices

### Appendix A: Command Reference

See `README.md` for complete command reference.

### Appendix B: Authentication Methods

See `AUTHENTICATION.md` for detailed authentication setup.

### Appendix C: Example Output

See `examples/` directory for sample generated Terraform files.

---

**Document Status:** ✅ Ready for Review  
**Next Review Date:** After ERD creation

