# Oroya - Video Sharing Platform Plan

## Project Overview

**Name:** Oroya  
**Type:** Mobile-first video sharing platform (YouTube-like)  
**Backend:** Go (Golang) — deployed separately at `https://api.oroya.xyz`  
**Frontend:** Svelte 5 + SvelteKit — deployed separately at `https://oroya.xyz`  
**Database:** Self-hosted Supabase (PostgreSQL for data; Auth for user management; Storage for file storage)  

**Base URLs:**
- Frontend: `https://oroya.xyz`
- Backend API: `https://api.oroya.xyz`

**Frontend Environment:**
```
PUBLIC_API_BASE_URL=https://api.oroya.xyz
```

**Deployment Rule:**
- `/backend` and `/frontend` are completely separate codebases.
- They are built and deployed independently.
- The frontend does not know anything about the database, Supabase service keys, or internal backend structure.

---

## Backend / Frontend Separation Rule

This is the most important rule of the project:

### Frontend = ONLY User Interface
The frontend is a thin client. It does not contain any business logic, database knowledge, or secrets.

**What the frontend DOES:**
- Render pages, components, buttons, forms, video player
- Collect user input (forms, clicks, file selection)
- Send HTTP requests to the Go backend API at `PUBLIC_API_BASE_URL`
- Receive responses and display them to the user
- Handle client-side navigation and UI state

**What the frontend DOES NOT do:**
- **Does NOT write to the database directly**
- **Does NOT know database table names, columns, or schema**
- **Does NOT contain Supabase service-role key or any secret API key**
- **Does NOT process videos**
- **Does NOT handle authentication logic** (only sends email/password to backend, stores returned token)
- **Does NOT generate signed URLs**
- **Does NOT query Supabase Storage directly** (only uses URLs given by backend)

### Backend = ALL Business Logic
The backend is the brain of the application. It handles everything except rendering UI.

**What the backend does:**
- **User System:** registration, login, logout, token refresh, token validation, profile management
- **Video System:** metadata CRUD, upload control, compression, HLS splitting, link management
- **Comment System:** create, delete, reply, list comments
- **Like System:** like/unlike videos, like/unlike comments
- **View System:** count video views
- **Search & Channels:** search indexing, channel info, subscriptions
- **All database queries** (the only service that talks to PostgreSQL)
- **Admin / monitoring dashboard**

### Supabase = Infrastructure Only
Supabase is used purely as infrastructure. The application logic lives in the Go backend.

- **PostgreSQL** = data persistence. Accessed ONLY by Go backend via service-role key.
- **Auth** = user identity provider. Backend validates tokens; frontend only stores them.
- **Storage** = file storage. Backend generates signed upload URLs. Files are served directly from CDN. Frontend never talks to Storage API directly.

---

## Design Philosophy

- **Color Palette:**
  - **Background:** `#F5F0EB` (warm, muted off-white — low brightness, easy on the eyes)
  - **Primary Accent:** `#E86A33` (warm orange for buttons, active states, icons)
  - **Secondary Accent:** `#D35400` (deeper orange for hover states)
  - **Text Primary:** `#2C2C2C` (dark gray, not pure black)
  - **Text Secondary:** `#6B6B6B` (medium gray for descriptions, timestamps)
  - **Surface/Cards:** `#FFFFFF` (slightly brighter than background for elevation)
  - **Borders:** `#E0D8D0` (subtle warm gray borders)
  - **Success:** `#27AE60`
  - **Error:** `#C0392B`

- **Typography:**
  - System font stack with Inter or Geist as primary
  - Highly legible, clean sans-serif
  - Large touch targets for mobile (min 44px)

- **Layout:**
  - Mobile-first responsive design
  - Bottom navigation bar on mobile (Home, Shorts, Upload, Subscriptions, You)
  - Sidebar/drawer on tablet/desktop
  - Video grid with 1 column on mobile, 2-3 on tablet, 4-5 on desktop

---

## Architecture

