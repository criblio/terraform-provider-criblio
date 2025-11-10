# Building CLI Tool as Separate Private Repository - Analysis

## Difficulty Assessment: **Medium to High** ⚠️

The main challenge is maintaining code reuse of the SDK while managing dependencies across repositories.

---

## ✅ PROS of Separate Repository

### 1. **Independent Release Cycle** ⭐
- **Benefit**: CLI tool can be released independently of Terraform provider
- **Impact**: Faster iteration, separate versioning (e.g., `criblio-import v2.0` while provider is `v1.5`)
- **Use Case**: CLI tool needs frequent updates for new features, provider is stable

### 2. **Access Control & Security** ⭐
- **Benefit**: Different teams can have different access levels
- **Impact**: CLI tool team doesn't need full provider repo access
- **Use Case**: CLI tool maintained by different team or external contractors

### 3. **Cleaner Repository Structure**
- **Benefit**: Provider repo stays focused on provider code
- **Impact**: Easier navigation, clearer ownership
- **Use Case**: Provider repo is large, adding tools makes it unwieldy

### 4. **Independent CI/CD Pipelines**
- **Benefit**: Separate build/test/deploy workflows
- **Impact**: Faster CI runs, independent deployment schedules
- **Use Case**: Different deployment targets (CLI binaries vs provider releases)

### 5. **License & Distribution Flexibility**
- **Benefit**: Can have different licenses or distribution models
- **Impact**: CLI tool could be open-source while provider is private (or vice versa)
- **Use Case**: Different business models for different products

---

## ❌ CONS of Separate Repository

### 1. **SDK Code Reuse Complexity** ⚠️ (CRITICAL)

**The Problem:**
- CLI tool needs `internal/sdk` for API client and authentication
- `internal/sdk` is currently in the provider repo with path `internal/sdk`
- Go modules can't import `internal/` paths from other repos

**Solutions & Trade-offs:**

#### Option A: Extract SDK to Separate Module (Recommended)
```go
// In provider repo: go.mod
module github.com/criblio/terraform-provider-criblio
require github.com/criblio/criblio-sdk-go v1.0.0

// In CLI repo: go.mod
module github.com/criblio/criblio-import
require github.com/criblio/criblio-sdk-go v1.0.0
```

**Difficulty**: **High**
- Requires refactoring SDK into separate Go module
- Need to update provider to use extracted SDK
- Version management between SDK, provider, and CLI
- Breaking change for provider repo

**Effort**: 2-3 weeks of refactoring

#### Option B: Use Go Replace Directive (Development Only)
```go
// In CLI repo: go.mod
replace github.com/criblio/terraform-provider-criblio => ../terraform-provider-criblio
require github.com/criblio/terraform-provider-criblio v0.0.0
```

**Difficulty**: **Low** (but not production-ready)
- Only works for local development
- Can't be used in CI/CD or distributed binaries
- Not a real solution

#### Option C: Copy/Vendor SDK Code
**Difficulty**: **Medium**
- Copy SDK code to CLI repo
- Maintain two copies of SDK
- Manual sync required on SDK updates
- Risk of divergence

**Effort**: Ongoing maintenance burden

#### Option D: Reimplement API Client
**Difficulty**: **Very High**
- Lose all code reuse benefits
- Reimplement authentication, API calls, models
- Defeats the purpose of choosing Go

**Effort**: 4-6 weeks, defeats design goals

### 2. **Dependency Management Complexity**

**The Challenge:**
- Provider and CLI need to stay in sync on SDK versions
- Breaking changes in SDK affect both repos
- Need coordination for releases

**Example Scenario:**
```
SDK v1.0.0 → Provider v1.5.0 uses it
SDK v1.1.0 → Adds new API endpoint
CLI needs v1.1.0 for new feature
Provider still on v1.0.0
→ Need to update provider OR maintain two SDK versions
```

### 3. **Development Workflow Friction**

**Issues:**
- Testing CLI changes requires SDK changes in separate repo
- Can't easily test end-to-end in single repo
- More complex local development setup
- Cross-repo PR reviews and coordination

**Workflow Example:**
```
1. Developer needs new API endpoint in CLI
2. SDK needs update (in provider repo)
3. Create PR in provider repo → wait for review → merge
4. Release new SDK version
5. Update CLI to use new SDK version
6. Test CLI changes
→ 3-5 day cycle vs 1 day in monorepo
```

### 4. **Versioning & Release Coordination**

