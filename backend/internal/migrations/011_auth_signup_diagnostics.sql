-- =====================================================================
-- Auth signup diagnostics
-- =====================================================================
-- Run these queries in the Supabase SQL editor if /auth/v1/signup returns:
-- "Database error saving new user".
--
-- The most common cause is a broken custom trigger on auth.users. The Oroya
-- backend creates public.profiles itself, so an auth.users -> profiles trigger
-- is not required by this backend.

-- 1) Inspect non-internal triggers on auth.users.
select
  t.tgname as trigger_name,
  pg_get_triggerdef(t.oid) as trigger_definition
from pg_trigger t
join pg_class c on c.oid = t.tgrelid
join pg_namespace n on n.oid = c.relnamespace
where n.nspname = 'auth'
  and c.relname = 'users'
  and not t.tgisinternal;

-- 2) If you see a custom trigger such as on_auth_user_created that inserts
-- into public.profiles and it is failing, remove it. Adjust names if your
-- diagnostic query shows different trigger/function names.
--
-- drop trigger if exists on_auth_user_created on auth.users;
-- drop function if exists public.handle_new_user();

-- 3) If this deployment only has SUPABASE_ANON_KEY and the key is kept only on
-- the backend, allow the backend to write profiles through PostgREST.
--
-- drop policy if exists profiles_insert_backend_anon on profiles;
-- create policy profiles_insert_backend_anon on profiles
--     for insert to anon, authenticated
--     with check (true);
--
-- drop policy if exists profiles_update_backend_anon on profiles;
-- create policy profiles_update_backend_anon on profiles
--     for update to anon, authenticated
--     using (true)
--     with check (true);
