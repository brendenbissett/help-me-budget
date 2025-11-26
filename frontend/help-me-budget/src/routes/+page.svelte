<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let user: { email: string; name: string; provider: string } | null = null;
	let error: string | null = null;
	let isLoading = false;
	let loadingProvider: string | null = null;

	onMount(async () => {
		// Check if user is already logged in
		try {
			const response = await fetch('/api/auth/me');
			if (response.ok) {
				const data = await response.json();
				user = data.user;
				// Redirect to dashboard if authenticated
				if (user) {
					goto('/dashboard');
				}
			}
		} catch (err) {
			console.error('Error checking auth status:', err);
		}
	});

	function handleLogin(provider: string) {
		// Set loading state
		isLoading = true;
		loadingProvider = provider;
		// Navigate to the SvelteKit login endpoint which will redirect to the OAuth provider
		window.location.href = `/api/auth/login/${provider}`;
	}

	async function handleLogout() {
		try {
			await fetch('/api/auth/logout', { method: 'POST' });
			user = null;
			window.location.href = '/';
		} catch (err) {
			console.error('Error logging out:', err);
		}
	}
</script>

<div class="min-h-screen bg-gray-50 flex">
	<!-- Left Side - Blue Section with Demo -->
	<div class="hidden lg:flex w-1/2 bg-blue-600 text-white flex-col justify-between p-12">
		<div>
			<div class="flex items-center gap-2 mb-16">
				<div class="w-8 h-8 bg-white rounded-lg flex items-center justify-center">
					<span class="text-blue-600 font-bold text-lg">💰</span>
				</div>
				<span class="text-2xl font-bold">Help Me Budget</span>
			</div>

			<!-- Demo Dashboard Cards -->
			<div class="space-y-6">
				<div class="bg-white/10 backdrop-blur-md rounded-2xl p-6 space-y-4">
					<div class="grid grid-cols-2 gap-4">
						<div>
							<div class="text-sm text-white/60 mb-1">Income</div>
							<div class="text-2xl font-bold">$24,908.00</div>
						</div>
						<div>
							<div class="text-sm text-white/60 mb-1">Expenses</div>
							<div class="text-2xl font-bold">$1,028.00</div>
						</div>
					</div>
				</div>

				<div class="bg-white rounded-xl p-4 text-gray-800">
					<div class="text-sm text-gray-600 mb-3">Recent Transactions</div>
					<div class="space-y-3">
						<div class="flex justify-between items-center">
							<span class="text-sm">Salary Deposit</span>
							<span class="text-sm font-semibold text-green-600">+$3,500.00</span>
						</div>
						<div class="flex justify-between items-center">
							<span class="text-sm">Rent Payment</span>
							<span class="text-sm font-semibold text-red-600">-$1,200.00</span>
						</div>
						<div class="flex justify-between items-center">
							<span class="text-sm">Groceries</span>
							<span class="text-sm font-semibold text-red-600">-$125.50</span>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div>
			<h2 class="text-4xl font-bold mb-4 leading-tight">Smart budgeting, made simple</h2>
			<p class="text-white/70 text-sm">Track your spending, set goals, and take control of your finances.</p>
		</div>
	</div>

	<!-- Right Side - Login Form -->
	<div class="w-full lg:w-1/2 flex items-center justify-center p-6">
		<div class="w-full max-w-md">
			<!-- Login View -->
				<div>
					<h1 class="text-4xl font-bold text-gray-900 mb-2">Sign up for an account</h1>
					<p class="text-gray-600 mb-8">Send, spend and save smarter</p>

					{#if error}
						<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-4 rounded-xl mb-6 text-sm">
							{error}
						</div>
					{/if}

					<div class="space-y-3 mb-6">
						<button
							on:click={() => handleLogin('google')}
							disabled={isLoading}
							class="w-full bg-white border-2 border-gray-300 hover:border-gray-400 disabled:opacity-70 disabled:cursor-not-allowed text-gray-800 font-semibold py-3 px-4 rounded-xl transition duration-200 flex items-center justify-center gap-3"
						>
							{#if loadingProvider === 'google' && isLoading}
								<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
							{:else}
								<svg class="w-5 h-5" viewBox="0 0 24 24">
									<circle cx="12" cy="12" r="10" fill="#4285F4"/>
									<text x="12" y="15" text-anchor="middle" fill="white" font-size="10" font-weight="bold">G</text>
								</svg>
							{/if}
							<span>{isLoading && loadingProvider === 'google' ? 'Signing up...' : 'Sign Up with Google'}</span>
						</button>

						<button
							on:click={() => handleLogin('facebook')}
							disabled={isLoading}
							class="w-full bg-white border-2 border-gray-300 hover:border-gray-400 disabled:opacity-70 disabled:cursor-not-allowed text-gray-800 font-semibold py-3 px-4 rounded-xl transition duration-200 flex items-center justify-center gap-3"
						>
							{#if loadingProvider === 'facebook' && isLoading}
								<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
							{:else}
								<svg class="w-5 h-5" viewBox="0 0 24 24">
									<rect x="2" y="2" width="20" height="20" fill="#1877F2" rx="3"/>
									<text x="12" y="15" text-anchor="middle" fill="white" font-size="10" font-weight="bold">f</text>
								</svg>
							{/if}
							<span>{isLoading && loadingProvider === 'facebook' ? 'Signing up...' : 'Sign Up with Facebook'}</span>
						</button>
					</div>

					<div class="relative mb-6">
						<div class="absolute inset-0 flex items-center">
							<div class="w-full border-t border-gray-300"></div>
						</div>
						<div class="relative flex justify-center text-sm">
							<span class="px-2 bg-white text-gray-600">Or with email</span>
						</div>
					</div>

					<div class="space-y-3 mb-6">
						<input
							type="text"
							placeholder="First name"
							class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-blue-600 focus:outline-none transition"
						/>
						<input
							type="text"
							placeholder="Last name"
							class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-blue-600 focus:outline-none transition"
						/>
						<input
							type="email"
							placeholder="Email"
							class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-blue-600 focus:outline-none transition"
						/>
						<input
							type="password"
							placeholder="Password"
							class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-blue-600 focus:outline-none transition"
						/>
					</div>

					<button
						class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 rounded-xl transition duration-200 mb-4"
					>
						Sign Up
					</button>

					<p class="text-center text-gray-600 text-sm">
						Already have an account? <span class="text-blue-600 font-semibold hover:cursor-pointer">Sign In</span>
					</p>

					<p class="text-center text-gray-500 text-xs mt-6">
						By creating an account, you agreeing to our <span class="text-gray-700 font-semibold">Privacy Policy</span> and
						<span class="text-gray-700 font-semibold">Terms and Conditions</span>
					</p>
				</div>
		</div>
	</div>
</div>
