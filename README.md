# fritzbox-cloudflare-dyndns-bridge

Update your domain using the dynamic DNS feature of an AVM FritzBox   

# Installation (Docker)

```
git clone https://github.com/lucasl0st/fritzbox-cloudflare-dyndns-bridge.git
```

docker-compose.yml
```
services:
  bridge:
    build: ./fritzbox-cloudflare-dyndns-bridge
    ports:
      - 80:8000
    restart: unless-stopped
    volumes:
      - ./config.json:/app/config.json
```

config.json
```
{
  "users": [
    {
      "username": "username",
      "password": "password",
      "allowedDomains": [
        "FQDN"
      ],
      "cloudflareApiKey": "APIKEY"
    }
  ]
}
```

Use this update-url in your FritzBox settings:   

Don't replace the variables inside the <>, these are replaced by the FritzBox!

```
http://DYN_HOSTNAME/dyn?domain=<domain>&username=<username>&password=<pass>&ipaddr=<ipaddr>&ip6addr=<ip6addr>&dualstack=<dualstack>&ip6lanprefix=<ip6lanprefix>
```