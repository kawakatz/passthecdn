# passthecdn

Remove subdomains which are behind CDN (CloudFlare, CloudFront, Fasty, Akamai).<br>
This tool compares resolved IP and IP ranges of CDNs.<br>
For Akamai, exec reverse DNS lookup and check the result contains "deploy.static.akamaitechnologies.com".<br>

install<br>
```➜  ~ go get -u -v github.com/kawakatz/passthecdn```

usage<br>
```➜  ~ cat subdomains.txt | passthecdn```

I'm thinking about the best way how to output the result including flexibility of options.
