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
