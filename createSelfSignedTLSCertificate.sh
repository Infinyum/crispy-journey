# -x509 for self-signed
# -nodes (no DES) to disable password protection for the new RSA key
# -subj to avoid interactive mode and specify the common name (CN) right away
openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout tls.key -out tls.crt -subj "/CN=crispy.local"