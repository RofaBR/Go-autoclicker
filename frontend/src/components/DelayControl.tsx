import { useState, useEffect } from 'react';
import { timeToMs } from '../utils/time';

interface DelayControlProps {
    delayMs: number;
    onChange: (delayMs: number) => void;
}

export const DelayControl = ({ delayMs, onChange }: DelayControlProps) => {
    const [hours, setHours] = useState<number>(0);
    const [minutes, setMinutes] = useState<number>(0);
    const [seconds, setSeconds] = useState<number>(0);

    useEffect(() => {
        const totalSeconds = Math.floor(delayMs / 1000);
        const h = Math.floor(totalSeconds / 3600);
        const m = Math.floor((totalSeconds % 3600) / 60);
        const s = totalSeconds % 60;

        setHours(h);
        setMinutes(m);
        setSeconds(s);
    }, [delayMs]);

    const handleChange = (h: number, m: number, s: number) => {
        const newDelayMs = timeToMs(h, m, s);
        onChange(newDelayMs);
    };

    const handleHoursChange = (value: string) => {
        const h = Math.max(0, parseInt(value) || 0);
        setHours(h);
        handleChange(h, minutes, seconds);
    };

    const handleMinutesChange = (value: string) => {
        const m = Math.max(0, Math.min(59, parseInt(value) || 0));
        setMinutes(m);
        handleChange(hours, m, seconds);
    };

    const handleSecondsChange = (value: string) => {
        const s = Math.max(0, Math.min(59, parseInt(value) || 0));
        setSeconds(s);
        handleChange(hours, minutes, s);
    };

    return (
        <div className="delay-control">
            <div className="input-group time-input">
                <label>H</label>
                <input
                    type="number"
                    value={hours}
                    onChange={(e) => handleHoursChange(e.target.value)}
                    min="0"
                />
            </div>
            <span className="time-separator">:</span>
            <div className="input-group time-input">
                <label>M</label>
                <input
                    type="number"
                    value={minutes}
                    onChange={(e) => handleMinutesChange(e.target.value)}
                    min="0"
                    max="59"
                />
            </div>
            <span className="time-separator">:</span>
            <div className="input-group time-input">
                <label>S</label>
                <input
                    type="number"
                    value={seconds}
                    onChange={(e) => handleSecondsChange(e.target.value)}
                    min="0"
                    max="59"
                />
            </div>
        </div>
    );
};