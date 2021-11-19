# vault-aws-credential-helper

The Vault AWS Credential Helper is a component that can be injected
into a task environment and be used as a credential helper process for
the AWS SDK.  More details about the AWS configuration can be found on
[this
page](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html).


You must export `VACH_VAULT_BLOB` to the environment as a path that
points to the JSON blob from the Vault AWS Secrets backend.

In your Vault Template file, you should specify the secret like below
to ensure the JSON winds up in the right shape:

```
{{secret "aws/creds/my-app" | toJSON}}
```
