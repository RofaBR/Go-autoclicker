export const timeToMs = (h: number, m: number, s: number, ms: number) => {
    return ((h * 3600000) + (m * 60000) + (s * 1000) + ms);
};

export const msToTime = (duration: number) => {
    const milliseconds = Math.floor(duration % 1000);
    const seconds = Math.floor((duration / 1000) % 60);
    const minutes = Math.floor((duration / (1000 * 60)) % 60);
    const hours = Math.floor((duration / (1000 * 60 * 60)));

    return { h: hours, m: minutes, s: seconds, ms: milliseconds };
}