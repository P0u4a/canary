from flask import Flask, request, jsonify
import canary
import io

app = Flask(__name__)

# TODO Replace with a thread-safe data structure
user_to_model = {}

@app.route('/process-voice', methods=['POST'])
def process_voice():
    voicedata = request.json.get('voicedata')
    username = request.json.get('username')
    new_user_model = canary.train_model(voicedata)

    user_to_model[username] = new_user_model

    # TODO add error handling
    return jsonify({'message': 'User model initialised'}), 200


@app.route('/analyse-voice', methods=['POST'])
def analyse_voice():
    voicedata = request.json.get('voicedata')
    username = request.json.get('username')

    user_model = user_to_model[username]

    similarity_score = canary.get_similarity(user_model, voicedata)

    status = similarity_score < -16000

    return jsonify({'verified': status}), 200


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=3001)