```
oroya/
├── backend/                    # Go (Golang) backend — completely separate
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go         # HTTP API server entry point
│   │   └── worker/
│   │       └── main.go         # Background video processing worker (FFmpeg)
│   ├── internal/
│   │   ├── config/             # Environment & config management
│   │   ├── api/                # HTTP handlers (REST API)
│   │   │   ├── auth.go
│   │   │   ├── videos.go
│   │   │   ├── comments.go
│   │   │   ├── channels.go
│   │   │   ├── search.go
│   │   │   ├── upload.go
│   │   │   └── admin.go        # Admin dashboard handlers
│   │   ├── middleware/         # Auth, CORS, logging, rate limiting
│   │   ├── models/             # Data structures / DTOs
│   │   ├── services/           # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   ├── video_service.go
│   │   │   ├── user_service.go
│   │   │   ├── comment_service.go
│   │   │   ├── search_service.go
│   │   │   ├── upload_service.go
│   │   │   └── admin_service.go
│   │   ├── repository/         # Supabase/Postgres data access layer ONLY
│   │   │   ├── auth_repo.go
│   │   │   ├── video_repo.go
│   │   │   ├── user_repo.go
│   │   │   └── comment_repo.go
│   │   ├── migrations/         # Database schema migrations (SQL files)
│   │   │   ├── 001_profiles.sql
│   │   │   ├── 002_videos.sql
│   │   │   ├── 003_video_assets.sql
│   │   │   ├── 004_video_likes.sql
│   │   │   ├── 005_comments.sql
│   │   │   ├── 006_comment_likes.sql
│   │   │   ├── 007_subscriptions.sql
│   │   │   ├── 008_views.sql
│   │   │   └── 009_video_links.sql
│   │   ├── worker/             # Background job processing
│   │   │   ├── ffmpeg.go       # FFmpeg HLS transcoding
│   │   │   ├── queue.go        # Job queue (in-memory or Redis)
│   │   │   └── processor.go    # Video processing orchestrator
│   │   └── utils/              # Helpers, validators, JWT validation
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Svelte frontend — completely separate
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/     # Reusable UI components
│   │   │   ├── stores/         # Svelte stores (auth, theme, player)
│   │   │   ├── api/            # API client — ONLY calls Go backend
│   │   │   │   ├── auth.js
│   │   │   │   ├── videos.js
│   │   │   │   ├── comments.js
│   │   │   │   ├── channels.js
│   │   │   │   └── search.js
│   │   │   └── utils/          # Helpers, formatters
│   │   ├── routes/
│   │   │   ├── +page.svelte    # Home / Feed
│   │   │   ├── watch/
│   │   │   │   └── [id]/
│   │   │   │       └── +page.svelte   # Video player page (HLS)
│   │   │   ├── upload/
│   │   │   │   └── +page.svelte       # Video upload page
│   │   │   ├── channel/
│   │   │   │   └── [id]/
│   │   │   │       └── +page.svelte   # User channel page
│   │   │   ├── search/
│   │   │   │   └── +page.svelte       # Search results
│   │   │   └── login/
│   │   │       └── +page.svelte       # Auth page
│   │   ├── app.html
│   │   └── app.css             # Global styles with orange/off-white theme
│   ├── static/
│   ├── .env.example            # PUBLIC_API_BASE_URL only
│   ├── svelte.config.js
│   ├── vite.config.js
│   └── package.json
│
└── oroya.md                    # This plan
```

**Important:** The frontend does not have a `migrations/` folder. It does not know table names, column names, or relationships. The backend owns the entire database schema.

---

## Database Schema (Backend Migrations Only)

**Important:** The database is ONLY accessed by the Go backend. Frontend NEVER talks directly to Postgres. RLS policies will still be enabled as defense-in-depth, but all application queries go through the Go API. The schema lives in `backend/internal/migrations/`.

### Tables

**1. `profiles`** (extends Supabase auth.users)

Stores application-level user information. Authentication (password hashes) is handled by Supabase Auth; this table stores profile data.

