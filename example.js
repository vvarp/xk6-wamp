import wamp from 'k6/x/wamp'
import {sleep} from 'k6'

export default function () {
    const client = new wamp.Client(
        "ws://127.0.0.1:8080/ws",
        {
            realm: "realm1",
        }
    )

    const subId = client.subscribe("com.example.topic1", {}, function (args, kwargs) {
        console.log(`Received args=${args} kwargs=${JSON.stringify(kwargs)}`)
    });

    const sessId = client.getSessionID()
    console.log(`sessionId=${sessId} subscriptionId=${subId}`)

    client.publish("com.example.topic2", {}, [1, 2, 3], {one: "1", two: "2"})

    sleep(Math.random() * 5 + 1)

    client.disconnect()
}
