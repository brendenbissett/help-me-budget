<script lang="ts">
	import { onMount } from 'svelte';

	let logs: any[] = [];
	let loading = true;
	let error: string | null = null;
	let currentPage = 0;
	let limit = 50;

	async function loadLogs() {
		try {
			loading = true;
			error = null;
			const offset = currentPage * limit;
			const response = await fetch(`/api/admin/audit-logs?limit=${limit}&offset=${offset}`);
			if (!response.ok) {
				throw new Error('Failed to load audit logs');
			}
			const data = await response.json();
			logs = data.logs || [];
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	function nextPage() {
		currentPage++;
		loadLogs();
	}

	function prevPage() {
		if (currentPage > 0) {
			currentPage--;
			loadLogs();
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	function getActionIcon(action: string): string {
		if (action.includes('deactivate')) return '🚫';
		if (action.includes('reactivate')) return '✅';
		if (action.includes('delete')) return '🗑️';
		if (action.includes('grant')) return '➕';
		if (action.includes('revoke')) return '➖';
		if (action.includes('kill')) return '🔒';
		return '📝';
	}

	function getActionColor(action: string): string {
		if (action.includes('deactivate')) return 'bg-yellow-100 text-yellow-800';
		if (action.includes('reactivate')) return 'bg-green-100 text-green-800';
		if (action.includes('delete')) return 'bg-red-100 text-red-800';
		if (action.includes('grant')) return 'bg-blue-100 text-blue-800';
		if (action.includes('revoke')) return 'bg-orange-100 text-orange-800';
		if (action.includes('kill')) return 'bg-purple-100 text-purple-800';
		return 'bg-gray-100 text-gray-800';
	}

	onMount(() => {
		loadLogs();
	});
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h2 class="text-3xl font-bold text-gray-900">Audit Logs</h2>
		<button
			on:click={loadLogs}
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

	<div class="bg-white rounded-lg shadow overflow-hidden">
		{#if loading}
			<div class="p-8 text-center">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-4 text-gray-600">Loading audit logs...</p>
			</div>
		{:else if logs.length === 0}
			<div class="p-8 text-center text-gray-500">No audit logs found.</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
							>
								Timestamp
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
							>
								Action
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
							>
								Resource
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
							>
								Actor
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
							>
								Details
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each logs as log}
							<tr class="hover:bg-gray-50">
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
									{formatDate(log.CreatedAt)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span
										class="px-2 inline-flex items-center text-xs leading-5 font-semibold rounded-full {getActionColor(
											log.Action
										)}"
									>
										<span class="mr-1">{getActionIcon(log.Action)}</span>
										{log.Action}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
									<div>
										<div class="font-medium">{log.ResourceType}</div>
										{#if log.ResourceID}
											<div class="text-xs text-gray-500 font-mono">{log.ResourceID}</div>
										{/if}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
									{#if log.ActorID}
										<div class="font-mono text-xs">{log.ActorID}</div>
									{:else}
										<span class="text-gray-400">System</span>
									{/if}
								</td>
								<td class="px-6 py-4 text-sm text-gray-500">
									{#if log.Details && Object.keys(log.Details).length > 0}
										<details class="cursor-pointer">
											<summary class="text-blue-600 hover:text-blue-800">View details</summary>
											<div class="mt-2 p-2 bg-gray-50 rounded text-xs font-mono overflow-x-auto">
												<pre>{JSON.stringify(log.Details, null, 2)}</pre>
											</div>
										</details>
									{:else}
										<span class="text-gray-400">No details</span>
									{/if}
									{#if log.IPAddress}
										<div class="text-xs text-gray-400 mt-1">IP: {log.IPAddress}</div>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination -->
			<div class="bg-gray-50 px-6 py-3 flex items-center justify-between border-t border-gray-200">
				<div class="text-sm text-gray-700">
					Page {currentPage + 1} • Showing {logs.length} log(s)
				</div>
				<div class="flex gap-2">
					<button
						on:click={prevPage}
						disabled={currentPage === 0 || loading}
						class="px-4 py-2 bg-white border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						← Previous
					</button>
					<button
						on:click={nextPage}
						disabled={logs.length < limit || loading}
						class="px-4 py-2 bg-white border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Next →
					</button>
				</div>
			</div>
		{/if}
	</div>

	<!-- Info Card -->
	<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
		<div class="flex items-start">
			<div class="text-2xl mr-3">ℹ️</div>
			<div class="flex-1">
				<h4 class="font-semibold text-blue-900 mb-2">About Audit Logs</h4>
				<ul class="text-sm text-blue-800 space-y-1">
					<li>• All administrative actions are automatically logged and cannot be deleted</li>
					<li>• Logs include the action type, affected resource, actor, IP address, and timestamp</li>
					<li>• Use audit logs to track who did what and when for security and compliance</li>
					<li>• Moderators can view audit logs but cannot modify users or sessions</li>
				</ul>
			</div>
		</div>
	</div>
</div>
