import librosa
import numpy as np
from sklearn import mixture

def extract_features(audio_file, fs=16000, wst=0.02, fpt=0.01, nbands=40, ncomp=20):
    # y: Speech signal, sr: Sampling rate
    y, sr = librosa.load(audio_file, sr=fs)

    # Window size in samples
    nfft = int(wst * fs) 
    # Frame period in samples
    fp = int(fpt * fs)    

    # Extract MFCC features
    mfcc = librosa.feature.mfcc(y=y, sr=sr, n_fft=nfft, hop_length=fp, n_mels=nbands, n_mfcc=ncomp).T

    return mfcc


def train_model(audio_file):
    mfcc = extract_features(audio_file)
    gmm = mixture.GaussianMixture(n_components=8, covariance_type='diag', random_state=1)
    gmm.fit(mfcc)
    return gmm


def get_similarity(model, audio_file):
    mfcc = extract_features(audio_file)
    probas = model.score_samples(mfcc)
    return np.sum(probas)
