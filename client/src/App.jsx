import { useState, useRef } from 'react';
import toast from 'react-hot-toast';

import './App.css';

const App = () => {
    const [isRecording, setIsRecording] = useState(false);
    const [blob, setBlob] = useState(null);
    const [blobURL, setBlobURL] = useState(null);
    const mediaRecorderRef = useRef(null);
    const audioChunks = useRef([]);
    const [token, setToken] = useState(null);

    // Start recording
    const startRecording = async () => {
        setIsRecording(true);
        audioChunks.current = [];

        // Request permission to access the microphone
        const stream = await navigator.mediaDevices.getUserMedia({
            audio: true,
        });

        // Create a new MediaRecorder instance
        mediaRecorderRef.current = new MediaRecorder(stream);

        // Handle the dataavailable event to collect audio chunks
        mediaRecorderRef.current.addEventListener('dataavailable', (event) => {
            audioChunks.current.push(event.data);
        });

        // Handle the stop event to create a blob from the recorded audio
        mediaRecorderRef.current.addEventListener('stop', () => {
            const audioBlob = new Blob(audioChunks.current, {
                type: 'audio/wav',
            });
            setBlob(audioBlob);
            const audioUrl = URL.createObjectURL(audioBlob);
            setBlobURL(audioUrl);
        });

        // Start the recording
        mediaRecorderRef.current.start();
    };

    // Stop recording
    const stopRecording = () => {
        setIsRecording(false);
        mediaRecorderRef.current.stop();
    };

    // Function to send the recorded audio blob to the API
    const sendVoiceData = async (blob, endpoint) => {
        const formData = new FormData();
        formData.append('voicedata', blob, 'voice-input.wav');

        try {
            const response = await fetch(`http://localhost:3000/${endpoint}`, {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                if (endpoint === 'signup') toast.success('Signed up!');
                if (endpoint === 'signin') {
                    const { accesstoken } = await response.json();
                    setToken(accesstoken);
                    toast.success('Signed in!');
                }
            } else {
                if (endpoint === 'signin') toast.error('Not authorized');
                if (endpoint === 'signup') toast.error('Something went wrong');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div className="voice-signup-container">
            <h2>Canary Demo</h2>
            <div className="button-container">
                <button onClick={startRecording} disabled={isRecording}>
                    Start Recording
                </button>
                <button onClick={stopRecording} disabled={!isRecording}>
                    Stop Recording
                </button>
            </div>
            {<audio src={blobURL} controls />}
            <button onClick={() => sendVoiceData(blob, 'signup')}>
                Sign Up
            </button>
            <button onClick={() => sendVoiceData(blob, 'signin')}>
                Sign In
            </button>
        </div>
    );
};

export default App;
