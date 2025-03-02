<script lang="ts">
    import moment from 'moment';
    import type { RunState } from "$lib/types/RunState";
    import { page } from '$app/state';
    import { goto, invalidate, invalidateAll, replaceState } from '$app/navigation';
    import { config } from '$lib/state.svelte.js';

    let {data} = $props();

    let runs: RunState[] = $state(data.runs); 
    $effect(() => {
        runs = data.runs
    });
    $inspect(runs);

    // refresh 
    async function refresh() {
        invalidateAll();
    }

    // Define available status options
    let status: RunState["status"] | "all" = $state(page.url.searchParams.get('status') || 'all') as RunState["status"] | "all"; 
    const statusOptions = ['all', 'pending', 'running', 'finished', 'errored'] as const; 
    $inspect(status)

    // Handle status change
    async function handleStatusChange(event: Event) {
        const target = event.target as HTMLSelectElement;
        status = target.value as typeof status;

        const params = new URLSearchParams(page.url.searchParams);
        params.set('status', status);
        await goto(`?${params.toString()}`);
    }

    async function onStart(runId: number) {
        const runUrl = `${config.apiServer}/runs/${runId}/start`;
        const res = await fetch(runUrl, { method: 'POST'});
        const data = await res.json();
        $inspect(data);
        await refresh();
    }
</script>

<button onclick={refresh} class="text-blue-600 hover:text-blue-800 hover:cursor-pointer" title="Refresh" aria-label="Refresh">
    Refresh
</button>
<div class="relative overflow-x-auto shadow-md sm:rounded-lg">
    <table class="w-full text-sm text-left">
        <thead class="text-xs uppercase bg-gray-100">
            <tr>
                <th scope="col" class="px-6 py-3 font-semibold">
                    Title
                </th>
                <th scope="col" class="px-6 py-3 font-semibold">
                    <select 
                        class="text-xs uppercase bg-gray-100 font-semibold outline-none cursor-pointer"
                        bind:value={status}
                        onchange={handleStatusChange}
                    >
                        {#each statusOptions as option}
                            <option value={option}>
                                Status: {option}
                            </option>
                        {/each}
                    </select>
                </th>
                <th scope="col" class="px-6 py-3 font-semibold">
                    Created
                </th>
                <th scope="col" class="px-6 py-3 font-semibold">
                    Finished
                </th>
                <th scope="col" class="px-6 py-3 font-semibold">
                    Actions
                </th>
            </tr>
        </thead>
        <tbody>
            {#each runs as run}
                <tr class="bg-white border-b border-b-gray-300 hover:bg-gray-50">
                    <td class="px-6 py-4">
                        <a href="/manager/runs/{run.id}" class="font-medium text-blue-600 hover:underline">{run.title}</a>
                    </td>
                    <td class="px-6 py-4">
                        {run.status}
                    </td>
                    <td class="px-6 py-4">
                        {moment(run.created_at).fromNow()}
                    </td>
                    <td class="px-6 py-4">
                        {run.status === 'finished' ? moment(run.finished_at).fromNow() : null}
                        {run.status === 'errored' ? 'Errored' : null}
                    </td>
                    <td class="px-6 py-4 flex gap-2">
                        {#if run.status === 'running'}
                            Running...
                        {/if}
                        {#if run.status === 'pending'}
                        <button 
                            onclick={() => onStart(run.id)}
                            class="text-green-600 hover:text-green-800 hover:cursor-pointer" 
                            title="Start" 
                            aria-label="Start Run"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </button>
                        {/if}
                        {#if run.status === 'finished'}
                        <button 
                            class="text-blue-600 hover:text-blue-800 hover:cursor-pointer" 
                            title="Download" 
                            aria-label="Download Result"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                            </svg>
                        </button>
                        {/if}
                        {#if run.status !== 'running'}
                        <button disabled class="text-gray-200  hover:cursor-not-allowed" title="Delete" aria-label="Delete Run">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                        {/if}
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>