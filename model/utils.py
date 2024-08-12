import os

def cleanup(path):
    try:
        os.remove(path)
        return None
    except Exception as e:
        return e
    