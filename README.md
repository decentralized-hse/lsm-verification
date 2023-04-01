# lsm-verification

Utility for hashing, signing, and validating [LSM database](https://github.com/ds-project-lseqdb/ds-project-public) replica entries.

### Usage: 
1. Generate private and public RSA keys, for example, by running
```bash
openssl genrsa -out ~/mykey.pem 2048
openssl rsa -in ~/mykey.pem -pubout > ~/mykey.pub
```

2. Build the project:
```bash
go build
```

3. Set environment variables:
```bash
export dbServerAddress="<address>:<port>"
export dbReplicaID="<replicaID>"
export rsaPublicKey=$(cat ~/mykey.pub)
export rsaPrivateKey=$(cat ~/mykey.pem)
```

3. Hash and sign your database replica entries:
    - Set `run_mode: "Sign"` in `config/config.yml`
    - Run `./lsm-verification`

4. Distribute your public key

5. Others can verify data from your replica by
    - Setting their `run_mode: "Validation"` in `config/config.yml`
    - Setting their `dbReplicaID` to your replica ID
    - Setting their `rsaPublicKey` to your public key
    - Running `./lsm-verification`

