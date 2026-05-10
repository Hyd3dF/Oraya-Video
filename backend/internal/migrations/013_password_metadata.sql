-- =====================================================================
-- Password metadata
-- =====================================================================
-- Do NOT store raw user passwords in public tables.
--
-- Supabase Auth stores password hashes internally in auth.users
-- (encrypted_password). The application should never copy or expose that
-- value. This migration only keeps safe metadata on public.profiles.

alter table public.profiles
  add column if not exists password_updated_at timestamptz;

create or replace function public.sync_profile_password_updated_at()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
begin
  if new.encrypted_password is distinct from old.encrypted_password then
    update public.profiles
       set password_updated_at = now(),
           updated_at = now()
     where id = new.id;
  end if;

  return new;
end;
$$;

drop trigger if exists on_auth_user_password_updated on auth.users;

create trigger on_auth_user_password_updated
after update of encrypted_password on auth.users
for each row
execute function public.sync_profile_password_updated_at();
