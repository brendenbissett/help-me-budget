<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import type { PageData, ActionData } from './$types';
	import type { BudgetEntry } from '$lib/server/budget/budgets';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedEntry = $state<BudgetEntry | null>(null);
	let createEntryType = $state<'income' | 'expense'>('income');

	// Separate entries by type
	const incomeEntries = $derived(
		(data.budget?.entries || []).filter((e) => e.entry_type === 'income' && e.is_active)
	);
	const expenseEntries = $derived(
		(data.budget?.entries || []).filter((e) => e.entry_type === 'expense' && e.is_active)
	);

	function openCreateModal(type: 'income' | 'expense') {
		createEntryType = type;
		showCreateModal = true;
	}

	function openEditModal(entry: BudgetEntry) {
		selectedEntry = entry;
		showEditModal = true;
	}

	function openDeleteModal(entry: BudgetEntry) {
		selectedEntry = entry;
		showDeleteModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedEntry = null;
	}

	$effect(() => {
		if (form?.success) {
			closeModals();
		}
	});

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD'
		}).format(amount);
	}

	function formatFrequency(frequency: string): string {
		const map: Record<string, string> = {
			once_off: 'Once-off',
			daily: 'Daily',
			weekly: 'Weekly',
			fortnightly: 'Fortnightly',
			monthly: 'Monthly',
			annually: 'Annually'
		};
		return map[frequency] || frequency;
	}

	function getHealthColor(score: number): string {
		if (score >= 80) return 'bg-green-100 text-green-700 border-green-200';
		if (score >= 60) return 'bg-blue-100 text-blue-700 border-blue-200';
		if (score >= 40) return 'bg-yellow-100 text-yellow-700 border-yellow-200';
		if (score >= 20) return 'bg-orange-100 text-orange-700 border-orange-200';
		return 'bg-red-100 text-red-700 border-red-200';
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-start">
		<div>
			<button
				onclick={() => goto('/dashboard/budgets')}
				class="text-blue-600 hover:text-blue-700 mb-2 flex items-center gap-1 text-sm"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Back to Budgets
			</button>
			<h1 class="text-3xl font-bold text-gray-900">{data.budget.name}</h1>
			{#if data.budget.description}
				<p class="text-gray-600 mt-1">{data.budget.description}</p>
			{/if}
		</div>
		<div class="flex gap-3">
			{#if data.budget.is_active}
				<span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
					Active
				</span>
			{:else}
				<span class="px-3 py-1 bg-gray-100 text-gray-600 rounded-full text-sm font-medium">
					Archived
				</span>
			{/if}
		</div>
	</div>

	<!-- Budget Summary Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<!-- Monthly Income -->
		<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
			<div class="flex items-center gap-3 mb-2">
				<div class="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center">
					<svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
				</div>
				<h3 class="text-sm font-medium text-gray-600">Monthly Income</h3>
			</div>
			<p class="text-2xl font-bold text-gray-900">
				{formatCurrency(data.summary.total_monthly_income)}
			</p>
			<p class="text-xs text-gray-500 mt-1">{data.summary.income_entries_count} entries</p>
		</div>

		<!-- Monthly Expenses -->
		<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
			<div class="flex items-center gap-3 mb-2">
				<div class="w-10 h-10 bg-red-100 rounded-xl flex items-center justify-center">
					<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
					</svg>
				</div>
				<h3 class="text-sm font-medium text-gray-600">Monthly Expenses</h3>
			</div>
			<p class="text-2xl font-bold text-gray-900">
				{formatCurrency(data.summary.total_monthly_expenses)}
			</p>
			<p class="text-xs text-gray-500 mt-1">{data.summary.expense_entries_count} entries</p>
		</div>

		<!-- Net Balance -->
		<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
			<div class="flex items-center gap-3 mb-2">
				<div class="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center">
					<svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
						/>
					</svg>
				</div>
				<h3 class="text-sm font-medium text-gray-600">Monthly Net</h3>
			</div>
			<p
				class="text-2xl font-bold"
				class:text-green-600={data.summary.monthly_surplus_deficit > 0}
				class:text-red-600={data.summary.monthly_surplus_deficit < 0}
				class:text-gray-900={data.summary.monthly_surplus_deficit === 0}
			>
				{formatCurrency(data.summary.monthly_surplus_deficit)}
			</p>
			<p class="text-xs text-gray-500 mt-1">
				{data.summary.monthly_surplus_deficit >= 0 ? 'Surplus' : 'Deficit'}
			</p>
		</div>
	</div>

	<!-- Budget Health -->
	<div class="bg-white rounded-2xl p-6 border-2 {getHealthColor(data.health.score)}">
		<div class="flex items-center justify-between mb-3">
			<h3 class="text-lg font-bold">Budget Health</h3>
			<span class="text-2xl font-bold">{data.health.score}/100</span>
		</div>
		<div class="w-full bg-gray-200 rounded-full h-3 mb-3">
			<div
				class="h-3 rounded-full transition-all"
				style="width: {data.health.score}%; background-color: {data.health.color}"
			></div>
		</div>
		<p class="text-sm">{data.health.message}</p>
	</div>

	<!-- Two-Column Layout: Income & Expenses -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<!-- Income Column -->
		<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
			<div class="flex justify-between items-center mb-6">
				<h2 class="text-xl font-bold text-gray-900">Income</h2>
				<button
					onclick={() => openCreateModal('income')}
					class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-xl text-sm font-medium flex items-center gap-2 transition"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					Add Income
				</button>
			</div>

			{#if incomeEntries.length === 0}
				<div class="text-center py-12">
					<div class="text-4xl mb-3">💰</div>
					<p class="text-gray-600 mb-4">No income entries yet</p>
					<button
						onclick={() => openCreateModal('income')}
						class="text-green-600 hover:text-green-700 font-medium"
					>
						Add your first income source
					</button>
				</div>
			{:else}
				<div class="space-y-3">
					{#each incomeEntries as entry}
						<div
							class="border-2 border-gray-200 rounded-xl p-4 hover:border-green-300 transition group"
						>
							<div class="flex justify-between items-start mb-2">
								<div class="flex-1">
									<h4 class="font-semibold text-gray-900">{entry.name}</h4>
									{#if entry.description}
										<p class="text-sm text-gray-600 mt-1">{entry.description}</p>
									{/if}
								</div>
								<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition">
									<button
										onclick={() => openEditModal(entry)}
										class="text-gray-400 hover:text-blue-600 transition p-1"
										title="Edit"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
											/>
										</svg>
									</button>
									<button
										onclick={() => openDeleteModal(entry)}
										class="text-gray-400 hover:text-red-600 transition p-1"
										title="Delete"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
											/>
										</svg>
									</button>
								</div>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-lg font-bold text-green-600">{formatCurrency(entry.amount)}</span>
								<span
									class="text-xs px-2 py-1 bg-gray-100 text-gray-700 rounded-full font-medium"
								>
									{formatFrequency(entry.frequency)}
								</span>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Expenses Column -->
		<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
			<div class="flex justify-between items-center mb-6">
				<h2 class="text-xl font-bold text-gray-900">Expenses</h2>
				<button
					onclick={() => openCreateModal('expense')}
					class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-xl text-sm font-medium flex items-center gap-2 transition"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					Add Expense
				</button>
			</div>

			{#if expenseEntries.length === 0}
				<div class="text-center py-12">
					<div class="text-4xl mb-3">💸</div>
					<p class="text-gray-600 mb-4">No expense entries yet</p>
					<button
						onclick={() => openCreateModal('expense')}
						class="text-red-600 hover:text-red-700 font-medium"
					>
						Add your first expense
					</button>
				</div>
			{:else}
				<div class="space-y-3">
					{#each expenseEntries as entry}
						<div class="border-2 border-gray-200 rounded-xl p-4 hover:border-red-300 transition group">
							<div class="flex justify-between items-start mb-2">
								<div class="flex-1">
									<h4 class="font-semibold text-gray-900">{entry.name}</h4>
									{#if entry.description}
										<p class="text-sm text-gray-600 mt-1">{entry.description}</p>
									{/if}
								</div>
								<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition">
									<button
										onclick={() => openEditModal(entry)}
										class="text-gray-400 hover:text-blue-600 transition p-1"
										title="Edit"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
											/>
										</svg>
									</button>
									<button
										onclick={() => openDeleteModal(entry)}
										class="text-gray-400 hover:text-red-600 transition p-1"
										title="Delete"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
											/>
										</svg>
									</button>
								</div>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-lg font-bold text-red-600">{formatCurrency(entry.amount)}</span>
								<span
									class="text-xs px-2 py-1 bg-gray-100 text-gray-700 rounded-full font-medium"
								>
									{formatFrequency(entry.frequency)}
								</span>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>

	<!-- Error Message -->
	{#if form?.error}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
			{form.error}
		</div>
	{/if}
</div>

<!-- Create Entry Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full max-h-[90vh] overflow-y-auto"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Escape' && closeModals()}
			role="dialog"
			aria-modal="true"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">
				Add {createEntryType === 'income' ? 'Income' : 'Expense'}
			</h2>

			<form method="POST" action="?/createEntry" use:enhance>
				<input type="hidden" name="entry_type" value={createEntryType} />

				<div class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
							Name <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="name"
							name="name"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="e.g., Salary, Rent, Groceries"
						/>
					</div>

					<div>
						<label for="amount" class="block text-sm font-medium text-gray-700 mb-2">
							Amount <span class="text-red-500">*</span>
						</label>
						<input
							type="number"
							id="amount"
							name="amount"
							step="0.01"
							min="0"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="0.00"
						/>
					</div>

					<div>
						<label for="frequency" class="block text-sm font-medium text-gray-700 mb-2">
							Frequency <span class="text-red-500">*</span>
						</label>
						<select
							id="frequency"
							name="frequency"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							<option value="monthly">Monthly</option>
							<option value="weekly">Weekly</option>
							<option value="fortnightly">Fortnightly</option>
							<option value="daily">Daily</option>
							<option value="annually">Annually</option>
							<option value="once_off">Once-off</option>
						</select>
					</div>

					<div>
						<label for="start_date" class="block text-sm font-medium text-gray-700 mb-2">
							Start Date <span class="text-red-500">*</span>
						</label>
						<input
							type="date"
							id="start_date"
							name="start_date"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="end_date" class="block text-sm font-medium text-gray-700 mb-2">
							End Date (Optional)
						</label>
						<input
							type="date"
							id="end_date"
							name="end_date"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="description" class="block text-sm font-medium text-gray-700 mb-2">
							Description (Optional)
						</label>
						<textarea
							id="description"
							name="description"
							rows="3"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="Additional details..."
						></textarea>
					</div>
				</div>

				<div class="flex gap-3 mt-6">
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 px-6 py-3 border border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="flex-1 px-6 py-3 text-white font-semibold rounded-xl transition"
						class:bg-green-600={createEntryType === 'income'}
						class:hover:bg-green-700={createEntryType === 'income'}
						class:bg-red-600={createEntryType === 'expense'}
						class:hover:bg-red-700={createEntryType === 'expense'}
					>
						Add {createEntryType === 'income' ? 'Income' : 'Expense'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Entry Modal -->
{#if showEditModal && selectedEntry}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full max-h-[90vh] overflow-y-auto"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Escape' && closeModals()}
			role="dialog"
			aria-modal="true"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Edit Entry</h2>

			<form method="POST" action="?/updateEntry" use:enhance>
				<input type="hidden" name="entry_id" value={selectedEntry.id} />

				<div class="space-y-4">
					<div>
						<label for="edit_name" class="block text-sm font-medium text-gray-700 mb-2">
							Name <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="edit_name"
							name="name"
							value={selectedEntry.name}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_amount" class="block text-sm font-medium text-gray-700 mb-2">
							Amount <span class="text-red-500">*</span>
						</label>
						<input
							type="number"
							id="edit_amount"
							name="amount"
							value={selectedEntry.amount}
							step="0.01"
							min="0"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_frequency" class="block text-sm font-medium text-gray-700 mb-2">
							Frequency <span class="text-red-500">*</span>
						</label>
						<select
							id="edit_frequency"
							name="frequency"
							value={selectedEntry.frequency}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							<option value="monthly">Monthly</option>
							<option value="weekly">Weekly</option>
							<option value="fortnightly">Fortnightly</option>
							<option value="daily">Daily</option>
							<option value="annually">Annually</option>
							<option value="once_off">Once-off</option>
						</select>
					</div>

					<div>
						<label for="edit_start_date" class="block text-sm font-medium text-gray-700 mb-2">
							Start Date <span class="text-red-500">*</span>
						</label>
						<input
							type="date"
							id="edit_start_date"
							name="start_date"
							value={selectedEntry.start_date}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_end_date" class="block text-sm font-medium text-gray-700 mb-2">
							End Date (Optional)
						</label>
						<input
							type="date"
							id="edit_end_date"
							name="end_date"
							value={selectedEntry.end_date || ''}
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_description" class="block text-sm font-medium text-gray-700 mb-2">
							Description (Optional)
						</label>
						<textarea
							id="edit_description"
							name="description"
							rows="3"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>{selectedEntry.description || ''}</textarea>
					</div>

					<div class="flex items-center gap-3">
						<input
							type="checkbox"
							id="edit_is_active"
							name="is_active"
							value="true"
							checked={selectedEntry.is_active}
							class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
						/>
						<label for="edit_is_active" class="text-sm font-medium text-gray-700">
							Active Entry
						</label>
					</div>
				</div>

				<div class="flex gap-3 mt-6">
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 px-6 py-3 border border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition"
					>
						Save Changes
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal && selectedEntry}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Escape' && closeModals()}
			role="dialog"
			aria-modal="true"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Delete Entry?</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to delete <strong>{selectedEntry.name}</strong>? This will mark the
				entry as inactive.
			</p>

			<form method="POST" action="?/deleteEntry" use:enhance>
				<input type="hidden" name="entry_id" value={selectedEntry.id} />

				<div class="flex gap-3">
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 px-6 py-3 border border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="flex-1 px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-xl transition"
					>
						Delete
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
