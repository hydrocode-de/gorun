import { config } from "$lib/state.svelte";
import type { ToolSpec } from "$lib/types/ToolSpec";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch }): Promise<{specs: ToolSpec[]}> => {
    const res = await fetch(`${config.apiServer}/specs`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });
    const data = await res.json();
    const specs = data.tools as ToolSpec[];

    return { specs };
}; 