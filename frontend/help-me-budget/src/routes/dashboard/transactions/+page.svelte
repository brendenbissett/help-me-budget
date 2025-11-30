<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedTransaction = $state<any>(null);

	// Filter state
	let filterAccountId = $state('');
	let filterCategoryId = $state('');
	let filterStartDate = $state('');
	let filterEndDate = $state('');

	function openCreateModal() {
		showCreateModal = true;
	}

	function openEditModal(transaction: any) {
		selectedTransaction = transaction;
		showEditModal = true;
	}

	function openDeleteModal(transaction: any) {
		selectedTransaction = transaction;
		showDeleteModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedTransaction = null;
	}

	function applyFilters() {
		const params = new URLSearchParams();
		if (filterAccountId) params.append('account_id', filterAccountId);
		if (filterCategoryId) params.append('category_id', filterCategoryId);
		if (filterStartDate) params.append('start_date', filterStartDate);
		if (filterEndDate) params.append('end_date', filterEndDate);

		const queryString = params.toString();
		goto(`/dashboard/transactions${queryString ? '?' + queryString : ''}`);
	}

	function clearFilters() {
		filterAccountId = '';
		filterCategoryId = '';
		filterStartDate = '';
		filterEndDate = '';
		goto('/dashboard/transactions');
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

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function getCategoryName(categoryId: string | null | undefined): string {
		if (!categoryId) return 'Uncategorized';
		const category = data.categories.find((c) => c.id === categoryId);
		return category?.name || 'Unknown';
	}

	function getAccountName(accountId: string): string {
		const account = data.accounts.find((a) => a.id === accountId);
		return account?.name || 'Unknown';
	}

	// Get today's date in YYYY-MM-DD format for default date
	const today = new Date().toISOString().split('T')[0];
</script>

<div class="space-y-6">
	<!-- Load Error Message -->
	{#if data.loadError}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
			{data.loadError}
		</div>
	{/if}

	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Transactions</h1>
			<p class="text-gray-600 mt-1">Track your income and expenses</p>
		</div>
		<div class="flex gap-3">
			<a
				href="/dashboard/transactions/review"
				class="bg-purple-600 hover:bg-purple-700 text-white font-semibold px-6 py-3 rounded-xl flex items-center gap-2 transition"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
					/>
				</svg>
				Review Matches
			</a>
			<button
				onclick={openCreateModal}
				class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-3 rounded-xl flex items-center gap-2 transition"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Add Transaction
			</button>
		</div>
	</div>

	<!-- Filters -->
	<div class="bg-white rounded-2xl p-6 border-2 border-gray-100">
		<h3 class="text-lg font-bold text-gray-900 mb-4">Filters</h3>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div>
				<label for="filter_account" class="block text-sm font-medium text-gray-700 mb-2">
					Account
				</label>
				<select
					id="filter_account"
					bind:value={filterAccountId}
					class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
				>
					<option value="">All Accounts</option>
					{#each data.accounts as account}
						<option value={account.id}>{account.name}</option>
					{/each}
				</select>
			</div>

			<div>
				<label for="filter_category" class="block text-sm font-medium text-gray-700 mb-2">
					Category
				</label>
				<select
					id="filter_category"
					bind:value={filterCategoryId}
					class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
				>
					<option value="">All Categories</option>
					{#each data.categories as category}
						<option value={category.id}>{category.name}</option>
					{/each}
				</select>
			</div>

			<div>
				<label for="filter_start_date" class="block text-sm font-medium text-gray-700 mb-2">
					From Date
				</label>
				<input
					type="date"
					id="filter_start_date"
					bind:value={filterStartDate}
					class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
				/>
			</div>

			<div>
				<label for="filter_end_date" class="block text-sm font-medium text-gray-700 mb-2">
					To Date
				</label>
				<input
					type="date"
					id="filter_end_date"
					bind:value={filterEndDate}
					class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
				/>
			</div>
		</div>

		<div class="flex gap-3 mt-4">
			<button
				onclick={applyFilters}
				class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition"
			>
				Apply Filters
			</button>
			<button
				onclick={clearFilters}
				class="px-4 py-2 border border-gray-300 text-gray-700 font-medium rounded-xl hover:bg-gray-50 transition"
			>
				Clear
			</button>
		</div>
	</div>

	<!-- Empty State -->
	{#if !data.transactions || data.transactions.length === 0}
		<div class="bg-white rounded-2xl p-12 text-center">
			<div class="text-6xl mb-4">💸</div>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">No transactions yet</h3>
			<p class="text-gray-600 mb-6">Start tracking your income and expenses</p>
			<button
				onclick={openCreateModal}
				class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-6 py-3 rounded-xl inline-flex items-center gap-2 transition"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Add Your First Transaction
			</button>
		</div>
	{:else}
		<!-- Transactions List -->
		<div class="bg-white rounded-2xl border-2 border-gray-100 overflow-hidden">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="bg-gray-50 border-b border-gray-200">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Date
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Description
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Category
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Account
							</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
								Amount
							</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each data.transactions as transaction}
							<tr class="hover:bg-gray-50 transition">
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
									{formatDate(transaction.transaction_date)}
								</td>
								<td class="px-6 py-4 text-sm text-gray-900">
									<div class="flex items-center gap-2">
										<span
											class="w-2 h-2 rounded-full"
											class:bg-green-500={transaction.transaction_type === 'income'}
											class:bg-red-500={transaction.transaction_type === 'expense'}
										></span>
										{transaction.description || 'No description'}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
									{getCategoryName(transaction.category_id)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
									{getAccountName(transaction.account_id)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<span
										class:text-green-600={transaction.transaction_type === 'income'}
										class:text-red-600={transaction.transaction_type === 'expense'}
									>
										{transaction.transaction_type === 'income' ? '+' : '-'}
										{formatCurrency(transaction.amount)}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end gap-2">
										<button
											onclick={() => openEditModal(transaction)}
											class="text-blue-600 hover:text-blue-700"
											title="Edit"
										>
											<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
												/>
											</svg>
										</button>
										<button
											onclick={() => openDeleteModal(transaction)}
											class="text-red-600 hover:text-red-700"
											title="Delete"
										>
											<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												/>
											</svg>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}

	<!-- Form Error Message -->
	{#if form?.error}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
			{form.error}
		</div>
	{/if}
</div>

<!-- Create Transaction Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
		onclick={closeModals}
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
		role="button"
		tabindex="0"
	>
		<div
			class="bg-white rounded-2xl p-6 max-w-lg w-full"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Add Transaction</h2>

			<form method="POST" action="?/create" use:enhance>
				<div class="space-y-4">
					<!-- Transaction Type -->
					<div>
						<label for="create_type" class="block text-sm font-medium text-gray-700 mb-2">
							Type <span class="text-red-500">*</span>
						</label>
						<select
							name="transaction_type"
							id="create_type"
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							<option value="expense">Expense</option>
							<option value="income">Income</option>
						</select>
					</div>

					<!-- Amount -->
					<div>
						<label for="create_amount" class="block text-sm font-medium text-gray-700 mb-2">
							Amount <span class="text-red-500">*</span>
						</label>
						<input
							type="number"
							name="amount"
							id="create_amount"
							step="0.01"
							min="0.01"
							required
							placeholder="0.00"
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Account -->
					<div>
						<label for="create_account" class="block text-sm font-medium text-gray-700 mb-2">
							Account <span class="text-red-500">*</span>
						</label>
						<select
							name="account_id"
							id="create_account"
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							<option value="">Select an account</option>
							{#each data.accounts as account}
								<option value={account.id}>{account.name}</option>
							{/each}
						</select>
					</div>

					<!-- Category -->
					<div>
						<label for="create_category" class="block text-sm font-medium text-gray-700 mb-2">
							Category
						</label>
						<select
							name="category_id"
							id="create_category"
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							<option value="">Uncategorized</option>
							{#each data.categories as category}
								<option value={category.id}>{category.name}</option>
							{/each}
						</select>
					</div>

					<!-- Date -->
					<div>
						<label for="create_date" class="block text-sm font-medium text-gray-700 mb-2">
							Date <span class="text-red-500">*</span>
						</label>
						<input
							type="date"
							name="transaction_date"
							id="create_date"
							value={today}
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Description -->
					<div>
						<label for="create_description" class="block text-sm font-medium text-gray-700 mb-2">
							Description
						</label>
						<input
							type="text"
							name="description"
							id="create_description"
							placeholder="What was this for?"
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Notes -->
					<div>
						<label for="create_notes" class="block text-sm font-medium text-gray-700 mb-2">
							Notes
						</label>
						<textarea
							name="notes"
							id="create_notes"
							rows="3"
							placeholder="Additional details..."
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						></textarea>
					</div>
				</div>

				<div class="flex gap-3 mt-6">
					<button
						type="submit"
						class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-xl transition"
					>
						Add Transaction
					</button>
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 border border-gray-300 text-gray-700 font-medium py-3 px-6 rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Transaction Modal -->
{#if showEditModal && selectedTransaction}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
		onclick={closeModals}
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
		role="button"
		tabindex="0"
	>
		<div
			class="bg-white rounded-2xl p-6 max-w-lg w-full"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Edit Transaction</h2>

			<form method="POST" action="?/update" use:enhance>
				<input type="hidden" name="transaction_id" value={selectedTransaction.id} />

				<div class="space-y-4">
					<!-- Transaction Type -->
					<div>
						<label for="edit_type" class="block text-sm font-medium text-gray-700 mb-2">
							Type <span class="text-red-500">*</span>
						</label>
						<select
							name="transaction_type"
							id="edit_type"
							value={selectedTransaction.transaction_type}
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							<option value="expense">Expense</option>
							<option value="income">Income</option>
						</select>
					</div>

					<!-- Amount -->
					<div>
						<label for="edit_amount" class="block text-sm font-medium text-gray-700 mb-2">
							Amount <span class="text-red-500">*</span>
						</label>
						<input
							type="number"
							name="amount"
							id="edit_amount"
							value={selectedTransaction.amount}
							step="0.01"
							min="0.01"
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Account -->
					<div>
						<label for="edit_account" class="block text-sm font-medium text-gray-700 mb-2">
							Account <span class="text-red-500">*</span>
						</label>
						<select
							name="account_id"
							id="edit_account"
							value={selectedTransaction.account_id}
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							{#each data.accounts as account}
								<option value={account.id}>{account.name}</option>
							{/each}
						</select>
					</div>

					<!-- Category -->
					<div>
						<label for="edit_category" class="block text-sm font-medium text-gray-700 mb-2">
							Category
						</label>
						<select
							name="category_id"
							id="edit_category"
							value={selectedTransaction.category_id || ''}
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						>
							<option value="">Uncategorized</option>
							{#each data.categories as category}
								<option value={category.id}>{category.name}</option>
							{/each}
						</select>
					</div>

					<!-- Date -->
					<div>
						<label for="edit_date" class="block text-sm font-medium text-gray-700 mb-2">
							Date <span class="text-red-500">*</span>
						</label>
						<input
							type="date"
							name="transaction_date"
							id="edit_date"
							value={selectedTransaction.transaction_date}
							required
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Description -->
					<div>
						<label for="edit_description" class="block text-sm font-medium text-gray-700 mb-2">
							Description
						</label>
						<input
							type="text"
							name="description"
							id="edit_description"
							value={selectedTransaction.description || ''}
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Notes -->
					<div>
						<label for="edit_notes" class="block text-sm font-medium text-gray-700 mb-2">
							Notes
						</label>
						<textarea
							name="notes"
							id="edit_notes"
							rows="3"
							value={selectedTransaction.notes || ''}
							class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500"
						></textarea>
					</div>
				</div>

				<div class="flex gap-3 mt-6">
					<button
						type="submit"
						class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-xl transition"
					>
						Save Changes
					</button>
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 border border-gray-300 text-gray-700 font-medium py-3 px-6 rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal && selectedTransaction}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
		onclick={closeModals}
		onkeydown={(e) => e.key === 'Escape' && closeModals()}
		role="button"
		tabindex="0"
	>
		<div
			class="bg-white rounded-2xl p-6 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Delete Transaction?</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to delete this transaction? This action cannot be undone.
			</p>

			<div class="bg-gray-50 border border-gray-200 rounded-xl p-4 mb-6">
				<div class="text-sm text-gray-600">
					<strong class="text-gray-900">{selectedTransaction.description || 'No description'}</strong>
				</div>
				<div class="text-sm text-gray-600 mt-1">
					Amount: <span
						class:text-green-600={selectedTransaction.transaction_type === 'income'}
						class:text-red-600={selectedTransaction.transaction_type === 'expense'}
						class="font-semibold"
					>
						{selectedTransaction.transaction_type === 'income' ? '+' : '-'}
						{formatCurrency(selectedTransaction.amount)}
					</span>
				</div>
				<div class="text-sm text-gray-600 mt-1">
					Date: {formatDate(selectedTransaction.transaction_date)}
				</div>
			</div>

			<form method="POST" action="?/delete" use:enhance>
				<input type="hidden" name="transaction_id" value={selectedTransaction.id} />

				<div class="flex gap-3">
					<button
						type="submit"
						class="flex-1 bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-6 rounded-xl transition"
					>
						Delete
					</button>
					<button
						type="button"
						onclick={closeModals}
						class="flex-1 border border-gray-300 text-gray-700 font-medium py-3 px-6 rounded-xl hover:bg-gray-50 transition"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
