<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let currentUser: any = null;
	let isAdmin = false;
	let loading = true;

	async function checkAuth() {
		try {
			const response = await fetch('/api/auth/me');
			if (response.ok) {
				const data = await response.json();
				if (!data.user) {
					// Session was killed - redirect to home
					goto('/');
					return false;
				}
				currentUser = data;
				return true;
			}
		} catch (error) {
			console.error('Failed to check admin access:', error);
		}
		return false;
	}

	onMount(async () => {
		const isAuthenticated = await checkAuth();
		if (isAuthenticated) {
			// Check if user has admin role by trying to fetch users
			const adminCheck = await fetch('/api/admin/users?limit=1');
			isAdmin = adminCheck.ok;

			// Periodically check session validity (every 30 seconds)
			const interval = setInterval(async () => {
				const stillAuthenticated = await checkAuth();
				if (!stillAuthenticated) {
					clearInterval(interval);
				}
			}, 30000);
		}
		loading = false;
		if (!isAdmin) {
			goto('/');
		}
	});

	const navItems = [
		{ href: '/admin', label: 'Dashboard', icon: '📊' },
		{ href: '/admin/users', label: 'Users', icon: '👥' },
		{ href: '/admin/sessions', label: 'Sessions', icon: '🔐' },
		{ href: '/admin/audit-logs', label: 'Audit Logs', icon: '📋' }
	];

	function isActive(href: string): boolean {
		if (href === '/admin') {
			return $page.url.pathname === '/admin';
		}
		return $page.url.pathname.startsWith(href);
	}
</script>

{#if loading}
	<div class="flex items-center justify-center min-h-screen">
		<div class="text-center">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-4 text-gray-600">Checking admin access...</p>
		</div>
	</div>
{:else if isAdmin}
	<div class="min-h-screen bg-gray-50">
		<!-- Header -->
		<header class="bg-white shadow">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
				<div class="flex justify-between items-center">
					<div class="flex items-center space-x-4">
						<h1 class="text-2xl font-bold text-gray-900">Admin Panel</h1>
					</div>
					<div class="flex items-center space-x-4">
						<span class="text-sm text-gray-600">
							{currentUser?.email}
						</span>
						<a href="/" class="text-sm text-blue-600 hover:text-blue-800">
							← Back to App
						</a>
					</div>
				</div>
			</div>
		</header>

		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="flex gap-8">
				<!-- Sidebar Navigation -->
				<nav class="w-64 flex-shrink-0">
					<div class="bg-white rounded-lg shadow p-4">
						<ul class="space-y-2">
							{#each navItems as item}
								<li>
									<a
										href={item.href}
										class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-colors {isActive(
											item.href
										)
											? 'bg-blue-600 text-white'
											: 'text-gray-700 hover:bg-gray-100'}"
									>
										<span class="text-xl">{item.icon}</span>
										<span class="font-medium">{item.label}</span>
									</a>
								</li>
							{/each}
						</ul>
					</div>
				</nav>

				<!-- Main Content -->
				<main class="flex-1">
					<slot />
				</main>
			</div>
		</div>
	</div>
{:else}
	<div class="flex items-center justify-center min-h-screen">
		<div class="text-center">
			<h1 class="text-4xl font-bold text-red-600 mb-4">Access Denied</h1>
			<p class="text-gray-600 mb-4">You don't have permission to access the admin panel.</p>
			<a href="/" class="text-blue-600 hover:text-blue-800 underline">Return to Home</a>
		</div>
	</div>
{/if}
