<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	interface User {
		email: string;
		name: string;
		provider: string;
		user_id: string;
		avatar_url: string;
	}

	let user: User | null = null;
	let loading = true;
	let showDropdown = false;
	let isAdmin = false;

	async function checkAuth() {
		try {
			const response = await fetch('/api/auth/me');
			if (response.ok) {
				const data = await response.json();
				if (!data.user) {
					// Session was killed - redirect to login
					goto('/');
					return false;
				}
				user = data.user;
				return true;
			} else {
				// Not authenticated, redirect to home
				goto('/');
				return false;
			}
		} catch (err) {
			console.error('Error checking auth status:', err);
			goto('/');
			return false;
		}
	}

	onMount(async () => {
		// Check if user is authenticated
		const isAuthenticated = await checkAuth();
		if (isAuthenticated) {
			// Check if user has admin role
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
	});

	async function handleLogout() {
		try {
			await fetch('/api/auth/logout', { method: 'POST' });
			await goto('/');
		} catch (err) {
			console.error('Error logging out:', err);
		}
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase();
	}
</script>

{#if loading}
	<div class="min-h-screen bg-gray-50 flex items-center justify-center">
		<div class="text-gray-600">Loading...</div>
	</div>
{:else if user}
	<div class="min-h-screen bg-gray-50 flex">
		<!-- Sidebar -->
		<div class="w-64 bg-white border-r border-gray-200">
			<div class="p-6">
				<div class="flex items-center gap-2 mb-8">
					<div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
						<span class="text-white font-bold">💰</span>
					</div>
					<span class="text-xl font-bold text-gray-900">Help Me Budget</span>
				</div>

				<nav class="space-y-2">
					<a href="/dashboard" class="flex items-center gap-3 px-4 py-3 rounded-lg bg-blue-50 text-blue-600 font-medium">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-3m0 0l7-4 7 4M5 9v10a1 1 0 001 1h12a1 1 0 001-1V9m-9 11l4-4m0 0l4 4m-4-4V3" />
						</svg>
						<span>Dashboard</span>
					</a>

					<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
						</svg>
						<span>Invoices</span>
					</button>

					<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
						</svg>
						<span>Messages</span>
						<span class="ml-auto bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">1</span>
					</button>

					<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-blue-600 font-medium">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m2.25-12C11.7 3 15 6.3 15 10.5" />
						</svg>
						<span>My Wallets</span>
					</button>

					<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<span>Activity</span>
					</button>

					<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
						</svg>
						<span>Analytics</span>
					</button>

					<div class="pt-4 mt-4 border-t border-gray-200 space-y-2">
						<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
							<span>Get Help</span>
						</button>

						<button type="button" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-50">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
							<span>Settings</span>
						</button>
					</div>
				</nav>
			</div>
		</div>

		<!-- Main Content -->
		<div class="flex-1 flex flex-col">
			<!-- Top Header -->
			<div class="bg-white border-b border-gray-200 px-8 py-4 flex items-center justify-between">
				<h1 class="text-2xl font-bold text-gray-900">My Wallet</h1>

				<div class="flex items-center gap-6">
					<button type="button" aria-label="Search" class="text-gray-600 hover:text-gray-900">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
						</svg>
					</button>

					<button type="button" aria-label="Notifications" class="text-gray-600 hover:text-gray-900 relative">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
						</svg>
						<span class="absolute top-0 right-0 w-3 h-3 bg-red-500 rounded-full"></span>
					</button>

					<div class="flex items-center gap-3 pl-6 border-l border-gray-200 relative">
						<div class="text-right">
							<p class="text-sm font-medium text-gray-900">{user.name}</p>
							<p class="text-xs text-gray-600">{user.provider}</p>
						</div>
						<div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-bold">
							{getInitials(user.name)}
						</div>
						<button type="button" on:click={() => { showDropdown = !showDropdown }} aria-label="User menu" class="text-gray-600 hover:text-gray-900">
							<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
							</svg>
						</button>

						{#if showDropdown}
							<div class="absolute top-full right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-50">
								{#if isAdmin}
									<a
										href="/admin"
										class="w-full text-left px-4 py-3 text-gray-700 hover:bg-gray-50 flex items-center gap-2 border-b border-gray-100"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
										</svg>
										<span>Admin Panel</span>
									</a>
								{/if}
								<button
									on:click={handleLogout}
									class="w-full text-left px-4 py-3 text-gray-700 hover:bg-gray-50 flex items-center gap-2 rounded-lg"
								>
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
									</svg>
									<span>Logout</span>
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Content Grid -->
			<div class="p-8 flex-1 overflow-y-auto">
				<div class="grid grid-cols-3 gap-8">
					<!-- Left Column -->
					<div class="col-span-2 space-y-8">
						<!-- Total Balance Card -->
						<div class="bg-white rounded-2xl p-8">
							<div class="flex justify-between items-start mb-8">
								<div>
									<p class="text-gray-600 mb-2">Total Balance</p>
									<h2 class="text-5xl font-bold text-gray-900">$56,476.00</h2>
									<p class="text-green-600 text-sm mt-2">📈 2.05% Feb 05, 2025</p>
								</div>
								<button type="button" aria-label="More options" class="text-gray-400 hover:text-gray-600">
									<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
										<path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
									</svg>
								</button>
							</div>
						</div>

						<!-- Card Lists -->
						<div>
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Card Lists <span class="text-gray-500 text-sm font-normal">2</span></h3>
							<div class="space-y-4">
								<!-- Card 1 -->
								<div class="bg-gradient-to-br from-blue-600 to-blue-700 rounded-2xl p-8 text-white">
									<div class="flex justify-between items-start mb-12">
										<div class="text-yellow-300 text-2xl">💳</div>
										<div class="text-right">
											<p class="text-xs text-blue-100">Budget</p>
										</div>
									</div>
									<div class="mb-8">
										<p class="text-blue-100 text-sm">Balance</p>
										<p class="text-3xl font-bold">$24,098.00</p>
									</div>
									<div class="flex justify-between items-end">
										<div>
											<p class="text-xs text-blue-100">Card Number</p>
											<p class="text-lg font-semibold">•••• •••• •••• 4242</p>
										</div>
										<p class="text-xl font-semibold">VISA</p>
									</div>
								</div>

								<!-- Card 2 -->
								<div class="bg-white border-2 border-gray-200 rounded-2xl p-8">
									<div class="flex justify-between items-start mb-12">
										<div class="text-gray-400 text-2xl">💳</div>
										<div class="text-right">
											<p class="text-xs text-gray-600">Budget</p>
										</div>
									</div>
									<div class="mb-8">
										<p class="text-gray-600 text-sm">Balance</p>
										<p class="text-3xl font-bold text-gray-900">$33,000.00</p>
									</div>
									<div class="flex justify-between items-end">
										<div>
											<p class="text-xs text-gray-600">Card Number</p>
											<p class="text-lg font-semibold text-gray-900">•••• •••• •••• 5555</p>
										</div>
										<p class="text-xl font-semibold text-gray-900">VISA</p>
									</div>
								</div>
							</div>

							<button class="w-full mt-6 border-2 border-blue-600 text-blue-600 font-semibold py-3 rounded-xl hover:bg-blue-50 transition">
								Manage Card
							</button>
						</div>

						<!-- Quick Links & Stats -->
						<div class="grid grid-cols-2 gap-8">
							<!-- Quick Links -->
							<div class="bg-white rounded-2xl p-6">
								<h3 class="text-lg font-semibold text-gray-900 mb-6">Quick Links</h3>
								<div class="space-y-4">
									<button class="w-full flex flex-col items-center gap-3 p-6 rounded-xl hover:bg-gray-50 transition">
										<div class="w-12 h-12 bg-blue-50 rounded-full flex items-center justify-center text-blue-600 text-xl">
											💰
										</div>
										<span class="text-sm font-medium text-gray-900">Deposit</span>
									</button>
									<button class="w-full flex flex-col items-center gap-3 p-6 rounded-xl hover:bg-gray-50 transition">
										<div class="w-12 h-12 bg-blue-50 rounded-full flex items-center justify-center text-blue-600 text-xl">
											📤
										</div>
										<span class="text-sm font-medium text-gray-900">Send</span>
									</button>
								</div>
							</div>

							<!-- Statistics -->
							<div class="bg-white rounded-2xl p-6">
								<h3 class="text-lg font-semibold text-gray-900 mb-6">Statistics</h3>
								<div class="space-y-4">
									<div class="flex items-end gap-2 h-32">
										<div class="flex-1 bg-gray-200 rounded h-16"></div>
										<div class="flex-1 bg-gray-200 rounded h-20"></div>
										<div class="flex-1 bg-blue-600 rounded h-32 relative">
											<span class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-900 text-white text-xs px-2 py-1 rounded">$5,100</span>
										</div>
										<div class="flex-1 bg-gray-200 rounded h-24"></div>
									</div>
									<div class="flex justify-between text-xs text-gray-600 pt-2">
										<span>Jan 10</span>
										<span>Jan 11</span>
										<span class="text-blue-600 font-semibold">Jan 12</span>
										<span>Jan 13</span>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Right Column - Notifications -->
					<div class="bg-white rounded-2xl p-6 h-fit">
						<div class="flex justify-between items-center mb-6">
							<h3 class="text-lg font-semibold text-gray-900">Notifications</h3>
							<button class="text-blue-600 text-sm font-semibold hover:underline">Mark all as read</button>
						</div>

						<div class="space-y-4">
							<div class="p-4 bg-gray-50 rounded-xl">
								<div class="flex gap-3">
									<div class="w-10 h-10 rounded-full bg-gradient-to-br from-gray-400 to-gray-600 flex-shrink-0"></div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-semibold text-gray-900">Tenner Stafford</p>
										<p class="text-xs text-gray-600 mt-1">You have sent <span class="font-semibold text-green-600">$200.00</span> to Tenner Stafford</p>
										<p class="text-xs text-gray-500 mt-1">2 mins ago</p>
									</div>
								</div>
							</div>

							<div class="p-4 bg-gray-50 rounded-xl">
								<div class="flex gap-3">
									<div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
										<svg class="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
											<path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
											<path fill-rule="evenodd" d="M4 5a2 2 0 012-2 1 1 0 000-2H3a1 1 0 00-1 1v12a1 1 0 001 1h10a1 1 0 001-1V4a1 1 0 00-1-1 1 1 0 000 2h1a2 2 0 012 2v10a2 2 0 01-2 2H6a2 2 0 01-2-2V5z" clip-rule="evenodd" />
										</svg>
									</div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-semibold text-gray-900">New Invoice Sent</p>
										<p class="text-xs text-gray-600 mt-1">You have sent a new invoice of <span class="font-semibold text-green-600">$4,567.00</span> to Birce Enterprises</p>
										<p class="text-xs text-gray-500 mt-1">5 mins ago</p>
									</div>
								</div>
							</div>

							<div class="p-4 bg-gray-50 rounded-xl">
								<div class="flex gap-3">
									<div class="w-10 h-10 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex-shrink-0"></div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-semibold text-gray-900">Cindy Lillibridge</p>
										<p class="text-xs text-gray-600 mt-1">You have received a new payment request from Cindy Lillibridge for <span class="font-semibold text-red-600">$800.00</span></p>
										<p class="text-xs text-gray-500 mt-1">1 hour ago</p>
										<div class="flex gap-2 mt-3">
											<button class="flex-1 text-xs font-semibold text-gray-600 py-2 border border-gray-300 rounded hover:bg-gray-100 transition">
												Decline
											</button>
											<button class="flex-1 text-xs font-semibold text-white bg-blue-600 py-2 rounded hover:bg-blue-700 transition">
												Pay Now
											</button>
										</div>
									</div>
								</div>
							</div>

							<div class="p-4 bg-green-50 rounded-xl border border-green-200">
								<div class="flex gap-3">
									<div class="text-2xl flex-shrink-0">✓</div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-semibold text-gray-900">Payment Received</p>
										<p class="text-xs text-gray-600 mt-1">Received a new payment <span class="font-semibold">$100</span> from Alesia Alexandra</p>
										<p class="text-xs text-gray-500 mt-1">18 hour ago</p>
									</div>
								</div>
							</div>
						</div>

						<button class="w-full mt-6 text-blue-600 text-sm font-semibold hover:underline">
							See all notifications
						</button>
					</div>
				</div>

				<!-- Currency Section -->
				<div class="mt-8 bg-white rounded-2xl p-6">
					<h3 class="text-lg font-semibold text-gray-900 mb-6">Currency</h3>
					<div class="space-y-4">
						<div class="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg">
							<div class="flex items-center gap-3">
								<span class="text-2xl">🇺🇸</span>
								<span class="font-medium text-gray-900">USD</span>
							</div>
							<span class="font-semibold text-gray-900">56,476.00 USD</span>
						</div>
						<div class="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg">
							<div class="flex items-center gap-3">
								<span class="text-2xl">🇪🇺</span>
								<span class="font-medium text-gray-900">EUR</span>
							</div>
							<span class="font-semibold text-gray-900">49,973.67 EUR</span>
						</div>
						<div class="flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg">
							<div class="flex items-center gap-3">
								<span class="text-2xl">🇬🇧</span>
								<span class="font-medium text-gray-900">GBP</span>
							</div>
							<span class="font-semibold text-gray-900">45,098.56 GBP</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
