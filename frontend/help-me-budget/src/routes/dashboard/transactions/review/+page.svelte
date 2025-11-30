<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	// State for teach mode modal
	let showTeachModal = $state(false);
	let selectedTransactionId = $state<string | null>(null);
	let selectedBudgetEntryId = $state<string>('');
	let createRules = $state(true);
	let amountTolerance = $state<number>(2.0);

	// State for suggestions
	let loadingSuggestions = $state<Record<string, boolean>>({});
	let suggestions = $state<Record<string, any>>({});

	// Helper to get account name
	function getAccountName(accountId: string): string {
		const account = data.accounts.find((a) => a.id === accountId);
		return account?.name || 'Unknown Account';
	}

	// Helper to get category name
	function getCategoryName(categoryId: string | null | undefined): string {
		if (!categoryId) return 'Uncategorized';
		const category = data.categories.find((c) => c.id === categoryId);
		return category?.name || 'Unknown';
	}

	// Open teach modal for a transaction
	function openTeachModal(transactionId: string) {
		selectedTransactionId = transactionId;
		selectedBudgetEntryId = '';
		createRules = true;
		amountTolerance = 2.0;
		showTeachModal = true;
	}

	// Close teach modal
	function closeTeachModal() {
		showTeachModal = false;
		selectedTransactionId = null;
		selectedBudgetEntryId = '';
	}

	// Load suggestions for a transaction
	async function loadSuggestions(transactionId: string) {
		loadingSuggestions[transactionId] = true;

		const formData = new FormData();
		formData.append('transaction_id', transactionId);

		try {
			const response = await fetch('?/getSuggestions', {
				method: 'POST',
				body: formData
			});

			const result = await response.json();
			if (result.data?.success) {
				suggestions[transactionId] = result.data.suggestions;
			}
		} catch (error) {
			console.error('Error loading suggestions:', error);
		} finally {
			loadingSuggestions[transactionId] = false;
		}
	}

	// Get confidence badge color
	function getConfidenceBadgeClass(level: string): string {
		return level === 'auto_high'
			? 'bg-green-100 text-green-800'
			: 'bg-yellow-100 text-yellow-800';
	}

	// Get confidence label
	function getConfidenceLabel(level: string): string {
		return level === 'auto_high' ? 'High Confidence' : 'Low Confidence';
	}
</script>

