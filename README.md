# lsm-verification

Usage: 
1. Generate private and public RSA keys, for example, by running
```bash
openssl genrsa -out ~/mykey.pem 2048
openssl rsa -in ~/mykey.pem -pubout > ~/mykey.pub
```

2. Build the project:
```bash
go build
```

3. Hash and sign your database replica:
```bash
lsm-verification hash-and-sign your_server_address your_replica_id ~/mykey.pem
```

4. Distribute your public key

5. Others can verify data from your replica by running
```bash
lsm-verification verify their_server_address your_replica_id ~/mykey.pub
```