| Column        | Type        | Constraints                 |
|---------------|-------------|-----------------------------|
| id            | uuid        | PK, refs auth.users         |
| real_name     | text        | user's real/full name       |
| username      | text        | unique, not null, in-app display handle |
| display_name  | text        | optional custom display name |
| email         | text        | unique, not null            |
| avatar_url    | text        | profile picture URL         |
| banner_url    | text        | profile banner URL          |
| bio           | text        | user biography              |
| login_type    | text        | default 'email', options: 'email', 'google' |
| created_at    | timestamptz | default now()               |
| updated_at    | timestamptz | default now()               |

**Registration Flow:**
- **Email registration:** User sends email + password + real_name + username to `POST /api/v1/auth/register`. Backend creates user in Supabase Auth (password is hashed by Supabase), then inserts profile row with `login_type = 'email'`.
- **Google registration:** User authenticates with Google OAuth. Backend receives Google profile info, creates user in Supabase Auth, inserts profile row with `login_type = 'google'` and `email` from Google account.

**2. `videos`**

| Column          | Type        | Constraints                 |
|-----------------|-------------|-----------------------------|
| id              | uuid        | PK, default gen_random_uuid() |
| owner_id        | uuid        | FK → profiles.id            |
| title           | text        | not null                    |
| description     | text        | video description           |
| thumbnail_url   | text        | (Supabase Storage CDN URL)  |
| duration_seconds| int         |                             |
| views_count     | bigint      | default 0                   |
| likes_count     | bigint      | default 0                   |
| visibility      | text        | default 'public'            |
| status          | text        | default 'processing', options: 'processing', 'ready', 'failed' |
| created_at      | timestamptz | default now()               |
| updated_at      | timestamptz | default now()               |

**3. `video_assets`** (HLS renditions / quality variants)

| Column        | Type        | Constraints                           |
|---------------|-------------|---------------------------------------|
| id            | uuid        | PK, default gen_random_uuid()         |
| video_id      | uuid        | FK → videos.id, on delete cascade     |
| quality       | text        | not null, e.g. '360p', '720p', '1080p'|
| playlist_url  | text        | not null (path to .m3u8 file)         |
| master_url    | text        | (path to master.m3u8 if applicable)   |
| width         | int         |                                       |
| height        | int         |                                       |
| bitrate       | int         | bits per second                       |
| size_bytes    | bigint      | total size of all segments            |
| created_at    | timestamptz | default now()                         |

**MVP Note:** For the initial MVP, only `720p` quality will be generated. The `video_assets` table is designed so that `360p` and `1080p` can be added later without schema changes.

**4. `video_links`** (links shown under video description)

| Column     | Type        | Constraints                 |
|------------|-------------|-----------------------------|
| id         | uuid        | PK, default gen_random_uuid() |
| video_id   | uuid        | FK → videos.id, on delete cascade |
| title      | text        | link display title          |
| url        | text        | link URL, not null          |
| sort_order | int         | default 0                   |
| created_at | timestamptz | default now()               |

**Managed by backend:** Adding, removing, and listing links under a video description are all handled by the backend API. The frontend only sends requests.

**5. `video_likes`**

| Column     | Type        | Constraints                 |
|------------|-------------|-----------------------------|
| id         | bigint      | PK, auto-increment          |
| video_id   | uuid        | FK → videos.id, on delete cascade |
| user_id    | uuid        | FK → profiles.id, on delete cascade |
| created_at | timestamptz | default now()               |
| **Unique** | (video_id, user_id) |                     |

**Managed by backend:** Like and unlike operations are handled by the backend. The frontend sends `POST /api/v1/videos/:id/like` and the backend toggles the like state in the database.

**6. `comments`**

| Column     | Type        | Constraints                 |
|------------|-------------|-----------------------------|
| id         | uuid        | PK, default gen_random_uuid() |
| video_id   | uuid        | FK → videos.id              |
| user_id    | uuid        | FK → profiles.id            |
| parent_id  | uuid        | FK → comments.id (replies)  |
| content    | text        | not null                    |
| likes_count| bigint      | default 0                   |
| created_at | timestamptz | default now()               |

