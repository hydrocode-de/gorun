import type { CitationFile } from "./CitationFile";
import type { InputParams } from "./InputParams";

export interface ToolSpec {
    id: string;
    name: string;
    title: string;
    description: string;
    parameters?: Record<string, ParameterSpec>;
    data?: Record<string, DataSpec>;
    citation?: CitationFile;
}

export interface ParameterSpec {
    description?: string;
    type: string;
    array?: boolean;
    default?: InputParams;
    optional?: boolean;
    values?: string[];
    min?: number;
    max?: number;
}

export interface DataSpec {
    path: string;
    description?: string;
    example?: string;
    extension: string[];
}