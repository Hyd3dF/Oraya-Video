-- =====================================================================
-- Oroya — initial database schema
-- =====================================================================
-- Single-file consolidation of 001..009 migrations. Run this once in the
-- Supabase SQL editor (or via psql) against the project database.
--
-- Column names and types here match exactly what the Go repository layer
-- expects (backend/internal/repository/*.go). Do not rename anything in
-- this file without updating the corresponding repository code.
--
-- The Go backend connects with the service-role key, which bypasses RLS.
-- RLS is still enabled on every table as defense-in-depth, with public
-- SELECT policies only where the API exposes the data anonymously.
-- =====================================================================

create extension if not exists pgcrypto;

-- ---------------------------------------------------------------------
-- 001 profiles
-- ---------------------------------------------------------------------
create table if not exists profiles (
    id            uuid primary key references auth.users(id) on delete cascade,
    real_name     text,
    username      text not null unique,
    display_name  text,
    email         text not null unique,
    avatar_url    text,
    banner_url    text,
    bio           text,
    login_type    text not null default 'email' check (login_type in ('email', 'google')),
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now()
);

create index if not exists idx_profiles_username on profiles (username);

alter table profiles enable row level security;
drop policy if exists profiles_select_public on profiles;
create policy profiles_select_public on profiles for select using (true);

-- ---------------------------------------------------------------------
-- 002 videos
-- ---------------------------------------------------------------------
create table if not exists videos (
    id               uuid primary key default gen_random_uuid(),
    owner_id         uuid not null references profiles(id) on delete cascade,
    title            text not null,
    description      text,
    thumbnail_url    text,
    duration_seconds int,
    views_count      bigint not null default 0,
    likes_count      bigint not null default 0,
    visibility       text not null default 'public' check (visibility in ('public', 'unlisted', 'private')),
    status           text not null default 'processing' check (status in ('processing', 'ready', 'failed')),
    created_at       timestamptz not null default now(),
    updated_at       timestamptz not null default now()
);

create index if not exists idx_videos_owner   on videos (owner_id);
create index if not exists idx_videos_status  on videos (status);
create index if not exists idx_videos_created on videos (created_at desc);

alter table videos enable row level security;
drop policy if exists videos_select_public on videos;
create policy videos_select_public on videos for select
    using (visibility = 'public' and status = 'ready');

-- ---------------------------------------------------------------------
-- 003 video_assets (HLS renditions)
-- ---------------------------------------------------------------------
create table if not exists video_assets (
    id           uuid primary key default gen_random_uuid(),
    video_id     uuid not null references videos(id) on delete cascade,
    quality      text not null,
    playlist_url text not null,
    master_url   text,
    width        int,
    height       int,
    bitrate      int,
    size_bytes   bigint,
    created_at   timestamptz not null default now()
);

create index if not exists idx_video_assets_video on video_assets (video_id);
-- ON CONFLICT (video_id, quality) target in repository.AddAsset uses this index.
create unique index if not exists uq_video_assets_quality on video_assets (video_id, quality);

alter table video_assets enable row level security;
drop policy if exists video_assets_select_public on video_assets;
create policy video_assets_select_public on video_assets for select using (true);

-- ---------------------------------------------------------------------
-- 004 video_likes
-- ---------------------------------------------------------------------
create table if not exists video_likes (
    id         bigserial primary key,
    video_id   uuid not null references videos(id) on delete cascade,
    user_id    uuid not null references profiles(id) on delete cascade,
    created_at timestamptz not null default now(),
    unique (video_id, user_id)
);

create index if not exists idx_video_likes_video on video_likes (video_id);
create index if not exists idx_video_likes_user  on video_likes (user_id);

alter table video_likes enable row level security;

-- ---------------------------------------------------------------------
-- 005 comments
-- ---------------------------------------------------------------------
create table if not exists comments (
    id          uuid primary key default gen_random_uuid(),
    video_id    uuid not null references videos(id) on delete cascade,
    user_id     uuid not null references profiles(id) on delete cascade,
    parent_id   uuid references comments(id) on delete cascade,
    content     text not null,
    likes_count bigint not null default 0,
    created_at  timestamptz not null default now()
);

create index if not exists idx_comments_video  on comments (video_id, created_at desc);
create index if not exists idx_comments_parent on comments (parent_id);

alter table comments enable row level security;
drop policy if exists comments_select_public on comments;
create policy comments_select_public on comments for select using (true);

-- ---------------------------------------------------------------------
-- 006 comment_likes
-- ---------------------------------------------------------------------
create table if not exists comment_likes (
    id         bigserial primary key,
    comment_id uuid not null references comments(id) on delete cascade,
    user_id    uuid not null references profiles(id) on delete cascade,
    created_at timestamptz not null default now(),
    unique (comment_id, user_id)
);

create index if not exists idx_comment_likes_comment on comment_likes (comment_id);

alter table comment_likes enable row level security;

-- ---------------------------------------------------------------------
-- 007 subscriptions
-- ---------------------------------------------------------------------
create table if not exists subscriptions (
    id            bigserial primary key,
    subscriber_id uuid not null references profiles(id) on delete cascade,
    channel_id    uuid not null references profiles(id) on delete cascade,
    created_at    timestamptz not null default now(),
    unique (subscriber_id, channel_id),
    check (subscriber_id <> channel_id)
);

create index if not exists idx_subs_channel    on subscriptions (channel_id);
create index if not exists idx_subs_subscriber on subscriptions (subscriber_id);

alter table subscriptions enable row level security;

-- ---------------------------------------------------------------------
-- 008 views
-- ---------------------------------------------------------------------
create table if not exists views (
    id         bigserial primary key,
    video_id   uuid not null references videos(id) on delete cascade,
    user_id    uuid references profiles(id) on delete set null,
    ip_hash    text,
    created_at timestamptz not null default now()
);

create index if not exists idx_views_video      on views (video_id, created_at desc);
create index if not exists idx_views_video_user on views (video_id, user_id);

alter table views enable row level security;

-- ---------------------------------------------------------------------
-- 009 video_links
-- ---------------------------------------------------------------------
create table if not exists video_links (
    id         uuid primary key default gen_random_uuid(),
    video_id   uuid not null references videos(id) on delete cascade,
    title      text,
    url        text not null,
    sort_order int not null default 0,
    created_at timestamptz not null default now()
);

create index if not exists idx_video_links_video on video_links (video_id, sort_order);

alter table video_links enable row level security;
drop policy if exists video_links_select_public on video_links;
create policy video_links_select_public on video_links for select using (true);
