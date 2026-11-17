# Search source authentication and provider upgrades

Existing `criblio_search_source` state retains its schema and token values when
upgrading. If `auth_tokens` was never configured, the provider continues to read
those credentials from Cribl. An unchanged configuration does not clear tokens,
and unrelated updates omit the authentication fields from API requests.

Imported sources receive the same protection. API responses may contain masked or
encrypted tokens; plaintext cannot be recovered from those responses.

To take ownership of legacy or imported credentials, configure the desired token
list and apply a resource update. A no-change apply cannot record ownership; if
configuring the existing list produces no update, removing it still preserves the
legacy credentials. Use an explicit empty list to clear them in that case.

To remove all credentials explicitly, configure
`auth_tokens = []` and apply. After configured tokens are written, removing the
attribute or setting it to `null` clears them. Sources newly created by this
provider also use this removal behavior.

Raw HTTP sends inline tokens through `authTokensExt` and clears both extended and
legacy lists when removing authentication. Cribl HTTP uses stored-secret IDs;
Splunk HEC uses inline tokens. Sensitive token values remain in Terraform state.

Ownership is stored in Terraform's provider-private state. Existing states without
this marker preserve API-managed tokens. Keep normal Terraform state and its
private data when moving resources; importing a resource resets ownership to the
protected, API-managed behavior until configured tokens are applied.
