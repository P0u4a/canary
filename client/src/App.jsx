import { useState, useRef, useEffect } from 'react';
import './App.css';

function App() {
    const [micAccessGranted, setMicAccessGranted] = useState(false);

    useEffect(() => {
        const getMicAccess = async () => {
            const micAccess = await navigator.permissions.query({
                name: 'microphone',
            });
            if (micAccess.state === 'granted') {
                setMicAccessGranted(true);
            }
        };

        getMicAccess();
    }, []);

    const [loginType, setLoginType] = useState('signup');
    const [isRecording, setIsRecording] = useState(false);
    const [recordedChunks, setRecordedChunks] = useState([]);
    const mediaRecorderRef = useRef(null);
    const [username, setUsername] = useState('');
    const [authStatus, setAuthStatus] = useState(false);

    const handleStartRecording = async () => {
        if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
            try {
                const stream = await navigator.mediaDevices.getUserMedia({
                    audio: true,
                });
                mediaRecorderRef.current = new MediaRecorder(stream);

                mediaRecorderRef.current.ondataavailable = (event) => {
                    if (event.data.size > 0) {
                        setRecordedChunks([event.data]);
                    }
                };

                mediaRecorderRef.current.start();
                setIsRecording(true);
            } catch (error) {
                console.error('Error accessing microphone:', error);
            }
        }
    };

    const handleStopRecording = () => {
        if (mediaRecorderRef.current) {
            mediaRecorderRef.current.stop();
            setIsRecording(false);
        }
    };

    const handlePlayRecording = () => {
        const blob = new Blob(recordedChunks, { type: 'audio/wav' });
        const url = URL.createObjectURL(blob);

        const audio = new Audio(url);
        audio.play();
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const blob = new Blob(recordedChunks, { type: 'audio/wav' });

        const formData = new FormData();
        formData.append('voiceData', blob, 'audio.wav');
        formData.append('username', username);

        const endpoint = `http://localhost:3000/${loginType}`;
        const res = await fetch(endpoint, {
            method: 'POST',
            body: formData,
        });

        if (res.status != 200) {
            alert('Something went wrong.');
            return;
        }

        const tokens = await res.json();
        // TODO store the tokens somewhere safe
        console.log(tokens);

        setAuthStatus(true);
    };

    return (
        <>
            <h1>Canary Auth Demo</h1>
            <p style={{ color: 'red' }}>
                {micAccessGranted === true ? '' : 'Please enable your mic'}
            </p>
            <form className="auth-form" onSubmit={handleSubmit}>
                <input
                    type="text"
                    value={username}
                    placeholder="Your username"
                    onChange={(e) => setUsername(e.currentTarget.value)}
                />
                <div className="audio-controls">
                    <button
                        type="button"
                        onClick={handleStartRecording}
                        disabled={isRecording}
                    >
                        Start
                    </button>
                    <button
                        type="button"
                        onClick={handleStopRecording}
                        disabled={!isRecording}
                    >
                        Stop
                    </button>
                    <button
                        type="button"
                        onClick={handlePlayRecording}
                        disabled={recordedChunks.length === 0}
                    >
                        Play
                    </button>
                </div>
                <div className="btn-group">
                    <button
                        style={
                            recordedChunks.length === 0
                                ? { pointerEvents: 'none' }
                                : { pointerEvents: 'auto' }
                        }
                        disabled={recordedChunks.length === 0}
                        type="submit"
                        onClick={() => setLoginType('signin')}
                    >
                        Sign In
                    </button>
                    <button
                        style={
                            recordedChunks.length === 0
                                ? { pointerEvents: 'none' }
                                : { pointerEvents: 'auto' }
                        }
                        disabled={recordedChunks.length === 0}
                        type="submit"
                        onClick={() => setLoginType('signup')}
                    >
                        Sign Up
                    </button>
                </div>
            </form>

            <div className="test">
                Status: {authStatus ? 'Signed in' : 'Not signed in'}
                <a href="http://localhost:3000/protected">Protected route</a>
            </div>
        </>
    );
}

export default App;
