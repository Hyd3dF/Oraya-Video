-- =====================================================================
-- Storage + video upload policies for anon-key backend deployments
-- =====================================================================
-- Run this in Supabase SQL Editor if video upload fails with:
-- "new row violates row-level security policy"
--
-- This project keeps SUPABASE_ANON_KEY on the Go backend only. The frontend
-- must never receive it. These policies allow the backend's anon-key requests
-- to create signed upload URLs, store processed HLS assets, and write video
-- metadata.

-- Ensure buckets exist.
insert into storage.buckets (id, name, public)
values
  ('raw-videos', 'raw-videos', false),
  ('hls-videos', 'hls-videos', true),
  ('thumbnails', 'thumbnails', true)
on conflict (id) do nothing;

-- Storage object policies, limited to Oroya buckets.
drop policy if exists oroya_storage_select on storage.objects;
create policy oroya_storage_select on storage.objects
  for select to anon, authenticated
  using (bucket_id in ('raw-videos', 'hls-videos', 'thumbnails'));

drop policy if exists oroya_storage_insert on storage.objects;
create policy oroya_storage_insert on storage.objects
  for insert to anon, authenticated
  with check (bucket_id in ('raw-videos', 'hls-videos', 'thumbnails'));

drop policy if exists oroya_storage_update on storage.objects;
create policy oroya_storage_update on storage.objects
  for update to anon, authenticated
  using (bucket_id in ('raw-videos', 'hls-videos', 'thumbnails'))
  with check (bucket_id in ('raw-videos', 'hls-videos', 'thumbnails'));

drop policy if exists oroya_storage_delete on storage.objects;
create policy oroya_storage_delete on storage.objects
  for delete to anon, authenticated
  using (bucket_id in ('raw-videos', 'hls-videos', 'thumbnails'));

-- Video metadata policies needed by POST /api/v1/videos and the worker.
drop policy if exists videos_insert_backend_anon on public.videos;
create policy videos_insert_backend_anon on public.videos
  for insert to anon, authenticated
  with check (true);

drop policy if exists videos_update_backend_anon on public.videos;
create policy videos_update_backend_anon on public.videos
  for update to anon, authenticated
  using (true)
  with check (true);

drop policy if exists videos_delete_backend_anon on public.videos;
create policy videos_delete_backend_anon on public.videos
  for delete to anon, authenticated
  using (true);

drop policy if exists video_assets_insert_backend_anon on public.video_assets;
create policy video_assets_insert_backend_anon on public.video_assets
  for insert to anon, authenticated
  with check (true);

drop policy if exists video_assets_update_backend_anon on public.video_assets;
create policy video_assets_update_backend_anon on public.video_assets
  for update to anon, authenticated
  using (true)
  with check (true);