**Managed by backend:**
- Adding comments: `POST /api/v1/videos/:id/comments`
- Deleting comments: `DELETE /api/v1/comments/:id`
- Replying to comments: `POST /api/v1/videos/:id/comments` with `parent_id`
- All database operations are performed by the backend.

**7. `comment_likes`**

| Column     | Type        | Constraints                 |
|------------|-------------|-----------------------------|
| id         | bigint      | PK, auto-increment          |
| comment_id | uuid        | FK → comments.id, on delete cascade |
| user_id    | uuid        | FK → profiles.id, on delete cascade |
| created_at | timestamptz | default now()               |
| **Unique** | (comment_id, user_id) |                   |

**Managed by backend:** Like and unlike on comments are handled by the backend. The frontend sends `POST /api/v1/comments/:id/like` and the backend toggles the state.

**8. `subscriptions`**

| Column      | Type        | Constraints                 |
|-------------|-------------|-----------------------------|
| id          | bigint      | PK                          |
| subscriber_id | uuid      | FK → profiles.id            |
| channel_id  | uuid        | FK → profiles.id            |
| created_at  | timestamptz | default now()               |
| **Unique**  | (subscriber_id, channel_id) |           |

**9. `views`** (for analytics / unique view tracking)

| Column     | Type        | Constraints                 |
|------------|-------------|-----------------------------|
| id         | bigint      | PK                          |
| video_id   | uuid        | FK → videos.id              |
| user_id    | uuid        | nullable (guests allowed)   |
| ip_hash    | text        |                             |
| created_at | timestamptz | default now()               |

---

## Backend API Design (Go)

### Tech Stack
- **Language:** Go 1.22+
- **Router:** `github.com/go-chi/chi/v5` (lightweight, fast)
- **Supabase Client:** `github.com/supabase-community/supabase-go` (backend-only, service role key)
- **Auth:** Supabase Auth (JWT validation in middleware — backend validates tokens, frontend only stores them)
- **Storage:** Supabase Storage (signed URLs for upload, public URLs for playback)
- **Video Processing:** FFmpeg CLI invoked from Go worker
- **Config:** `github.com/joho/godotenv` for env vars

### API Endpoints

#### Auth (User System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| POST   | `/api/v1/auth/register`          | No            | Register new user (email + password)|
| POST   | `/api/v1/auth/login`             | No            | Login user, return JWT token        |
| POST   | `/api/v1/auth/logout`            | Yes           | Logout user, invalidate token       |
| POST   | `/api/v1/auth/refresh`           | No            | Refresh Supabase session token      |
| GET    | `/api/v1/auth/me`                | Yes           | Get current authenticated user      |
| POST   | `/api/v1/auth/google`            | No            | Google OAuth login/register         |

#### Upload (Video System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| POST   | `/api/v1/videos/upload-url`      | Yes           | Request signed Storage upload URL   |
| POST   | `/api/v1/videos`                 | Yes           | Create video metadata after upload  |

#### Videos (Video System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| GET    | `/api/v1/videos`                 | No            | List videos (feed, paginated)       |
| GET    | `/api/v1/videos/:id`             | No            | Get single video details + assets + links |
| PUT    | `/api/v1/videos/:id`             | Yes           | Update video metadata (owner only)  |
| DELETE | `/api/v1/videos/:id`             | Yes           | Delete own video                    |
| POST   | `/api/v1/videos/:id/view`        | No            | Register a view (debounced)         |
| POST   | `/api/v1/videos/:id/like`        | Yes           | Toggle like on video                |

#### Video Links (Video System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| POST   | `/api/v1/videos/:id/links`       | Yes           | Add a link under video description  |
| DELETE | `/api/v1/videos/:id/links/:linkId` | Yes         | Remove a link from video            |
| GET    | `/api/v1/videos/:id/links`       | No            | List links under a video            |

#### Comments (Comment System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| GET    | `/api/v1/videos/:id/comments`    | No            | List comments (paginated, with replies) |
| POST   | `/api/v1/videos/:id/comments`    | Yes           | Post comment or reply (set parent_id) |
| DELETE | `/api/v1/comments/:id`           | Yes           | Delete own comment                  |
| POST   | `/api/v1/comments/:id/like`      | Yes           | Toggle like on comment              |

