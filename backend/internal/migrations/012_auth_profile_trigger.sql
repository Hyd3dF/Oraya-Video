-- =====================================================================
-- Auth signup profile trigger
-- =====================================================================
-- Fixes Supabase Auth signup failures caused by a broken/missing
-- auth.users -> public.profiles trigger.
--
-- Run this in the Supabase SQL editor. It creates a profile row from the
-- metadata sent by the Go backend during /auth/v1/signup:
--   data.username
--   data.real_name

create or replace function public.handle_new_user()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
declare
  v_username text;
  v_real_name text;
begin
  v_username := nullif(new.raw_user_meta_data ->> 'username', '');
  v_real_name := nullif(new.raw_user_meta_data ->> 'real_name', '');

  if v_username is null then
    v_username := 'user_' || substr(replace(new.id::text, '-', ''), 1, 12);
  end if;

  insert into public.profiles (
    id,
    email,
    username,
    real_name,
    display_name,
    login_type
  )
  values (
    new.id,
    new.email,
    v_username,
    v_real_name,
    coalesce(v_real_name, v_username),
    'email'
  )
  on conflict (id) do update set
    email = excluded.email,
    username = excluded.username,
    real_name = excluded.real_name,
    display_name = excluded.display_name,
    updated_at = now();

  return new;
end;
$$;

drop trigger if exists on_auth_user_created on auth.users;

create trigger on_auth_user_created
after insert on auth.users
for each row
execute function public.handle_new_user();
