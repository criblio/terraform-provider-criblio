# Cobra vs Alternatives: CLI Framework Comparison

## Why Cobra Was Selected

Cobra is the **industry-standard CLI framework** for Go, used by major projects:
- Kubernetes (`kubectl`)
- Docker (`docker`)
- GitHub CLI (`gh`)
- Helm (`helm`)
- Hugo (`hugo`)
- Terraform (`terraform`)

This makes it the **de facto standard** for professional Go CLI tools.

---

## Common Alternatives to Cobra

### 1. **Cobra** ⭐ (Selected)

**What it is:**
- Full-featured CLI framework with command/subcommand support
- Part of the `spf13` ecosystem (Cobra + Viper for config)
- Mature, battle-tested, widely adopted

**Pros:**
- ✅ **Industry standard** - Used by Kubernetes, Docker, GitHub CLI
- ✅ **Rich feature set** - Commands, subcommands, flags, persistent flags, aliases
- ✅ **Excellent documentation** - Extensive docs and examples
- ✅ **Viper integration** - Seamless config management (env vars, files, flags)
- ✅ **Auto-generated help** - Automatic `--help` and `-h` support
- ✅ **Shell completion** - Built-in bash/zsh/fish completion generation
- ✅ **Large community** - Easy to find examples and get help
- ✅ **Mature & stable** - Production-ready, well-maintained

**Cons:**
- ⚠️ **Learning curve** - More complex than simpler alternatives
- ⚠️ **Boilerplate** - More code required for simple CLIs
- ⚠️ **Dependency** - Adds external dependency (though widely used)

**Best for:**
- Professional CLI tools with multiple commands
- Tools that need subcommands (like `kubectl get`, `kubectl apply`)
- Tools requiring complex flag handling
- Production-grade applications

**Example:**
```go
var rootCmd = &cobra.Command{
    Use:   "criblio-import",
    Short: "Import Cribl configuration to Terraform",
    Long:  "A tool to import existing Cribl configuration...",
}

var importCmd = &cobra.Command{
    Use:   "import",
    Short: "Import configuration",
    Run: func(cmd *cobra.Command, args []string) {
        // Import logic
    },
}
```

---

### 2. **urfave/cli (v2/v3)** 

**What it is:**
- Simpler, more functional approach to CLI building
- Less boilerplate than Cobra
- Popular for smaller tools

**Pros:**
- ✅ **Simpler API** - Less boilerplate, more intuitive
- ✅ **Functional style** - Action-based commands
- ✅ **Good for simple CLIs** - Less overhead
- ✅ **Active development** - Well-maintained

**Cons:**
- ⚠️ **Less feature-rich** - Fewer built-in features
- ⚠️ **Smaller ecosystem** - Less community support
- ⚠️ **No built-in completion** - Need to implement yourself
- ⚠️ **Less structured** - Harder to organize complex CLIs

**Best for:**
- Simple single-command tools
- Tools with flat command structure
- Quick prototypes

**Example:**
```go
app := &cli.App{
    Name:  "criblio-import",
    Usage: "Import Cribl configuration",
    Action: func(c *cli.Context) error {
        // Import logic
        return nil
    },
}
```

**Comparison:**
- **Cobra**: Better for complex, multi-command tools
- **urfave/cli**: Better for simple, single-purpose tools

---

### 3. **kingpin**

**What it is:**
- Declarative CLI framework
- Inspired by Python's `argparse`
- Focus on type-safe flag parsing

**Pros:**
- ✅ **Type-safe** - Strong typing for flags
- ✅ **Declarative** - Clean, readable syntax
- ✅ **Good error messages** - Helpful validation errors
- ✅ **Lightweight** - Minimal dependencies

**Cons:**
- ⚠️ **Less popular** - Smaller community than Cobra
- ⚠️ **Less features** - Fewer built-in capabilities
- ⚠️ **Different paradigm** - May feel unfamiliar

**Best for:**
- Tools with complex flag validation
- Type-safe flag parsing needs
- Smaller tools

**Example:**
```go
var (
    output = kingpin.Flag("output", "Output directory").Required().String()
    dryRun = kingpin.Flag("dry-run", "Preview mode").Bool()
)
```

---

### 4. **go-flags**

**What it is:**
- Struct-based flag parsing
- Uses struct tags for configuration
- Similar to `flag` package but more powerful

**Pros:**
- ✅ **Struct-based** - Clean, organized code
- ✅ **Type-safe** - Compile-time safety
- ✅ **Good for config** - Natural fit for config structs

**Cons:**
- ⚠️ **Less flexible** - Harder for dynamic commands
- ⚠️ **Smaller community** - Less support
- ⚠️ **No subcommands** - Limited command structure

**Best for:**
- Simple tools with well-defined flags
- Config-driven applications

---

### 5. **Standard Library: `flag`**

**What it is:**
- Go's built-in flag package
- No external dependencies
- Minimal and simple

**Pros:**
- ✅ **No dependencies** - Part of standard library
- ✅ **Simple** - Easy to understand
- ✅ **Lightweight** - Minimal overhead

**Cons:**
- ❌ **Very limited** - No subcommands, no help generation
- ❌ **Manual work** - Need to implement help, validation yourself
- ❌ **Not scalable** - Hard to maintain for complex CLIs

**Best for:**
- Very simple tools with 1-2 flags
- Internal scripts
- When you want zero dependencies

