# dp ![Build status](https://travis-ci.org/Bo0mer/dp.svg?branch=master)


Sniff HTTP communication between two applications. Useful for reverse engineering, debugging or just for fun.

## Example usage

Let's say you want to monitor the communication between a local CLI, called `cf`, and a remote located at https://target.com. This could be done via:
```
dp -target https://target.com -addr localhost:8080
```

And now, tell the `cf` CLI that your remote server is located at `localhost:8080`. For example:

```
cf api http://localhost:8080
Setting api endpoint to localhost:8080...
OK
```

Now observe how `dp` has printed the request and response headers and body to its stdout.
It will automatically try to guess how to format the content of the body, if
no explicit format is provided using the `format` flag.

## Additional options
If the payload of the request/response bodies is formatted in JSON:
```
dp -target https://target.com -format json
```

If the remote host is using TLS, but its certificate is not valid for some reason you can use the `-insecure` flag.
Please be careful!
```
dp -target https://invalid.cert -insecure
```

### Supported payload formats (for pretty printing)
* auto (try to guess, includes all of the below)
* json
* plain

## Full usage help
```
Usage of dp:
  -addr string
    	Address to bind to. (default "localhost:8080")
  -cert string
    	Certificate file to use for TLS.
  -format string
    	Attempt to format payloads as. (default "auto")
  -header value
    	Header to add. Must be in Name:value format.
  -insecure
    	Please do not!
  -key string
    	Key file to use for TLS.
  -skip-body
    	Don't dump body.
  -target string
    	Target to proxy to. (default "https://example.com")
```
