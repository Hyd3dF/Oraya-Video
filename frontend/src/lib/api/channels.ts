import { apiGet, apiPost } from './client';
import type { ChannelView, SearchResults, SubscribeResponse } from '$lib/types';

export async function getChannel(id: string): Promise<ChannelView> {
	return apiGet<ChannelView>(`/api/v1/channels/${id}`);
}

export async function subscribeToChannel(id: string): Promise<SubscribeResponse> {
	return apiPost<SubscribeResponse>(`/api/v1/channels/${id}/subscribe`);
}

export async function search(q: string, limit = 25): Promise<SearchResults> {
	return apiGet<SearchResults>(`/api/v1/search?q=${encodeURIComponent(q)}&limit=${limit}`);
}
