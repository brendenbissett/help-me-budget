<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	// Helper to format currency
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	// Helper to format percentage
	function formatPercent(value: number): string {
		return `${value.toFixed(1)}%`;
	}

	// Helper to get variance color class
	function getVarianceClass(variance: number): string {
		if (variance > 0) return 'text-green-600'; // Under budget (good)
		if (variance < 0) return 'text-red-600'; // Over budget (bad)
		return 'text-gray-600'; // Exactly on budget
	}

	// Helper to get variance label
	function getVarianceLabel(variance: number): string {
		if (variance > 0) return 'Under Budget';
		if (variance < 0) return 'Over Budget';
		return 'On Budget';
	}

	// Group spending trends by month for display
	function groupTrendsByMonth() {
		const grouped: Record<string, { month: string; categories: any[] }> = {};

		data.spendingTrends.forEach((trend) => {
			if (!grouped[trend.month]) {
				grouped[trend.month] = {
					month: trend.month,
					categories: []
				};
			}
			grouped[trend.month].categories.push(trend);
		});

		return Object.values(grouped).sort((a, b) => b.month.localeCompare(a.month));
	}

	const monthlyTrends = groupTrendsByMonth();

	// Calculate projection summary stats
	const projectionSummary = data.cashFlowProjection.reduce(
		(acc, day) => {
			acc.totalIncome += day.projected_income;
			acc.totalExpenses += day.projected_expenses;
			if (day.projected_balance < acc.lowestBalance) {
				acc.lowestBalance = day.projected_balance;
				acc.lowestBalanceDate = day.date;
			}
			return acc;
		},
		{
			totalIncome: 0,
			totalExpenses: 0,
			lowestBalance: data.totalBalance,
			lowestBalanceDate: ''
		}
	);

	const endingBalance =
		data.cashFlowProjection.length > 0
			? data.cashFlowProjection[data.cashFlowProjection.length - 1].projected_balance
			: data.totalBalance;
</script>

