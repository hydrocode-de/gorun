export interface ToolSpec {
    id: string;
    name: string;
    title: string;
    description: string;
    parameters?: Record<string, ParameterSpec>;
    data?: Record<string, DataSpec>;
}

export interface ParameterSpec {
    description?: string;
    type: string;
    array?: boolean;
    default?: boolean;
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