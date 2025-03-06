export const bytesToSize = (bytes: number): string => {
    let short = bytes;
    let suffix = 'B';
    if (short > 1024) {
        short = short / 1024;
        suffix = 'KB';
    }
    if (short > 1024) {
        short = short / 1024;
        suffix = 'MB';
    }
    if (short > 1024) {
        short = short / 1024;
        suffix = 'GB';
    }

    return `${short.toFixed(1)} ${suffix}`
}