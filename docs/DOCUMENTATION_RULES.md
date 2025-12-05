# Documentation Rules for starknet.go

These are the strict rules that must be followed when creating or updating documentation for starknet.go. These rules were established after discovering AI-hallucinated struct definitions that would have caused compilation errors.

## 1. Source Code Verification - MANDATORY

### Every Type Definition Must Be Verified
- **NEVER** write struct definitions from memory or by guessing
- **ALWAYS** read the actual source code file from the repository
- **ALWAYS** include a GitHub source link for every type definition
- **Format**: `**Source:** [filename.go:LX-LY](https://github.com/NethermindEth/starknet.go/blob/main/path/to/filename.go#LX-LY)`

### Verification Process
1. Identify the source file containing the type definition
2. Read the ACTUAL file from `/Users/aayushgiri/work/starknet.go/` or fetch from GitHub
3. Copy the struct definition EXACTLY as it appears in source code
4. Include ALL fields with correct:
   - Field names (exact capitalization)
   - Field types (exact type names including package prefixes)
   - JSON tags (exact tag values including `omitempty` where present)
   - Comments (copy helpful comments from source)
5. Provide GitHub link to the exact line numbers

### What to Verify
- ✅ Struct field names
- ✅ Struct field types
- ✅ JSON tags
- ✅ Method signatures (receiver type, parameters, return types)
- ✅ Function signatures
- ✅ Parameter names and types
- ✅ Return value types
- ✅ Package imports in examples

## 2. Method Signature Verification

### Every Method Must Be Verified
- Read the actual method implementation from source
- Verify receiver type: `(p *Provider)` vs `(provider *Provider)`
- Verify parameter names and types match exactly
- Verify return types match exactly
- Include source link to method implementation

### Example
```go
// CORRECT - Verified against source
func (provider *Provider) StorageProof(
    ctx context.Context,
    storageProofInput StorageProofInput,
) (*StorageProofResult, error)
```
**Source:** [contract.go:L245-L262](https://github.com/NethermindEth/starknet.go/blob/main/rpc/contract.go#L245-L262)

## 3. Usage Examples Must Compile

### All Code Examples Must Be Valid Go Code
- Examples must use the CORRECT struct field names from verified source
- Examples must import the correct packages
- Examples must use the correct types
- Examples must be complete (not pseudo-code)

### Testing Examples
Before writing an example:
1. Verify all struct definitions it uses
2. Verify all method signatures it calls
3. Ensure imports are correct
4. Use actual field names from structs

### Environment Variables
All examples MUST use environment variables for RPC URLs:
- Use `github.com/joho/godotenv` to load `.env` files
- Always call `godotenv.Load()` at the start of `main()`
- Check if the environment variable is set and provide clear error message
- Example `.env` file format:
  ```bash
  STARKNET_RPC_URL=https://starknet-sepolia.g.alchemy.com/starknet/version/rpc/v0_9/YOUR_API_KEY
  ```
- Required package installation note:
  ```bash
  go get github.com/joho/godotenv
  ```

### Example Structure
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/NethermindEth/juno/core/felt"
    "github.com/NethermindEth/starknet.go/rpc"
)

func main() {
    // Get RPC URL from environment variable
    rpcURL := os.Getenv("STARKNET_RPC_URL")
    if rpcURL == "" {
        log.Fatal("STARKNET_RPC_URL environment variable is not set")
    }

    // Create provider
    provider, err := rpc.NewProvider(context.Background(), rpcURL)
    if err != nil {
        log.Fatal(err)
    }

    // Use actual struct with VERIFIED field names
    input := rpc.StorageProofInput{
        BlockID:           rpc.BlockID{Tag: "latest"},
        ContractAddresses: []*felt.Felt{contractAddr},
    }

    // Call with actual verified method signature
    result, err := provider.StorageProof(context.Background(), input)
    if err != nil {
        log.Fatal(err)
    }

    // Access using VERIFIED field names
    fmt.Printf("Result: %s\n", result.GlobalRoots.BlockHash)
}
```

## 4. Documentation Structure

### Method Documentation Page Structure

**ONLY SECTIONS REQUIRED (in this exact order):**

1. **Title + Description Paragraph**
   - Method name as H1 (e.g., `# BlockNumber`)
   - 2-3 sentence description explaining what it does, important details, edge cases
   - Include information like block numbering, error conditions, etc.

2. **Method Signature**
   - Code block with verified signature
   - **MUST include source link**: `**Source:** [filename.go:LX-LY](https://github.com/...)`

