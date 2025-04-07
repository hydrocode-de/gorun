<script lang="ts">
    import type { PageProps } from "./$types";
    import type { ToolSpec } from "$lib/types/ToolSpec";

    let { data }: PageProps = $props();
    let specs: ToolSpec[] = $state(data.specs);
</script>

<div class="p-4">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Specs</h1>
    </div>

    <div class="space-y-4">
        {#each specs as spec}
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