# xk6-redis

This is a [WAMP protocol](https://wamp-proto.org) client library for [k6](https://github.com/loadimpact/k6),
implemented as an extension using the [xk6](https://github.com/k6io/xk6) system.

It is based on [nexus WAMP client library](https://github.com/gammazero/nexus).

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  go get -u github.com/k6io/xk6/cmd/xk6
  ```

2. Build the binary:
  ```shell
  xk6 build v0.31.1 --with github.com/vvarp/xk6-wamp
  ```

## Example test script

```javascript
import wamp from 'k6/x/wamp';
import { sleep } from 'k6';

export default function () {
    const client = new wamp.Client(
        "ws://127.0.0.1:9000",
        {
            realm: "default",
        }
    );

    const subId = client.subscribe("test", {}, function (args, kwargs) {
        console.log(args, kwargs)
    });

    const sessId = client.getSessionID();
    console.log(`Subscription ID: ${sessId} ${subId}`)

    sleep(Math.random() * 2);

    client.disconnect()
}
```