#### Channels & Search
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| GET    | `/api/v1/channels/:id`           | No            | Get channel profile + videos        |
| POST   | `/api/v1/channels/:id/subscribe` | Yes           | Toggle subscription                 |
| GET    | `/api/v1/search`                 | No            | Search videos & channels            |

#### Me (User System)
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| GET    | `/api/v1/me`                     | Yes           | Get current user profile            |
| PUT    | `/api/v1/me`                     | Yes           | Update current user profile         |

#### Admin / Dashboard
| Method | Endpoint                         | Auth Required | Description                         |
|--------|----------------------------------|---------------|-------------------------------------|
| GET    | `/admin`                         | No (or Basic) | Admin dashboard HTML page           |
| GET    | `/api/v1/admin/health`           | No            | System health check                 |
| GET    | `/api/v1/admin/stats`            | Yes (Admin)   | Platform statistics                 |
| GET    | `/api/v1/admin/queue`            | Yes (Admin)   | Video processing queue status       |
| GET    | `/api/v1/admin/worker-status`    | Yes (Admin)   | Worker health & job count           |
| GET    | `/api/v1/admin/storage-status`   | Yes (Admin)   | Storage bucket usage                |

### Middleware
- CORS
- JWT validation (extract from `Authorization: Bearer <token>`, validate against Supabase JWKS)
- Rate limiting (per IP & per user)
- Request logging
- Admin role check (for admin endpoints)

---

## Admin / Dashboard Panel

The backend provides a lightweight admin panel at `/admin` (or `/dashboard`) for monitoring the system.

**Accessible Information:**
- **API Base URL:** Current backend URL configuration
- **Frontend Connection Info:** Allowed CORS origins, frontend URL
- **Upload Status:** Recent uploads, success/failure rates, average upload time
- **Worker Status:** Is the worker running? Current job being processed, jobs in queue, completed jobs, failed jobs
- **Video Processing Queue:** List of videos waiting to be processed, currently processing, completed, failed
- **Storage Bucket Status:**
  - `raw-videos` bucket: file count, total size
  - `hls-videos` bucket: file count, total size
  - `thumbnails` bucket: file count, total size
- **System Health Check:**
  - Database connection status
  - Supabase Auth connection status
  - Supabase Storage connection status
  - Worker process status
  - Disk space / memory usage

**Tech:** The admin panel can be a simple server-rendered HTML page served by the Go backend (using Go templates) or a minimal embedded JS bundle. It does not need to be a full SPA.

---

## Video System (HLS Only — No Single MP4 Playback)

**Rule:** Videos will NEVER be played as a single MP4 file. All video playback is done via HLS (HTTP Live Streaming) with segmented `.ts` files and `.m3u8` playlists. This enables adaptive bitrate, faster seeking, and better caching.

### Upload Flow

1. **User selects video file** on `/upload`
2. **Frontend requests upload permission** from Go backend:
   - `POST /api/v1/videos/upload-url`
   - Backend validates user, checks upload limits, generates a signed Supabase Storage URL for the `raw-videos` bucket
   - Backend returns: `{ upload_url, path, token, expires_at }`
3. **Frontend uploads raw video file DIRECTLY to Supabase Storage** using the signed URL
   - This bypasses the Go backend entirely — saves backend bandwidth and memory
   - Frontend shows upload progress
4. **After upload completes, frontend notifies backend:**
   - `POST /api/v1/videos` with `{ title, description, visibility, storage_path, duration_seconds }`
   - Backend inserts row into `videos` table with `status = 'processing'`
   - Backend enqueues a processing job for the worker
5. **Backend returns video ID** to frontend. Frontend shows "Processing..." status.

### Processing Flow (Backend Worker + FFmpeg)

