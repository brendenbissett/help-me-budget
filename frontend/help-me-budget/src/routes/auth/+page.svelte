<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { signInWithEmail, signUpWithEmail, signInWithOAuth } from '$lib/supabase.client';

	// Check URL for initial mode and verification status
	let mode = $state<'login' | 'signup'>(($page.url.searchParams.get('mode') || 'login') as 'login' | 'signup');
	let emailVerified = $derived($page.url.searchParams.get('verified') === 'true');

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let fullName = $state('');
	let loading = $state(false);
	let error = $state('');
	let successMessage = $state('');

	async function handleSubmit() {
		loading = true;
		error = '';
		successMessage = '';

		if (mode === 'signup') {
			// Validation for signup
			if (password !== confirmPassword) {
				error = 'Passwords do not match';
				loading = false;
				return;
			}

			if (password.length < 6) {
				error = 'Password must be at least 6 characters';
				loading = false;
				return;
			}

			try {
				await signUpWithEmail(email, password, {
					full_name: fullName
				});
				successMessage = 'Account created! Please check your email to confirm your account.';
				// Clear form
				email = '';
				password = '';
				confirmPassword = '';
				fullName = '';
			} catch (err: any) {
				error = err.message || 'An error occurred during signup';
			} finally {
				loading = false;
			}
		} else {
			// Login
			try {
				await signInWithEmail(email, password);
				goto('/dashboard');
			} catch (err: any) {
				error = err.message || 'An error occurred during login';
				loading = false;
			}
		}
	}

	async function handleOAuth(provider: 'google' | 'facebook') {
		loading = true;
		error = '';

		try {
			await signInWithOAuth(provider);
			// OAuth will redirect automatically
		} catch (err: any) {
			error = err.message || 'An error occurred during OAuth';
			loading = false;
		}
	}

	function toggleMode() {
		mode = mode === 'login' ? 'signup' : 'login';
		error = '';
		successMessage = '';
		// Update URL without reload
		const url = new URL(window.location.href);
		url.searchParams.set('mode', mode);
		window.history.replaceState({}, '', url);
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<!-- Logo/Header -->
		<div>
			<div class="flex justify-center mb-6">
				<div class="flex items-center gap-2">
					<div class="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
						<span class="text-white font-bold text-xl">💰</span>
					</div>
					<span class="text-2xl font-bold text-gray-900">Help Me Budget</span>
				</div>
			</div>

			<!-- Mode Toggle Tabs -->
			<div class="flex border-b border-gray-200">
				<button
					type="button"
					onclick={() => { mode = 'login'; error = ''; successMessage = ''; }}
					class="flex-1 py-4 px-1 text-center border-b-2 font-medium text-sm transition-colors {mode === 'login' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
				>
					Login
				</button>
				<button
					type="button"
					onclick={() => { mode = 'signup'; error = ''; successMessage = ''; }}
					class="flex-1 py-4 px-1 text-center border-b-2 font-medium text-sm transition-colors {mode === 'signup' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
				>
					Sign Up
				</button>
			</div>
		</div>

		<!-- Email Verified Success Message -->
		{#if emailVerified}
			<div class="rounded-md bg-green-50 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-green-800">Email verified successfully!</h3>
						<div class="mt-2 text-sm text-green-700">
							<p>You can now sign in with your credentials below.</p>
						</div>
					</div>
				</div>
			</div>
		{/if}

		<!-- Success Message -->
		{#if successMessage}
			<div class="rounded-md bg-green-50 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<p class="text-sm text-green-700">{successMessage}</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Error Message -->
		{#if error}
			<div class="rounded-md bg-red-50 p-4">
				<div class="text-sm text-red-700">{error}</div>
			</div>
		{/if}

		<!-- OAuth Buttons -->
		<div class="space-y-3">
			<button
				type="button"
				onclick={() => handleOAuth('google')}
				disabled={loading}
				class="w-full flex items-center justify-center gap-3 py-3 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
				<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
					<path d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z" />
				</svg>
				Continue with Google
			</button>

			<button
				type="button"
				onclick={() => handleOAuth('facebook')}
				disabled={loading}
				class="w-full flex items-center justify-center gap-3 py-3 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
				<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
					<path fill-rule="evenodd" d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z" clip-rule="evenodd" />
				</svg>
				Continue with Facebook
			</button>
		</div>

		<!-- Divider -->
		<div class="relative">
			<div class="absolute inset-0 flex items-center">
				<div class="w-full border-t border-gray-300"></div>
			</div>
			<div class="relative flex justify-center text-sm">
				<span class="px-2 bg-gray-50 text-gray-500">Or continue with email</span>
			</div>
		</div>

		<!-- Email/Password Form -->
		<form class="space-y-4" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
			{#if mode === 'signup'}
				<div>
					<label for="full-name" class="block text-sm font-medium text-gray-700 mb-1">
						Full Name <span class="text-gray-400">(optional)</span>
					</label>
					<input
						id="full-name"
						name="full-name"
						type="text"
						autocomplete="name"
						bind:value={fullName}
						class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
						placeholder="John Doe"
					/>
				</div>
			{/if}

			<div>
				<label for="email-address" class="block text-sm font-medium text-gray-700 mb-1">
					Email address
				</label>
				<input
					id="email-address"
					name="email"
					type="email"
					autocomplete="email"
					required
					bind:value={email}
					class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
					placeholder="you@example.com"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-gray-700 mb-1">
					Password
				</label>
				<input
					id="password"
					name="password"
					type="password"
					autocomplete={mode === 'signup' ? 'new-password' : 'current-password'}
					required
					bind:value={password}
					class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
					placeholder={mode === 'signup' ? 'Min 6 characters' : '••••••••'}
				/>
			</div>

			{#if mode === 'signup'}
				<div>
					<label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-1">
						Confirm Password
					</label>
					<input
						id="confirm-password"
						name="confirm-password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={confirmPassword}
						class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
						placeholder="••••••••"
					/>
				</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
			>
				{loading ? (mode === 'signup' ? 'Creating account...' : 'Signing in...') : (mode === 'signup' ? 'Create Account' : 'Sign In')}
			</button>
		</form>

		<!-- Footer Links -->
		<div class="text-center text-sm text-gray-600">
			<a href="/" class="font-medium text-blue-600 hover:text-blue-500">
				← Back to Home
			</a>
		</div>
	</div>
</div>
