# frontend

## Purpose

The frontend server pod is react application react web application server.
User can upload photo, and view the touring log and related photo on map.

## Available Scripts

For development.

```
npm start
```

For production build.

```
npm run build
```

## Initialization

## Create sealed secret manifest of AWS credential

Store credential to file.
Then, Remove command history from `~/.bash_history`.

```bash
echo -n 'foo' > access_key_id.txt
echo -n 'bar' > secret_access_key.txt
echo -n 'baz' > s3_bucket.txt
```

Install kubeseal cli and create controller resource. [github](https://github.com/bitnami-labs/sealed-secrets)
Create sealed secret manifest from certificate.

```bash
kubectl create secret -n web generic web-credential \
  --from-file=access_key_id=./access_key_id.txt \
  --from-file=secret_access_key=./secret_access_key.txt \
  --from-file=s3_bucket=./s3_bucket.txt \
  -o yaml --dry-run=client >secret1.yml
kubeseal -o yaml <secret1.yml >sealedsecret1.yml
rm secret1.yml access_key_id.txt secret_access_key.txt s3_bucket.txt
```

## Copy sealed secret manifest to web manifest

Copy sealed secret manifest to the part of web_server.yml.

## Usage

Apply web_server.yml.
