import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, parent }) => {
    const parentData = await parent();
    const run = parentData.runs.find(run => run.id === Number(params.runId));

    return run ;
} 