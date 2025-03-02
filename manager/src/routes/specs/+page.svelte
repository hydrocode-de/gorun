<script lang="ts">
    import {config } from "$lib/state.svelte";
    import { tools } from "../state.svelte";

    $inspect(tools);
    let isLoading = $state(false);

    function refresh() {
        isLoading = true;
        fetch(`${config.apiServer}/specs`)
        .then(res => res.json())
        .then(data => {

            tools.count = data.count;
            tools.specs = [...data.tools];
            tools.lastUpdated = new Date();
            isLoading = false;
        })

    }
</script>

<div class="p-4">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Specs</h1>
        <button 
            class="p-2 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer"
            aria-label="refresh"
            class:animate-spin={isLoading}
            disabled={isLoading}
            onclick={refresh}>
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
                <path d="M21 3v5h-5"/>
            </svg>
        </button>
    </div>

    <div class="space-y-4">
        {#each tools.specs as spec}
            <a href="/manager/specs/{spec.id}" class="block">
                <div class="p-4 bg-white rounded-lg shadow hover:shadow-md transition-shadow border border-gray-200">
                    <h2 class="text-lg font-semibold text-gray-900">{spec.title}</h2>
                    <h4 class="text-md font-semibold text-gray-500">ID: {spec.id}</h4>
                    <p class="mt-2 text-gray-600">{spec.description}</p>
                    <div class="mt-2 flex items-center text-sm text-blue-500">
                        View details →
                    </div>
                </div>
            </a>
        {/each}
    </div>
</div>