3. **Parameters**
   - List each parameter with type and description
   - Format: `- paramName (Type): Description`

4. **Returns**
   - List each return value with type and description
   - Format: `- Type: Description`

5. **Type Structure** (REQUIRED if method returns a custom type/struct)
   - Add section IMMEDIATELY after Returns section
   - MUST include GitHub source link to type definition
   - Show complete struct/type definition
   - Explain all fields with comments
   - Provide example showing what a created instance contains
   - Show how to access fields
   - Format:
     ```markdown
     ## TypeName Structure

     **Source:** [file.go:LX-LY](github-link)

     The returned TypeName instance contains:

     ```go
     type TypeName struct {
         Field1 Type1  // Description
         Field2 Type2  // Description
     }
     ```

     **Example of a created instance:**

     ```go
     result, err := SomeMethod(...)
     // result.Field1 -> value
     // result.Field2 -> value
     ```
     ```

6. **Usage Example**
   - ONE complete, runnable code example
   - MUST use godotenv pattern for environment variables
   - MUST include `.env` file setup instructions
   - MUST include package installation instructions

7. **Error Handling**
   - Show proper error handling with actual errors from source code
   - Use switch statement for multiple error types

8. **Common Use Cases**
   - Bullet list of common scenarios where this method is used
   - Each bullet should link to examples page: `See [example](/docs/rpc/examples#method-name)`
   - Include final link to full examples: `For complete examples, see [RPC Examples - MethodName](/docs/rpc/examples#method-name)`

**FORBIDDEN SECTIONS:**
- ❌ NO "RPC Specification" section
- ❌ NO "Important Notes" section (put in description paragraph instead)
- ❌ NO inline code examples in "Common Use Cases" (link to examples page)
- ❌ NO type definitions in method pages (only in dedicated type pages)
- ❌ NO "Expected Output" section (examples should be self-documenting)
- ❌ NO "Related Methods" or "Related Functions" section

**Example Structure:**
```markdown
# MethodName

Description paragraph explaining what the method does, including important details about behavior, edge cases, etc.

## Method Signature

\`\`\`go
func (provider *Provider) MethodName(ctx context.Context, param Type) (ReturnType, error)
\`\`\`

**Source:** [file.go:LX-LY](https://github.com/...)

## Parameters

- `param` (Type): Description

## Returns

- `ReturnType`: Description
- `error`: Error description

## Usage Example

\`\`\`go
// Complete example with godotenv
\`\`\`

Create a `.env` file:
\`\`\`bash
STARKNET_RPC_URL=...
\`\`\`

Install godotenv:
\`\`\`bash
go get github.com/joho/godotenv
\`\`\`

## Error Handling

\`\`\`go
// Error handling example
\`\`\`

## Common Use Cases

This method is commonly used for:

- **Use Case 1**: Description. See [example](/docs/rpc/examples#method-name).
- **Use Case 2**: Description. See [example](/docs/rpc/examples#method-name).

For complete examples, see [RPC Examples - MethodName](/docs/rpc/examples#method-name).

## Related Methods

- [RelatedMethod](/docs/rpc/methods/related-method) - Description
```

### Type Definition Section
```markdown
### TypeName

Brief description from source code comments.

\`\`\`go
type TypeName struct {
    // Copy comments from source if helpful
    FieldName FieldType `json:"field_name,omitempty"`
    // ...
}
\`\`\`

**Source:** [filename.go:LX-LY](https://github.com/...)
```

## 5. Examples Organization

### Central Examples File
- All comprehensive examples go in `/docs/docs/pages/docs/rpc/examples.mdx`
- Use HTML span with id for anchors: `## <span id="method-name">MethodName Examples</span>`
- Each method gets its own section with multiple example scenarios
- Examples must be complete, runnable code

### Method Page Examples
- Method documentation pages include ONE basic usage example
- Link to central examples page for more: `[RPC Examples - MethodName](/docs/rpc/examples#method-name)`

### Example Categories
Each method should have examples for:
1. Basic usage
2. Different parameter combinations
3. Error handling
4. Common use cases
5. Advanced scenarios

## 6. Cross-References and Links

### Internal Links
- Link to related methods
- Link to type definition pages
- Link to examples page
- Use relative paths: `/docs/rpc/methods/method-name`

### External Links
- **ALWAYS** link to GitHub source code for verification
- Use `main` branch: `https://github.com/NethermindEth/starknet.go/blob/main/...`
- Include line numbers: `#L245-L262`

## 7. Sidebar Configuration

### Adding New Pages
1. Edit `/docs/sidebar.ts`
2. Add entry in appropriate section
3. Use exact path: `/docs/section/page-name`
4. Examples should be at END of section

