<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<div class="min-h-screen flex flex-col">
	<nav class="bg-gray-800 text-white p-4 shadow-lg">
		<div class="container mx-auto flex items-center justify-between">
			<div class="text-xl font-bold">GoRun</div>
			<div class="space-x-4">
				<a href="/manager" class="hover:text-gray-300">Home</a>
				<a href="/manager/specs" class="hover:text-gray-300">Tools</a>
			</div>
		</div>
	</nav>

	<main class="flex-1 overflow-auto">
		<div class="container mx-auto p-4">
			{@render children()}
		</div>
	</main>
</div>
