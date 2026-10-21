CREATE TABLE IF NOT EXISTS public.seed (
    is_seed BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO public.seed (is_seed) VALUES (FALSE);

GRANT SELECT, INSERT, UPDATE, DELETE ON public.seed TO gas_backoffice_service;
