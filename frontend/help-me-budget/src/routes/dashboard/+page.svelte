<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric'
		});
	}

	function getHealthColorClass(color: string): string {
		const colorMap: Record<string, string> = {
			green: 'bg-green-500',
			blue: 'bg-blue-500',
			yellow: 'bg-yellow-500',
			orange: 'bg-orange-500',
			red: 'bg-red-500'
		};
		return colorMap[color] || 'bg-gray-500';
	}

	const summary = data.summary;
</script>

{#if data.error}
	<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl mb-6">
		{data.error}
	</div>
{/if}

{#if !summary}
	<div class="bg-white rounded-2xl p-12 text-center">
		<div class="text-6xl mb-4">📊</div>
		<h3 class="text-xl font-semibold text-gray-900 mb-2">No data available</h3>
		<p class="text-gray-600 mb-6">Start by adding accounts, budgets, and transactions</p>
		<div class="flex gap-4 justify-center">
			<a
				href="/dashboard/accounts"
				class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-3 rounded-xl transition"
			>
				Add Account
			</a>
			<a
				href="/dashboard/budgets"
				class="border border-gray-300 text-gray-700 font-medium px-6 py-3 rounded-xl hover:bg-gray-50 transition"
			>
				Create Budget
			</a>
		</div>
	</div>
{:else}
	<div class="space-y-6">
		<!-- Page Header -->
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
			<p class="text-gray-600 mt-1">Your financial overview at a glance</p>
		</div>

		<!-- Top Stats Row -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-6">
			<!-- Total Balance -->
			<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
				<div class="flex items-center gap-3 mb-2">
					<div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm text-gray-600">Total Balance</p>
						<p class="text-2xl font-bold text-gray-900">{formatCurrency(summary.total_balance)}</p>
					</div>
				</div>
				<p class="text-xs text-gray-500 mt-2">Across {summary.account_count} account{summary.account_count !== 1 ? 's' : ''}</p>
			</div>

			<!-- Month Income -->
			<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
				<div class="flex items-center gap-3 mb-2">
					<div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm text-gray-600">Month Income</p>
						<p class="text-2xl font-bold text-green-600">{formatCurrency(summary.month_to_date_income)}</p>
					</div>
				</div>
				<p class="text-xs text-gray-500 mt-2">Budget: {formatCurrency(summary.budgeted_monthly_income)}</p>
			</div>

			<!-- Month Expenses -->
			<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
				<div class="flex items-center gap-3 mb-2">
					<div class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm text-gray-600">Month Expenses</p>
						<p class="text-2xl font-bold text-red-600">{formatCurrency(summary.month_to_date_expenses)}</p>
					</div>
				</div>
				<p class="text-xs text-gray-500 mt-2">Budget: {formatCurrency(summary.budgeted_monthly_expense)}</p>
			</div>

			<!-- Net -->
			<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
				<div class="flex items-center gap-3 mb-2">
					<div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
						<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
							/>
						</svg>
					</div>
					<div>
						<p class="text-sm text-gray-600">Month Net</p>
						<p
							class="text-2xl font-bold"
							class:text-green-600={summary.month_to_date_net > 0}
							class:text-red-600={summary.month_to_date_net < 0}
							class:text-gray-900={summary.month_to_date_net === 0}
						>
							{formatCurrency(summary.month_to_date_net)}
						</p>
					</div>
				</div>
				<p class="text-xs text-gray-500 mt-2">
					{summary.month_to_date_net > 0 ? 'Surplus' : summary.month_to_date_net < 0 ? 'Deficit' : 'Break even'}
				</p>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Left Column (2/3 width) -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Budget Health -->
				{#if summary.budget_health_score > 0}
					<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
						<h3 class="text-lg font-bold text-gray-900 mb-4">Budget Health</h3>
						<div class="flex items-center gap-4 mb-4">
							<div class="flex-1">
								<div class="flex justify-between text-sm mb-2">
									<span class="font-medium text-gray-700">{summary.budget_health_status}</span>
									<span class="text-gray-600">{summary.budget_health_score}/100</span>
								</div>
								<div class="h-3 bg-gray-200 rounded-full overflow-hidden">
									<div
										class="h-full {getHealthColorClass(summary.budget_health_color)} transition-all"
										style="width: {summary.budget_health_score}%"
									></div>
								</div>
							</div>
						</div>
						<p class="text-sm text-gray-600">{summary.budget_health_message}</p>
					</div>
				{/if}

				<!-- Recent Transactions -->
				<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
					<div class="flex justify-between items-center mb-4">
						<h3 class="text-lg font-bold text-gray-900">Recent Transactions</h3>
						<a href="/dashboard/transactions" class="text-sm text-blue-600 hover:text-blue-700 font-medium">
							View All →
						</a>
					</div>

					{#if summary.recent_transactions.length === 0}
						<div class="text-center py-8 text-gray-500">
							<p class="mb-2">No transactions yet</p>
							<a href="/dashboard/transactions" class="text-blue-600 hover:text-blue-700 text-sm">
								Add your first transaction
							</a>
						</div>
					{:else}
						<div class="space-y-3">
							{#each summary.recent_transactions as transaction}
								<div class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition">
									<div class="flex items-center gap-3">
										<div
											class="w-10 h-10 rounded-lg flex items-center justify-center"
											class:bg-green-100={transaction.transaction_type === 'income'}
											class:bg-red-100={transaction.transaction_type === 'expense'}
										>
											<span class:text-green-600={transaction.transaction_type === 'income'} class:text-red-600={transaction.transaction_type === 'expense'}>
												{transaction.transaction_type === 'income' ? '↓' : '↑'}
											</span>
										</div>
										<div>
											<p class="font-medium text-gray-900">
												{transaction.description || 'No description'}
											</p>
											<p class="text-sm text-gray-500">{formatDate(transaction.transaction_date)}</p>
										</div>
									</div>
									<span
										class="font-semibold"
										class:text-green-600={transaction.transaction_type === 'income'}
										class:text-red-600={transaction.transaction_type === 'expense'}
									>
										{transaction.transaction_type === 'income' ? '+' : '-'}{formatCurrency(transaction.amount)}
									</span>
								</div>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Spending by Category -->
				{#if summary.spending_by_category.length > 0}
					<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
						<h3 class="text-lg font-bold text-gray-900 mb-4">Spending by Category</h3>
						<div class="space-y-3">
							{#each summary.spending_by_category.slice(0, 5) as category}
								<div>
									<div class="flex justify-between text-sm mb-1">
										<span class="font-medium text-gray-700">{category.category_name}</span>
										<span class="text-gray-600">{formatCurrency(category.total_amount)} ({category.percentage.toFixed(1)}%)</span>
									</div>
									<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
										<div
											class="h-full {category.color ? 'bg-blue-500' : 'bg-gray-400'}"
											style="width: {category.percentage}%"
										></div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Right Column (1/3 width) -->
			<div class="space-y-6">
				<!-- Quick Actions -->
				<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
					<h3 class="text-lg font-bold text-gray-900 mb-4">Quick Actions</h3>
					<div class="space-y-3">
						<a
							href="/dashboard/transactions"
							class="block w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-4 rounded-xl text-center transition"
						>
							+ Add Transaction
						</a>
						<a
							href="/dashboard/budgets"
							class="block w-full border border-gray-300 text-gray-700 font-medium py-3 px-4 rounded-xl text-center hover:bg-gray-50 transition"
						>
							View Budgets
						</a>
						<a
							href="/dashboard/accounts"
							class="block w-full border border-gray-300 text-gray-700 font-medium py-3 px-4 rounded-xl text-center hover:bg-gray-50 transition"
						>
							Manage Accounts
						</a>
					</div>
				</div>

				<!-- Upcoming Bills -->
				{#if summary.upcoming_bills.length > 0}
					<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
						<h3 class="text-lg font-bold text-gray-900 mb-4">Upcoming Bills</h3>
						<div class="space-y-3">
							{#each summary.upcoming_bills.slice(0, 5) as bill}
								<div class="flex justify-between items-center p-3 rounded-lg hover:bg-gray-50 transition">
									<div>
										<p class="font-medium text-gray-900">{bill.name}</p>
										<p class="text-sm text-gray-500">Due: {formatDate(bill.due_date)}</p>
									</div>
									<span class="font-semibold text-gray-900">{formatCurrency(bill.amount)}</span>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
