<script lang="ts">
	import { onMount } from 'svelte';

	let stats = {
		totalUsers: 0,
		activeSessions: 0,
		recentLogs: 0
	};
	let loading = true;

	onMount(async () => {
		try {
			// Fetch stats from API
			const [usersRes, sessionsRes, logsRes] = await Promise.all([
				fetch('/api/admin/users?limit=1'),
				fetch('/api/admin/sessions'),
				fetch('/api/admin/audit-logs?limit=10')
			]);

			if (usersRes.ok) {
				const usersData = await usersRes.json();
				stats.totalUsers = usersData.users?.length || 0;
			}

			if (sessionsRes.ok) {
				const sessionsData = await sessionsRes.json();
				stats.activeSessions = sessionsData.count || 0;
			}

			if (logsRes.ok) {
				const logsData = await logsRes.json();
				stats.recentLogs = logsData.logs?.length || 0;
			}
		} catch (error) {
			console.error('Failed to load stats:', error);
		} finally {
			loading = false;
		}
	});
</script>

<div class="space-y-6">
	<h2 class="text-3xl font-bold text-gray-900">Dashboard</h2>

	{#if loading}
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			{#each [1, 2, 3] as _}
				<div class="bg-white rounded-lg shadow p-6 animate-pulse">
					<div class="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
					<div class="h-8 bg-gray-200 rounded w-1/3"></div>
				</div>
			{/each}
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<!-- Total Users Card -->
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-600">Total Users</p>
						<p class="text-3xl font-bold text-gray-900 mt-2">{stats.totalUsers}</p>
					</div>
					<div class="text-4xl">👥</div>
				</div>
				<a href="/admin/users" class="text-sm text-blue-600 hover:text-blue-800 mt-4 inline-block">
					View all users →
				</a>
			</div>

			<!-- Active Sessions Card -->
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-600">Active Sessions</p>
						<p class="text-3xl font-bold text-gray-900 mt-2">{stats.activeSessions}</p>
					</div>
					<div class="text-4xl">🔐</div>
				</div>
				<a href="/admin/sessions" class="text-sm text-blue-600 hover:text-blue-800 mt-4 inline-block">
					Manage sessions →
				</a>
			</div>

			<!-- Recent Logs Card -->
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-600">Recent Audit Logs</p>
						<p class="text-3xl font-bold text-gray-900 mt-2">{stats.recentLogs}</p>
					</div>
					<div class="text-4xl">📋</div>
				</div>
				<a href="/admin/audit-logs" class="text-sm text-blue-600 hover:text-blue-800 mt-4 inline-block">
					View logs →
				</a>
			</div>
		</div>
	{/if}

	<!-- Quick Actions -->
	<div class="bg-white rounded-lg shadow p-6">
		<h3 class="text-xl font-bold text-gray-900 mb-4">Quick Actions</h3>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<a
				href="/admin/users"
				class="flex items-center space-x-3 p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
			>
				<span class="text-2xl">👥</span>
				<span class="font-medium text-gray-700">Manage Users</span>
			</a>
			<a
				href="/admin/sessions"
				class="flex items-center space-x-3 p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
			>
				<span class="text-2xl">🔐</span>
				<span class="font-medium text-gray-700">View Sessions</span>
			</a>
			<a
				href="/admin/audit-logs"
				class="flex items-center space-x-3 p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
			>
				<span class="text-2xl">📋</span>
				<span class="font-medium text-gray-700">Audit Logs</span>
			</a>
			<a
				href="/"
				class="flex items-center space-x-3 p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
			>
				<span class="text-2xl">🏠</span>
				<span class="font-medium text-gray-700">Back to App</span>
			</a>
		</div>
	</div>
</div>