<div class="p-6">
	<!-- Header -->
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Reports & Insights</h1>
		<p class="text-gray-600 mt-2">Understand your spending patterns and financial trends</p>
	</div>

	<!-- Load Error Message -->
	{#if data.loadError}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl mb-6">
			{data.loadError}
		</div>
	{/if}

	<!-- Top Summary Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
		<!-- Current Balance -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-600 mb-1">Current Balance</p>
					<p class="text-2xl font-bold text-gray-900">{formatCurrency(data.totalBalance)}</p>
				</div>
				<div class="bg-blue-100 p-3 rounded-lg">
					<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				</div>
			</div>
		</div>

		<!-- Projected Balance (90 days) -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-600 mb-1">Projected (90d)</p>
					<p
						class={`text-2xl font-bold ${endingBalance >= data.totalBalance ? 'text-green-600' : 'text-red-600'}`}
					>
						{formatCurrency(endingBalance)}
					</p>
					<p class="text-xs text-gray-500 mt-1">
						{endingBalance >= data.totalBalance ? '+' : ''}{formatCurrency(endingBalance - data.totalBalance)}
					</p>
				</div>
				<div class="bg-purple-100 p-3 rounded-lg">
					<svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
						/>
					</svg>
				</div>
			</div>
		</div>

		<!-- Lowest Projected Balance -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm text-gray-600 mb-1">Lowest Point (90d)</p>
					<p
						class={`text-2xl font-bold ${projectionSummary.lowestBalance >= 0 ? 'text-gray-900' : 'text-red-600'}`}
					>
						{formatCurrency(projectionSummary.lowestBalance)}
					</p>
					<p class="text-xs text-gray-500 mt-1">
						{new Date(projectionSummary.lowestBalanceDate).toLocaleDateString()}
					</p>
				</div>
				<div class="bg-orange-100 p-3 rounded-lg">
					<svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
						/>
					</svg>
				</div>
			</div>
		</div>
	</div>

	<!-- Two Column Layout -->
	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<!-- Left Column (2/3) -->
		<div class="lg:col-span-2 space-y-6">
			<!-- Budget Variance -->
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Budget vs Actual (This Month)</h2>

				{#if data.budgetVariance.length === 0}
					<p class="text-gray-500 text-center py-8">
						No budget entries found. Create a budget to see variance analysis.
					</p>
				{:else}
					<div class="space-y-3">
						{#each data.budgetVariance as item}
							<div class="border-b border-gray-100 pb-3 last:border-0">
								<div class="flex items-center justify-between mb-2">
									<div>
										<p class="font-medium text-gray-900">{item.entry_name}</p>
										<p class="text-sm text-gray-500">{item.category}</p>
									</div>
									<div class="text-right">
										<p class={`font-semibold ${getVarianceClass(item.variance)}`}>
											{getVarianceLabel(item.variance)}
										</p>
										<p class="text-sm text-gray-600">
											{formatCurrency(item.actual)} / {formatCurrency(item.budgeted)}
										</p>
									</div>
								</div>
								<div class="w-full bg-gray-200 rounded-full h-2">
									<div
										class={`h-2 rounded-full ${item.actual > item.budgeted ? 'bg-red-500' : 'bg-green-500'}`}
										style={`width: ${Math.min((item.actual / item.budgeted) * 100, 100)}%`}
									></div>
								</div>
								<p class="text-xs text-gray-500 mt-1">
									{formatPercent(Math.abs(item.variance_pct))}
									{item.variance > 0 ? 'under' : item.variance < 0 ? 'over' : 'on'} budget
								</p>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Cash Flow Projection -->
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Cash Flow Projection (90 Days)</h2>

				{#if data.cashFlowProjection.length === 0}
					<p class="text-gray-500 text-center py-8">
						No projection available. Add budget entries to see cash flow forecast.
					</p>
				{:else}
					<!-- Summary Stats -->
					<div class="grid grid-cols-3 gap-4 mb-6 pb-6 border-b border-gray-200">
						<div>
							<p class="text-sm text-gray-600">Total Income</p>
							<p class="text-lg font-semibold text-green-600">
								{formatCurrency(projectionSummary.totalIncome)}
							</p>
						</div>
						<div>
							<p class="text-sm text-gray-600">Total Expenses</p>
							<p class="text-lg font-semibold text-red-600">
								{formatCurrency(projectionSummary.totalExpenses)}
							</p>
						</div>
						<div>
							<p class="text-sm text-gray-600">Net Change</p>
							<p
								class={`text-lg font-semibold ${projectionSummary.totalIncome - projectionSummary.totalExpenses >= 0 ? 'text-green-600' : 'text-red-600'}`}
							>
								{formatCurrency(projectionSummary.totalIncome - projectionSummary.totalExpenses)}
							</p>
						</div>
					</div>

					<!-- Simple Text-Based Chart -->
					<div class="space-y-2 max-h-96 overflow-y-auto">
						{#each data.cashFlowProjection.filter((_, i) => i % 7 === 0 || i === data.cashFlowProjection.length - 1) as day}
							<div class="flex items-center gap-3">
								<p class="text-xs text-gray-600 w-24">{new Date(day.date).toLocaleDateString()}</p>
								<div class="flex-1 bg-gray-100 rounded h-6 relative">
									<div
										class={`h-6 rounded ${day.projected_balance >= 0 ? 'bg-green-500' : 'bg-red-500'}`}
										style={`width: ${Math.max(Math.abs(day.projected_balance) / Math.max(...data.cashFlowProjection.map((d) => Math.abs(d.projected_balance))), 0.05) * 100}%`}
									></div>
								</div>
								<p class={`text-sm font-medium w-28 text-right ${day.projected_balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
									{formatCurrency(day.projected_balance)}
								</p>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<!-- Right Column (1/3) -->
		<div class="space-y-6">
			<!-- Top Expenses -->
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Top Expenses (This Month)</h2>

				{#if data.topExpenses.length === 0}
					<p class="text-gray-500 text-center py-8">No expenses tracked this month.</p>
				{:else}
					<div class="space-y-3">
						{#each data.topExpenses as expense}
							<div>
								<div class="flex items-center justify-between mb-1">
									<p class="font-medium text-gray-900">{expense.category_name}</p>
									<p class="text-sm font-semibold text-gray-900">
										{formatCurrency(expense.total_amount)}
									</p>
								</div>
								<div class="w-full bg-gray-200 rounded-full h-2 mb-1">
									<div
										class="bg-red-500 h-2 rounded-full"
										style={`width: ${expense.percentage}%`}
									></div>
								</div>
								<div class="flex items-center justify-between text-xs text-gray-500">
									<span>{expense.count} transactions</span>
									<span>{formatPercent(expense.percentage)}</span>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Spending Trends (Last 6 Months) -->
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Spending Trends (6 Months)</h2>

				{#if monthlyTrends.length === 0}
					<p class="text-gray-500 text-center py-8">No spending data available.</p>
				{:else}
					<div class="space-y-4">
						{#each monthlyTrends.slice(0, 6) as monthData}
							<div>
								<p class="font-medium text-gray-900 mb-2">
									{new Date(monthData.month + '-01').toLocaleDateString('en-US', {
										month: 'long',
										year: 'numeric'
									})}
								</p>
								{#each monthData.categories.slice(0, 3) as category}
									<div class="flex items-center justify-between text-sm mb-1">
										<span class="text-gray-600">{category.category}</span>
										<span class="font-medium text-gray-900">{formatCurrency(category.amount)}</span>
									</div>
								{/each}
								<p class="text-xs text-gray-500 mt-1">
									Total: {formatCurrency(
										monthData.categories.reduce((sum, cat) => sum + cat.amount, 0)
									)}
								</p>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>
