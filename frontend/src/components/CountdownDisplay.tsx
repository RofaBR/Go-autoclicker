import { msToTime } from '../utils/time';

interface CountdownDisplayProps {
    countdownMs: number;
}

export const CountdownDisplay = ({ countdownMs }: CountdownDisplayProps) => {
    const { h, m, s, ms } = msToTime(countdownMs);

    const formatTime = (hours: number, minutes: number, seconds: number, milliseconds: number): string => {
        const msStr = milliseconds.toString().padStart(3, '0');
        if (hours > 0) {
            return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}.${msStr}`;
        }
        return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}.${msStr}`;
    };

    const isLowTime = countdownMs < 2000;

    return (
        <div className={`countdown-display ${isLowTime ? 'low-time' : ''}`}>
            <span className="countdown-label">Next click:</span>
            <span className="countdown-time">{formatTime(h, m, s, ms)}</span>
        </div>
    );
}