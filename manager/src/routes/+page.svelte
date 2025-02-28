<script lang="ts">
    import { config } from './state.svelte.js';
    
    let apiStatus: boolean | null = null; 
    
    async function testApiConnection() {
        try {
            const response = await fetch(`${config.apiServer}/health`);
            const text = await response.text();
            apiStatus = text.trim() === 'OK';
        } catch (error) {
            apiStatus = false;
        }
    }
</script>

<div class="max-w-3xl mx-auto pt-8">
    <h1 class="text-3xl font-bold mb-6">Get Started</h1>
    
    <div class="mb-8">
        <label for="apiServer" class="block text-sm font-medium text-gray-700 mb-2">
            GoRun API endpoint
        </label>
        <div class="relative flex items-center">
            <input 
                type="text" 
                id="apiServer"
                bind:value={config.apiServer}
                class="flex-1 px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors text-lg pr-24"
                placeholder="Enter API server URL..."
            />
            <button 
                on:click={testApiConnection}
                class="absolute right-12 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors text-sm"
            >
                Test
            </button>
            <div class="absolute right-4 w-4 h-4 rounded-full transition-colors"
                class:bg-gray-300={apiStatus === null}
                class:bg-green-500={apiStatus === true}
                class:bg-red-500={apiStatus === false}
            ></div>
        </div>
        <p class="mt-2 text-sm text-gray-500">
            Status: {apiStatus === null ? 'Not tested' : apiStatus ? 'Connected' : 'Connection failed'}
        </p>
    </div>

    <div class="mt-12">
        <h2 class="text-2xl font-bold mb-4">HowTo</h2>
        <p class="text-gray-700">This section will explain everything at one point</p>
    </div>
</div>



