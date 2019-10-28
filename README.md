#### Reverse proxy

It can host a website or api on embedded server. 
All original url's are translated into server host address. 

You can host for example https://en.wikipedia.org on your server, 
perform search and get requests - all presented on you site. 

Installation is simple as
```
go get https://github.com/blunext/reverse-proxy
```
You can use help:
```
 ./Proxy -help
```
You can run it:
```
./Proxy -proxy https://localhost:9013 \
    -target https://en.wikipedia.org \
    -parseWithoutSchema=true \
    -certFile /Users/tom/.ssh/localhost-ssl/localhost.crt \
    -keyFile /Users/tom/.ssh/localhost-ssl/localhost.key
``` 
- Certs are not necessary when you specify http scheme.
- parseWithoutSchema:
    - flag set into true "https://en.wikipedia.org" and "en.wikipedia.org" will be changed
    - flag set into false "https://en.wikipedia.org" will be changed while "en.wikipedia.org" will not
        
Limitations: modern websites use different scripts from external urls. 
Those external resources will not be processed and parsed so site may not work correctly.   

The project was made for fun. I tested it on my work on on Jira and Confluence with success but it's not production ready.  