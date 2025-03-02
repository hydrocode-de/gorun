import type { ToolSpec } from "$lib/types/ToolSpec";

// Update your state definition
export const tools = $state({
    specs: [] as ToolSpec[],
    lastUpdated: null as Date | null, 
    count: 0,
});
export const specs: ToolSpec[] = [];