**Example:**
```go
output := flag.String("output", "", "Output directory")
flag.Parse()
```

---

### 6. **mitchellh/cli**

**What it is:**
- HashiCorp's CLI framework
- Used by Terraform, Packer, Consul
- Simple, opinionated approach

**Pros:**
- ✅ **HashiCorp standard** - Used by Terraform ecosystem
- ✅ **Simple** - Less boilerplate than Cobra
- ✅ **Good for HashiCorp tools** - Consistent with Terraform

**Cons:**
- ⚠️ **Less features** - Fewer capabilities than Cobra
- ⚠️ **HashiCorp-specific** - Less general-purpose
- ⚠️ **Smaller ecosystem** - Less community support

**Best for:**
- Tools in HashiCorp ecosystem
- Simple multi-command tools

---

## Comparison Matrix

| Feature | Cobra | urfave/cli | kingpin | go-flags | flag | mitchellh/cli |
|---------|-------|------------|---------|----------|------|---------------|
| **Subcommands** | ✅ Excellent | ✅ Good | ✅ Good | ❌ Limited | ❌ No | ✅ Good |
| **Help Generation** | ✅ Auto | ✅ Auto | ✅ Auto | ⚠️ Manual | ❌ Manual | ✅ Auto |
| **Shell Completion** | ✅ Built-in | ❌ No | ⚠️ Limited | ❌ No | ❌ No | ❌ No |
| **Config Integration** | ✅ Viper | ⚠️ Manual | ⚠️ Manual | ⚠️ Manual | ❌ No | ⚠️ Manual |
| **Learning Curve** | ⚠️ Medium | ✅ Easy | ✅ Easy | ✅ Easy | ✅ Very Easy | ✅ Easy |
| **Boilerplate** | ⚠️ Medium | ✅ Low | ✅ Low | ✅ Low | ✅ Very Low | ✅ Low |
| **Community** | ✅ Huge | ✅ Large | ⚠️ Medium | ⚠️ Small | ✅ Built-in | ⚠️ Medium |
| **Production Ready** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ⚠️ Limited | ✅ Yes |
| **Best For** | Complex CLIs | Simple CLIs | Type-safe | Struct-based | Minimal | HashiCorp |

---

## Why Cobra for `criblio-import`?

### 1. **Multi-Command Structure** ⭐
The tool needs multiple commands:
- `criblio-import import` - Main import command
- `criblio-import version` - Version info
- `criblio-import help` - Help

Cobra excels at organizing subcommands.

### 2. **Complex Flag Handling** ⭐
The tool has many flags:
- `--output`, `--include`, `--exclude`, `--dry-run`
- `--bearer-token`, `--client-id`, `--workspace-id`
- Authentication flags, filtering flags, execution flags

Cobra's persistent flags and flag groups make this manageable.

### 3. **Viper Integration** ⭐
The tool needs config from multiple sources:
- Environment variables (`CRIBL_BEARER_TOKEN`)
- Config files (`~/.cribl/credentials`)
- Command-line flags

Cobra + Viper is the standard solution for this.

### 4. **Professional Tool** ⭐
This is a production tool for end users. Cobra provides:
- Professional help output
- Shell completion (bash/zsh/fish)
- Consistent UX with other major tools

### 5. **Ecosystem Alignment** ⭐
Terraform uses Cobra (via HashiCorp's CLI framework). Using Cobra keeps the tool familiar to Terraform users.

### 6. **Future Extensibility** ⭐
As the tool grows, Cobra makes it easy to add:
- New subcommands
- Plugin system
- Advanced features

---

## Real-World Usage Examples

### Kubernetes (`kubectl`)
```bash
kubectl get pods          # Subcommand: get
kubectl apply -f file.yaml # Subcommand: apply
kubectl --help            # Auto-generated help
```

### Docker
```bash
docker run ...            # Subcommand: run
docker build ...          # Subcommand: build
docker --help             # Auto-generated help
```

### GitHub CLI (`gh`)
```bash
gh repo clone ...         # Subcommand: repo
gh pr create ...         # Subcommand: pr
gh --help                # Auto-generated help
```

All use Cobra, demonstrating its production-readiness.

---

## Alternative Consideration: urfave/cli

**When urfave/cli might be better:**
- Simple single-command tool
- Minimal flag requirements
- Quick prototype

**Why not for `criblio-import`:**
- ❌ Tool has multiple commands (import, version, help)
- ❌ Complex flag structure (20+ flags)
- ❌ Needs config management (Viper integration)
- ❌ Production tool requiring professional UX

---

## Conclusion

**Cobra is the right choice** for `criblio-import` because:

1. ✅ **Industry standard** - Used by major tools users are familiar with
2. ✅ **Feature-rich** - Handles complex CLI requirements
3. ✅ **Viper integration** - Seamless config management
4. ✅ **Production-ready** - Battle-tested at scale
5. ✅ **Future-proof** - Easy to extend as tool grows
6. ✅ **Professional UX** - Auto-help, completion, consistent interface

**Alternatives considered:**
- **urfave/cli**: Too simple for multi-command tool
- **kingpin**: Less popular, smaller ecosystem
- **go-flags**: Limited subcommand support
- **flag**: Too basic for production tool
- **mitchellh/cli**: Good but less features than Cobra

The slight learning curve and additional boilerplate are worth it for the professional-grade CLI experience Cobra provides.

