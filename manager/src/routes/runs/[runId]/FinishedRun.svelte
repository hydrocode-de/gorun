<script lang="ts">
    import { bytesToSize } from "$lib/helper";
    import moment from "moment";

    import { config } from "$lib/state.svelte";
    import type { RunState } from "$lib/types/RunState";
    import type { ResultFile } from "$lib/types/ResultFile";
    interface $$Props {
        run: RunState;
        files: ResultFile[];
    }

    let { run, files }: $$Props = $props();
    let activeTab = $state('files'); // or 'files' or 'logs'

    async function downloadFile(fileName: string) {
        const response = await fetch(`${config.apiServer}/runs/${run.id}/results/${fileName}`, {
            headers: {
                'X-User-ID': config.auth.user.id
            }
        });
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    }
</script>

<div class="w-full">
    <!-- Tab Navigation -->
    <div class="border-b border-gray-200">
        <nav class="-mb-px flex" aria-label="Tabs">
            <button
                class={`
                    w-24 py-2 px-3 text-center border-b-2 font-medium text-sm
                    ${activeTab === 'files' 
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                `}
                onclick={() => activeTab = 'files'}
            >
                Files
            </button>
            <button
                class={`
                    w-24 py-2 px-3 text-center border-b-2 font-medium text-sm
                    ${activeTab === 'details' 
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                `}
                onclick={() => activeTab = 'details'}
            >
                Details
            </button>
            <button
                class={`
                    w-24 py-2 px-3 text-center border-b-2 font-medium text-sm
                    ${activeTab === 'logs' 
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                `}
                onclick={() => activeTab = 'logs'}
            >
                Logs
            </button>
        </nav>
    </div>

    <div class="mt-4">
        {#if activeTab === 'details'}
            <div>
            <code>
                <pre>{JSON.stringify(run, null, 2)}</pre>
            </code>
            </div>
        {:else if activeTab === 'files'}
            <div>
                <table class="w-full text-sm text-left">
                    <thead class="text-xs uppercase bg-gray-50">
                        <tr>
                            <th scope="col" class="px-6 py-3">File name</th>
                            <th scope="col" class="px-6 py-3">Size</th>
                            <th scope="col" class="px-6 py-3">Last modified</th>
                            <th scope="col" class="px-6 py-3">Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each files as file}
                        <tr class="bg-white border-b border-b-gray-300 hover:bg-gray-50">
                            <td class="px-6 py-4">{file.name}</td>
                            <td class="px-6 py-4">{bytesToSize(file.size)}</td>
                            <td class="px-6 py-4">{moment(file.lastModified).fromNow()}</td>
                            <td class="px-6 py-4">
                                <button 
                                    onclick={() => downloadFile(file.name)}
                                    class="text-blue-600 hover:text-blue-800 hover:cursor-pointer"
                                    title="Download" 
                                    aria-label="Download Result"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                    </svg>
                                </button> 
                            </td>
                        </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {:else if activeTab === 'logs'}
            <div>
                <h2 class="text-lg font-semibold">StdOut</h2>
                
                <h2 class="text-lg font-semibold">StdErr</h2>
            </div>
        {/if}
    </div>
</div>
