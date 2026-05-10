<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { getUploadUrl, createVideo } from '$lib/api/videos';
	import { ApiError, friendlyApiMessage } from '$lib/api/client';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	type UploadStep = 'form' | 'uploading' | 'creating' | 'done' | 'error';

	let step = $state<UploadStep>('form');
	let title = $state('');
	let description = $state('');
	let file = $state<File | null>(null);
	let error = $state('');
	let uploadedVideoId = $state('');

	onMount(() => {
		if (!auth.isLoggedIn) {
			goto('/login');
		}
	});

	async function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			file = input.files[0];
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		if (!title.trim()) {
			error = 'Title is required';
			return;
		}
		if (!file) {
			error = 'Please select a video file';
			return;
		}

		try {
			step = 'uploading';

			const uploadData = await getUploadUrl(file.name);

			const uploadResponse = await fetch(uploadData.upload_url, {
				method: 'PUT',
				body: file,
				headers: {
					'Content-Type': file.type || 'application/octet-stream'
				}
			});
			if (!uploadResponse.ok) {
				throw new Error('Video file upload failed. Please try again.');
			}

			step = 'creating';

			const video = await createVideo({
				title: title.trim(),
				description: description.trim() || undefined,
				storage_path: uploadData.storage_path,
				duration_seconds: 0,
				visibility: 'public'
			});

			uploadedVideoId = video.id;
			step = 'done';
		} catch (e) {
			step = 'error';
			if (e instanceof ApiError) {
				error = friendlyApiMessage(e, 'Upload failed. Please try again.');
			} else if (e instanceof Error) {
				error = e.message;
			} else {
				error = 'Upload failed. Please try again.';
			}
		}
	}

	function reset() {
		step = 'form';
		title = '';
		description = '';
		file = null;
		error = '';
		uploadedVideoId = '';
	}
</script>

<div class="page-container">
	<div class="upload-card">
		<h1 class="page-title">Upload Video</h1>

		{#if !auth.isLoggedIn}
			<p>Please <a href="/login">login</a> to upload videos.</p>

		{:else if step === 'form'}
			<form onsubmit={handleSubmit}>
				<div class="form-group">
					<label for="title">Title</label>
					<input id="title" type="text" bind:value={title} placeholder="Video title" required />
				</div>
				<div class="form-group">
					<label for="description">Description</label>
					<textarea id="description" bind:value={description} placeholder="Video description (optional)" rows="3"></textarea>
				</div>
				<div class="form-group">
					<label for="file">Video File</label>
					<div class="file-input-wrapper">
						<input
							id="file"
							type="file"
							accept="video/*"
							onchange={handleFileSelect}
							required
						/>
						{#if file}
							<span class="file-name">{file.name} ({(file.size / (1024 * 1024)).toFixed(1)} MB)</span>
						{/if}
					</div>
				</div>

				{#if error}
					<p class="form-error">{error}</p>
				{/if}

				<button type="submit" class="btn btn-primary">
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
						<polyline points="17 8 12 3 7 8" />
						<line x1="12" y1="3" x2="12" y2="15" />
					</svg>
					Upload
				</button>
			</form>

		{:else if step === 'uploading'}
			<div class="step-feedback">
				<LoadingSpinner message="Uploading video file..." />
			</div>

		{:else if step === 'creating'}
			<div class="step-feedback">
				<LoadingSpinner message="Creating video record..." />
			</div>

		{:else if step === 'done'}
			<div class="step-feedback">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--success)" stroke-width="2">
					<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
					<polyline points="22 4 12 14.01 9 11.01" />
				</svg>
				<h2>Upload Complete!</h2>
				<p>Your video is being processed. It will be available shortly.</p>
				<div class="done-actions">
					<a href={`/video/${uploadedVideoId}`} class="btn btn-primary">View Video</a>
					<button class="btn btn-secondary" onclick={reset}>Upload Another</button>
				</div>
			</div>

		{:else if step === 'error'}
			<div class="step-feedback">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--danger)" stroke-width="2">
					<circle cx="12" cy="12" r="10" />
					<line x1="15" y1="9" x2="9" y2="15" />
					<line x1="9" y1="9" x2="15" y2="15" />
				</svg>
				<h2>Upload Failed</h2>
				<p>{error}</p>
				<div class="done-actions">
					<button class="btn btn-primary" onclick={reset}>Try Again</button>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.upload-card {
		max-width: 600px;
		margin: 0 auto;
		background: var(--bg-card);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: 32px;
	}

	.upload-card form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.file-input-wrapper {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.file-input-wrapper input[type='file'] {
		padding: 8px;
		cursor: pointer;
	}

	.file-name {
		font-size: 13px;
		color: var(--text-secondary);
	}

	.step-feedback {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		padding: 40px 16px;
		text-align: center;
	}

	.step-feedback h2 {
		font-size: 20px;
		font-weight: 600;
	}

	.step-feedback p {
		font-size: 14px;
		color: var(--text-secondary);
		max-width: 400px;
	}

	.done-actions {
		display: flex;
		gap: 12px;
		margin-top: 16px;
	}
</style>