6. **Worker picks up job** from queue (in-memory channel or Redis list)
7. **Worker downloads raw video** from Supabase Storage (`raw-videos` bucket) to a local temp directory
8. **Worker runs FFmpeg** to compress and convert to HLS:
   ```bash
   ffmpeg -i raw_video.mp4 \
     -vf "scale=-2:720,format=yuv420p" \
     -c:v libx264 -b:v 2000k -maxrate 2500k -bufsize 4000k \
     -c:a aac -b:a 128k -ar 48000 \
     -f hls -hls_time 6 -hls_playlist_type vod \
     -hls_segment_filename "segment_%03d.ts" \
     -master_pl_name master.m3u8 \
     playlist.m3u8
   ```
   - **MVP Quality:** 720p only
   - **Target Bitrate:** ~1700k-2500k (variable, capped at 2500k)
   - **Audio:** AAC 128k
   - **Segment Duration:** 6 seconds
   - **Output Files:**
     - `master.m3u8` — master playlist (for future multi-quality)
     - `playlist.m3u8` — 720p variant playlist
     - `segment_000.ts`, `segment_001.ts`, ... — actual video segments
9. **Worker uploads all HLS output files** to Supabase Storage (`hls-videos` bucket):
   - Storage path: `hls/{video_id}/720p/`
   - MIME types: `application/vnd.apple.mpegurl` for `.m3u8`, `video/mp2t` for `.ts`
   - Bucket policy: public read (so CDN can serve segments directly)
10. **Worker inserts row into `video_assets`:**
    - `quality: '720p'`
    - `playlist_url: 'https://oroya.xyz/storage/v1/object/public/hls-videos/{video_id}/720p/playlist.m3u8'`
    - `master_url: 'https://oroya.xyz/storage/v1/object/public/hls-videos/{video_id}/720p/master.m3u8'`
    - `width: 1280`, `height: 720`, `bitrate: 2000000`, `size_bytes: ...`
11. **Worker updates `videos` table:** `status = 'ready'`
12. **Worker cleans up** temp files and optionally deletes raw video from `raw-videos` bucket
13. **Frontend polls `GET /api/v1/videos/:id`** and sees `status = 'ready'`

### Playback Flow

14. **User opens `/watch/{id}`**
15. **Frontend fetches video details** from `GET /api/v1/videos/:id`
16. **Response includes:** video metadata + `video_assets` array with HLS URLs + `video_links` array
17. **Frontend uses HLS.js** (or native Safari HLS) to load the `.m3u8` playlist:
    - For MVP: only one `720p` asset
    - HLS.js attaches to `<video>` element and handles segment loading automatically
18. **Video segments stream directly from Supabase Storage / CDN**
    - The Go backend is NOT involved in serving video bytes
    - Backend ONLY handles: auth validation, metadata, permissions, view counting, comments, likes
    - Large video files NEVER pass through the Go backend

---

## Frontend Design (Svelte)

### Tech Stack
- **Framework:** Svelte 5 (runes) + SvelteKit
- **Build Tool:** Vite
- **Styling:** Plain CSS with CSS custom properties (no heavy UI framework)
- **Icons:** Lucide icons (SVG, lightweight)
- **Video Player:** Native HTML5 `<video>` + **HLS.js** for HLS playback
- **State:** Svelte stores (writable/derived) — no external state library needed
- **API:** Native `fetch` with a thin wrapper → calls Go backend at `PUBLIC_API_BASE_URL` ONLY

**What the frontend knows:**
- `PUBLIC_API_BASE_URL=https://api.oroya.xyz`
- API endpoint paths (e.g., `/api/v1/videos`, `/api/v1/auth/login`)
- How to render UI components
- How to store user JWT token in memory / localStorage

