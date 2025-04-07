<script lang="ts">
    import DataInput from "$lib/components/DataInput.svelte";
    import ParameterInput from "$lib/components/ParameterInput.svelte";
    import type { RemoteFile } from "$lib/types/TempFile";
    import { config } from "$lib/state.svelte";
    import type { PageProps } from "./$types";
    import { goto } from "$app/navigation";
    import ClInfo from "./CLInfo.svelte";

    let { data }: PageProps = $props();
    // $inspect(data);

    let parameterValues: {[name: string]: any} = $state({});
    let dataValues: {[name: string]: RemoteFile} = $state({});
    $inspect(parameterValues);
    $inspect(dataValues);

    let currentTab: 'parameters' | 'cli' = $state('parameters');

    function updateParameterValues(name: string, value: any) {
        parameterValues = {...parameterValues, [name]: value};
    }

    let parameterAreValid = $derived(
        !data.parameters ||
        Object.keys(data.parameters)
        .map(name => parameterValues[name] !== null && parameterValues[name] !== undefined && parameterValues[name] !== '')
        .reduce((a, b) => a && b, true)
    );

    let dataAreValid = $derived(
        !data.data ||
        Object.keys(data.data)
        .map(name => dataValues[name] !== null && dataValues[name] !== undefined)
        .reduce((a, b) => a && b, true)
    );

    let allValid = $derived(parameterAreValid && dataAreValid);

    let dockerImage = $derived(data.id!.split('::')[0]);
    let toolName = $derived(data.id!.split('::')[1]);

    function startRun() {
        const [dockerImage, toolName, ...o] = data.id!.split('::');
        if (!dockerImage || !toolName || data.name !== toolName || o.length > 0) {
            console.error(`Invalid tool slug: ${data.id}`) 
            return 
        }

        const payload = ({
            name: toolName,
            docker_image: dockerImage,
            parameters: {...parameterValues},
            data: Object.fromEntries(Object.entries(dataValues).map(([name, conf]) => ([name, conf.path])))
        })

        fetch(`${config.apiServer}/runs`, {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(res => res.json())
        .then(response => { 
            console.log(response);
            goto('/manager/runs');
        })
    }
</script>


<div>
    <h1 class="text-2xl font-bold text-gray-900">{data.title}</h1>
    <h4 class="text-md font-semibold text-gray-500">ID: {data.id}</h4>
    
    <p class="mt-2 text-gray-600">{data.description}</p>

    <div class="flex mt-4 ml-6">
        <button 
        class="px-4 pt-2 text-sm font-medium border-b-1 transition-colors duration-200" 
        onclick={() => currentTab = 'parameters'}
        class:border-indigo-500={currentTab === 'parameters'}
        class:text-indigo-600={currentTab === 'parameters'}
        class:border-transparent={currentTab !== 'parameters'}
        class:text-gray-500={currentTab !== 'parameters'}
        class:hover:text-gray-700={currentTab !== 'parameters'}
        class:hover:border-gray-300={currentTab !== 'parameters'}
        >
        Parameters
    </button>
    <button
        class="px-4 pt-2 text-sm font-medium border-b-1 transition-colors duration-200"
        onclick={() => currentTab = 'cli'}
        class:border-indigo-500={currentTab === 'cli'}
        class:text-indigo-600={currentTab === 'cli'}
        class:border-transparent={currentTab !== 'cli'}
        class:text-gray-500={currentTab !== 'cli'}
        class:hover:text-gray-700={currentTab !== 'cli'}
        class:hover:border-gray-300={currentTab !== 'cli'}
    >
        API Access</button>
    </div>

    {#if currentTab === 'parameters'}
        <div class="p-3 rounded-lg border border-gray-200 shadow-md mb-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-3">Parameters</h2>
                {#if data.parameters}    
                    {#each Object.entries(data.parameters) as [name, parameter]}
                    <ParameterInput {parameter} {name} oninput={value => updateParameterValues(name, value)} />
                {/each}
            {:else}
                    <p class="mt-2 text-gray-600">Tool {data.title} has no parameters defined.</p>
            {/if}
            {#if parameterAreValid}
                <p class="mt-2 text-green-500">All parameters are valid</p>
            {:else}
                <p class="mt-2 text-red-600">Some required parameters are not yet set.</p>
            {/if}
        </div>

        <div class="p-3 rounded-lg border border-gray-200 shadow-md my-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-3">Data</h2>
            {#if data.data}
                {#each Object.entries(data.data) as [name, dataSpec]}
                    <DataInput {name} data={dataSpec} onupload={f => f ? dataValues[name] = {...f} : delete dataValues[name]} /> 
                {/each}
                {#if dataAreValid}
                        <p class="mt-2 text-green-500">All data is valid</p>
                {:else}
                        <p class="mt-2 text-red-600">Some required data is not yet set.</p>
                {/if}
            {:else}
                <p class="mt-2 text-gray-600">Tool {data.title} does not require any data</p>
            {/if}
        </div>
    {:else if currentTab === 'cli'}
        <ClInfo id={data.id!} name={toolName} image={dockerImage} {parameterValues} {dataValues} />
    {/if}

</div>
{#if allValid}
<button class="w-full px-3 py-2 bg-green-500 text-white rounded-lg shadow-md hover:bg-green-600 transition-colors cursor-pointer" onclick={startRun}>
    Create
</button>
{/if}


