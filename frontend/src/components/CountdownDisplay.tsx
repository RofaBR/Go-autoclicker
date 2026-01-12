import { msToTime } from '../utils/time';

interface CountdownDisplayProps {
    countdownMs: number;
}

export const CountdownDisplay = ({ countdownMs }: CountdownDisplayProps) => {
    const { h, m, s } = msToTime(countdownMs);

    const formatTime = (hours: number, minutes: number, seconds: number): string => {
        if (hours > 0) {
            return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
        }
        return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    };

    const isLowTime = countdownMs < 2000;

    return (
        <div className={`countdown-display ${isLowTime ? 'low-time' : ''}`}>
            <span className="countdown-label">Next click:</span>
            <span className="countdown-time">{formatTime(h, m, s)}</span>
        </div>
    );
}