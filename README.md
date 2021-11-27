# smartthings-influx

A simple program to bring to Influx your SmartThings data through the SmartThings API. No SmartApp installation needed.

## Getting started

If you have Docker and Docker Compose, you can use this getting started, in just 3 steps.

- Step 1: Create an SmartThings API token

Go to [SmartThings API Token](https://account.smartthings.com/tokens) page and create a token. Make sure you give full access to devices.

- Step 2: Place the token at the `docker-compose-config.yml` file

- Step 3: Bring the stack up and see your Grafana chart

Run Docker Compose:
```
$ docker-compose up --build
```

Go to [Grafana inteface](http://localhost:3000) and log with user `admin` and password `password`.

There is already a pre-provisioned Grafana dashboard to hold your SmartThings data. In Grafana go to Dashboards->Manage and there click on Smartthings.

Have fun!

## Usage

Create the `.smartthings-influx.yaml` file either at your home folder or at the folder where you run the program:

```yaml
apitoken: <put your SmartThings API token here>
monitor:
  - temperatureMeasurement
  - illuminanceMeasurement
  - relativeHumidityMeasurement
period: 120
influxurl: http://localhost:8086
influxuser: user
influxpassword: password
influxdatabase: database
```

Run the monitor option:
```
$ smartthings-influx monitor
```

