export interface ResultFile {
    name: string;
    relPath: string;
    absPath: string;
    size: number;
    lastModified?: Date;
}