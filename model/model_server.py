from flask import Flask, request, jsonify
import threading
import canary
from utils import cleanup



# -16000 is the value arrived at by doing some tests against positive and negative voice matches
# In future iterations a more robust calculation should be used
SIMILARITY_THRESHOLD = -16000
TEMP_FILE_NAME = 'audio.wav'
SESSION_TYPE = 'filesystem'

app = Flask(__name__)
app.config.from_object(__name__)

user_to_model = {}

@app.route('/init-model', methods=['POST'])
def process_voice():
    voice_data = request.files['voicedata']
    username = request.form.get('username')

    audio_bytes = voice_data.read()
    with open(TEMP_FILE_NAME, 'wb') as f:
        f.write(audio_bytes)

    try:
        new_user_model = canary.train_model(TEMP_FILE_NAME)
    except Exception as e:
        threading.Thread(target=cleanup, args=(TEMP_FILE_NAME,)).start()
        return jsonify({'message': e}), 500
    
    user_to_model[username] = new_user_model

    threading.Thread(target=cleanup, args=(TEMP_FILE_NAME,)).start()

    return jsonify({'message': 'initialised user model'}), 200


@app.route('/verify-voice', methods=['POST'])
def analyse_voice():
    voice_data = request.files['voicedata']
    username = request.form.get('username')

    user_model = user_to_model[username]

    audio_bytes = voice_data.read()
    with open(TEMP_FILE_NAME, 'wb') as f:
        f.write(audio_bytes)

    try:
        similarity_score = canary.get_similarity(user_model, TEMP_FILE_NAME)
    except Exception as e:
        threading.Thread(target=cleanup, args=(TEMP_FILE_NAME,)).start()
        return jsonify({'message': e}), 500

    threading.Thread(target=cleanup, args=(TEMP_FILE_NAME,)).start()

    status = 1 if similarity_score < SIMILARITY_THRESHOLD else 0

    return jsonify({'verified': status}), 200


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=3001)
