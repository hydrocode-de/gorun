<script lang="ts">
    import moment from "moment";
    import type { PageProps } from "./$types";

    let {data}: PageProps = $props();
    $inspect(data);
</script>

<div>
    <h1 class="text-2xl font-bold text-gray-900">{data.title}</h1>
    <p class="mt-2 text-gray-600">{data.description}</p>

    {#if data.status === 'running'}
        <div class="mt-2 text-sm text-gray-600">
            Running since {moment(data.started_at).fromNow()}
        </div>
        {:else if data.status === 'finished'}
            <div class="mt-2 text-sm text-gray-600">
                Finished {moment(data.finished_at).fromNow()} {data.finished_at}
            </div>
            <i>Inspect the results</i>
        {:else if data.status === 'errored'}
            <div class="mt-2 text-sm text-gray-600">
                Errored {moment(data.finished_at).fromNow()}
            </div>
            <p class="mt-2 text-sm text-red-500">{data.error_message}</p>
        {:else if data.status === 'pending'}
            <button 
                class="w-full px-3 py-2 bg-green-500 text-white rounded-lg shadow-md hover:bg-green-600 transition-colors cursor-pointer" 
                onclick={() => console.log('start')}
            >
                Start
            </button>
        {/if}
</div>