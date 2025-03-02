import { config } from "$lib/state.svelte";
import type { RunState } from "$lib/types/RunState";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ url, fetch }): Promise<{runs: RunState[]}> => {
    let status = url.searchParams.get('status');
    if (status === 'all') {
        status = '';
    }
    let backendUrl = `${config.apiServer}/runs`
    if (status !== '') {
        backendUrl += `?status=${status}`
    }
    console.log(backendUrl);
    const res = await fetch(backendUrl, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    const data = await res.json()
    const runs = data.runs as RunState[]

    return { runs };
}