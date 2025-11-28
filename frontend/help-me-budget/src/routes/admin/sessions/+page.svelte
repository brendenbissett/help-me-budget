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

	async function killSession(sessionKey: string, userName: string) {
		if (!confirm(`Are you sure you want to terminate ${userName}'s session?`)) {
			return;
		}

		try {
			actionInProgress = true;
			const response = await fetch(`/api/admin/sessions/${encodeURIComponent(sessionKey)}`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to kill session');
			}

			await loadSessions();
			alert(`${userName}'s session terminated successfully`);
		} catch (err: any) {
			alert('Error: ' + err.message);
		} finally {
			actionInProgress = false;
		}
	}

	function formatDate(dateString: string) {
		const date = new Date(dateString);
		return date.toLocaleString();
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
								<div class="flex items-center space-x-3 mb-2">
									<div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-bold">
										{session.name ? session.name.charAt(0).toUpperCase() : '?'}
									</div>
									<div>
										<h4 class="font-semibold text-gray-900">{session.name || 'Unknown User'}</h4>
										<p class="text-sm text-gray-600">{session.email}</p>
									</div>
									<span class="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full">
										Active
									</span>
								</div>
								<div class="text-sm text-gray-600 space-y-1 ml-13">
									<div><span class="font-semibold">Provider:</span> {session.provider}</div>
									<div><span class="font-semibold">Last Sign In:</span> {formatDate(session.login_at)}</div>
									<div class="text-xs text-gray-500 font-mono">{session.user_id}</div>
								</div>
							</div>
							<button
								on:click={() => killSession(session.key, session.name || session.email)}
								disabled={actionInProgress}
								class="ml-4 px-3 py-1.5 bg-red-600 text-white text-sm rounded hover:bg-red-700 disabled:opacity-50 transition-colors flex-shrink-0"
							>
								Terminate Session
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
					<li>• Shows all users who have logged in (managed by Supabase)</li>
					<li>• Terminating a session will immediately sign out the user from all devices</li>
					<li>• Sessions are managed by Supabase authentication system</li>
					<li>• Users will need to re-authenticate after their session is terminated</li>
				</ul>
			</div>
		</div>
	</div>
</div>
