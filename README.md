# Canary Auth

Biometric authentication protocol using voice samples.

> 🚧 Please do not use this for anything real. It's just a fun experiment.

## Design

The auth flow uses JWT, with the access and refresh token flow.

On Sign up, the user's voice MFCC features are extracted, and their passphrase is transcribed and hashed.

On Sign in, the MFCC features from the user's input voice is compared against their saved features, if the similarity is greater than 98% (accounting for model errors so not 100%) and the transcribed passphrase hashes match, then a acess, refresh token pair is generated and sent to the client, authenticating the user.

Upon subsequent requests, the client will sent the access token in the Authorization header of their API requests to the server. If the access token is expired, the client will have to sent a request to the server at the `/refresh` route, sending the refresh token in the Authorization header of this request. The server will then validate this refresh token and send a new access token. If the refresh token is also expired, then the user is required to log in again to receive a new token pair.

## Limitations

-   You'll have to say your passphrase out loud, problematic when in public
-   Probably will not work accuratley in noisy environments
-   Voice data is not currrently encrypted when in transit from client to server, though it could be
-   It's possible to spoof the user's voice (relatively easy these days) and then brute force the passphrase
