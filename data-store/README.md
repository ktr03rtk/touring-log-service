# data store

## Purpose

The data-store server pod subscribes touring log, then compress to gzip and store to AWS S3.

## Initialization

## Create sealed secret manifest of mqtt client certificate and AWS credential

Store credential to file.
Then, Remove command history from `~/.bash_history`.

```bash
echo -n 'foo' > access_key_id.txt
echo -n 'bar' > secret_access_key.txt
echo -n 'baz' > bucket.txt
```

Install kubeseal cli and create controller resource. [github](https://github.com/bitnami-labs/sealed-secrets)
Create sealed secret manifest from certificate.

```bash
kubectl create secret -n data-store generic data-store-credential \
  --from-file=access_key_id=./access_key_id.txt \
  --from-file=secret_access_key=./secret_access_key.txt \
  --from-file=bucket=./bucket.txt \
  -o yaml --dry-run=client >secret1.yml
kubeseal -o yaml <secret1.yml >sealedsecret1.yml
rm secret1.yml access_key_id.txt secret_access_key.txt bucket.txt
```

```bash
kubectl create secret -n data-store generic mqtt-client-certificate \
  --from-file=ca.crt \
  --from-file=client.crt \
  --from-file=client.key \
  -o yaml --dry-run=client >secret2.yml
kubeseal -o yaml <secret2.yml >sealedsecret2.yml
rm secret2.yml
```

## Copy sealed secret manifest to mqtt manifest

Copy sealed secret manifest to the part of mqtt_server.yml

## Usage

Apply data_store.yml and execute pub/sub command.
