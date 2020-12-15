# passthecdn

Remove subdomains which are behind CDN (CloudFlare, CloudFront, Fasty).<br>
This tool compares resolved IP and IP ranges of CDNs.<br>

install<br>
```➜  ~ go get -u -v github.com/kawakatz/passthecdn```

usage<br>
```➜  ~ cat subdomains.txt | passthecdn```
