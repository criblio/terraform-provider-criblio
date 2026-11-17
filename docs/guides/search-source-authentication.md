# Search source authentication and provider upgrades

Existing `criblio_search_source` state retains its schema and token values when
upgrading. If `auth_tokens` was never configured, the provider continues to read
those credentials from Cribl. An unchanged configuration does not clear tokens,
and unrelated updates preserve authentication using the current API representation,
because the v1 API replaces the complete source. Updates stop before writing if
the API masks credentials that the provider cannot safely preserve.

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

Raw HTTP sends inline tokens and stored-secret references through `authTokensExt` and clears both extended and
legacy lists when removing authentication. Cribl HTTP uses stored-secret IDs;
Splunk HEC supports inline tokens and stored-secret references. Sensitive token values remain in Terraform state.

For Raw HTTP and Splunk HEC, configure exactly one of `token` (an inline value) or `token_secret`
(a stored text-secret ID) in each entry. Existing inline configurations continue
to work. For example:

```hcl
auth_tokens = [{
  token_secret = "my-http-text-secret"
  description  = "Stored authentication credential"
}]
```

The provider sends `authType = "secret"` for stored-secret references and keeps
the reference when reading existing sources. Resolved runtime values returned by
the API do not replace the reference or become inline tokens in Terraform state.

Ownership is stored in Terraform's provider-private state. Existing states without
this marker preserve API-managed tokens. Keep normal Terraform state and its
private data when moving resources; importing a resource resets ownership to the
protected, API-managed behavior until configured tokens are applied.

Raw HTTP stored-secret authentication requires Cribl newer than 4.19.2. The
referenced text secret must be resolvable by the Search ingest group (`search`);
a secret in another Worker Group is not automatically available there. The
provider cannot verify availability from the reference ID alone.

When Cribl masks a global secret ID as `****/******`, refresh retains a known
configured ID. Imported sources without a known ID need access to the secret
folder or explicitly configured credentials before they can be safely updated.
