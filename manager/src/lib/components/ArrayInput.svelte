<script lang="ts">
    import type { InputArray } from "$lib/types/InputParams";
    import type { ParameterSpec } from "$lib/types/ToolSpec";

    interface ArrayInputProps {
        parameter: ParameterSpec
        value: InputArray | null
        oninput: (value: InputArray) => void
    }
    let {parameter, value, oninput}: ArrayInputProps = $props()

    let arrayValue: InputArray = $state(value || (parameter.default && typeof parameter.default === 'object' ? parameter.default as InputArray : []))
    let newItemValue = $state('');

    function addItem() {
        // Check if the value is empty based on its type
        if (parameter.type === 'string' && newItemValue.trim() === '') return;
        if (parameter.type === 'integer' || parameter.type === 'float') {
            if (newItemValue === '' || isNaN(Number(newItemValue))) return;
        }
        if (parameter.type === 'boolean' && newItemValue === '') return;
        if (parameter.type === 'datetime' && newItemValue === '') return;
        if (parameter.type === 'enum' && newItemValue === '') return;
        
        let typedValue: string | number | boolean | Date | null;
        
        switch (parameter.type) {
            case 'string':
                typedValue = newItemValue;
                break;
            case 'integer':
            case 'float':
                typedValue = Number(newItemValue);
                break;
            case 'boolean':
                typedValue = newItemValue.toLowerCase() === 'true';
                break;
            case 'datetime':
                typedValue = new Date(newItemValue);
                break;
            default:
                typedValue = newItemValue;
        }
        
        arrayValue = [...arrayValue, typedValue] as InputArray;
        oninput(arrayValue);
        newItemValue = '';
    }

    function removeItem(index: number) {
        arrayValue = arrayValue.filter((_, i) => i !== index) as InputArray;
        oninput(arrayValue);
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            addItem();
        }
    }
</script>

<div class="flex items-center gap-2">
    <div class="flex flex-wrap gap-1 max-w-md">
        {#each arrayValue as item, index}
            <div class="inline-flex items-center bg-gray-100 rounded-full px-2 py-1 text-sm">
                <span class="truncate max-w-[100px]">{String(item)}</span>
                <button 
                    class="ml-1 text-gray-500 hover:text-red-500"
                    onclick={() => removeItem(index)}
                    aria-label="Remove item"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                    </svg>
                </button>
            </div>
        {/each}
    </div>
    
    <div class="flex-grow relative">
        {#if parameter.type === 'string'}
            <input 
                type="text" 
                bind:value={newItemValue}
                placeholder="Add new item"
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                onkeydown={handleKeyDown}
            />
        {:else if parameter.type === 'integer' || parameter.type === 'float'}
            <input 
                type="number" 
                bind:value={newItemValue}
                placeholder="Add new number"
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                onkeydown={handleKeyDown}
            />
        {:else if parameter.type === 'boolean'}
            <select 
                bind:value={newItemValue}
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
            >
                <option value="">Select value</option>
                <option value="true">True</option>
                <option value="false">False</option>
            </select>
        {:else if parameter.type === 'enum'}
            <select 
                bind:value={newItemValue}
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
            >
                <option value="">Select value</option>
                {#each parameter.values! as value}
                    <option value={value}>{value}</option>
                {/each}
            </select>
        {:else if parameter.type === 'datetime'}
            <input 
                type="datetime-local" 
                bind:value={newItemValue}
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                onkeydown={handleKeyDown}
            />
        {:else}
            <input 
                type="text" 
                bind:value={newItemValue}
                placeholder="Add new item"
                class="w-full px-3 py-1.5 pr-10 border border-gray-200 rounded-md shadow-sm focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                onkeydown={handleKeyDown}
            />
        {/if}
        <button 
            class="absolute right-2 top-1/2 transform -translate-y-1/2 text-indigo-600 hover:text-indigo-800"
            onclick={addItem}
            aria-label="Add item"
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
        </button>
    </div>
</div>



