# heatpump2mqtt systemd service

Here is a sample file of systemd service definition in a [heatpump2mqtt.service](heatpump2mqtt.service) file.

### Install service

```console
cp heatpump2mqtt.service /etc/systemd/system/heatpump2mqtt.service
```

### Enable service

```console
systemctl daemon-reload
systemctl enable heatpump2mqtt.service
```
