# CLI Tool Design Flow Diagram

## Overview

This document provides visual diagrams and detailed explanations of how the `criblio-import` CLI tool is designed and works.

---

## 1. High-Level Architecture Flow

```mermaid
graph TB
    subgraph "Code Generation (Build Time)"
        A[openapi.yml<br/>Base OpenAPI Schema] --> B[entity-mapping-overlay.yml<br/>Schema Modifications]
        B --> C[Speakeasy Generator]
        C --> D[internal/sdk/<br/>Generated API Client]
        C --> E[internal/provider/<br/>Generated Terraform Provider]
    end
    
    subgraph "CLI Tool Runtime"
        F[User Input<br/>Flags/Env/Config] --> G[criblio-import CLI]
        G --> H[Authentication<br/>CriblTerraformHook]
        H --> I[SDK Client<br/>internal/sdk]
        I --> J[Cribl API<br/>/bulk/diag/download]
        J --> K[Archive Download<br/>.tar.gz]
        K --> L[Extract & Parse<br/>YAML Files]
        L --> M[Resource Converters<br/>YAML → HCL]
        M --> N[Terraform Files<br/>.tf Output]
    end
    
    D -.->|Import| I
    E -.->|Reference| M
    
    style A fill:#e1f5ff
    style D fill:#c8e6c9
    style E fill:#c8e6c9
    style G fill:#fff9c4
    style I fill:#c8e6c9
    style N fill:#ffccbc
```

---

## 2. Detailed Component Flow

```mermaid
graph LR
    subgraph "1. Initialization"
        A1[main.go<br/>Cobra CLI] --> A2[initViper<br/>Config Management]
        A2 --> A3[Load Config<br/>Flags > Env > File]
    end
    
    subgraph "2. Authentication"
        B1[Parse Auth Flags] --> B2[CriblTerraformHook<br/>Auto Auth Handler]
        B2 --> B3{Auth Method?}
        B3 -->|Bearer| B4[Bearer Token]
        B3 -->|OAuth| B5[OAuth Flow]
        B3 -->|On-Prem| B6[Username/Password]
        B3 -->|File| B7[~/.cribl/credentials]
        B4 --> B8[SDK Client Initialized]
        B5 --> B8
        B6 --> B8
        B7 --> B8
    end
    
    subgraph "3. API Interaction"
        C1[SDK Client] --> C2[Call API<br/>/bulk/diag/download]
        C2 --> C3[Download Archive<br/>.tar.gz]
        C3 --> C4[Extract Archive]
    end
    
    subgraph "4. Processing"
        D1[Parse YAML Files] --> D2[Detect Resource Type]
        D2 --> D3[Route to Converter]
        D3 --> D4[Source Converter]
        D3 --> D5[Destination Converter]
        D3 --> D6[Pipeline Converter]
        D3 --> D7[Other Converters]
        D4 --> D8[Generate HCL]
        D5 --> D8
        D6 --> D8
        D7 --> D8
    end
    
    subgraph "5. Output"
        E1[HCL Generation] --> E2[Format with<br/>terraform fmt]
        E2 --> E3[Write Files<br/>modules/ or resources/]
        E3 --> E4[Optional: Auto-Import<br/>terraform import]
    end
    
    A3 --> B1
    B8 --> C1
    C4 --> D1
    D8 --> E1
    E4 --> E5[Complete]
    
    style A1 fill:#fff9c4
    style B2 fill:#c8e6c9
    style C1 fill:#c8e6c9
    style D8 fill:#ffccbc
    style E3 fill:#ffccbc
```

---

## 3. Code Generation Flow (Build Time)

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant OpenAPI as openapi.yml
    participant Overlay as entity-mapping-overlay.yml
    participant Speakeasy as Speakeasy Generator
    participant SDK as internal/sdk/
    participant Provider as internal/provider/
    participant CLI as CLI Tool
    
    Dev->>OpenAPI: Update API Schema
    Dev->>Overlay: Add Overlay Rules
    Dev->>Speakeasy: Run: make build-speakeasy
    Speakeasy->>OpenAPI: Read Base Schema
    Speakeasy->>Overlay: Apply Overlays
    Speakeasy->>SDK: Generate API Client
    Speakeasy->>Provider: Generate Terraform Provider
    SDK->>CLI: Available for Import
    Provider->>CLI: Reference for Schemas
