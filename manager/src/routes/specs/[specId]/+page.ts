import { tools } from "../../state.svelte";
import type { PageLoad } from "./$types";

export const load: PageLoad = ({ params }) => {
    const spec = tools.specs.find(spec => spec.id === params.specId);
    return spec;
}