import { useState, useEffect } from 'react';
import { timeToMs } from '../utils/time';

interface DelayControlProps {
    delayMs: number;
    onChange: (delayMs: number) => void;
}

interface TimeParts {
    h: number;
    m: number;
    s: number;
    ms: number;
}

export const DelayControl = ({ delayMs, onChange }: DelayControlProps) => {
    const [parts, setParts] = useState<TimeParts>({ h: 0, m: 0, s: 0, ms: 0 });

    useEffect(() => {
        const h = Math.floor(delayMs / 3600000);
        const m = Math.floor((delayMs % 3600000) / 60000);
        const s = Math.floor((delayMs % 60000) / 1000);
        const ms = delayMs % 1000;

        setParts({ h, m, s, ms });
    }, [delayMs]);

    const timeFields: { label: string; key: keyof TimeParts; max?: number }[] = [
        { label: 'H', key: 'h' },
        { label: 'M', key: 'm', max: 59 },
        { label: 'S', key: 's', max: 59 },
        { label: 'MS', key: 'ms', max: 999 }
    ];

    const handleFieldChange = (key: keyof TimeParts, valueStr: string, max?: number) => {
        let value = parseInt(valueStr) || 0;

        if (value < 0) value = 0;
        if (max !== undefined && value > max) value = max;

        const newParts = { ...parts, [key]: value };

        setParts(newParts);

        const totalMs = timeToMs(newParts.h, newParts.m, newParts.s, newParts.ms);
        onChange(totalMs);
    };

    return (
        <div className="delay-control">
            {timeFields.map((field, index) => (
                <div key={field.key} style={{ display: 'flex', alignItems: 'center' }}>
                    <div className="input-group time-input">
                        <label>{field.label}</label>
                        <input
                            type="number"
                            value={parts[field.key]}
                            onChange={(e) => handleFieldChange(field.key, e.target.value, field.max)}
                            min="0"
                            max={field.max}
                        />
                    </div>
                    {index < timeFields.length - 1 && (
                        <span className="time-separator">:</span>
                    )}
                </div>
            ))}
        </div>
    );
};