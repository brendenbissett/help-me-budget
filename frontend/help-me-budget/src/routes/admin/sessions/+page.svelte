<script lang="ts">
	import { onMount } from 'svelte';

	let sessions: any[] = [];
	let loading = true;
	let error: string | null = null;
	let actionInProgress = false;

	async function loadSessions() {
		try {
			loading = true;
			error = null;
			const response = await fetch('/api/admin/sessions');
			if (!response.ok) {
				throw new Error('Failed to load sessions');
			}
			const data = await response.json();
			sessions = data.sessions || [];
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	async function killSession(sessionKey: string) {
		if (!confirm(`Are you sure you want to kill session ${sessionKey}?`)) {
			return;
		}

		try {
			actionInProgress = true;
			const response = await fetch(`/api/admin/sessions/${encodeURIComponent(sessionKey)}`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error('Failed to kill session');
			}

			await loadSessions();
			alert('Session killed successfully');
		} catch (err: any) {
			alert('Error: ' + err.message);
		} finally {
			actionInProgress = false;
		}
	}

	onMount(() => {
		loadSessions();
	});
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h2 class="text-3xl font-bold text-gray-900">Session Management</h2>
		<button
			on:click={loadSessions}
			disabled={loading}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
		>
			{loading ? 'Loading...' : '🔄 Refresh'}
		</button>
	</div>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-800 rounded-lg p-4">
			<strong>Error:</strong>
			{error}
		</div>
	{/if}

	<div class="bg-white rounded-lg shadow p-6">
		<div class="flex items-center justify-between mb-4">
			<h3 class="text-lg font-semibold text-gray-900">Active Sessions</h3>
			<span class="text-sm text-gray-600">{sessions.length} active session(s)</span>
		</div>

		{#if loading}
			<div class="p-8 text-center">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-4 text-gray-600">Loading sessions...</p>
			</div>
		{:else if sessions.length === 0}
			<div class="p-8 text-center text-gray-500">No active sessions found.</div>
		{:else}
			<div class="space-y-4">
				{#each sessions as session}
					<div class="border border-gray-200 rounded-lg p-4 hover:border-blue-500 transition-colors">
						<div class="flex items-start justify-between">
							<div class="flex-1">
								<div class="flex items-center space-x-2 mb-2">
									<span class="text-sm font-mono font-semibold text-gray-900">{session.key}</span>
									<span class="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full">
										Active
									</span>
								</div>
								<div class="text-sm text-gray-600 space-y-1">
									<div class="font-mono bg-gray-50 p-2 rounded text-xs overflow-x-auto">
										{session.value}
									</div>
								</div>
							</div>
							<button
								on:click={() => killSession(session.key)}
								disabled={actionInProgress}
								class="ml-4 px-3 py-1.5 bg-red-600 text-white text-sm rounded hover:bg-red-700 disabled:opacity-50 transition-colors"
							>
								Kill Session
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Info Card -->
	<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
		<div class="flex items-start">
			<div class="text-2xl mr-3">ℹ️</div>
			<div class="flex-1">
				<h4 class="font-semibold text-blue-900 mb-2">About Session Management</h4>
				<ul class="text-sm text-blue-800 space-y-1">
					<li>• Active sessions represent users currently logged in to the application</li>
					<li>• Killing a session will immediately log out the user</li>
					<li>
						• Sessions are stored in Redis and automatically expire after 24 hours of inactivity
					</li>
					<li>• All session kills are logged in the audit logs</li>
				</ul>
			</div>
		</div>
	</div>
</div>
