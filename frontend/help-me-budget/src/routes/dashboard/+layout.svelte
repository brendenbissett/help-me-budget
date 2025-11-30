<script lang="ts">
	import { signOut } from '$lib/supabase.client';
	import { onMount } from 'svelte';
	import DashboardSidebar from '$lib/components/DashboardSidebar.svelte';

	let { data, children } = $props();
	let showDropdown = $state(false);
	let isAdmin = $state(false);

	onMount(async () => {
		// Check if user has admin role
		try {
			const response = await fetch('/api/user/roles');
			if (response.ok) {
				const rolesData = await response.json();
				isAdmin = rolesData.is_admin || false;
			}
		} catch (err) {
			console.error('Error checking admin status:', err);
		}
	});

	async function handleLogout() {
		try {
			// signOut() will handle redirect and cache invalidation
			await signOut();
		} catch (err) {
			console.error('Error logging out:', err);
		}
	}

	function getInitials(name: string | null | undefined): string {
		if (!name) return '?';
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase();
	}

	function getUserName(): string {
		return (
			data.user?.user_metadata?.full_name ||
			data.user?.user_metadata?.name ||
			data.user?.email?.split('@')[0] ||
			'User'
		);
	}

	function getProvider(): string {
		return data.user?.app_metadata?.provider || 'email';
	}
</script>

{#if data.user}
	<div class="min-h-screen bg-gray-50 flex">
		<!-- Sidebar -->
		<DashboardSidebar />

		<!-- Main Content -->
		<div class="flex-1 flex flex-col">
			<!-- Top Header -->
			<div class="bg-white border-b border-gray-200 px-8 py-4 flex items-center justify-between">
				<h1 class="text-2xl font-bold text-gray-900">My Wallet</h1>

				<div class="flex items-center gap-6">
					<button type="button" aria-label="Search" class="text-gray-600 hover:text-gray-900">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
							/>
						</svg>
					</button>

					<button
						type="button"
						aria-label="Notifications"
						class="text-gray-600 hover:text-gray-900 relative"
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
							/>
						</svg>
						<span class="absolute top-0 right-0 w-3 h-3 bg-red-500 rounded-full"></span>
					</button>

					<div class="flex items-center gap-3 pl-6 border-l border-gray-200 relative">
						<div class="text-right">
							<p class="text-sm font-medium text-gray-900">{getUserName()}</p>
						</div>
						<div
							class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-bold"
						>
							{getInitials(getUserName())}
						</div>
						<button
							type="button"
							onclick={() => {
								showDropdown = !showDropdown;
							}}
							aria-label="User menu"
							class="text-gray-600 hover:text-gray-900"
						>
							<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
								<path
									fill-rule="evenodd"
									d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>

						{#if showDropdown}
							<div
								class="absolute top-full right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-50"
							>
								{#if isAdmin}
									<a
										href="/admin"
										class="w-full text-left px-4 py-3 text-gray-700 hover:bg-gray-50 flex items-center gap-2 border-b border-gray-100"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
											/>
										</svg>
										<span>Admin Panel</span>
									</a>
								{/if}
								<button
									onclick={handleLogout}
									class="w-full text-left px-4 py-3 text-gray-700 hover:bg-gray-50 flex items-center gap-2 rounded-lg"
								>
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
										/>
									</svg>
									<span>Logout</span>
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Content Area -->
			<div class="p-8 flex-1 overflow-y-auto">
				{@render children()}
			</div>
		</div>
	</div>
{/if}
