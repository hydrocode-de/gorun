import { config } from "$lib/state.svelte";
import type { ResultFile } from "$lib/types/ResultFile";
import type { RunState } from "$lib/types/RunState";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, parent, fetch }): Promise<{run: RunState | undefined, files: ResultFile[]}> => {
    const parentData = await parent();
    const run = parentData.runs.find(run => run.id === Number(params.runId));

    let files: ResultFile[];
    if (run) {
        const res = await fetch(`${config.apiServer}/runs/${run.id}/results`, {
            headers: {
                'X-User-ID': config.auth.user.id
            }
        }).then(resp => {
            if (resp.ok) {
                return resp.json();
            } else {
                throw new Error(`Failed to fetch results for run ${run.id}`);
            }
        }).catch(error => {
            console.log(`Failed to fetch results for run ${run.id}`, error);
            return {
                files: []
            }
        });

        files = res.files;
    } else {
        files = [];
    }

    return {
        run,
        files 
    };
} 