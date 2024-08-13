# Canary Auth 🔊

Biometric authentication protocol using your voice.

> 🚧 Not meant for use in production

## Design

The auth flow uses JWT, with the [access and refresh token pattern](https://www.baeldung.com/cs/access-refresh-tokens).

On sign up, the user's voice is passed to the `canary` model, which extracts the [MFCC features](https://en.wikipedia.org/wiki/Mel-frequency_cepstrum) from the audio file and trains a [Gaussian Mixture Model](https://scikit-learn.org/stable/modules/mixture.html) (GMM) on them.
This model is then linked to the username of the client in a dictionary on the `canary` server. Each user has their own Gaussian model that is responsible for recognising their voice.

On sign in, the user's GMM is fetched from the dictionary using their username as the key. The MFCC features of the voice used to sign in are extracted, and passed to the model that then generates a similarity score between its training data and the input data.
The score is then tested against a reasonable threshold value, and if the score is less than this threshold, the user is considered authenticated. An access and refresh token pair is then sent to the client.

Upon subsequent requests, the client will sent the access token in the Authorization header of their API requests to the server. If the access token is expired, the client will have to sent a request to the server at the `/refresh` route, sending the refresh token in the Authorization header of this request. The server will then validate this refresh token and send a new access token. If the refresh token is also expired, then the user is required to log in again to receive a new token pair.

## Limitations

-   You'll have to say your passphrase out loud, problematic when in public
-   Probably will not work accuratley in noisy environments
-   Voice data is not currrently encrypted when in transit from client to server
-   It's possible to spoof the user's voice or get it close enough (relatively easy these days)
