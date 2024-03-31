# Migrating from InfluxDB v1 to v2

Support for InfluxDB v2 was added in March 2024. This guide will help you upgrade from InfluxDB v1 to InfluxDB v2.

It assumes default configuration. Make changes as you see fit for your setup.

## Step 1: Migrating the data

First we recommend stopping your monitoring system.

We will be migrating the data from `data/influxdb` to `data/influxdb2`. Your v1 setup will stay functional with all its existing data.

For doing that create the `docker-migration-compose.yaml` file that will be responsible for the migration. Put in it the file content as below and adjust environment variables according to your setup:

```
version: '2'
services:
  influxdb:
    image: influxdb:2
    ports:
      - '8086:8086'
    volumes:
      - ./data/influxdb:/var/lib/influxdb
      - ./data/influxdb2:/var/lib/influxdb2
    environment:
      - INFLUXDB_DB=SmartThings
      - INFLUXDB_ADMIN_USER=admin
      - INFLUXDB_ADMIN_PASSWORD=password
      - DOCKER_INFLUXDB_INIT_MODE=upgrade
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=password
      - DOCKER_INFLUXDB_INIT_ORG=org
      - DOCKER_INFLUXDB_INIT_BUCKET=SmartThings
      - DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=token
```

We are now ready to create the InfluxDB2 database and migrate the data. 

To do so, run docker compose:
```
docker-compose -f docker-migration-compose.yaml up
```

Keep the docker compose running for the next step.

## Step 2: Renaming the database

First and foremost, let's check that the migration went well. Go to Influx2 US at http://localhost:8086

Login with user and password set at the docker compose file (`admin/password` by default).

Go to data explorer and you should see two buckets: `SmartThings` and `SmartThings/augoten`. The migration happened to the `SmartThings/autogen`. Validate the data by clicking on it and querying some metric.

From here you have two options. Either reconfigure your setup to use the new bucket, or delete the `SmartThings` one and rename `SmartThings/autogen` to `SmartThings`. 

In this guide we are considering the delete/renaming path. Do that at `Data Explorer/Buckets`.

Now you can stop your migration compose.

## Step 3: Update your setup to use Influx v2

On your configuration file change the db properties deleting the Influx v1 as below:
```
...
influxurl: http://influxdb:8086
influxuser: admin
influxpassword: password
influxdatabase: SmartThings
...
```

And adding the Influx v2:
```
...
database:
  type: influxdbv2
  url: http://influxdb:8086
  token: token
  org: org
  bucket: SmartThings
...
```

On your docker compose file `docker-compose.yaml` you will need to update your InfluxDB service from:
```
...
  influxdb:
    image: influxdb:1.8
    ports:
      - '8086:8086'
    volumes:
      - ./data/influxdb:/var/lib/influxdb
    environment:
      - INFLUXDB_DB=SmartThings
      - INFLUXDB_ADMIN_USER=admin
      - INFLUXDB_ADMIN_PASSWORD=password
...
```

To:
```
...
  influxdb:
    image: influxdb:2
    ports:
      - '8086:8086'
    volumes:
      - ./data/influxdb2:/var/lib/influxdb2
    environment:
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=password
      - DOCKER_INFLUXDB_INIT_ORG=org
      - DOCKER_INFLUXDB_INIT_BUCKET=SmartThings/autogen
      - DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=token
...
```
You should be able to run your setup again. But you still need to update Grafana configuration

## Step 4: Updating Grafana

Log in to your Grafana and update its database connection at data sources. There will be a connection named "InfluxDB".

Change `Query language` from `InfluxQL` to `Flux`. 

At `InfluxDB Details` fill `Organization` and `Token`. Defaults are `org` and `token`.

Set `Default bucket` to `SmartThings`. 

Save and test should work well. 

Your dashboards won't work given that Influx changed its query language. If you want a head start, import the dashboard as a json [here](/grafana-provisioning/dashboards/smartthings.json)

At times I needed to enter each panel edit, wait a few seconds, refresh a few times until they worked properly.

Now you are set!

