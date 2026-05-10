-- Makes source uploads playable when FFmpeg is not installed.
-- The backend still uses HLS when FFmpeg is available.

update storage.buckets
set public = true
where id = 'raw-videos';
