<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedBudget = $state<any>(null);

	function openCreateModal() {
		showCreateModal = true;
	}

	function openEditModal(budget: any) {
		selectedBudget = budget;
		showEditModal = true;
	}

	function openDeleteModal(budget: any) {
		selectedBudget = budget;
		showDeleteModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedBudget = null;
	}

	function viewBudget(budgetId: string) {
		goto(`/dashboard/budgets/${budgetId}`);
	}

	$effect(() => {
		if (form?.success) {
			closeModals();
		}
	});

	const activeBudgets = $derived((data.budgets || []).filter(b => b.is_active));
	const inactiveBudgets = $derived((data.budgets || []).filter(b => !b.is_active));

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
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
			<h1 class="text-3xl font-bold text-gray-900">Budgets</h1>
			<p class="text-gray-600 mt-1">Plan and track your income and expenses</p>
		</div>
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
			Create Budget
		</button>
	</div>

	<!-- Empty State -->
	{#if !data.budgets || data.budgets.length === 0}
		<div class="bg-white rounded-2xl p-12 text-center">
			<div class="text-6xl mb-4">📊</div>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">No budgets yet</h3>
			<p class="text-gray-600 mb-6">
				Create your first budget to start planning your finances
			</p>
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
				Create Your First Budget
			</button>
		</div>
	{:else}
		<!-- Active Budgets -->
		{#if activeBudgets.length > 0}
			<div>
				<h2 class="text-xl font-bold text-gray-900 mb-4">Active Budgets</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each activeBudgets as budget}
						<div
							class="bg-white rounded-2xl p-6 border-2 border-blue-200 hover:border-blue-400 transition cursor-pointer group"
							onclick={() => viewBudget(budget.id)}
						>
							<!-- Header -->
							<div class="flex justify-between items-start mb-4">
								<div class="flex-1">
									<h3 class="text-xl font-bold text-gray-900 mb-1">{budget.name}</h3>
									{#if budget.description}
										<p class="text-sm text-gray-600">{budget.description}</p>
									{/if}
								</div>
								<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition">
									<button
										onclick={(e) => {
											e.stopPropagation();
											openEditModal(budget);
										}}
										class="text-gray-400 hover:text-blue-600 transition p-1"
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
										onclick={(e) => {
											e.stopPropagation();
											openDeleteModal(budget);
										}}
										class="text-gray-400 hover:text-red-600 transition p-1"
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
							</div>

							<!-- Status Badge -->
							<div class="mb-4">
								<span
									class="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium"
								>
									<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
										<path
											fill-rule="evenodd"
											d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
											clip-rule="evenodd"
										/>
									</svg>
									Active
								</span>
							</div>

							<!-- Footer -->
							<div class="pt-4 border-t border-gray-100">
								<p class="text-xs text-gray-500">Created {formatDate(budget.created_at)}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Inactive Budgets -->
		{#if inactiveBudgets.length > 0}
			<div>
				<h2 class="text-xl font-bold text-gray-900 mb-4">Archived Budgets</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each inactiveBudgets as budget}
						<div
							class="bg-white rounded-2xl p-6 border-2 border-gray-200 opacity-60 hover:opacity-100 transition cursor-pointer group"
							onclick={() => viewBudget(budget.id)}
						>
							<!-- Header -->
							<div class="flex justify-between items-start mb-4">
								<div class="flex-1">
									<h3 class="text-xl font-bold text-gray-900 mb-1">{budget.name}</h3>
									{#if budget.description}
										<p class="text-sm text-gray-600">{budget.description}</p>
									{/if}
								</div>
								<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition">
									<button
										onclick={(e) => {
											e.stopPropagation();
											openEditModal(budget);
										}}
										class="text-gray-400 hover:text-blue-600 transition p-1"
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
								</div>
							</div>

							<!-- Status Badge -->
							<div class="mb-4">
								<span
									class="inline-flex items-center gap-1 px-3 py-1 bg-gray-100 text-gray-600 rounded-full text-sm font-medium"
								>
									Archived
								</span>
							</div>

							<!-- Footer -->
							<div class="pt-4 border-t border-gray-100">
								<p class="text-xs text-gray-500">Created {formatDate(budget.created_at)}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}

	<!-- Error Message -->
	{#if form?.error}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
			{form.error}
		</div>
	{/if}
</div>

<!-- Create Budget Modal -->
{#if showCreateModal}
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
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Create New Budget</h2>

			<form method="POST" action="?/create" use:enhance>
				<div class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
							Budget Name
						</label>
						<input
							type="text"
							id="name"
							name="name"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="e.g., Monthly Budget 2025"
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
							placeholder="What's this budget for?"
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
						class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition"
					>
						Create Budget
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Budget Modal -->
{#if showEditModal && selectedBudget}
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
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Edit Budget</h2>

			<form method="POST" action="?/update" use:enhance>
				<input type="hidden" name="id" value={selectedBudget.id} />

				<div class="space-y-4">
					<div>
						<label for="edit_name" class="block text-sm font-medium text-gray-700 mb-2">
							Budget Name
						</label>
						<input
							type="text"
							id="edit_name"
							name="name"
							value={selectedBudget.name}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_description" class="block text-sm font-medium text-gray-700 mb-2">
							Description
						</label>
						<textarea
							id="edit_description"
							name="description"
							rows="3"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>{selectedBudget.description || ''}</textarea>
					</div>

					<div class="flex items-center gap-3">
						<input
							type="checkbox"
							id="edit_is_active"
							name="is_active"
							value="true"
							checked={selectedBudget.is_active}
							class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
						/>
						<label for="edit_is_active" class="text-sm font-medium text-gray-700">
							Active Budget
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
{#if showDeleteModal && selectedBudget}
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
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Delete Budget?</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to archive <strong>{selectedBudget.name}</strong>? This will mark
				the budget as inactive.
			</p>

			<form method="POST" action="?/delete" use:enhance>
				<input type="hidden" name="id" value={selectedBudget.id} />

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
						Archive
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
