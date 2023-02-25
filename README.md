# cleanmyarr
A lightweight utility to delete movies and shows from Radarr and Sonarr after specified time. 

**NOTE: Currently it only supports Radarr. Development for Sonarr is underway.**

Cleanmyarr uses Radarr and Sonarr API to communicate with them and get the list of movies and episodes.
It checks the age of movie/Show files in Radarr/Sonarr.

If age of file is more than the expiry age (deleteAfterDays) specified in configuration file (config.yaml), It deletes the file and movie/show from Radarr/Sonarr. It performs this job at every maintenance cycle which is configured in configuration file.

Some Movies/Shows can also be marked to ignore. In this case, movies/shows marked to ignore will not be deleted even if their age is more than expiry time. Movies can be marked to ignore by adding a tag to it which is specified in configuration file (config.yaml).


It also sends the notification through configured channels to user. Currently supported channels are Email (SMTP), Gotify and Telegram.

Status of the last activity is stored in /config/status.yaml file, where it stores information like,

- Last time maintenance was run (lastMaintenanceDate)
- Next time at which maintenance will run (nextMaintenanceDate)
- Movies deleted when last maintenance was run (deletedMovies)
- Shows deleted when last maintenance was run (deletedShows)
- Movies which are marked to ignore (ignoredMovies)
- Shows which are marked to ignore (ignoredShows)
- Movies which are marked for deletion and will be deleted on next maintenance (moviesMarkedForDeletion)
- Shows which are marked for deletion and will be deleted on next maintenance (showsMarkedForDeletion)

## How to install

To install cleanmyarr on docker run below command

```
mkdir $HOME/cleanmyarr/config
chown 1000:1000 $HOME/cleanmyarr/config
chmod 755 $HOME/cleanmyarr/config

docker run -d \
    --name cleanmyarr \
    -v $HOME/cleanmyarr/config:/config \
    --net mynetwork \
    --restart=unless-stopped \
    linuxshots/cleanmyarr:latest
```

With `docker-compose.yml`

```
version: "3.9"
services:
  radarr:
    container_name: radarr
    image: lscr.io/linuxserver/radarr:latest
    networks:
      - mynetwork
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=UTC
    ports:
      - 7878:7878
    volumes:
      - radarr-config:/config
      - torrent-downloads:/downloads
    restart: "unless-stopped"

  cleanmyarr:
    depends_on:
      - radarr
    image: linuxshots/cleanmyarr:latest
    container_name: cleanmyarr
    networks:
      - mynetwork
    volumes:
      - /path/to/config:/config
    restart: unless-stopped

volumes:
  cleanmyarr-config:
  radarr-config:
  torrent-downloads:

networks:
  mynetwork:
    external: true
```

## Configuration

A configuration file `config.yaml` is created in `/config` directory inside container. By default, Configuration file has everything disabled and values set to default. It is a YAML file which is very easy to understand and easy to configure.

Find sample configuration file [here](sample-config.yaml).

**`maintenanceCycle`** This defines how often files must be scanned for deletion. Movies/Shows are deleted only during maintenance. Valid values are `daily`, `every3days`, `weekly`, `bimonthly` and `monthly`. 

**`deleteAfterDays`** This is time-to-live of a movie/show. If a movie/show has aged more than this, It will be deleted, unless its marked to ignore. Valid values are any non-decimal (non-fractional) number more than 0. *Default value is 90* **days.** 

**`ignoreTag`** This can have any string value. Any movie/show tagged with this string will be ignored and will not be deleted even after its age is more than age of expiry. *Default value is `cma-donotdelete`* 

**`notificationChannel`** This section has three different notification channels, `smtp`, `gotify` and `telegram`. These three must be enabled by setting `enabled: true`. They need creadentials and some other information specific to them for them to work. Check sample configuration file [here](sample-config.yaml). 

**`radarr`** This section contains configuration related to radarr. Cleanmyarr will check movies only if it is enabled by setting `enabled: true` in its configuration. Some other configurations are Radarr url, Base64 encoded API key to communicate with Radarr and whether notification is enabled for radarr.

**`sonarr`** This section contains configuration related to sonarr. Cleanmyarr will check shows only if it is enabled by setting `enabled: true` in its configuration. Some other configurations are sonarr url, Base64 encoded API key to communicate with sonarr and whether notification is enabled for sonarr.

👉🏾 [SAMPLE CONFIGURATION FILE](sample-config.yaml)

All the credential and sensitive values must be added after base64 encoding it. You can base64 encode any string using below command in Linux terminal.

```
echo 'str1ng-to_enc0de' | base64 -w 0
```

## Setup telegram bot for telegram notification

Article link to be updated.

