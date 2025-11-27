<script lang="ts">
	import { onMount } from 'svelte';

	let users: any[] = [];
	let loading = true;
	let error: string | null = null;
	let actionInProgress = false;

	async function loadUsers() {
		try {
			loading = true;
			error = null;
			const response = await fetch('/api/admin/users');
			if (!response.ok) {
				throw new Error('Failed to load users');
			}
			const data = await response.json();
			users = data.users || [];
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	async function deactivateUser(userId: string, email: string) {
		if (!confirm(`Are you sure you want to deactivate ${email}?`)) {
			return;
		}

		const reason = prompt('Reason for deactivation:');
		if (!reason) return;

		try {
			actionInProgress = true;
			const response = await fetch(`/api/admin/users/${userId}/deactivate`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ reason })
			});

			if (!response.ok) {
				throw new Error('Failed to deactivate user');
			}

			await loadUsers();
			alert('User deactivated successfully');
		} catch (err: any) {
			alert('Error: ' + err.message);
		} finally {
			actionInProgress = false;
		}
	}

	async function reactivateUser(userId: string, email: string) {
		if (!confirm(`Are you sure you want to reactivate ${email}?`)) {
			return;
		}

		try {
			actionInProgress = true;
			const response = await fetch(`/api/admin/users/${userId}/reactivate`, {
				method: 'POST'
			});

			if (!response.ok) {
				throw new Error('Failed to reactivate user');
			}

			await loadUsers();
			alert('User reactivated successfully');
		} catch (err: any) {
			alert('Error: ' + err.message);
		} finally {
			actionInProgress = false;
		}
	}

	async function deleteUser(userId: string, email: string) {
		if (!confirm(`⚠️ WARNING: Are you sure you want to PERMANENTLY DELETE ${email}?`)) {
			return;
		}

		if (!confirm('This action cannot be undone. Proceed with deletion?')) {
			return;
		}

		const reason = prompt('Reason for deletion:');
		if (!reason) return;

		try {
			actionInProgress = true;
			const response = await fetch(`/api/admin/users/${userId}`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ reason })
			});

			if (!response.ok) {
				throw new Error('Failed to delete user');
			}

			await loadUsers();
			alert('User deleted successfully');
		} catch (err: any) {
			alert('Error: ' + err.message);
		} finally {
			actionInProgress = false;
		}
	}

	function formatDate(dateString: string | null): string {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	onMount(() => {
		loadUsers();
	});
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h2 class="text-3xl font-bold text-gray-900">User Management</h2>
		<button
			on:click={loadUsers}
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
				<p class="mt-4 text-gray-600">Loading users...</p>
			</div>
		{:else if users.length === 0}
			<div class="p-8 text-center text-gray-500">No users found.</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								User
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Roles
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Status
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Last Login
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each users as { user, roles }}
							<tr class={user.IsActive ? '' : 'bg-gray-50'}>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0 h-10 w-10">
											{#if user.AvatarURL}
												<img class="h-10 w-10 rounded-full" src={user.AvatarURL} alt="" />
											{:else}
												<div
													class="h-10 w-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold"
												>
													{user.Name?.[0]?.toUpperCase() || '?'}
												</div>
											{/if}
										</div>
										<div class="ml-4">
											<div class="text-sm font-medium text-gray-900">{user.Name}</div>
											<div class="text-sm text-gray-500">{user.Email}</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex gap-2 flex-wrap">
										{#each roles as role}
											<span
												class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800"
											>
												{role.Name}
											</span>
										{/each}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									{#if user.IsActive}
										<span
											class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
										>
											Active
										</span>
									{:else}
										<span
											class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800"
										>
											Inactive
										</span>
									{/if}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
									{formatDate(user.LastLoginAt)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
									<div class="flex gap-2">
										{#if user.IsActive}
											<button
												on:click={() => deactivateUser(user.ID, user.Email)}
												disabled={actionInProgress}
												class="text-yellow-600 hover:text-yellow-900 disabled:opacity-50"
											>
												Deactivate
											</button>
										{:else}
											<button
												on:click={() => reactivateUser(user.ID, user.Email)}
												disabled={actionInProgress}
												class="text-green-600 hover:text-green-900 disabled:opacity-50"
											>
												Reactivate
											</button>
										{/if}
										<button
											on:click={() => deleteUser(user.ID, user.Email)}
											disabled={actionInProgress}
											class="text-red-600 hover:text-red-900 disabled:opacity-50"
										>
											Delete
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
