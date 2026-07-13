ALTER TABLE public.audit_events RENAME COLUMN actor_id TO author_id;
ALTER TABLE public.audit_events RENAME COLUMN actor_role TO author_role;

ALTER TABLE public.outbox_event_bus RENAME COLUMN actor_id TO author_id;
ALTER TABLE public.outbox_event_bus RENAME COLUMN actor_role TO author_role;
