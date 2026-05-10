import { apiGet, apiPost, apiDelete } from './client';
import type { Comment, CommentListResponse, LikeResponse } from '$lib/types';

export async function getComments(
	videoId: string,
	limit = 30,
	offset = 0
): Promise<CommentListResponse> {
	return apiGet<CommentListResponse>(
		`/api/v1/videos/${videoId}/comments?limit=${limit}&offset=${offset}`
	);
}

export async function addComment(
	videoId: string,
	content: string,
	parentId?: string
): Promise<Comment> {
	const body: { content: string; parent_id?: string } = { content };
	if (parentId) body.parent_id = parentId;
	return apiPost<Comment>(`/api/v1/videos/${videoId}/comments`, body);
}

export async function deleteComment(id: string): Promise<void> {
	return apiDelete<void>(`/api/v1/comments/${id}`);
}

export async function likeComment(id: string): Promise<LikeResponse> {
	return apiPost<LikeResponse>(`/api/v1/comments/${id}/like`);
}
