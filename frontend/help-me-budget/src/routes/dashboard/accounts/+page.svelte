<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedAccount = $state<any>(null);

	// Account type icons and colors
	const accountTypeConfig = {
		checking: { icon: '🏦', color: 'from-blue-500 to-blue-600', label: 'Checking' },
		savings: { icon: '💰', color: 'from-green-500 to-green-600', label: 'Savings' },
		credit_card: { icon: '💳', color: 'from-purple-500 to-purple-600', label: 'Credit Card' },
		cash: { icon: '💵', color: 'from-yellow-500 to-yellow-600', label: 'Cash' },
		investment: { icon: '📈', color: 'from-indigo-500 to-indigo-600', label: 'Investment' }
	};

	function formatCurrency(amount: number, currency: string = 'USD'): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: currency
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function openCreateModal() {
		showCreateModal = true;
	}

	function openEditModal(account: any) {
		selectedAccount = account;
		showEditModal = true;
	}

	function openDeleteModal(account: any) {
		selectedAccount = account;
		showDeleteModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedAccount = null;
	}

	// Calculate totals by account type
	$effect(() => {
		if (form?.success) {
			closeModals();
			// Optionally show success toast
		}
	});

	const totalBalance = $derived(
		data.accounts.reduce((sum, account) => sum + (account.is_active ? account.balance : 0), 0)
	);
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Accounts</h1>
			<p class="text-gray-600 mt-1">Manage your bank accounts, credit cards, and cash</p>
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
			Add Account
		</button>
	</div>

	<!-- Total Balance Card -->
	<div class="bg-gradient-to-br from-blue-600 to-blue-700 rounded-2xl p-8 text-white">
		<p class="text-blue-100 text-sm mb-2">Total Balance</p>
		<h2 class="text-5xl font-bold">{formatCurrency(totalBalance)}</h2>
		<p class="text-blue-100 text-sm mt-3">{data.accounts.filter(a => a.is_active).length} active accounts</p>
	</div>

	<!-- Accounts Grid -->
	{#if data.accounts.length === 0}
		<div class="bg-white rounded-2xl p-12 text-center">
			<div class="text-6xl mb-4">🏦</div>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">No accounts yet</h3>
			<p class="text-gray-600 mb-6">Get started by adding your first account</p>
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
				Add Your First Account
			</button>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each data.accounts as account}
				{@const config = accountTypeConfig[account.account_type]}
				<div
					class="bg-white rounded-2xl p-6 border-2 {account.is_active
						? 'border-gray-200'
						: 'border-red-200 opacity-60'}"
				>
					<!-- Card Header -->
					<div class="flex justify-between items-start mb-4">
						<div class="w-12 h-12 bg-gradient-to-br {config.color} rounded-full flex items-center justify-center text-2xl">
							{config.icon}
						</div>
						<div class="flex gap-2">
							<button
								onclick={() => openEditModal(account)}
								class="text-gray-400 hover:text-blue-600 transition"
								title="Edit account"
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
								onclick={() => openDeleteModal(account)}
								class="text-gray-400 hover:text-red-600 transition"
								title="Delete account"
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

					<!-- Account Details -->
					<div class="space-y-3">
						<div>
							<p class="text-sm text-gray-500">{config.label}</p>
							<h3 class="text-xl font-bold text-gray-900">{account.name}</h3>
						</div>

						<div>
							<p class="text-sm text-gray-500">Balance</p>
							<p class="text-2xl font-bold text-gray-900">
								{formatCurrency(account.balance, account.currency)}
							</p>
						</div>

						<div class="pt-3 border-t border-gray-100 flex justify-between text-xs text-gray-500">
							<span>Created {formatDate(account.created_at)}</span>
							{#if !account.is_active}
								<span class="text-red-600 font-semibold">Inactive</span>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Error Message -->
	{#if form?.error}
		<div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-xl">
			{form.error}
		</div>
	{/if}
</div>

<!-- Create Account Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Add New Account</h2>

			<form method="POST" action="?/create" use:enhance>
				<div class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
							Account Name
						</label>
						<input
							type="text"
							id="name"
							name="name"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="e.g., Chase Checking"
						/>
					</div>

					<div>
						<label for="account_type" class="block text-sm font-medium text-gray-700 mb-2">
							Account Type
						</label>
						<select
							id="account_type"
							name="account_type"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							{#each Object.entries(accountTypeConfig) as [value, config]}
								<option {value}>{config.icon} {config.label}</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="balance" class="block text-sm font-medium text-gray-700 mb-2">
							Current Balance
						</label>
						<input
							type="number"
							id="balance"
							name="balance"
							step="0.01"
							value="0"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="0.00"
						/>
					</div>

					<div>
						<label for="currency" class="block text-sm font-medium text-gray-700 mb-2">
							Currency
						</label>
						<select
							id="currency"
							name="currency"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							<option value="USD">USD - US Dollar</option>
							<option value="EUR">EUR - Euro</option>
							<option value="GBP">GBP - British Pound</option>
							<option value="CAD">CAD - Canadian Dollar</option>
							<option value="AUD">AUD - Australian Dollar</option>
						</select>
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
						Add Account
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Account Modal -->
{#if showEditModal && selectedAccount}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Edit Account</h2>

			<form method="POST" action="?/update" use:enhance>
				<input type="hidden" name="id" value={selectedAccount.id} />

				<div class="space-y-4">
					<div>
						<label for="edit_name" class="block text-sm font-medium text-gray-700 mb-2">
							Account Name
						</label>
						<input
							type="text"
							id="edit_name"
							name="name"
							value={selectedAccount.name}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_account_type" class="block text-sm font-medium text-gray-700 mb-2">
							Account Type
						</label>
						<select
							id="edit_account_type"
							name="account_type"
							value={selectedAccount.account_type}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							{#each Object.entries(accountTypeConfig) as [value, config]}
								<option {value} selected={value === selectedAccount.account_type}>
									{config.icon} {config.label}
								</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="edit_balance" class="block text-sm font-medium text-gray-700 mb-2">
							Current Balance
						</label>
						<input
							type="number"
							id="edit_balance"
							name="balance"
							value={selectedAccount.balance}
							step="0.01"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label for="edit_currency" class="block text-sm font-medium text-gray-700 mb-2">
							Currency
						</label>
						<select
							id="edit_currency"
							name="currency"
							value={selectedAccount.currency}
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						>
							<option value="USD">USD - US Dollar</option>
							<option value="EUR">EUR - Euro</option>
							<option value="GBP">GBP - British Pound</option>
							<option value="CAD">CAD - Canadian Dollar</option>
							<option value="AUD">AUD - Australian Dollar</option>
						</select>
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
{#if showDeleteModal && selectedAccount}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Delete Account?</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to delete <strong>{selectedAccount.name}</strong>? This will mark
				the account as inactive.
			</p>

			<form method="POST" action="?/delete" use:enhance>
				<input type="hidden" name="id" value={selectedAccount.id} />

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