```

---

## 4. CLI Tool Execution Flow (Runtime)

```mermaid
sequenceDiagram
    participant User
    participant CLI as criblio-import
    participant Viper as Config Manager
    participant Auth as CriblTerraformHook
    participant SDK as internal/sdk
    participant API as Cribl API
    participant Converter as Resource Converters
    participant HCL as HCL Generator
    participant Terraform as Terraform CLI
    
    User->>CLI: Run: criblio-import --output ./configs
    CLI->>Viper: Load Config (Flags > Env > File)
    Viper-->>CLI: Config Values
    CLI->>Auth: Initialize SDK Client
    Auth->>Auth: Determine Auth Method
    Auth->>SDK: Create Authenticated Client
    CLI->>SDK: Call GetDiagDownload()
    SDK->>API: GET /bulk/diag/download
    API-->>SDK: Archive (.tar.gz)
    SDK-->>CLI: Archive Data
    CLI->>CLI: Extract Archive
    CLI->>CLI: Parse YAML Files
    loop For Each Resource
        CLI->>Converter: Convert YAML to HCL
        Converter->>HCL: Generate Resource Block
        HCL-->>CLI: HCL String
    end
    CLI->>CLI: Write .tf Files
    CLI->>Terraform: terraform fmt (optional)
    opt Auto-Import Enabled
        CLI->>Terraform: terraform import (for each resource)
    end
    CLI-->>User: Success: Files Generated
```

---

## 5. Data Flow Diagram

```mermaid
graph TD
    subgraph "Input Sources"
        I1[Command Line Flags]
        I2[Environment Variables<br/>CRIBL_*]
        I3[Config File<br/>~/.cribl/credentials]
    end
    
    subgraph "Configuration Layer"
        C1[Viper Config Manager]
        C2[Priority Resolution<br/>Flags > Env > File]
    end
    
    subgraph "Authentication Layer"
        A1[CriblTerraformHook]
        A2[Auto-Detect Auth Method]
        A3[Bearer Token / OAuth / On-Prem / File]
    end
    
    subgraph "API Layer"
        API1[SDK Client<br/>internal/sdk]
        API2[Cribl API Endpoints]
        API3[Archive Download]
    end
    
    subgraph "Processing Layer"
        P1[Archive Extraction]
        P2[YAML Parsing]
        P3[Resource Type Detection]
        P4[Field Mapping<br/>camelCase → snake_case]
        P5[HCL Generation]
    end
    
    subgraph "Output Layer"
        O1[Terraform Files<br/>.tf]
        O2[Module Structure]
        O3[Optional: Auto-Import]
    end
    
    I1 --> C1
    I2 --> C1
    I3 --> C1
    C1 --> C2
    C2 --> A1
    A1 --> A2
    A2 --> A3
    A3 --> API1
    API1 --> API2
    API2 --> API3
    API3 --> P1
    P1 --> P2
    P2 --> P3
    P3 --> P4
    P4 --> P5
    P5 --> O1
    O1 --> O2
    O2 --> O3
    
    style C1 fill:#fff9c4
    style A1 fill:#c8e6c9
    style API1 fill:#c8e6c9
    style P5 fill:#ffccbc
    style O1 fill:#ffccbc
```

---

## 6. Resource Conversion Flow

```mermaid
graph TB
    subgraph "Input"
        A[YAML Config File<br/>from diag bundle]
    end
    
    subgraph "Detection"
        B[Parse YAML Structure]
        C{Resource Type?}
        C -->|source| D[Source Converter]
        C -->|destination| E[Destination Converter]
        C -->|pipeline| F[Pipeline Converter]
        C -->|routes| G[Routes Converter]
        C -->|pack| H[Pack Converter]
        C -->|global_var| I[GlobalVar Converter]
    end
    
    subgraph "Conversion"
        D --> J[Field Mapper]
        E --> J
        F --> J
        G --> J
        H --> J
        I --> J
        J --> K[Type Conversion]
        K --> L[Schema Validation]
    end
    
    subgraph "Output"
        L --> M[HCL Block Generation]
        M --> N[Format with hclwrite]
        N --> O[Write to .tf File]
    end
    
    A --> B
    B --> C
    
    style J fill:#fff9c4
    style M fill:#ffccbc
    style O fill:#ffccbc
```

---

## 7. File Structure Flow

```
User runs: criblio-import --output ./terraform-configs

Input:
├── Flags: --output, --include, --bearer-token, etc.
├── Env Vars: CRIBL_BEARER_TOKEN, CRIBL_WORKSPACE_ID, etc.
└── Config: ~/.cribl/credentials.ini

Processing:
├── Viper merges config (Flags > Env > File)
├── CriblTerraformHook authenticates
├── SDK calls /bulk/diag/download
├── Archive extracted to temp directory
└── YAML files parsed and converted

Output Structure:
terraform-configs/
├── main.tf                 # Provider configuration
├── variables.tf            # Sensitive variables (optional)
├── README.md              # Generated documentation
└── modules/
    ├── sources/
    │   └── main.tf        # All source resources
    ├── destinations/
    │   └── main.tf        # All destination resources
    ├── pipelines/
    │   └── main.tf        # All pipeline resources
    ├── routes/
    │   └── main.tf        # All routes resources
    └── packs/
        └── main.tf        # All pack resources
