# smartthings-influx

A simple program to bring to Influx your SmartThings data through the SmartThings API. No SmartApp installation needed.

## Getting started

If you have Docker and Docker Compose, you get started in just 3 steps.

- Step 1: Create an SmartThings API token

Go to [SmartThings API Token](https://account.smartthings.com/tokens) page and create a token. Make sure you give full access to devices.

- Step 2: Either place the token at the `docker-compose-config.yml` file or set the APITOKEN environment variable `export APITOKEN=YOUR-TOKEN-HERE`

- Step 3: Bring the stack up and see your Grafana chart

Run Docker Compose:
```
$ UID=$(id -u) GID=$(id -g) docker-compose up
```

Go to [Grafana inteface](http://localhost:3000) and log with user `admin` and password `password`.

There is already a pre-provisioned Grafana dashboard to hold your SmartThings data. In Grafana go to Dashboards->Manage and there click on Smartthings.

Have fun!

## Running locally

1. Download the latest version of `smartthings-influx` that is compatible to your platform [here](https://github.com/eargollo/smartthings-influx/releases)

1. Go to [SmartThings API Token](https://account.smartthings.com/tokens) page and create a token. Make sure you give full access to devices.

1. Create the configuration `.smartthings-influx.yaml` file either at your home folder or at the folder where you run the program. Here is a sample of the configuration content:

```yaml
apitoken: <put your SmartThings API token here or export APITOKEN env var>
monitor:
  - light
  - temperatureMeasurement
  - illuminanceMeasurement
  - relativeHumidityMeasurement
period: 120
influxurl: http://localhost:8086
influxuser: user
influxpassword: password
influxdatabase: database
valuemap:
  switch: 
    off: 0
    on: 1
```

4. Put the API token you generated at step 2 at the configuration file (alternativelly you can export it as the `APITOKEN` environment variable)

1. Now we can test the setup by listing your devices running the `list` command:

```
$ ./smartthings-influx list
Using config file: /Users/eduardoargollo/src/eargollo/smartthings-influx/.smartthings-influx.yaml
0: f33840a1-f835-41ff-b8f8-b8c95d768363, c2c-rgbw-color-bulb, Living Room Floor Lamp
   | switch
   | switchLevel
   | colorControl
   | colorTemperature
   | refresh
   | healthCheck
1: 27118afd-2f98-425c-8211-899eb596edad, GE Wall Switch, Garbage Disposal
   | switch
   | refresh
2: da74302c-5653-406a-8770-3323500a41fb, c2c-rgbw-color-bulb, Office Floor Lamp
   | switch
   | switchLevel
   | colorControl
   | colorTemperature
   | refresh
   | healthCheck
...
21: 99e5de18-3fbc-4769-a83e-e466a8564f6d, GE Wall Switch, Closet Light
   | switch
   | refresh
```

If you are getting a similar list then the one above, `smarththings-influx` is working properly. You can even use the data on this list to enhance your configuration and monitor more devices.

If this command is not working, check your APIKEY is correct, with permissions, and placed correctly at the file (replacing the `<put your SmartThings API token here or export APITOKEN env var>` [comment](https://github.com/eargollo/smartthings-influx/blob/master/docker-compose-config.yaml#L1) )


6. Now you will need to install InfluxDB. The program does not support InfluxDB 2, so install the latest InfluxDB version 1 (at the time of this writing it is the version 1.8.10). There may be packages for your computer platform. For instance you can install it on Mac using homebrew with `brew install influxdb@1.8.10`

1. Run InfluxDB and note the URL, port, user and password for the database installation

1. Update the `.smartthings-influx.yaml` with the InfluxDB configuraiton (Example for a local InfluxDB [here](https://github.com/eargollo/smartthings-influx/blob/master/docker-compose-config.yaml#L9-L12))

1. Now you can run the monitor command and it should work without errors: `./smartthings-influx monitor` . This means data is being loaded to InfluxDB

1. Install [Grafana](https://grafana.com/), run it and configure it to access your InfluxDB. 

You should now be able to see the datapoints and create your Grafana charts.

Note, all of this has been made and pre-set with the Docker compose at this repository and this could serve you as a guide. Look at the [docker compose](https://github.com/eargollo/smartthings-influx/blob/master/docker-compose.yml) file to validate the components and at (grafana-provisioning folder)[] for the Grafana configuraiton used at the docker compose version. It also has an initial dashboard.

Have fun!
