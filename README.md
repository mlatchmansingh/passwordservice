# passwordservice

## To build
```
go build
```

## To run
```
./PasswordService
```

## To test
```
$Passwd=@{                                                                  
>> password='angryMonkey'
>> }

Invoke-WebRequest -Uri http://localhost:8080/hash -Method Post -Body ($Passwd | ConvertTo-Json) -ContentType "application/json"
```

Five seconds after requesting a password you can get the hashed value by invoking

```
Invoke-WebRequest -Uri http://localhost:8080/hash/1 -Method Get
```


## Example

```
Invoke-WebRequest -Uri http://localhost:8080/hash/1 -Method Get


StatusCode        : 200
StatusDescription : OK
Content           : ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0 
                    Fd7su2A-gf7Q==
RawContent        : HTTP/1.1 200 OK
                    Content-Length: 88
                    Content-Type: text/plain; charset=utf-8
                    Date: Sun, 16 Jan 2022 18:13:35 GMT

                    ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0
                    Fd7su2A-g...
Forms             : {}
Headers           : {[Content-Length, 88], [Content-Type, text/plain; charset=utf-8], [Date,   
                    Sun, 16 Jan 2022 18:13:35 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 88
```

## Stats
Get the statistics for POST requests
```
Invoke-WebRequest -Uri http://localhost:8080/stats -Method Get
```

## Shutdown
To shutdown the server and prevent new requests
```
Invoke-WebRequest -Uri http://localhost:8080/shutdown -Method Post
```
This call may take some time to complete if there are in-flight requests that must complete first
