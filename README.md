# Echo Alert - Alertmanager proof of concept

Echo Alert or `echoalert` is a small web service that consists of a number of webhook receivers that will be defined as receivers in the Alertmanager configuration. The `echoalert` web service can also respond with a fixed HTTP status code by providing it as a query parameter in the URL.

When alerts are triggered they'll be processed by Alertmanager and by matching on labels configured on the alert rules the alerts will be sent to different webhook receivers. In this case the same web service will handle all webhook receiver endpoints.

You can see the different alerts sent to the `echoalert` service by browsing to `http://localhost:8080`. This is a simple view of all the alerts, as well as counters on the number of received HTTP POST requests on every endpoint. 

## Components

* Prometheus
* Alertmanager
* Blackbox Exporter
* Echo Alert service

We'll define a couple of [alert rules](prometheus/blackbox-rules.yml) in Prometheus, this is for alerts triggered on Blackbox Exporter probe pings. This is to test the `echoalert` endpoint that can return a fixed HTTP status code. Please see the `targets` list of the `blackbox-http-checks` job config in the [Prometheus configuration](prometheus/prometheus.yml) file.

## Run the proof of concept

In the root dir of this repository run:
  ```
  docker-compose up -d
  ```

Now you should be able to reach the following endpoints:

* Prometheus: [`http://localhost:9090`](`http://localhost:9090`)
* Alertmanager: [`http://localhost:9093`](`http://localhost:9093`)
* `echoalert`: [`http://localhost:8080`](`http://localhost:8080`)

## Testing the services

* Start a browser and browse to [`http://localhost:8080`](`http://localhost:8080`), you'll see the Blackbox Exporter alert rules triggered soon, one of the is the target that returns a HTTP `401` status code.
* Fire off the alerts defined in the [alert send script](scripts/send_alerts.sh). 
* Test the `echoalert` customizable endpoint by running:
```
curl http://localhost:8080/shapeshifter?body=notok&code=400

HTTP/1.1 400 Bad Request
Connection: close

notok
```

```
curl http://localhost:8080/shapeshifter?body=ok&code=200

HTTP/1.1 200 OK
Connection: close

ok
```