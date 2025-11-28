<script lang="ts">
	import { goto } from '$app/navigation';

	function handleGetStarted() {
		goto('/auth?mode=signup');
	}

	function handleLogin() {
		goto('/auth?mode=login');
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

	<!-- Right Side - Marketing Content -->
	<div class="w-full lg:w-1/2 flex items-center justify-center p-6">
		<div class="w-full max-w-md">
			<div class="flex items-center gap-2 mb-8 lg:hidden">
				<div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
					<span class="text-white font-bold text-lg">💰</span>
				</div>
				<span class="text-2xl font-bold text-gray-900">Help Me Budget</span>
			</div>

			<h1 class="text-4xl font-bold text-gray-900 mb-2">Take control of your finances</h1>
			<p class="text-gray-600 mb-8 text-lg">Track spending, set goals, and budget smarter with our simple, powerful tools.</p>

			<div class="space-y-4 mb-8">
				<div class="flex items-start gap-3">
					<div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
						<svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
					</div>
					<div>
						<h3 class="font-semibold text-gray-900">Track every transaction</h3>
						<p class="text-gray-600 text-sm">Automatically categorize and monitor all your spending in one place.</p>
					</div>
				</div>

				<div class="flex items-start gap-3">
					<div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
						<svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
					</div>
					<div>
						<h3 class="font-semibold text-gray-900">Set realistic goals</h3>
						<p class="text-gray-600 text-sm">Create budgets that work for you and watch your progress in real-time.</p>
					</div>
				</div>

				<div class="flex items-start gap-3">
					<div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
						<svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
					</div>
					<div>
						<h3 class="font-semibold text-gray-900">Save for what matters</h3>
						<p class="text-gray-600 text-sm">Plan for the future with smart savings strategies and insights.</p>
					</div>
				</div>
			</div>

			<button
				onclick={handleGetStarted}
				class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-xl transition duration-200 mb-4"
			>
				Get Started Free
			</button>

			<p class="text-center text-gray-600 text-sm">
				Already have an account?
				<button onclick={handleLogin} class="text-blue-600 font-semibold hover:underline">
					Sign In
				</button>
			</p>

			<p class="text-center text-gray-500 text-xs mt-6">
				By creating an account, you agree to our <span class="text-gray-700 font-semibold">Privacy Policy</span> and
				<span class="text-gray-700 font-semibold">Terms and Conditions</span>
			</p>
		</div>
	</div>
</div>
