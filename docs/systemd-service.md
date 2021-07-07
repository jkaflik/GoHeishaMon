# GoHeishamon systemd service

Here is a sample file of systemd service definition in a [goheishamon.service](goheishamon.service) file.

### Install service

```console
cp goheishamon.service /etc/systemd/system/goheishamon.service
```

### Enable service

```console
systemctl daemon-reload
systemctl enable goheishamon.service
```