<div class="p-6">
	<!-- Header -->
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Match Review</h1>
		<p class="text-gray-600 mt-2">
			Review and match unmatched transactions to your budget entries
		</p>
	</div>

	<!-- Bulk Auto-Match Button -->
	{#if data.unmatchedTransactions.length > 0}
		<div class="mb-6">
			<form method="POST" action="?/bulkAutoMatch" use:enhance>
				<button
					type="submit"
					class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium transition-colors"
				>
					Auto-Match All ({data.unmatchedTransactions.length})
				</button>
			</form>
		</div>
	{/if}

	<!-- No Active Budget Warning -->
	{#if !data.budget}
		<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
			<div class="flex items-start">
				<svg
					class="w-5 h-5 text-yellow-600 mt-0.5 mr-3"
					fill="currentColor"
					viewBox="0 0 20 20"
				>
					<path
						fill-rule="evenodd"
						d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
						clip-rule="evenodd"
					/>
				</svg>
				<div>
					<h3 class="font-semibold text-yellow-900">No Active Budget</h3>
					<p class="text-yellow-700 text-sm mt-1">
						You need an active budget with entries to match transactions. Please create a budget
						first.
					</p>
					<a
						href="/dashboard/budgets"
						class="text-yellow-800 underline hover:text-yellow-900 text-sm mt-2 inline-block"
					>
						Go to Budgets →
					</a>
				</div>
			</div>
		</div>
	{/if}

	<!-- Unmatched Transactions List -->
	{#if data.unmatchedTransactions.length === 0}
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
			<svg
				class="w-16 h-16 text-gray-300 mx-auto mb-4"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">All Caught Up!</h3>
			<p class="text-gray-600">All your transactions are matched to budget entries.</p>
			<a
				href="/dashboard/transactions"
				class="mt-4 inline-block bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg font-medium transition-colors"
			>
				View Transactions
			</a>
		</div>
	{:else}
		<div class="space-y-4">
			{#each data.unmatchedTransactions as transaction}
				<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
					<!-- Transaction Details -->
					<div class="flex items-start justify-between mb-4">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
								<span
									class={transaction.transaction_type === 'income'
										? 'text-green-600 text-xl'
										: 'text-red-600 text-xl'}
								>
									{transaction.transaction_type === 'income' ? '↓' : '↑'}
								</span>
								<div>
									<h3 class="font-semibold text-gray-900">
										{transaction.description || 'No description'}
									</h3>
									<p class="text-sm text-gray-600">
										{transaction.transaction_date} • {getAccountName(transaction.account_id)} •
										{getCategoryName(transaction.category_id)}
									</p>
								</div>
							</div>
						</div>
						<div class="text-right">
							<p
								class={`text-lg font-bold ${transaction.transaction_type === 'income' ? 'text-green-600' : 'text-red-600'}`}
							>
								${transaction.amount.toFixed(2)}
							</p>
						</div>
					</div>

					<!-- Actions -->
					<div class="flex items-center gap-3">
						<button
							onclick={() => loadSuggestions(transaction.id)}
							disabled={loadingSuggestions[transaction.id]}
							class="text-blue-600 hover:text-blue-700 font-medium text-sm disabled:opacity-50"
						>
							{loadingSuggestions[transaction.id] ? 'Loading...' : 'Find Matches'}
						</button>
						<button
							onclick={() => openTeachModal(transaction.id)}
							class="text-purple-600 hover:text-purple-700 font-medium text-sm"
						>
							Teach Match
						</button>
					</div>

					<!-- Suggestions (if loaded) -->
					{#if suggestions[transaction.id]?.suggestions?.length > 0}
						<div class="mt-4 pt-4 border-t border-gray-200">
							<h4 class="font-medium text-gray-900 mb-3">Suggested Matches:</h4>
							<div class="space-y-2">
								{#each suggestions[transaction.id].suggestions as suggestion}
									<div
										class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
									>
										<div class="flex-1">
											<div class="flex items-center gap-2 mb-1">
												<p class="font-medium text-gray-900">{suggestion.budget_entry.name}</p>
												<span
													class={`px-2 py-0.5 rounded-full text-xs font-medium ${getConfidenceBadgeClass(suggestion.confidence_level)}`}
												>
													{suggestion.confidence_score}% - {getConfidenceLabel(
														suggestion.confidence_level
													)}
												</span>
											</div>
											<p class="text-sm text-gray-600">
												${suggestion.budget_entry.amount.toFixed(2)} •
												{suggestion.budget_entry.frequency}
											</p>
											<div class="flex flex-wrap gap-1 mt-2">
												{#each suggestion.match_reasons as reason}
													<span class="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded">
														{reason}
													</span>
												{/each}
											</div>
										</div>
										<form
											method="POST"
											action="?/teach"
											use:enhance
											class="ml-4"
										>
											<input type="hidden" name="transaction_id" value={transaction.id} />
											<input
												type="hidden"
												name="budget_entry_id"
												value={suggestion.budget_entry.id}
											/>
											<input type="hidden" name="create_rules" value="true" />
											<input type="hidden" name="amount_tolerance" value="2.0" />
											<button
												type="submit"
												class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm font-medium transition-colors"
											>
												Link & Learn
											</button>
										</form>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Teach Mode Modal -->
{#if showTeachModal && selectedTransactionId}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Teach Match</h2>
			<p class="text-gray-600 mb-6">
				Link this transaction to a budget entry and optionally create matching rules for future
				transactions.
			</p>

			<form method="POST" action="?/teach" use:enhance onsubmit={closeTeachModal}>
				<input type="hidden" name="transaction_id" value={selectedTransactionId} />

				<!-- Budget Entry Selection -->
				<div class="mb-4">
					<label for="budget_entry_id" class="block text-sm font-medium text-gray-700 mb-2">
						Budget Entry
					</label>
					<select
						id="budget_entry_id"
						name="budget_entry_id"
						bind:value={selectedBudgetEntryId}
						required
						class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					>
						<option value="">Select a budget entry...</option>
						{#if data.budget?.entries}
							{#each data.budget.entries as entry}
								<option value={entry.id}>
									{entry.name} (${entry.amount.toFixed(2)} - {entry.frequency})
								</option>
							{/each}
						{/if}
					</select>
				</div>

				<!-- Create Rules Checkbox -->
				<div class="mb-4">
					<label class="flex items-center">
						<input
							type="checkbox"
							name="create_rules"
							bind:checked={createRules}
							value="true"
							class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<span class="ml-2 text-sm text-gray-700">
							Create matching rules for future transactions
						</span>
					</label>
				</div>

				<!-- Amount Tolerance -->
				{#if createRules}
					<div class="mb-6">
						<label for="amount_tolerance" class="block text-sm font-medium text-gray-700 mb-2">
							Amount Tolerance (±$)
						</label>
						<input
							type="number"
							id="amount_tolerance"
							name="amount_tolerance"
							bind:value={amountTolerance}
							step="0.01"
							min="0"
							class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
						<p class="text-xs text-gray-500 mt-1">
							How much the amount can vary (e.g., $2.00 means within ±$2)
						</p>
					</div>
				{/if}

				<!-- Action Buttons -->
				<div class="flex gap-3">
					<button
						type="button"
						onclick={closeTeachModal}
						class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium transition-colors"
					>
						Link {createRules ? '& Learn' : ''}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
