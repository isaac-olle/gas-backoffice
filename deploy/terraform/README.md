# Terraform database resources

Terraform manages the DigitalOcean PostgreSQL cluster users and database. Database
object privileges are managed by the application's SQL migrations, which run in
the app's pre-deploy job using the cluster administrator credentials supplied by
DigitalOcean. That account is required to grant privileges on the database and
`public` schema. Terraform does not need a PostgreSQL provider connection.

## One-time transition from `postgresql_grant`

After checking out the change, run this from the GAS infrastructure repository
in `terragrunt/base/backoffice`, with the environment that currently tracks the
grants. Terragrunt selects the remote state backend; the current dev backend key
is `dev/backoffice/terraform.tfstate` and its Terraform workspace is `default`.
For dev:

```sh
cd gas/terragrunt/base/backoffice
ENV=dev terragrunt state rm \
  'postgresql_grant.gas_backoffice_admin_database_privileges' \
  'postgresql_grant.gas_backoffice_admin_schema_privileges' \
  'postgresql_grant.gas_backoffice_service_database_privileges' \
  'postgresql_grant.gas_backoffice_service_tables_privileges' \
  'postgresql_grant.gas_backoffice_service_sequences_privileges'
ENV=dev terragrunt state list
```

`terraform state rm` forgets these resources without issuing PostgreSQL `REVOKE`
statements. Run it before `terraform plan`/`apply` with this configuration, so
Terraform cannot propose destroying the grants it no longer manages. It does not
modify the users or database in Terraform state.

The first deployment applies the new permissions migration. It is safe if the
existing grants are already present. The migration gives the externally managed
`gas_admin` role the migration-level database and schema privileges, and grants
the app user the existing CRUD and sequence access. Default table privileges are
configured for objects created by `doadmin`, the account used by the pre-deploy
job.

## One-time database-user transition

The old `gas_backoffice_admin` DigitalOcean user is no longer managed here. If it
is still in this Terraform state, forget it without deleting it; object ownership
may still depend on that role. The service user resource keeps its Terraform
address and is replaced with `gas-backoffice-user` using create-before-destroy.
Run this from `deploy/terraform` against the same workspace/backend before
planning:

```sh
terraform state rm 'digitalocean_database_user.gas_backoffice_admin'
```

The old admin user will remain in DigitalOcean until it is deliberately retired.
