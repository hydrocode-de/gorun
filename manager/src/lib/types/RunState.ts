export interface RunState {
    id: number;
    name: string,
    title: string,
    description: string,
    dockerImage: string,
    mounts?: {
        [containerPath: string]: string
    },
    parameters?: {
        [name: string]: string | number | boolean | Date | Object
    },
    data?: {
        [name: string]: string
    },
    status: "pending" | "running" | "finished" | "errored",
    created_at: Date,
    started_at?: Date,
    finished_at?: Date,
    error_message?: string
}