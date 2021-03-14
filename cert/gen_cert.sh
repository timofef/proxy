#!/bin/sh

openssl req -new -key $2/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 3650 -CA $2/ca.crt -CAkey $2/ca.key -set_serial "$3" > $4$1.crt
