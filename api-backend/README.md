# api backend

## Purpose

The api-backend server returns gps log, and photos.

## Initialization

## Create sealed secret manifest of mysql secrets

Store credential to file.
Then, Remove command history from `~/.bash_history`.

```bash
echo -n 'foo' > mysql_root_password.txt
echo -n 'bar' > mysql_user.txt
echo -n 'baz' > mysql_password.txt
echo -n 'foobar' > access_key_id.txt
echo -n 'barbaz' > secret_access_key.txt
echo -n 'foobaz' > bucket.txt
echo -n 'foobarbaz' > jwt_secret.txt
```

Install kubeseal cli and create controller resource. [github](https://github.com/bitnami-labs/sealed-secrets)
Create sealed secret manifest from certificate.

```bash
kubectl create secret -n database generic mysql-secrets \
  --from-file=mysql_root_password=./mysql_root_password.txt \
  --from-file=mysql_user=./mysql_user.txt \
  --from-file=mysql_password=./mysql_password.txt \
  -o yaml --dry-run=client >secret1.yml
kubeseal -o yaml <secret1.yml >sealedsecret1.yml
kubectl create secret -n api generic api-secrets \
  --from-file=mysql_user=./mysql_user.txt \
  --from-file=mysql_password=./mysql_password.txt \
  --from-file=access_key_id=./access_key_id.txt \
  --from-file=secret_access_key=./secret_access_key.txt \
  --from-file=bucket=./bucket.txt \
  --from-file=jwt_secret=./jwt_secret.txt \
  -o yaml --dry-run=client >secret2.yml
kubeseal -o yaml <secret2.yml >sealedsecret2.yml
rm secret1.yml secret2.yml mysql_root_password.txt mysql_user.txt mysql_password.txt \
  access_key_id.txt secret_access_key.txt bucket.txt jwt_secret.txt
```

## Copy sealed secret manifest to api-backend manifest

Copy sealed secret manifest to the part of database.yml and api-backend.yml.