### Example Sidebar Entry
```typescript
{
  text: "MethodName",
  link: "/docs/rpc/methods/method-name"
}
```

## 8. Common Pitfalls to Avoid

### ❌ NEVER DO THIS:
- Write struct definitions without reading source code
- Guess field names or types
- Use AI to generate struct definitions
- Skip verification "because it looks right"
- Assume types are similar to other languages
- Copy examples from other SDKs without verification

### ✅ ALWAYS DO THIS:
- Read actual source code files
- Verify every field name
- Verify every type
- Include source links
- Test that examples would compile
- Copy JSON tags exactly

## 9. Quality Checklist

Before considering documentation complete:

- [ ] All struct definitions verified against source code
- [ ] Source links provided for all types
- [ ] Method signature verified against source
- [ ] All examples use correct field names
- [ ] Examples would compile if copy-pasted
- [ ] Imports are correct
- [ ] Parameters documented
- [ ] Return values documented
- [ ] Edge cases noted
- [ ] Related methods linked
- [ ] Added to sidebar if new page
- [ ] No hallucinated fields or types
- [ ] Comments are accurate and helpful

## 10. Example: Complete Verification Process

### Step 1: Identify Source Files
```bash
# Find where StorageProof is defined
grep -r "func.*StorageProof" /Users/aayushgiri/work/starknet.go/rpc/
# Found in: rpc/contract.go
```

### Step 2: Read Source Code
```bash
# Read the actual implementation
Read rpc/contract.go lines 245-262
Read rpc/types_contract.go for type definitions
```

### Step 3: Verify Types
```go
// From rpc/types_contract.go:66-78
type StorageProofInput struct {
    BlockID              BlockID                `json:"block_id"`
    ClassHashes          []*felt.Felt           `json:"class_hashes,omitempty"`
    ContractAddresses    []*felt.Felt           `json:"contract_addresses,omitempty"`
    ContractsStorageKeys []ContractStorageKeys  `json:"contracts_storage_keys,omitempty"`
}
```

### Step 4: Document with Source Links
```markdown
### StorageProofInput

\`\`\`go
type StorageProofInput struct {
    BlockID              BlockID                `json:"block_id"`
    ClassHashes          []*felt.Felt           `json:"class_hashes,omitempty"`
    ContractAddresses    []*felt.Felt           `json:"contract_addresses,omitempty"`
    ContractsStorageKeys []ContractStorageKeys  `json:"contracts_storage_keys,omitempty"`
}
\`\`\`

**Source:** [types_contract.go:L66-L78](https://github.com/NethermindEth/starknet.go/blob/main/rpc/types_contract.go#L66-L78)
```

### Step 5: Create Verified Example
```go
// Use VERIFIED field names from step 3
input := rpc.StorageProofInput{
    BlockID:           rpc.BlockID{Tag: "latest"},  // ✅ Correct field name
    ContractAddresses: []*felt.Felt{addr},          // ✅ Correct field name
}
// NOT:
input := rpc.StorageProofInput{
    BlockID:         rpc.BlockID{Tag: "latest"},
    ContractAddress: addr,  // ❌ WRONG - field doesn't exist!
}
```

## 11. Human-Toned Documentation

### Writing Style
- Write in clear, conversational English
- Explain WHY something is useful, not just WHAT it does
- Include practical use cases
- Add helpful comments in code examples
- Be concise but complete

### Example
❌ **Bad** (robotic):
```markdown
This method retrieves storage proofs.
```

✅ **Good** (human):
```markdown
Get merkle paths in one of the state tries: global state, classes, or individual contract storage. A single request can query for any mix of the three types of storage proofs (classes, contracts, and storage).
```

## 12. Version Control

- Always verify against the `main` branch
- GitHub source links should point to `main` branch
- Note the RPC version in specification section
- Update documentation when source code changes

## 13. Error Messages to Watch For

If you see these during development:
- `undefined: FieldName` - Field doesn't exist, verify struct definition
- `cannot use X as Y` - Type mismatch, verify types
- Build errors in examples - Example uses wrong types/fields

## Summary

The core principle is simple: **NEVER TRUST YOUR MEMORY OR AI GENERATION. ALWAYS VERIFY AGAINST ACTUAL SOURCE CODE.**

Every single struct, method, field, and type must be verified by reading the actual source code file and cross-referenced with a GitHub link. This is not optional.

This documentation represents real code that developers will copy and paste. If it doesn't compile, we've failed them and damaged the project's reputation.
