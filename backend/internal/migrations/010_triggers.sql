-- =====================================================================
-- Triggers that keep denormalized counters in sync with their child tables.
--
-- Run this AFTER init.sql. Backend code only inserts/deletes rows in
-- video_likes / comment_likes / views — these triggers update the parent
-- row's counter automatically. This way the Go backend never has to do
-- read-modify-write on a counter (race-condition-prone over REST).
-- =====================================================================

-- ---------------------------------------------------------------------
-- videos.likes_count <- video_likes
-- ---------------------------------------------------------------------
create or replace function update_video_likes_count() returns trigger language plpgsql as $$
begin
  if tg_op = 'INSERT' then
    update videos set likes_count = likes_count + 1 where id = new.video_id;
  elsif tg_op = 'DELETE' then
    update videos set likes_count = greatest(likes_count - 1, 0) where id = old.video_id;
  end if;
  return null;
end; $$;

drop trigger if exists trg_video_likes_count on video_likes;
create trigger trg_video_likes_count
after insert or delete on video_likes
for each row execute function update_video_likes_count();

-- ---------------------------------------------------------------------
-- comments.likes_count <- comment_likes
-- ---------------------------------------------------------------------
create or replace function update_comment_likes_count() returns trigger language plpgsql as $$
begin
  if tg_op = 'INSERT' then
    update comments set likes_count = likes_count + 1 where id = new.comment_id;
  elsif tg_op = 'DELETE' then
    update comments set likes_count = greatest(likes_count - 1, 0) where id = old.comment_id;
  end if;
  return null;
end; $$;

drop trigger if exists trg_comment_likes_count on comment_likes;
create trigger trg_comment_likes_count
after insert or delete on comment_likes
for each row execute function update_comment_likes_count();

-- ---------------------------------------------------------------------
-- videos.views_count <- views
--
-- Backend deduplicates before inserting (per user/IP within the last hour),
-- so every inserted row corresponds to one counted view.
-- ---------------------------------------------------------------------
create or replace function update_video_views_count() returns trigger language plpgsql as $$
begin
  update videos set views_count = views_count + 1 where id = new.video_id;
  return null;
end; $$;

drop trigger if exists trg_video_views_count on views;
create trigger trg_video_views_count
after insert on views
for each row execute function update_video_views_count();
