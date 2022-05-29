# MQTT Server

## Purpose

The pub/sub server pod receives touring log from raspberry pi, and send log to servers run by home k8s.

## Initialization

### Create mqtt server certificate

Set your mqtt server domain to `IP` value in `script/create_cert.sh` and execute the script.

## Create sealed secret manifest of mqtt server certificate

Install kubeseal cli and create controller resource. [github](https://github.com/bitnami-labs/sealed-secrets)
Create sealed secret manifest from certificate.

```bash
kubectl create secret -n mqtt generic mqtt-server-certificate \
  --from-file=ca.crt \
  --from-file=server.crt \
  --from-file=server.key \
  -o yaml --dry-run=client >secret.yml
kubeseal -o yaml <secret.yml >sealedsecret.yml
rm secret.yml
```

## Copy sealed secret manifest to mqtt manifest

Copy sealed secret manifest to the part of mqtt_server.yml

## Usage

Apply mqtt_server.yml and execute pub/sub command.

```bash
mosquitto_pub -h $IP -p $SERVER_PORT -t topic/subtopic --cafile ca.crt --cert client.crt --key client.key -m "message"
```

```bash
mosquitto_sub -h $IP -p $SERVER_PORT -t topic/subtopic --cafile ca.crt --cert client.crt --key client.key
```