```

---

## 8. Authentication Flow

```mermaid
graph TD
    A[CLI Starts] --> B[Read Config<br/>Viper]
    B --> C{Credentials Found?}
    C -->|Flags| D[Use Flag Values]
    C -->|Env Vars| E[Use CRIBL_* vars]
    C -->|Config File| F[Read ~/.cribl/credentials]
    C -->|None| G[Error: Auth Required]
    
    D --> H[CriblTerraformHook]
    E --> H
    F --> H
    
    H --> I{Detect Auth Type}
    I -->|Bearer Token| J[Bearer Auth]
    I -->|Client ID/Secret| K[OAuth Flow]
    I -->|On-Prem URL| L[Basic Auth]
    I -->|Credentials File| M[Read Profile]
    
    J --> N[SDK Client Ready]
    K --> N
    L --> N
    M --> N
    
    N --> O[API Calls Available]
    
    style H fill:#c8e6c9
    style N fill:#c8e6c9
```

---

## 9. Error Handling Flow

```mermaid
graph TD
    A[Operation Starts] --> B{Success?}
    B -->|Yes| C[Continue]
    B -->|No| D{Error Type?}
    
    D -->|Auth Error| E[Show Auth Help<br/>List Available Methods]
    D -->|API Error| F[Show API Error<br/>Suggest Troubleshooting]
    D -->|File Error| G[Show Permission Error<br/>Suggest Fix]
    D -->|Validation Error| H[Show Resource Errors<br/>List Failed Resources]
    D -->|Network Error| I[Show Network Error<br/>Check Connectivity]
    
    E --> J[Exit with Code 1]
    F --> J
    G --> J
    H --> J
    I --> J
    
    C --> K[Success: Files Generated]
    
    style D fill:#ffccbc
    style J fill:#ffccbc
    style K fill:#c8e6c9
```

---

## 10. Integration Points

### SDK Integration
```
CLI Tool
  ↓ imports
internal/sdk
  ├── criblio.go (Main Client)
  ├── models/operations/ (API Endpoints)
  └── models/shared/ (Data Types)
```

### Provider Integration
```
CLI Tool
  ↓ references
internal/provider
  ├── *_resource.go (Schemas)
  └── *_resource_sdk.go (Type Mappings)
```

### Terraform Integration
```
CLI Tool
  ↓ generates
Terraform Files (.tf)
  ↓ imports
Terraform State
  ↓ manages
Cribl Resources
```

---

## 11. Key Design Decisions

### 1. **Code Reuse**
- ✅ Import `internal/sdk` directly
- ✅ Reuse authentication (`CriblTerraformHook`)
- ✅ Reference provider schemas for validation

### 2. **Configuration Management**
- ✅ Viper for multi-source config
- ✅ Priority: Flags > Env > File
- ✅ Supports all auth methods

### 3. **Modular Conversion**
- ✅ Separate converter per resource type
- ✅ Reusable field mapper
- ✅ Type-safe HCL generation

### 4. **Output Flexibility**
- ✅ Module structure (default)
- ✅ Flat resources (option)
- ✅ Optional auto-import

---

## 12. Example Execution Flow

### Step-by-Step Example

```bash
# 1. User runs command
criblio-import --output ./configs --include sources,destinations

# 2. Viper loads config
#    - Checks flags: output=./configs, include=sources,destinations
#    - Checks env: CRIBL_BEARER_TOKEN, CRIBL_WORKSPACE_ID
#    - Checks file: ~/.cribl/credentials

# 3. Authentication
#    - CriblTerraformHook detects: Bearer token from env
#    - Creates authenticated SDK client

# 4. API Call
#    - SDK calls: GET /bulk/diag/download
#    - Receives: archive.tar.gz

# 5. Processing
#    - Extract archive
#    - Parse YAML files
#    - Filter: only sources and destinations
#    - Convert YAML → HCL

# 6. Output
#    - Write modules/sources/main.tf
#    - Write modules/destinations/main.tf
#    - Write main.tf (provider config)

# 7. Result
#    ✅ Success: Files generated in ./configs/
```

---

## Summary

The CLI tool follows a **clean, modular architecture**:

1. **Code Generation** (Build Time): Speakeasy generates SDK and provider from OpenAPI
2. **Configuration** (Runtime): Viper manages multi-source config
3. **Authentication** (Runtime): CriblTerraformHook handles all auth methods
4. **API Interaction** (Runtime): SDK client calls Cribl API
5. **Processing** (Runtime): Converters transform YAML to HCL
6. **Output** (Runtime): Generate organized Terraform files

**Key Benefits:**
- ✅ Maximum code reuse (SDK, auth, schemas)
- ✅ Type safety (generated types)
- ✅ Flexible configuration (flags, env, files)
- ✅ Modular design (easy to extend)
- ✅ Production-ready (error handling, validation)

