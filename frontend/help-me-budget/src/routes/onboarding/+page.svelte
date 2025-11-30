<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let currentStep = $state(1);
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	// Step 1: Welcome
	let userName = $derived(
		data.user?.user_metadata?.full_name ||
			data.user?.user_metadata?.name ||
			data.user?.email?.split('@')[0] ||
			'there'
	);

	// Step 2: First Account
	let accountName = $state('');
	let accountType = $state('checking');
	let accountBalance = $state('0');

	// Step 3: Categories choice
	let useDefaultCategories = $state(true);

	async function nextStep() {
		if (currentStep === 1) {
			currentStep = 2;
		} else if (currentStep === 2) {
			await createAccount();
		} else if (currentStep === 3) {
			await setupCategories();
		}
	}

	function previousStep() {
		if (currentStep > 1) {
			currentStep--;
		}
	}

	async function createAccount() {
		if (!accountName.trim()) {
			errorMessage = 'Please enter an account name';
			return;
		}

		isSubmitting = true;
		errorMessage = '';

		try {
			const response = await fetch('/api/onboarding/account', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: accountName,
					account_type: accountType,
					balance: parseFloat(accountBalance) || 0,
					currency: 'USD'
				})
			});

			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.message || 'Failed to create account');
			}

			currentStep = 3;
		} catch (err: any) {
			errorMessage = err.message || 'Failed to create account';
		} finally {
			isSubmitting = false;
		}
	}

	async function setupCategories() {
		isSubmitting = true;
		errorMessage = '';

		try {
			if (useDefaultCategories) {
				const response = await fetch('/api/onboarding/categories', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' }
				});

				if (!response.ok) {
					const error = await response.json();
					throw new Error(error.message || 'Failed to setup categories');
				}
			}

			// Complete onboarding and redirect to dashboard
			await goto('/dashboard');
		} catch (err: any) {
			errorMessage = err.message || 'Failed to setup categories';
		} finally {
			isSubmitting = false;
		}
	}

	function skipOnboarding() {
		goto('/dashboard');
	}

	const accountTypeConfig = {
		checking: { icon: '🏦', label: 'Checking Account' },
		savings: { icon: '💰', label: 'Savings Account' },
		credit_card: { icon: '💳', label: 'Credit Card' },
		cash: { icon: '💵', label: 'Cash' },
		investment: { icon: '📈', label: 'Investment Account' }
	};
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
	<div class="max-w-2xl w-full">
		<!-- Progress Steps -->
		<div class="mb-8">
			<div class="flex justify-between items-center">
				{#each [1, 2, 3] as step}
					<div class="flex-1 {step < 3 ? 'mr-2' : ''}">
						<div
							class="h-2 rounded-full transition-all {step <= currentStep
								? 'bg-blue-600'
								: 'bg-white'}"
						></div>
					</div>
				{/each}
			</div>
			<div class="flex justify-between mt-2 text-sm text-gray-600">
				<span class={currentStep === 1 ? 'font-semibold text-blue-600' : ''}>Welcome</span>
				<span class={currentStep === 2 ? 'font-semibold text-blue-600' : ''}>First Account</span>
				<span class={currentStep === 3 ? 'font-semibold text-blue-600' : ''}>Categories</span>
			</div>
		</div>

		<!-- Main Card -->
		<div class="bg-white rounded-2xl shadow-xl p-8 md:p-12">
			<!-- Step 1: Welcome -->
			{#if currentStep === 1}
				<div class="text-center">
					<div class="text-6xl mb-6">👋</div>
					<h1 class="text-4xl font-bold text-gray-900 mb-4">Welcome, {userName}!</h1>
					<p class="text-xl text-gray-600 mb-8">
						Let's get you started with Help Me Budget. We'll help you set up your first account
						and organize your finances in just a few steps.
					</p>

					<div class="bg-blue-50 rounded-xl p-6 mb-8">
						<h3 class="font-semibold text-gray-900 mb-3">What you'll do:</h3>
						<ul class="text-left space-y-2 text-gray-700">
							<li class="flex items-start gap-2">
								<svg class="w-6 h-6 text-blue-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
								</svg>
								<span>Add your first bank account or credit card</span>
							</li>
							<li class="flex items-start gap-2">
								<svg class="w-6 h-6 text-blue-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
								</svg>
								<span>Choose how to organize your expenses and income</span>
							</li>
							<li class="flex items-start gap-2">
								<svg class="w-6 h-6 text-blue-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
								</svg>
								<span>Start tracking your budget</span>
							</li>
						</ul>
					</div>

					<div class="flex gap-4 justify-center">
						<button
							onclick={skipOnboarding}
							class="px-6 py-3 text-gray-600 hover:text-gray-800 font-medium transition"
						>
							Skip for now
						</button>
						<button
							onclick={nextStep}
							class="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition"
						>
							Let's Get Started
						</button>
					</div>
				</div>
			{/if}

			<!-- Step 2: First Account -->
			{#if currentStep === 2}
				<div>
					<div class="text-center mb-8">
						<div class="text-5xl mb-4">🏦</div>
						<h2 class="text-3xl font-bold text-gray-900 mb-2">Add Your First Account</h2>
						<p class="text-gray-600">
							This could be your checking account, savings, credit card, or even cash on hand.
						</p>
					</div>

					<div class="space-y-6">
						<div>
							<label for="account_name" class="block text-sm font-medium text-gray-700 mb-2">
								Account Name
							</label>
							<input
								type="text"
								id="account_name"
								bind:value={accountName}
								placeholder="e.g., Chase Checking, Savings Account"
								class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-lg"
							/>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-700 mb-3">Account Type</label>
							<div class="grid grid-cols-2 gap-3">
								{#each Object.entries(accountTypeConfig) as [value, config]}
									<label class="cursor-pointer">
										<input
											type="radio"
											name="account_type"
											{value}
											bind:group={accountType}
											class="sr-only peer"
										/>
										<div
											class="p-4 border-2 border-gray-200 rounded-xl peer-checked:border-blue-600 peer-checked:bg-blue-50 hover:bg-gray-50 transition flex items-center gap-3"
										>
											<span class="text-3xl">{config.icon}</span>
											<span class="font-medium text-gray-900">{config.label}</span>
										</div>
									</label>
								{/each}
							</div>
						</div>

						<div>
							<label for="account_balance" class="block text-sm font-medium text-gray-700 mb-2">
								Current Balance (Optional)
							</label>
							<div class="relative">
								<span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-lg">$</span>
								<input
									type="number"
									id="account_balance"
									bind:value={accountBalance}
									step="0.01"
									placeholder="0.00"
									class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-lg"
								/>
							</div>
							<p class="text-sm text-gray-500 mt-2">
								You can leave this as $0 and update it later
							</p>
						</div>

						{#if errorMessage}
							<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
								{errorMessage}
							</div>
						{/if}
					</div>

					<div class="flex gap-4 mt-8">
						<button
							onclick={previousStep}
							disabled={isSubmitting}
							class="flex-1 px-6 py-3 border-2 border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition disabled:opacity-50"
						>
							Back
						</button>
						<button
							onclick={nextStep}
							disabled={isSubmitting || !accountName.trim()}
							class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isSubmitting ? 'Creating...' : 'Continue'}
						</button>
					</div>
				</div>
			{/if}

			<!-- Step 3: Categories -->
			{#if currentStep === 3}
				<div>
					<div class="text-center mb-8">
						<div class="text-5xl mb-4">📊</div>
						<h2 class="text-3xl font-bold text-gray-900 mb-2">Organize Your Money</h2>
						<p class="text-gray-600">
							Categories help you track where your money comes from and where it goes.
						</p>
					</div>

					<div class="space-y-6">
						<!-- Option 1: Use Default Categories -->
						<label class="cursor-pointer">
							<input
								type="radio"
								name="category_choice"
								checked={useDefaultCategories}
								onchange={() => (useDefaultCategories = true)}
								class="sr-only peer"
							/>
							<div
								class="p-6 border-2 border-gray-200 rounded-xl peer-checked:border-blue-600 peer-checked:bg-blue-50 hover:bg-gray-50 transition"
							>
								<div class="flex items-start gap-4">
									<div class="flex-shrink-0 w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center text-2xl">
										✨
									</div>
									<div class="flex-1">
										<h3 class="font-bold text-gray-900 mb-2 text-lg">
											Use Recommended Categories
										</h3>
										<p class="text-gray-600 mb-3">
											We'll set up common categories like Groceries, Rent, Salary, and more. You can
											customize them later.
										</p>
										<div class="flex flex-wrap gap-2">
											<span
												class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium"
												>Housing</span
											>
											<span
												class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium"
												>Food</span
											>
											<span
												class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium"
												>Transportation</span
											>
											<span
												class="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium"
												>Salary</span
											>
											<span
												class="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium"
												>Freelance</span
											>
											<span class="px-3 py-1 bg-gray-200 text-gray-600 rounded-full text-sm"
												>+13 more</span
											>
										</div>
									</div>
								</div>
							</div>
						</label>

						<!-- Option 2: Start from Scratch -->
						<label class="cursor-pointer">
							<input
								type="radio"
								name="category_choice"
								checked={!useDefaultCategories}
								onchange={() => (useDefaultCategories = false)}
								class="sr-only peer"
							/>
							<div
								class="p-6 border-2 border-gray-200 rounded-xl peer-checked:border-blue-600 peer-checked:bg-blue-50 hover:bg-gray-50 transition"
							>
								<div class="flex items-start gap-4">
									<div class="flex-shrink-0 w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center text-2xl">
										🎨
									</div>
									<div class="flex-1">
										<h3 class="font-bold text-gray-900 mb-2 text-lg">
											Start from Scratch
										</h3>
										<p class="text-gray-600">
											Skip the defaults and create your own categories from the dashboard.
										</p>
									</div>
								</div>
							</div>
						</label>

						{#if errorMessage}
							<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
								{errorMessage}
							</div>
						{/if}
					</div>

					<div class="flex gap-4 mt-8">
						<button
							onclick={previousStep}
							disabled={isSubmitting}
							class="flex-1 px-6 py-3 border-2 border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition disabled:opacity-50"
						>
							Back
						</button>
						<button
							onclick={nextStep}
							disabled={isSubmitting}
							class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition disabled:opacity-50"
						>
							{isSubmitting ? 'Setting up...' : 'Complete Setup'}
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="text-center mt-6 text-sm text-gray-600">
			Step {currentStep} of 3
		</div>
	</div>
</div>
