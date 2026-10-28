# Authorization and subscription access

This document describes the authorization model used by the backoffice. It is
the reference for the relationship between system roles, subscription roles,
memberships, JWT claims, and subscription endpoints.

## Two authorization scopes

Authorization has two separate scopes:

1. **System scope** answers what the user can do across the whole backoffice.
2. **Subscription scope** answers what the user can do inside one subscription.

These scopes must not be collapsed into one role. A user can have one global
system role and different subscription memberships as the subscription context
changes.

## System roles

The system role is global to the user:

| System role | Meaning |
| --- | --- |
| `system_administrator` | Has global visibility and can perform system-level operations. |
| `subscription_user` | Has no global access; access comes from subscription memberships. |

`is_system` is not an authorization role. It describes that a record is
managed or protected by the system.

User provenance is stored on the `users` row itself:

```text
users.author_id
users.author_role
```

These fields identify who created the user and which role the creator had.
They are authorship metadata, not subscription access and not tenancy.

System administrators do not need a membership for every subscription. Their
global role grants access implicitly. This avoids creating one membership per
administrator and subscription and keeps the membership relation scoped to
users who need subscription access.

## Subscription roles and memberships

A membership belongs to the relationship between a user and a subscription:

```text
subscription_members(subscription_id, user_id, subscription_role)
```

The current subscription role is:

```text
user
```

The role is kept on the membership because future subscription-scoped roles
may include `tenant`, `viewer`, or `auditor`.

Membership existence and membership role are distinct concepts:

- membership existence grants access to the subscription;
- membership role controls what the user can do inside it.

System administrators may access subscriptions without a membership. General
users must have a membership for the requested subscription.

## JWT contract

JWTs expose the global role as `system_role`. The old generic `role` claim is
not part of the contract.

### System administrator token

```json
{
  "sub": "user-system-admin",
  "system_role": "system_administrator"
}
```

The token does not contain `memberships` or `subscription_ids`. The list of
subscriptions is deliberately not copied into a global administrator token.

### General user token

```json
{
  "sub": "user-general",
  "system_role": "subscription_user",
  "memberships": [
    {
      "subscription_id": "subscription-a",
      "subscription_role": "user"
    }
  ]
}
```

The token contains the memberships used by scoped authorization. Membership
changes take effect in a newly issued or refreshed token.

Migration `20260101000018_rename_general_user_to_subscription_user` converts the legacy database
values `admin` and `tenant` to `system_administrator` and `subscription_user`.
New persistence and API code must use the public names documented here.

## Subscription listing

`GET /v1/backoffice/subscriptions` is available to both system roles.

- `system_administrator`: returns all non-deleted subscriptions.
- `subscription_user`: returns only subscriptions whose IDs are present in the
  user's memberships.

The endpoint remains available for a consistent server-side view even though
general users already carry their membership IDs in the JWT.

## Scoped subscription resources

Routes containing `:subscription_id` use the following authorization rule:

1. A system administrator is allowed by the global role.
2. A general user must have a membership for the requested subscription.
3. The membership role is used for subscription-scoped permissions.

This rule belongs in the authorization middleware or policy layer, rather
than being repeated as ad-hoc role checks in every handler.

## Events and projections

Subscription membership changes are domain events. The outbox delivers those
events to the user projection, which maintains `users.memberships` as a read
model containing the subscription ID and membership role.

The projection is not the source of truth. `subscription_members` is the
subscription aggregate's persistence model; `users.memberships` exists to
support authentication and read access efficiently.

## Deliberate non-goals

The current model does not include:

- subscription administrator roles;
- a change-member-role use case;
- ownership transfer;
- memberships for system administrators;
- a system administrator role embedded in a subscription membership.

