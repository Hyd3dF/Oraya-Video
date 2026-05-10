export interface Profile {
	id: string;
	real_name: string;
	username: string;
	display_name?: string;
	email: string;
	avatar_url?: string;
	banner_url?: string;
	bio?: string;
	login_type?: string;
	created_at?: string;
	updated_at?: string;
}

export interface Video {
	id: string;
	owner_id: string;
	title: string;
	description?: string;
	thumbnail_url?: string;
	duration_seconds: number;
	views_count: number;
	likes_count: number;
	visibility: string;
	status: string;
	created_at: string;
	updated_at: string;
}

export interface VideoAsset {
	id: string;
	video_id: string;
	quality: string;
	playlist_url: string;
	master_url?: string;
	width: number;
	height: number;
	bitrate: number;
	size_bytes: number;
	created_at: string;
}

export interface VideoLink {
	id: string;
	video_id: string;
	title: string;
	url: string;
	sort_order: number;
	created_at: string;
}

export interface Comment {
	id: string;
	video_id: string;
	user_id: string;
	parent_id?: string | null;
	content: string;
	likes_count: number;
	created_at: string;
	user?: Profile;
}

export interface Subscription {
	id: number;
	subscriber_id: string;
	channel_id: string;
	created_at: string;
}

export interface AuthResponse {
	access_token: string;
	refresh_token: string;
	expires_in: number;
	token_type: string;
	user: Profile;
}

export interface VideoListResponse {
	videos: Video[];
}

export interface VideoDetailResponse {
	video: Video;
	assets: VideoAsset[];
	links: VideoLink[];
}

export interface CommentListResponse {
	comments: Comment[];
}

export interface LinksResponse {
	links: VideoLink[];
}

export interface ChannelView {
	profile: Profile;
	videos: Video[];
	subscriber_count: number;
	is_subscribed: boolean;
}

export interface SearchResults {
	query: string;
	videos: Video[];
	channels: Profile[];
}

export interface UploadURLResponse {
	upload_url: string;
	storage_path: string;
	token: string;
	expires_at: number;
}

export interface LikeResponse {
	liked: boolean;
}

export interface SubscribeResponse {
	channel_id: string;
	subscribed: boolean;
}

export interface ViewResponse {
	status: string;
}

export interface ApiError {
	error: string;
	message: string;
}
