import type { PageLoad } from "./$types";

export const load: PageLoad = ({ params, parent }) => {
    return parent().then(({ specs }) => {
        const spec = specs.find(spec => spec.id === params.specId);
        return { spec };
    });
};