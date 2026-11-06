# Terraform application resources

Terraform manages the DigitalOcean App Platform application and its runtime
database user. The database and the `gas_admin` permissions are managed centrally
in the `gas` repository. Database object privileges are applied by the
application's SQL migrations, which run in the app's pre-deploy job.

This module must not perform an automatic `destroy` of the existing database.

## One-time database state transition

Before running the plan, remove only the database resource from Terraform state:

```sh
terraform state rm digitalocean_database_db.gas_backoffice
```

This removes the resource from state without deleting the database from
DigitalOcean. Then validate the configuration with:

```sh
terragrunt plan
```

The database and `gas_admin` permissions remain centralized in the `gas`
repository; `gas-backoffice` manages only the runtime user and the application.

## One-time transition from `postgresql_grant`

After checking out the change, run this from the GAS infrastructure repository
in `terragrunt/base/backoffice`, with the environment that currently tracks the
grants. Terragrunt selects the remote state backend; the current dev backend key
is `dev/backoffice/terraform.tfstate` and its Terraform workspace is `default`.
For dev:

```sh
cd gas/terragrunt/base/backoffice
ENV=dev terragrunt state rm \
  'postgresql_grant.gas_backoffice_database_privileges' \
  'postgresql_grant.gas_backoffice_schema_privileges' \
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
the app user the existing CRUD and sequence access. Default privileges are
configured for objects created by `gas_admin`, the account used by the
pre-deploy job.
