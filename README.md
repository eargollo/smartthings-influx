# smartthings-influx

A simple program to bring to Influx your SmartThings data through the SmartThings API. No SmartApp installation needed.

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

