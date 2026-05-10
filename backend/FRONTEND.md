# Frontend Integration Guide

This is the **only** document the frontend developer needs to read.
The frontend never sees Supabase keys, the database, or any secret.
It talks to one URL: this backend.

---

## 1. Connect

The frontend needs a single environment variable:

```env
PUBLIC_API_BASE_URL=https://api.oroya.xyz
```

For local development:

```env
PUBLIC_API_BASE_URL=http://localhost:8080
```

Get the live value by hitting `GET /api/v1/info` on a running backend — it returns
the canonical URL plus the full endpoint catalogue. The admin dashboard
(`/admin`) shows this prominently and has a one-click copy button.

---

## 2. The token

When a user logs in or registers, the backend returns:

```json
{
  "access_token":  "eyJhbGciOi...",
  "refresh_token": "v1.MzM...",
  "expires_in":    3600,
  "token_type":    "bearer",
  "user": { "id": "...", "username": "...", "email": "..." }
}
```

Store `access_token` (memory or localStorage) and send it on every protected
request:

```
Authorization: Bearer <access_token>
```

When `access_token` expires (typically after 1 hour), call:

```
POST /api/v1/auth/refresh
{ "refresh_token": "<refresh_token>" }
```

The response has the same shape — replace both tokens.

---

## 3. Endpoints

All endpoints are namespaced under `/api/v1`.
Auth-required endpoints are marked **🔒**.

### Auth

| Method | Path | Body | Description |
|---|---|---|---|
| POST | `/auth/register` | `{email, password, real_name, username}` | Creates auth user + profile |
| POST | `/auth/login` | `{email, password}` | Returns access_token + refresh_token |
| POST | `/auth/refresh` | `{refresh_token}` | Issues new tokens |
| POST | `/auth/logout` 🔒 | — | Invalidates the session |
| GET | `/auth/me` 🔒 | — | Current user profile |

### Profile

| Method | Path | Body | Description |
|---|---|---|---|
| GET | `/me` 🔒 | — | Same as `/auth/me` |
| PUT | `/me` 🔒 | `{real_name?, username?, display_name?, avatar_url?, banner_url?, bio?}` | Partial update |

### Videos

| Method | Path | Body | Description |
|---|---|---|---|
| GET | `/videos?limit=24&offset=0` | — | Paginated public feed (only `status=ready`, `visibility=public`) |
| GET | `/videos/{id}` | — | Returns `{video, assets: [...], links: [...]}` |
| POST | `/videos` 🔒 | `{title, description, visibility, storage_path, duration_seconds}` | Created after upload completes |
| PUT | `/videos/{id}` 🔒 owner | `{title?, description?, visibility?}` | Owner-only update |
| DELETE | `/videos/{id}` 🔒 owner | — | Owner-only delete |
| POST | `/videos/{id}/view` | — | Records a view (debounced server-side) |
| POST | `/videos/{id}/like` 🔒 | — | Returns `{liked: bool}` |
| POST | `/videos/upload-url` 🔒 | `{filename}` | Returns `{upload_url, storage_path, expires_at}` — see Upload Flow |
| GET | `/videos/{id}/links` | — | List video links |
| POST | `/videos/{id}/links` 🔒 owner | `{title, url, sort_order}` | Add a link |
| DELETE | `/videos/{id}/links/{linkId}` 🔒 owner | — | Remove a link |

### Comments

| Method | Path | Body | Description |
|---|---|---|---|
| GET | `/videos/{id}/comments?limit=30&offset=0` | — | Paginated comment list |
| POST | `/videos/{id}/comments` 🔒 | `{content, parent_id?}` | Top-level or reply |
| DELETE | `/comments/{id}` 🔒 owner | — | Delete own comment |
| POST | `/comments/{id}/like` 🔒 | — | Returns `{liked: bool}` |

### Channels

| Method | Path | Body | Description |
|---|---|---|---|
| GET | `/channels/{id}` | — | Returns `{profile, videos, subscriber_count, is_subscribed}` |
| POST | `/channels/{id}/subscribe` 🔒 | — | Returns `{subscribed: bool}` |

### Search

| Method | Path | Body | Description |
|---|---|---|---|
| GET | `/search?q=<query>&limit=25` | — | Returns `{query, videos: [...], channels: [...]}` |

### System

| Method | Path | Description |
|---|---|---|
| GET | `/info` | This document, machine-readable |
| GET | `/admin/health` | Public health check |
| GET | `/healthz` | Trivial liveness probe (text "ok") |

---

## 4. Upload flow (the only multi-step flow)

```
                ┌─────────────┐
   1. POST     │             │
   ─────────▶  │  Backend    │
   /videos/    │             │
   upload-url  └─────────────┘
                      │
                      │ returns { upload_url, storage_path }
                      ▼
                ┌─────────────┐
   2. PUT      │  Supabase   │  ◀── frontend uploads raw bytes
   ─────────▶  │  Storage    │      directly here, bypassing backend
                └─────────────┘
                      │
                      │ (file lives in raw-videos bucket)
                      ▼
                ┌─────────────┐
   3. POST     │             │
   ─────────▶  │  Backend    │  records video row, status="processing"
   /videos     │             │
                └─────────────┘
                      │
                      │ background worker transcodes to HLS,
                      │ uploads segments, sets status="ready"
                      ▼
   4. Poll  GET /videos/{id} until status="ready"
            assets[0].playlist_url is the HLS URL — feed to HLS.js
```

---

## 5. Errors

All error responses share this shape:

```json
{ "error": "<code>", "message": "<human-readable>" }
```

Common codes: `bad_request`, `unauthorized`, `forbidden`, `not_found`,
`invalid_input`, `username_taken`, `rate_limited`, `internal`.

HTTP status codes follow the obvious convention (200/201/204 for success;
400 for client error; 401/403 for auth; 404 not found; 429 rate limit;
5xx for backend failure).

---

## 6. CORS

The backend allows the origins listed in `CORS_ALLOWED_ORIGINS`. Local dev
on `http://localhost:5173` and production `https://oroya.xyz` are allowed
out of the box.

---

## 7. Things the frontend must NOT do

- Talk to Supabase directly (no `@supabase/supabase-js`)
- Store any Supabase key (anon, service-role, or JWT secret)
- Connect to Postgres
- Call Storage APIs other than the signed `upload_url` returned in step 1 of the upload flow
- Implement business logic (likes, view counts, dedup) — the backend does this

If you find yourself wanting to do any of the above, you're in the wrong layer.
