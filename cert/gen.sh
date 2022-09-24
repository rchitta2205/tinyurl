rm *.pem

# 1. Generate CA's private key and self-signed certificate.
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=CA/ST=Ontario/L=Cambridge/O=Teleport/OU=Assessment/CN=Cert Authority/emailAddress=rahulchitta22594@gmail.com"

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=CA/ST=Ontario/L=Cambridge/O=Teleport/OU=Assessment/CN=server/emailAddress=rahulchitta22594@gmail.com"

# 3. Use CA's private key to sign web server's CSR and produce server's signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

# 4. Generate client (Alice) private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-alice-key.pem -out client-alice-req.pem -subj "/C=CA/ST=Ontario/L=Cambridge/O=Teleport/OU=Assessment/CN=alice/emailAddress=rahulchitta22594@gmail.com"

# 5. Use CA's private key to sign client (Alice's) CSR and produce the signed certificate
openssl x509 -req -in client-alice-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-alice-cert.pem -extfile client-alice-ext.cnf

# 6. Generate client (Bob) private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-bob-key.pem -out client-bob-req.pem -subj "/C=CA/ST=Ontario/L=Cambridge/O=Teleport/OU=Assessment/CN=bob/emailAddress=rahulchitta22594@gmail.com"

# 7. Use CA's private key to sign client (Bobs's) CSR and produce the signed certificate
openssl x509 -req -in client-bob-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-bob-cert.pem  -extfile client-bob-ext.cnf

# 8. Generate unauthorized client private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-unauth-key.pem -out client-unauth-req.pem -subj "/C=CA/ST=Ontario/L=Cambridge/O=Teleport/OU=Assessment/CN=unauth/emailAddress=rahulchitta22594@gmail.com"

# 9. Use CA's private key to sign unauthorized client's CSR and produce the signed certificate
openssl x509 -req -in client-unauth-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-unauth-cert.pem -extfile client-unauth-ext.cnf