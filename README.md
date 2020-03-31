[![Build Status](https://travis-ci.org/proway2/kvserver.svg?branch=master)](https://travis-ci.org/proway2/kvserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/proway2/kvserver)](https://goreportcard.com/report/github.com/proway2/kvserver)

# kvserver
Simple yet fully functional in-memory key-value storage server based on HTTP protocol with elements purged based on TTL. All operations run in constant time.    
Operations provided by server:

- storing/updating value by its key
- getting value by its key
- deleting value by its key
- key-value element is cleaned up when expired (works automatically).

# Features

- all operations have time complexity of ```O(1)```, i.e. always run in constant time.    
- hits TTL as much accurate as it's possible.    
- lower CPU cycles consumption during approximation and idle.    
- TTL approximation's divider is always 2, i.e. next check time = current time + (time to the next element to purge)/2.

# Installation
Clone and run ```go install``` in project folder.

# Usage
Command line arguments:
```bash
$ kvserver -h
Usage of kvserver:
  -addr string
    	IP address to bind to (default "127.0.0.1")
  -port int
    	port to listen to (default 8080)
  -ttl uint
    	element's (key-value) lifetime in the storage, secs. (default 60)
```
# API
Base URL ```http://<host>:<port>/key/<key_name>```, where ```<key_name>``` - is the name of the key to be stored. Key and its value are always string.
## Storing/Updating value by its key
_HTTP method_: ```POST```    
_Request's parameter name_: ```value```    
_Success code_: ```200```    
_Error code_: no code, this is a crash.    
_Note_: TTL is reset for any subsequent requests for the same key.

## Getting value by its key
_HTTP method_: ```GET```    
_Request's parameter name_: no parameter is needed.    
_Success code_: ```200```, response's body contains string value for the key.    
_Error code_: ```404```

## Deleting value by its key
_HTTP method_: ```POST```    
_Request's parameter name_: no parameter is needed.    
_Success code_: ```200```    
_Error code_: ```404```

When error is occured code ```400``` is returned by server.

# Tests
Run ```run_test.sh``` in project folder.

# License
GPL v3