**What the frontend does NOT know:**
- Database table names or columns
- Supabase project URL or anon key (it doesn't use Supabase client libraries at all)
- How videos are processed
- How signed URLs are generated
- How likes, comments, or views are counted

### Key Pages & Components

**1. Layout (`+layout.svelte`)**
- Warm off-white background (`#F5F0EB`)
- Orange accent top bar with search, logo, user avatar
- Bottom nav bar (mobile) / Left sidebar (desktop)

**2. Home Feed (`/`)**
- Vertical scrollable video grid
- Video cards: thumbnail, duration badge, title, channel name, views + time ago
- Pull-to-refresh gesture
- Infinite scroll pagination

**3. Video Player Page (`/watch/[id]`)**
- Full-width video player (16:9)
- **HLS Playback:** Uses HLS.js to load `.m3u8` playlist URL received from backend
- Custom controls: play/pause, seek, volume, fullscreen, quality selector (future)
- Video info: title, views, likes, upload date, channel info, subscribe button
- Like button (orange when active)
- **Video links:** Links displayed under the description (fetched from backend)
- Comments section below video (with reply, like, delete)
- Related videos list

**4. Upload Page (`/upload`)**
- Drag & drop video file area
- Thumbnail selector (auto-generate or upload)
- Title, description inputs
- Visibility toggle (Public / Unlisted / Private)
- Upload progress bar (orange)
- Processing status indicator after upload completes
- Direct upload to Supabase Storage using signed URL received from Go backend

**5. Channel Page (`/channel/[id]`)**
- Banner + avatar + channel info
- Tabs: Videos, Shorts, About
- Subscribe button (orange)
- Video grid for this channel

**6. Search Page (`/search?q=...`)**
- Search bar with orange focus ring
- Filter tabs: Videos, Channels
- Results list with thumbnails

**7. Auth Page (`/login`)**
- Sign in / Sign up forms
- Email/password sent to backend (`POST /api/v1/auth/login`, `POST /api/v1/auth/register`)
- Google OAuth button (redirects to backend `/api/v1/auth/google`)
- Clean, minimal design with orange CTA buttons

### Global Components
- `VideoCard.svelte` — reusable thumbnail card
- `HlsPlayer.svelte` — HTML5 video wrapper with HLS.js integration
- `BottomNav.svelte` — mobile navigation
- `Sidebar.svelte` — desktop drawer
- `LikeButton.svelte` — animated heart/thumb icon (works for videos and comments)
- `CommentThread.svelte` — recursive comment + replies + like/delete actions
- `VideoLinkList.svelte` — links displayed under video description
- `Spinner.svelte` — orange loading indicator
- `Toast.svelte` — notification system

---

## Backend Responsibilities Summary

The Go backend is the single source of truth for all application logic. It manages four core systems:

### 1. User System
- Registration (email + password, or Google OAuth)
- Login / Logout / Token refresh
- Profile CRUD (real name, username, avatar, bio, etc.)
- JWT token validation on every protected request
- Password hashing is delegated to Supabase Auth; backend never stores raw passwords

### 2. Video System
- Upload permission control (signed URL generation)
- Video metadata CRUD (title, description, thumbnail, visibility)
- Video link management (add/remove links under description)
- Video status tracking (`processing` → `ready` / `failed`)
- HLS transcoding via FFmpeg worker
- View counting

### 3. Comment System
- Create top-level comments
- Reply to comments (nested threads)
- Delete own comments
- List comments with pagination
- Comment like/unlike counting

### 4. Like System
- Like / unlike videos
- Like / unlike comments
- Count aggregation (updates `likes_count` on videos and comments)

### 5. Database Operations
- All SQL queries are executed by the backend repository layer
- Frontend never sees a database connection string, table name, or column name
- If the schema changes, only the backend needs to be updated and redeployed

---

## Performance Targets

- **Lighthouse Score:** 95+ on mobile
- **First Contentful Paint:** < 1.5s
- **Time to Interactive:** < 3.5s on 4G
- **Backend Response Time:** p95 < 100ms for cached endpoints
- **Bundle Size:** < 150 KB gzipped initial JS
- **No React, No Node.js backend** — pure Go + compiled Svelte
- **Video Streaming:** Direct from CDN/Storage — zero backend bandwidth for video files

---

## Security Considerations

- **NO secret keys in frontend.**
  - No Supabase service-role key.
  - No Supabase anon key (frontend doesn't use Supabase client at all).
  - No internal API keys.
  - Frontend only knows `PUBLIC_API_BASE_URL`.
- **Authentication Flow:**
  1. User submits email/password to `POST /api/v1/auth/login`
  2. Backend validates via Supabase Auth and returns a JWT access token
  3. Frontend stores token (memory or localStorage)
  4. Every subsequent request includes `Authorization: Bearer <token>` header
  5. Backend validates the token on every protected endpoint
  6. Backend performs the database operation
- **Row Level Security (RLS)** enabled on all Supabase tables as defense-in-depth.
- **Frontend NEVER queries Postgres or Supabase Storage directly.** All access goes through Go API.
- **Signed URLs:** Supabase Storage upload URLs are time-limited and signed by backend.
- **Video uploads authenticated + authorized** via Go backend.
- **Rate limiting** on uploads, comments, likes, auth endpoints, and admin endpoints.
- **Input sanitization** for comments and user-generated content (XSS prevention).
- **CORS** restricted to `https://oroya.xyz` and local dev.
- **FFmpeg execution:** Worker runs in isolated environment; validate input file types before processing; strict temp file cleanup.

---

## Deployment Plan

### Backend (Go)
- Two binaries:
  - `oroya-server`: HTTP API server (listens on port 8080 or via reverse proxy)
  - `oroya-worker`: Background video processing worker
- Compile: `go build -o oroya-server ./cmd/server` and `go build -o oroya-worker ./cmd/worker`
- Deploy via Docker (Alpine Linux base + FFmpeg installed, ~50MB image)
- Serve behind Nginx or Caddy as reverse proxy at `https://api.oroya.xyz`
- Environment variables for Supabase credentials (never in repo)

### Frontend (Svelte)
- `npm run build` → static files in `build/` or `dist/`
- Deploy to `https://oroya.xyz` (served by Nginx/Caddy) or CDN
- All API calls point to `https://api.oroya.xyz/api/v1/*`
- `PUBLIC_API_BASE_URL=https://api.oroya.xyz` set at build time

### Supabase (Self-Hosted)
- Already running at `https://oroya.xyz`
- Configure Storage buckets:
  - `raw-videos`: for uploaded source files (private, signed access)
  - `hls-videos`: for processed HLS playlists and segments (public read, CDN-friendly)
  - `thumbnails`: for video thumbnails and avatars (public read)
- Set bucket policies and RLS rules

---

## Next Steps (Implementation Order)

1. **Initialize Go backend:** create `go.mod`, router, config, Supabase client setup
2. **Initialize Svelte frontend:** create SvelteKit project with Vite, set up global CSS theme
3. **Set up Supabase schema:** run SQL migrations in backend to create tables + RLS policies + Storage buckets
4. **Build Auth flow:** `register`, `login`, `logout`, `refresh`, `me`, `google` endpoints + frontend forms
5. **Build Upload Infrastructure:** signed URL endpoint + direct Storage upload from frontend
6. **Build FFmpeg Worker:** HLS transcoding pipeline, queue system, asset registration
7. **Build Home Feed:** API endpoint + frontend grid + video cards
8. **Build HLS Video Player Page:** watch page + HLS.js integration + like/comment/link display
9. **Build Comments System:** comment API + nested UI + reply + like + delete
10. **Build Search & Channel pages**
11. **Build Admin Dashboard:** system health, queue status, storage stats
12. **Polish & Optimize:** lazy loading, skeleton screens, caching, PWA manifest, worker monitoring

---

## Notes

- The off-white background `#F5F0EB` should be applied to `body` and all root containers.
- Orange should be used sparingly for primary actions only (CTAs, active states, logo).
- Avoid pure white `#FFFFFF` except for card surfaces to create subtle elevation.
- Use CSS `prefers-reduced-motion` for accessibility.
- Keep Svelte components small and focused; leverage runes for reactive state.
- **HLS.js** is required for non-Safari browsers. Safari supports HLS natively.
- For MVP, `video_assets` will contain only `720p`. The schema and worker are designed to easily add `360p` and `1080p` later.
- **Never stream video bytes through Go.** Backend is for metadata and control only.
- **Frontend has zero knowledge of database schema.** If the schema changes, only the backend needs updating.
