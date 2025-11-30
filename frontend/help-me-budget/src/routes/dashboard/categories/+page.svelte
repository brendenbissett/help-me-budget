<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedCategory = $state<any>(null);
	let activeTab = $state<'income' | 'expense'>('expense');

	// Common icons for categories
	const iconOptions = [
		'home', 'car', 'shopping-cart', 'zap', 'heart', 'film', 'shopping-bag',
		'user', 'book', 'shield', 'repeat', 'coffee', 'more-horizontal',
		'dollar-sign', 'briefcase', 'trending-up', 'gift', 'plus-circle',
		'music', 'phone', 'tool', 'truck', 'umbrella', 'scissors'
	];

	// Common colors for categories
	const colorOptions = [
		'#FF6B6B', '#4ECDC4', '#95E1D3', '#F38181', '#AA96DA', '#FCBAD3',
		'#A8D8EA', '#FFCFDF', '#C7CEEA', '#B5EAD7', '#FFD3B6', '#FFA8B6',
		'#D4A5A5', '#52B788', '#74C69D', '#95D5B2', '#B7E4C7', '#D8F3DC'
	];

	function openCreateModal(type: 'income' | 'expense') {
		activeTab = type;
		showCreateModal = true;
	}

	function openEditModal(category: any) {
		selectedCategory = category;
		showEditModal = true;
	}

	function openDeleteModal(category: any) {
		selectedCategory = category;
		showDeleteModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedCategory = null;
	}

	$effect(() => {
		if (form?.success) {
			closeModals();
		}
	});

	const incomeCategories = $derived(
		data.categories.filter(c => c.category_type === 'income' && c.is_active)
	);
	const expenseCategories = $derived(
		data.categories.filter(c => c.category_type === 'expense' && c.is_active)
	);
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Categories</h1>
			<p class="text-gray-600 mt-1">Organize your income and expenses</p>
		</div>
	</div>

	<!-- Empty State with Seed Option -->
	{#if data.categories.length === 0}
		<div class="bg-white rounded-2xl p-12 text-center">
			<div class="text-6xl mb-4">📊</div>
			<h3 class="text-xl font-semibold text-gray-900 mb-2">No categories yet</h3>
			<p class="text-gray-600 mb-6">Start with our recommended categories or create your own</p>
			<div class="flex justify-center gap-4">
				<form method="POST" action="?/seedDefaults" use:enhance>
					<button
						type="submit"
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
						Use Default Categories
					</button>
				</form>
				<button
					onclick={() => openCreateModal('expense')}
					class="border-2 border-blue-600 text-blue-600 font-semibold px-6 py-3 rounded-xl inline-flex items-center gap-2 hover:bg-blue-50 transition"
				>
					Create Custom Category
				</button>
			</div>
		</div>
	{:else}
		<!-- Tabs -->
		<div class="bg-white rounded-2xl p-2 flex gap-2">
			<button
				onclick={() => (activeTab = 'expense')}
				class="flex-1 py-3 px-6 rounded-xl font-semibold transition {activeTab === 'expense'
					? 'bg-red-50 text-red-600'
					: 'text-gray-600 hover:bg-gray-50'}"
			>
				Expenses ({expenseCategories.length})
			</button>
			<button
				onclick={() => (activeTab = 'income')}
				class="flex-1 py-3 px-6 rounded-xl font-semibold transition {activeTab === 'income'
					? 'bg-green-50 text-green-600'
					: 'text-gray-600 hover:bg-gray-50'}"
			>
				Income ({incomeCategories.length})
			</button>
		</div>

		<!-- Categories Grid -->
		{#if activeTab === 'expense'}
			<div>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-xl font-bold text-gray-900">Expense Categories</h2>
					<button
						onclick={() => openCreateModal('expense')}
						class="bg-red-600 hover:bg-red-700 text-white font-semibold px-4 py-2 rounded-lg flex items-center gap-2 transition"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 4v16m8-8H4"
							/>
						</svg>
						Add Expense Category
					</button>
				</div>

				{#if expenseCategories.length === 0}
					<div class="bg-white rounded-xl p-8 text-center text-gray-500">
						No expense categories yet. Click "Add Expense Category" to create one.
					</div>
				{:else}
					<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
						{#each expenseCategories as category}
							<div
								class="bg-white rounded-xl p-4 border-2 border-gray-200 hover:border-red-300 transition cursor-pointer group"
							>
								<div class="flex items-start justify-between mb-3">
									<div
										class="w-12 h-12 rounded-full flex items-center justify-center text-xl"
										style="background-color: {category.color || '#D4A5A5'}20"
									>
										{#if category.icon}
											<span>{category.icon}</span>
										{:else}
											<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
												/>
											</svg>
										{/if}
									</div>
									<div class="opacity-0 group-hover:opacity-100 transition flex gap-1">
										<button
											onclick={() => openEditModal(category)}
											class="text-gray-400 hover:text-blue-600 p-1"
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
											onclick={() => openDeleteModal(category)}
											class="text-gray-400 hover:text-red-600 p-1"
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
								<h3 class="font-semibold text-gray-900 text-sm">{category.name}</h3>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<div>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-xl font-bold text-gray-900">Income Categories</h2>
					<button
						onclick={() => openCreateModal('income')}
						class="bg-green-600 hover:bg-green-700 text-white font-semibold px-4 py-2 rounded-lg flex items-center gap-2 transition"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 4v16m8-8H4"
							/>
						</svg>
						Add Income Category
					</button>
				</div>

				{#if incomeCategories.length === 0}
					<div class="bg-white rounded-xl p-8 text-center text-gray-500">
						No income categories yet. Click "Add Income Category" to create one.
					</div>
				{:else}
					<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
						{#each incomeCategories as category}
							<div
								class="bg-white rounded-xl p-4 border-2 border-gray-200 hover:border-green-300 transition cursor-pointer group"
							>
								<div class="flex items-start justify-between mb-3">
									<div
										class="w-12 h-12 rounded-full flex items-center justify-center text-xl"
										style="background-color: {category.color || '#52B788'}20"
									>
										{#if category.icon}
											<span>{category.icon}</span>
										{:else}
											<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
												/>
											</svg>
										{/if}
									</div>
									<div class="opacity-0 group-hover:opacity-100 transition flex gap-1">
										<button
											onclick={() => openEditModal(category)}
											class="text-gray-400 hover:text-blue-600 p-1"
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
											onclick={() => openDeleteModal(category)}
											class="text-gray-400 hover:text-red-600 p-1"
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
								<h3 class="font-semibold text-gray-900 text-sm">{category.name}</h3>
							</div>
						{/each}
					</div>
				{/if}
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

<!-- Create Category Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full max-h-[90vh] overflow-y-auto"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">
				Add {activeTab === 'income' ? 'Income' : 'Expense'} Category
			</h2>

			<form method="POST" action="?/create" use:enhance>
				<input type="hidden" name="category_type" value={activeTab} />

				<div class="space-y-4">
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
							Category Name
						</label>
						<input
							type="text"
							id="name"
							name="name"
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="e.g., Groceries"
						/>
					</div>

					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">Color</label>
						<div class="grid grid-cols-6 gap-2">
							{#each colorOptions as color}
								<label class="cursor-pointer">
									<input
										type="radio"
										name="color"
										value={color}
										class="sr-only peer"
									/>
									<div
										class="w-10 h-10 rounded-lg border-2 border-transparent peer-checked:border-gray-900 peer-checked:scale-110 transition"
										style="background-color: {color}"
									></div>
								</label>
							{/each}
						</div>
					</div>

					<div>
						<label for="icon" class="block text-sm font-medium text-gray-700 mb-2">
							Icon (emoji or text)
						</label>
						<input
							type="text"
							id="icon"
							name="icon"
							maxlength="2"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="e.g., 🛒"
						/>
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
						Add Category
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Category Modal -->
{#if showEditModal && selectedCategory}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full max-h-[90vh] overflow-y-auto"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Edit Category</h2>

			<form method="POST" action="?/update" use:enhance>
				<input type="hidden" name="id" value={selectedCategory.id} />
				<input type="hidden" name="category_type" value={selectedCategory.category_type} />

				<div class="space-y-4">
					<div>
						<label for="edit_name" class="block text-sm font-medium text-gray-700 mb-2">
							Category Name
						</label>
						<input
							type="text"
							id="edit_name"
							name="name"
							value={selectedCategory.name}
							required
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
					</div>

					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">Color</label>
						<div class="grid grid-cols-6 gap-2">
							{#each colorOptions as color}
								<label class="cursor-pointer">
									<input
										type="radio"
										name="color"
										value={color}
										checked={color === selectedCategory.color}
										class="sr-only peer"
									/>
									<div
										class="w-10 h-10 rounded-lg border-2 border-transparent peer-checked:border-gray-900 peer-checked:scale-110 transition"
										style="background-color: {color}"
									></div>
								</label>
							{/each}
						</div>
					</div>

					<div>
						<label for="edit_icon" class="block text-sm font-medium text-gray-700 mb-2">
							Icon (emoji or text)
						</label>
						<input
							type="text"
							id="edit_icon"
							name="icon"
							value={selectedCategory.icon || ''}
							maxlength="2"
							class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
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
{#if showDeleteModal && selectedCategory}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
		onclick={closeModals}
	>
		<div
			class="bg-white rounded-2xl p-8 max-w-md w-full"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Delete Category?</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to delete <strong>{selectedCategory.name}</strong>? This will mark
				the category as inactive.
			</p>

			<form method="POST" action="?/delete" use:enhance>
				<input type="hidden" name="id" value={selectedCategory.id} />

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
