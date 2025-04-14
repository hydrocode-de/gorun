<script lang="ts">
    import type { CitationFile } from "$lib/types/CitationFile";
    import hljs from "highlight.js/lib/core";
    import json from "highlight.js/lib/languages/json";
    import "highlight.js/styles/atom-one-light.css";
    import { onMount } from "svelte";

    onMount(() => {
        hljs.registerLanguage("json", json);
        hljs.highlightAll();
    })

    interface CitationInfoProps {
        citation: CitationFile
    }
    let { citation }: CitationInfoProps = $props();
    $inspect(citation);

    let bibtexCode = $derived(`@software{${citation.Title.toLowerCase().replace(/[^a-z0-9]/g, '_')},
    author = {${citation.Authors?.map(author => 
        author.IsPerson ? 
            `${author.Person?.Family}, ${author.Person?.GivenNames}` : 
            author.Entity?.Name
    ).join(' and ')}},
    title = {${citation.Title}},
    version = {${citation.Version}},
    year = {${new Date().getFullYear()}},
    url = {${citation.RepositoryCode?.Host ? `https://${citation.RepositoryCode.Host}${citation.RepositoryCode.Path}` : ''}},
    keywords = {${citation.Keywords?.join(', ')}}
}`);

    let apaCitation = $derived(`${citation.Authors?.map(author => 
        author.IsPerson ? 
            `${author.Person?.Family}, ${author.Person?.GivenNames?.charAt(0)}.` : 
            author.Entity?.Name
    ).join(', ')} (${new Date().getFullYear()}). ${citation.Title} (Version ${citation.Version}) [Computer software]. ${citation.RepositoryCode?.Host ? `https://${citation.RepositoryCode.Host}${citation.RepositoryCode.Path}` : ''}`);
</script>

<div class="p-3 rounded-lg border border-gray-200 shadow-md mb-6">
    <h2 class="text-lg font-semibold text-gray-900 mb-3">{citation.Title}</h2>

    <p class="mt-2 text-gray-600">
        This section provides information on how to cite the tool.
    </p>

    {#if citation.Authors && citation.Authors.length > 0}
        <h4 class="mt-6 font-semibold text-gray-900">Authors</h4>
        <ul class="mt-2 space-y-4">
            {#each citation.Authors as author}
                {#if author.IsPerson && author.Person}
                    <li class="text-sm">
                        <div class="font-medium text-gray-900">
                            {author.Person.GivenNames} {author.Person.Family}
                        </div>
                        {#if author.Person.Affiliation}
                            <div class="text-gray-600">
                                {author.Person.Affiliation}
                            </div>
                        {/if}
                        {#if author.Person.Orcid}
                            <div class="mt-1">
                                <a 
                                    href={author.Person.Orcid} 
                                    target="_blank" 
                                    rel="noopener noreferrer"
                                    class="text-indigo-600 hover:text-indigo-800 text-xs flex items-center gap-1"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M12 0C5.372 0 0 5.372 0 12s5.372 12 12 12 12-5.372 12-12S18.628 0 12 0zm-2.092 17.264c-1.191 0-2.232-.532-2.232-1.247 0-.694.961-1.264 2.232-1.264s2.231.57 2.231 1.264c0 .715-1.04 1.247-2.231 1.247zm0-10.264c-1.191 0-2.232-.532-2.232-1.247 0-.694.961-1.264 2.232-1.264s2.231.57 2.231 1.264c0 .715-1.04 1.247-2.231 1.247zm4.184 5.264c-1.191 0-2.232-.532-2.232-1.247 0-.694.961-1.264 2.232-1.264s2.231.57 2.231 1.264c0 .715-1.04 1.247-2.231 1.247z"/>
                                    </svg>
                                    ORCID
                                </a>
                            </div>
                        {/if}
                    </li>
                {/if}
            {/each}
        </ul>
    {/if}

    {#if citation.Keywords && citation.Keywords.length > 0}
        <h4 class="mt-6 font-semibold text-gray-900">Keywords</h4>
        <div class="mt-2 flex flex-wrap gap-1.5">
            {#each citation.Keywords as keyword}
                <span class="px-1.5 py-0.5 bg-indigo-100 text-indigo-800 text-xs font-medium rounded-full">
                    {keyword}
                </span>
            {/each}
        </div>
    {/if}

    <h4 class="mt-6 font-semibold text-gray-900">APA</h4>
    <div class="mt-2 p-2 shadow-md relative">
        <pre class="text-wrap"><code class="language-json">{apaCitation}</code></pre>
        <button 
            aria-label="Copy to clipboard"
            class="absolute top-3 right-3 p-1 text-gray-500 hover:text-gray-400 transition-colors z-10 w-8 h-8 flex items-center justify-center border border-gray-400 hover:border-gray-300 hover:cursor-pointer rounded"
            onclick={() => navigator.clipboard.writeText(apaCitation)}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" />
                <path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z" />
            </svg>
        </button>
    </div>

    <h4 class="mt-6 font-semibold text-gray-900">BibTeX</h4>
    <div class="mt-2 p-2 shadow-md relative">
        <pre class="text-wrap"><code class="language-json">{bibtexCode}</code></pre>
        <button 
            aria-label="Copy to clipboard"
            class="absolute top-3 right-3 p-1 text-gray-500 hover:text-gray-400 transition-colors z-10 w-8 h-8 flex items-center justify-center border border-gray-400 hover:border-gray-300 hover:cursor-pointer rounded"
            onclick={() => navigator.clipboard.writeText(bibtexCode)}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" />
                <path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z" />
            </svg>
        </button>
    </div>
</div>