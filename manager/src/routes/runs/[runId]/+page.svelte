<script lang="ts">
    import moment from "moment";
    import type { PageProps } from "./$types";
    import { bytesToSize } from "$lib/helper";
    import FinishedRun from "./FinishedRun.svelte";

    let { data }: PageProps = $props();
    let run = data.run;
    let files = data.files;
    $inspect(data);

</script>

{#if run}
<div>
    <h1 class="text-2xl font-bold text-gray-900">{run.title}</h1>
    <p class="mt-2 text-gray-600">{run.description}</p>

    {#if run.status === 'running'}
        <div class="mt-2 text-sm text-gray-600">
            Running since {moment(run.started_at).fromNow()}
        </div>
    {:else if run.status === 'finished'}
        <div class="flex flex-row justify-between mt-2 text-sm text-gray-600">
            <span>Finished {moment(run.finished_at).fromNow()}</span>
            <span>{files.length} results ({bytesToSize(files.map(f => f.size).reduce((a, b) => a + b, 0))})</span> 
        </div>
        <FinishedRun {run} {files} />
    {:else if run.status === 'errored'}
        <div class="mt-2 text-sm text-gray-600">
            Errored {moment(run.finished_at).fromNow()}
        </div>
        <p class="mt-2 text-sm text-red-500">{run.error_message}</p>
    {:else if run.status === 'pending'}
        <button 
            disabled
            class="w-full px-3 py-2 bg-green-500 text-white rounded-lg shadow-md hover:bg-green-600 transition-colors cursor-pointer" 
            onclick={() => console.log('start')}
        >
            Start
        </button>
    {/if}
</div>
{:else}
<div class="flex flex-col items-center justify-center">
    <div class="text-lg font-bold text-gray-900">No run found. Try to refresh the page, this should not happen.</div>
</div>
{/if}