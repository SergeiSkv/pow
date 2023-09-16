# Proof-of-Work


### PoW Algorithm Choice

I've chosen the simple SHA-256 based hashing algorithm for the PoW. Clients have to find a nonce such that the SHA-256 hash of the nonce concatenated with a provided challenge string starts with a certain prefix (e.g., "0000"). This is computationally hard to find but easy to verify, making it suitable for a PoW system.