**Complexity:**
- Three version numbers to manage: SDK, Provider, CLI
- Need semantic versioning strategy
- Release coordination across repos
- Dependency graph management

**Example:**
```
SDK: v1.2.0
Provider: v1.5.0 (depends on SDK v1.2.0)
CLI: v2.1.0 (depends on SDK v1.2.0)

SDK v1.3.0 released with breaking changes
→ Both Provider and CLI need updates
→ Coordinated release or CLI stuck on old SDK
```

### 5. **Documentation & Discovery**

**Issues:**
- Users need to find two separate repos
- Documentation split across repos
- Less discoverable than single repo
- More complex contribution process

---

## 🔧 Implementation Approaches

### Approach 1: Extract SDK First (Recommended if Going Separate)

**Steps:**
1. Create new repo: `criblio-sdk-go`
2. Move `internal/sdk` → `criblio-sdk-go`
3. Convert to proper Go module with semantic versioning
4. Update provider to depend on SDK module
5. CLI repo depends on same SDK module

**Timeline**: 3-4 weeks
**Difficulty**: High
**Benefit**: Clean architecture, reusable SDK

### Approach 2: Monorepo with Workspaces

**Steps:**
1. Keep everything in provider repo
2. Use Go workspaces (Go 1.18+)
3. CLI tool in `tools/criblio-import/` (current structure)
4. Both depend on `internal/sdk`

**Timeline**: Already done
**Difficulty**: Low
**Benefit**: Maximum code reuse, simple development

### Approach 3: Git Submodules (Not Recommended)

**Steps:**
1. Provider repo has SDK
2. CLI repo uses git submodule to reference provider
3. Import SDK via submodule path

**Timeline**: 1 week
**Difficulty**: Medium
**Benefit**: Separate repos
**Drawback**: Submodules are painful, versioning issues

---

## 📊 Comparison Matrix

| Factor | Separate Repo | Monorepo (Current) |
|--------|---------------|-------------------|
| **Code Reuse** | ⚠️ Complex (need SDK extraction) | ✅ Direct import |
| **Development Speed** | ⚠️ Slower (cross-repo coordination) | ✅ Fast (single repo) |
| **Release Independence** | ✅ Independent | ⚠️ Coupled |
| **Access Control** | ✅ Fine-grained | ⚠️ All-or-nothing |
| **CI/CD Complexity** | ⚠️ Two pipelines | ✅ Single pipeline |
| **Version Management** | ⚠️ Three versions | ✅ Simpler |
| **Local Development** | ⚠️ Multi-repo setup | ✅ Simple |
| **Documentation** | ⚠️ Split | ✅ Centralized |

---

## 💡 Recommendation

### **Keep in Monorepo (Current Approach)** ✅

**Reasons:**
1. **Maximum Code Reuse**: Direct import of `internal/sdk` with zero friction
2. **Faster Development**: Single repo, single PR, single CI pipeline
3. **Simpler Maintenance**: One version to manage, one place for docs
4. **Design Goals Met**: The design spec explicitly chose Go for code reuse - monorepo maximizes this

### **When to Consider Separate Repo:**

Only if you have:
- ✅ Different teams maintaining provider vs CLI
- ✅ Different release cadences (CLI weekly, provider monthly)
- ✅ Different access control requirements
- ✅ Willingness to extract SDK to separate module (3-4 week effort)
- ✅ Need for different licensing/distribution models

### **Hybrid Approach (Future):**

If you later need separation:
1. Extract SDK to `criblio-sdk-go` module (public or private)
2. Both provider and CLI depend on SDK module
3. Keep provider and CLI in separate repos
4. SDK becomes shared dependency

This gives you separation while maintaining code reuse.

---

## 🎯 Difficulty Summary

| Task | Difficulty | Effort | Risk |
|------|-----------|--------|------|
| **Extract SDK to module** | High | 3-4 weeks | Medium |
| **Update provider to use extracted SDK** | Medium | 1 week | Low |
| **Set up CLI in separate repo** | Low | 1 week | Low |
| **Maintain version sync** | Medium | Ongoing | Medium |
| **Total (if extracting SDK)** | **High** | **5-6 weeks** | **Medium** |

**Without SDK extraction**: Not feasible (can't import `internal/` paths)

---

## 📝 Conclusion

**Building as separate repo is feasible but requires significant upfront work** (SDK extraction). The monorepo approach is simpler, faster, and better aligned with the design goals of maximum code reuse.

**Recommendation**: Start with monorepo, extract SDK later if separation becomes necessary.

