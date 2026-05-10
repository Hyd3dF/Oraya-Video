import { apiGet, apiPost, apiPut, apiDelete } from './client';
import type {
	Video,
	VideoListResponse,
	VideoDetailResponse,
	LinksResponse,
	VideoLink,
	UploadURLResponse,
	LikeResponse,
	ViewResponse
} from '$lib/types';

export async function getVideos(limit = 24, offset = 0): Promise<VideoListResponse> {
	return apiGet<VideoListResponse>(`/api/v1/videos?limit=${limit}&offset=${offset}`);
}

export async function getMyVideos(limit = 100, offset = 0): Promise<VideoListResponse> {
	return apiGet<VideoListResponse>(`/api/v1/me/videos?limit=${limit}&offset=${offset}`);
}

export async function getVideo(id: string): Promise<VideoDetailResponse> {
	return apiGet<VideoDetailResponse>(`/api/v1/videos/${id}`);
}

export async function createVideo(data: {
	title: string;
	description?: string;
	visibility?: string;
	storage_path: string;
	duration_seconds?: number;
	thumbnail_url?: string;
}): Promise<Video> {
	return apiPost<Video>('/api/v1/videos', data);
}

export async function updateVideo(
	id: string,
	data: { title?: string; description?: string; visibility?: string }
): Promise<Video> {
	return apiPut<Video>(`/api/v1/videos/${id}`, data);
}

export async function deleteVideo(id: string): Promise<void> {
	return apiDelete<void>(`/api/v1/videos/${id}`);
}

export async function viewVideo(id: string): Promise<ViewResponse> {
	return apiPost<ViewResponse>(`/api/v1/videos/${id}/view`);
}

export async function likeVideo(id: string): Promise<LikeResponse> {
	return apiPost<LikeResponse>(`/api/v1/videos/${id}/like`);
}

export async function getUploadUrl(filename: string): Promise<UploadURLResponse> {
	return apiPost<UploadURLResponse>('/api/v1/videos/upload-url', { filename });
}

export async function getVideoLinks(id: string): Promise<LinksResponse> {
	return apiGet<LinksResponse>(`/api/v1/videos/${id}/links`);
}

export async function addVideoLink(
	videoId: string,
	data: { title: string; url: string; sort_order: number }
): Promise<VideoLink> {
	return apiPost<VideoLink>(`/api/v1/videos/${videoId}/links`, data);
}

export async function deleteVideoLink(videoId: string, linkId: string): Promise<void> {
	return apiDelete<void>(`/api/v1/videos/${videoId}/links/${linkId}`);
}
