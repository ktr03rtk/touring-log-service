# VPN client

## Purpose

This VPN client cronjob download vpn config file from S3 bucket and execute openvpn command to VPN connection with the VPN server of the public cloud.

## Create sealed secret manifest of AWS credential

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
kubectl create secret -n vpn generic vpn-client-credential \
  --from-file=access_key_id=./access_key_id.txt \
  --from-file=secret_access_key=./secret_access_key.txt \
  --from-file=bucket=./bucket.txt \
  -o yaml --dry-run=client >secret.yml
kubeseal -o yaml <secret.yml >sealedsecret.yml
rm secret.yml access_key_id.txt secret_access_key.txt bucket.txt
```

## Copy sealed secret manifest to vpn client manifest

Copy sealed secret manifest to the part of vpn_client.yml

## Usage

Apply vpn_client.yml and execute pub/sub command.

Create AWS resource for VPN server. (the terraform github repository is private. Sorry.)
