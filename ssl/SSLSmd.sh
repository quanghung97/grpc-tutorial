openssl genrsa -out server.key 2048

openssl req -nodes -new -x509 -sha256 -days 1825 -config ./cert.conf -extensions 'req_ext' -key server.key -out server.crt
