<script lang="ts">
    import type { GeoJsonObject } from 'geojson';
    import L from 'leaflet';
    import { JSONEditor, Mode, type Content } from "svelte-jsoneditor";


    interface StructEditorProps {
        value: { [key: string]: any };
        oninput: (value: { [key: string]: any }) => void;
    }

    let { value, oninput }: StructEditorProps = $props();
    let mode: 'geojson' | 'editor' | 'dropzone' = $state('editor');
    let geojsonDetected = $state(false);
    
    // let jsonValue = $state(JSON.parse(JSON.stringify(value)));
    let geojsonValue: L.GeoJSON | null = $state(null);
    let map: L.Map | null = $state(null);

    //function handleJsonChange(newValue: {json: any | undefined, text: string | undefined}) {
    function handleJsonChange(newValue: Content) {
        if ('json' in newValue) {
            const newJson = newValue.json ? { ...newValue.json } : {};
            value = newJson;
            oninput(newJson);
        } else if (newValue.text) {
            const newJson = JSON.parse(newValue.text);
            value = newJson;
            oninput(newJson);
        }
    }

    function handleFileChange(event: Event) {
        const file = (event.target as HTMLInputElement).files?.[0]
        if (file) {
            readFile(file);
        }
    }

    function readFile(file: File) {
        const reader = new FileReader();
            reader.onload = e => {
                const content = e.target?.result as string;
                const json = JSON.parse(content);
                //jsonValue = json;
                oninput(json);
                mode = 'editor';
            }
            reader.readAsText(file);
    }

    function handleFileDrop(event: DragEvent) {
        console.log('handleFileDrop');
        event.preventDefault();
        const file = event.dataTransfer?.files[0];
        if (file) {
            readFile(file);
        }
    }

    $effect(() => {
        //console.log('value effect', value)
        if (value && value.type && (value.type === 'FeatureCollection' || value.type === 'Feature')) {
            geojsonDetected = true;
            geojsonValue = L.geoJson(value as GeoJsonObject);
        } else {
            geojsonDetected = false;
        }
    });

    function initMap(){
            if (!map) {
                map = L.map('map').setView([51.505, -0.09], 13);
                L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                }).addTo(map);
            }
            if (geojsonValue) {
                (geojsonValue as L.GeoJSON).addTo(map);
                map.fitBounds(geojsonValue.getBounds());
            }
    }

    function destroyMap(){
        if (map) {
            map.remove();
            map = null;
        }
    }

    $effect(() => {
        if (mode === 'geojson') {
            initMap();
        } else {
            destroyMap();
        }
    })
</script>

<svelte:head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
    integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
    crossorigin=""/>
</svelte:head>

<div class="w-full">
    <div class="flex border-b border-gray-200 mb-2">
        <button
            class="px-4 py-2 text-sm font-medium border-b-2 transition-colors duration-200"
            class:border-indigo-500={mode === 'editor'}
            class:text-indigo-600={mode === 'editor'}
            class:border-transparent={mode !== 'editor'}
            class:text-gray-500={mode !== 'editor'}
            class:hover:text-gray-700={mode !== 'editor'}
            class:hover:border-gray-300={mode !== 'editor'}
            onclick={() => mode = 'editor'}
        >
            JSON
        </button>
        <button
            class="px-4 py-2 text-sm font-medium border-b-2 transition-colors duration-200"
            class:border-indigo-500={mode === 'geojson' && geojsonDetected}
            class:text-indigo-600={mode === 'geojson' && geojsonDetected}
            class:border-transparent={mode !== 'geojson' || !geojsonDetected}
            class:text-gray-500={mode !== 'geojson' || !geojsonDetected}
            class:hover:text-gray-700={mode !== 'geojson' || !geojsonDetected}
            class:hover:border-gray-300={mode !== 'geojson' || !geojsonDetected}
            onclick={() => mode = 'geojson'}
            disabled={!geojsonDetected}
        >
            {geojsonDetected ? 'GeoJSON' : 'no GeoJSON detected'}
        </button>
        <!-- <button
            class="px-4 py-2 text-sm font-medium border-b-2 transition-colors duration-200"
            class:border-indigo-500={mode === 'dropzone'}
            class:text-indigo-600={mode === 'dropzone'}
            class:border-transparent={mode !== 'dropzone'}
            class:text-gray-500={mode !== 'dropzone'}
            class:hover:text-gray-700={mode !== 'dropzone'}
            class:hover:border-gray-300={mode !== 'dropzone'}
            onclick={() => mode = 'dropzone'}
        >
            Upload JSON
        </button> -->
    </div>

    {#if mode === 'editor'}
        <div class="w-full h-[300px] border border-gray-200 rounded-md shadow-sm focus-within:ring-1 focus-within:ring-indigo-500 focus-within:border-indigo-500">
            <JSONEditor
                content={{json: value}}
                onChange={handleJsonChange}
                statusBar={false}
                mode={Mode.text}
            />
        </div>
    {/if}

    {#if mode === 'geojson'}
        <div id="map" class="w-full h-[300px]"></div>
    {/if}

    {#if mode === 'dropzone'}
        <div class="w-full h-[300px] border border-gray-200 rounded-md shadow-sm focus-within:ring-1 focus-within:ring-indigo-500 focus-within:border-indigo-500"    >
            <div 
                aria-label="Dropzone for JSON file"
                role="button"
                tabindex="0"
                ondrop={handleFileDrop} 
                class="flex flex-col items-center justify-center h-full p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-indigo-500 transition-colors duration-200"
            >
                <svg class="w-12 h-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <button 
                    class="px-4 py-2 bg-white border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    onclick={() => document.getElementById('fileInput')?.click()}
                >
                    Drop JSON file here or click to upload
                </button>
            </div>
            <input type="file" id="fileInput" accept=".json" class="hidden" onchange={handleFileChange} />
        </div>
    {/if}
</div>