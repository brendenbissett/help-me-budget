<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	$: errorMessage = $page.url.searchParams.get('error') || 'Unknown authentication error occurred';
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Authentication Error
			</h2>
			<div class="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
				<p class="text-sm text-red-800">
					{errorMessage}
				</p>
			</div>
			<p class="mt-4 text-center text-sm text-gray-600">
				Please try logging in again. If the problem persists, contact support.
			</p>
		</div>
		<div class="mt-8 space-y-6">
			<button
				type="button"
				onclick={() => goto('/auth')}
				class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
			>
				Try Again
			</button>
			<button
				type="button"
				onclick={() => goto('/')}
				class="group relative w-full flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
			>
				Return to Home
			</button>
		</div>
	</div>
</div>
