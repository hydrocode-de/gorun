<script lang="ts">
    import { logout, initializeAuth } from '$lib/auth.svelte';
    import { config } from '$lib/state.svelte';
	import '../app.css';

	let { children } = $props();
    
    initializeAuth();
	$inspect(config);
</script>

<div class="min-h-screen flex flex-col">
	<nav class="bg-gray-800 text-white p-4 shadow-lg">
		<div class="container mx-auto flex items-center justify-between">
			<div class="text-xl font-bold">GoRun</div>
			<div class="space-x-4">
				<a href="/manager" class="hover:text-gray-300">Home</a>
				<a href="/manager/specs" class="hover:text-gray-300">Tools</a>
				<a href="/manager/runs" class="hover:text-gray-300">Runs</a>
				{#if !config.refreshToken}
					<a href="/manager/login" class="hover:text-gray-300">Login</a>
				{:else}
					<button onclick={() => logout()} class="hover:text-gray-300">Logout</button>
				{/if}
			</div>
		</div>
	</nav>

	<main class="flex-1 overflow-auto">
		<div class="container mx-auto p-4">
			{@render children()}
		</div>
	</main>
</div>
