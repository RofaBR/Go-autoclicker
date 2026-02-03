import {useState, useEffect, useRef} from 'react';
import './App.css';
import {AddPoint, GetPoints, RecordCoordinates, UpdatePointDelay, RemovePoint, StartClicker, StopClicker, SetClickerMode} from "../wailsjs/go/main/App";
import { DelayControl } from './components/DelayControl';
import { CountdownDisplay } from './components/CountdownDisplay';
import { EventsOn } from '../wailsjs/runtime/runtime';

interface ClickPoint {
    ID: number;
    X: number;
    Y: number;
    Delay: number;
}

function App() {
    const [points, setPoints] = useState<ClickPoint[]>([]);
    const [isRecording, setIsRecording] = useState<number | null>(null);
    const [isRunning, setIsRunning] = useState<boolean>(false);
    const [countdowns, setCountdowns] = useState<Record<number, number>>({});
    const [isSequential, setIsSequential] = useState<boolean>(false);
    const [currentPointIndex, setCurrentPointIndex] = useState<number>(0);
    const sequentialIndexRef = useRef<number>(0);
    const sequentialCountdownRef = useRef<number>(0);

    useEffect(() => {
        loadPoints();
        const unsubscribe = EventsOn('clicker:stopped', () => {
            setIsRunning(false);
        });

        return () => {
            if (unsubscribe) unsubscribe();
        };
    }, []);

    useEffect(() => {
        if (isRunning && points.length > 0) {
            if (isSequential) {
                sequentialIndexRef.current = 0;
                sequentialCountdownRef.current = points[0].Delay;
                setCurrentPointIndex(0);
                setCountdowns({ [points[0].ID]: points[0].Delay });

                const interval = setInterval(() => {
                    sequentialCountdownRef.current -= 100;

                    if (sequentialCountdownRef.current <= 0) {
                        sequentialIndexRef.current = (sequentialIndexRef.current + 1) % points.length;
                        const nextPoint = points[sequentialIndexRef.current];
                        sequentialCountdownRef.current = nextPoint.Delay;
                        setCurrentPointIndex(sequentialIndexRef.current);
                        setCountdowns({ [nextPoint.ID]: nextPoint.Delay });
                    } else {
                        const currentPoint = points[sequentialIndexRef.current];
                        setCountdowns({ [currentPoint.ID]: sequentialCountdownRef.current });
                    }
                }, 100);

                return () => clearInterval(interval);
            } else {
                const initialCountdowns: Record<number, number> = {};
                points.forEach(point => {
                    initialCountdowns[point.ID] = point.Delay;
                });
                setCountdowns(initialCountdowns);

                const interval = setInterval(() => {
                    setCountdowns(prevCountdowns => {
                        const newCountdowns = { ...prevCountdowns };
                        points.forEach(point => {
                            if (newCountdowns[point.ID] !== undefined) {
                                newCountdowns[point.ID] -= 100;
                                if (newCountdowns[point.ID] <= 0) {
                                    newCountdowns[point.ID] = point.Delay;
                                }
                            }
                        });
                        return newCountdowns;
                    });
                }, 100);

                return () => clearInterval(interval);
            }
        } else {
            setCountdowns({});
            setCurrentPointIndex(0);
        }
    }, [isRunning, points, isSequential]);

    const loadPoints = async () => {
        try {
            const loadedPoints = await GetPoints();
            setPoints(loadedPoints || []);
        } catch (error) {
            console.error("Failed to load points:", error);
        }
    };

    const handleAddPoint = async () => {
        try {
            await AddPoint();
            await loadPoints();
        } catch (error) {
            console.error("Failed to add point:", error);
        }
    };

    const handleRecordCoordinates = async (id: number) => {
        setIsRecording(id);
        try {
            await RecordCoordinates(id);
            await loadPoints();
        } catch (error) {
            console.error("Failed to record coordinates:", error);
        } finally {
            setIsRecording(null);
        }
    };

    const handleDelayChange = async (id: number, delay: number) => {
        try {
            await UpdatePointDelay(id, delay);
            await loadPoints();
        } catch (error) {
            console.error("Failed to update delay:", error);
        }
    };

    const handleRemovePoint = async (id: number) => {
        try {
            await RemovePoint(id);
            setPoints(points.filter(p => p.ID !== id));
        } catch (error) {
            console.error("Failed to remove point:", error);
        }
    };

    const handleModeChange = async (sequential: boolean) => {
        try {
            await SetClickerMode(sequential);
            setIsSequential(sequential);
        } catch (error) {
            console.error("Failed to set mode:", error);
        }
    };

    const handleStart = async () => {
        try {
            await StartClicker();
            setIsRunning(true);
        } catch (error) {
            console.error("Failed to start clicker:", error);
        }
    };

    const handleStop = async () => {
        try {
            await StopClicker();
            setIsRunning(false);
        } catch (error) {
            console.error("Failed to stop clicker:", error);
        }
    };

    return (
        <div className="container">
            <header className="app-header">
                <h1>Auto Clicker</h1>
                <div className="hotkey-hint">
                    <span>Emergency Stop:</span>
                    <kbd>F10</kbd>
                </div>
                <div className="header-actions">
                    <div className="mode-switcher">
                        <button
                            className={`mode-btn ${!isSequential ? 'active' : ''}`}
                            onClick={() => handleModeChange(false)}
                            disabled={isRunning}
                        >
                            Parallel
                        </button>
                        <button
                            className={`mode-btn ${isSequential ? 'active' : ''}`}
                            onClick={() => handleModeChange(true)}
                            disabled={isRunning}
                        >
                            Sequential
                        </button>
                    </div>
                    <button className="btn btn-primary" onClick={handleAddPoint}>
                        + Add Point
                    </button>
                    {isRunning ? (
                        <button className="btn btn-stop" onClick={handleStop}>
                            Stop
                        </button>
                    ) : (
                        <button
                            className="btn btn-start"
                            onClick={handleStart}
                            disabled={points.length === 0}
                        >
                            Start
                        </button>
                    )}
                </div>
            </header>

            <div className="points-list">
                {points.length === 0 ? (
                    <div className="empty-state">
                        <p>No points created yet.</p>
                        <small>Click the button above to start.</small>
                    </div>
                ) : (
                    points.map((point, index) => (
                        <div key={point.ID} className={`point-card ${isRunning && isSequential && index === currentPointIndex ? 'active' : ''}`}>

                            <div className="point-info">
                                <span className="point-id">#{point.ID}</span>

                                <div className="input-group">
                                    <label>X</label>
                                    <input
                                        type="number"
                                        value={point.X}
                                        readOnly
                                    />
                                </div>

                                <div className="input-group">
                                    <label>Y</label>
                                    <input
                                        type="number"
                                        value={point.Y}
                                        readOnly
                                    />
                                </div>
                            </div>

                            <div className="point-actions">
                                <button
                                    className={`btn btn-record ${isRecording === point.ID ? 'recording' : ''}`}
                                    onClick={() => handleRecordCoordinates(point.ID)}
                                    disabled={isRecording !== null && isRecording !== point.ID}
                                >
                                    {isRecording === point.ID ? 'Click anywhere...' : 'Set Coords'}
                                </button>

                                <DelayControl
                                    delayMs={point.Delay}
                                    onChange={(newDelay) => handleDelayChange(point.ID, newDelay)}
                                />

                                {isRunning && countdowns[point.ID] !== undefined && (
                                    <CountdownDisplay countdownMs={countdowns[point.ID]} />
                                )}
                                {!isRunning && (
                                    <button
                                        className="btn-icon-delete"
                                        onClick={() => handleRemovePoint(point.ID)}
                                        title="Remove point"
                                    >
                                        ✕
                                    </button>
                                )}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    )
}

export default App
