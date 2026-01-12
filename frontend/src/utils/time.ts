export const timeToMs = (h: number, m: number, s: number) => {
    return ((h * 3600000) + (m * 60000) + (s * 1000))
};

export const msToTime = (duration: number) => {
    const seconds = Math.floor((duration / 1000) % 60);
    const minutes = Math.floor((duration / (1000 * 60)) % 60);
    const hours = Math.floor((duration / (1000 * 60 * 60)));

    return { h: hours, m: minutes, s: seconds};
};