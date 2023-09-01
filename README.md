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

```bash
docker run -d \
    --name cleanmyarr \
    --net mynetwork \
    --restart=unless-stopped \
    --env CMA_MAINTENANCE_CYCLE="bimonthly" \
    --env CMA_DELETE_AFTER_DAYS=90 \
    --env CMA_ENABLE_GOTIFY_NOTIFICATION=true \
    --env CMA_GOTIFY_URL=gotify.example.com \
    --env CMA_GOTIFY_ENCODED_APP_TOKEN="dGgxc2lzbjB0QSQzY3IzdAo=" \
    --env CMA_MONITOR_RADARR=true \
    --env CMA_RADARR_URL=radarr.example.com \
    --env CMA_RADARR_ENCODED_API_KEY="dGhpc2lzbm90YW5hcGlrZXkK" \
    --env CMA_RADARR_ENABLE_NOTIFICATION=true \
    linuxshots/cleanmyarr:latest
```

With `docker-compose.yml`

```yaml
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
    envitonment:
      - CMA_MAINTENANCE_CYCLE="bimonthly"
      - CMA_DELETE_AFTER_DAYS=90
      - CMA_ENABLE_GOTIFY_NOTIFICATION=true
      - CMA_GOTIFY_URL=gotify.example.com
      - CMA_GOTIFY_ENCODED_APP_TOKEN="dGgxc2lzbjB0QSQzY3IzdAo="
      - CMA_MONITOR_RADARR=true
      - CMA_RADARR_URL=radarr.example.com
      - CMA_RADARR_ENCODED_API_KEY="dGhpc2lzbm90YW5hcGlrZXkK"
      - CMA_RADARR_ENABLE_NOTIFICATION=true
    restart: unless-stopped

volumes:
  radarr-config:
  torrent-downloads:

networks:
  mynetwork:
    external: true
```

## Parameters

Cleanmyarr can be configured using parametes passed at runtime as in above example.

| Parameter | Function | Default value | Possible values |
| :----: | --- | --- | --- |
| CMA_MAINTENANCE_CYCLE | Frequency of maintenance/cleanup | Daily | `Daily` , `Every3Days`, `Weekly`, `Bimonthly`, `Monthly` |
| CMA_DELETE_AFTER_DAYS | Expiry age of movies/series in days | 90 | Any whole number greater than 0 |
| CMA_IGNORE_TAG | Tag to mark a movie/series to ignore during maintenence | cma-donotedelete | Any string |
| CMA_ENABLE_EMAIL_NOTIFICATION | Enable email notification | false | `true`, `false`, `yes`, `no` |
| CMA_SMTP_SERVER | SMTP server to use for sending email notification | smtp.gmail.com | Valid smtp server hostname |
| CMA_SMTP_PORT | SMTP server TLS port | 587 | Valid smtp port number |
| CMA_SMTP_USERNAME | Username to authenticate to SMTP server | example@gmail.com | valid username string |
| CMA_SMTP_ENCODED_PASSWORD | Base64 encoded password to authenticate to SMTP server | | Valid encoded password |
| CMA_SMTP_FROM_EMAIL | Email id from which notification will be sent | example@gmail.com | Valid email id |
| CMA_SMTP_TO_EMAILS | Email ids to which notification will be sent | alert@example.com | Valid comma seperated email ids |
| CMA_SMTP_CC_EMAILS | CC email ids | | Valid comma seperated email ids |
| CMA_SMTP_BCC_EMAILS | BCC email ids | | Valid comma seperated email ids |
| CMA_ENABLE_GOTIFY_NOTIFICATION | Enable Gotify notification | false | `true`, `false`, `yes`, `no` |
| CMA_GOTIFY_URL | Gotify api endpoint URL | gotify.local | Any valid hostname string |
| CMA_GOTIFY_ENCODED_APP_TOKEN | Base64 encoded application token from Gotify | | Any base64 string |
| CMA_GOTIFY_PRIORITY | Gotify notification priority | 5 | Any number |
| CMA_ENABLE_TELEGRAM_NOTIFICATION | Enable telegram notification | false | `true`, `false`, `yes`, `no` |
| CMA_TELEGRAM_ENCODED_BOT_TOKEN | Base64 enccoded telegram bot token | | Any base64 string |
| CMA_TELEGRAM_CHAT_ID | Chat id of telegram | 000000000 | Any valid chat id |
| CMA_MONITOR_RADARR | Monitor radarr for maintenance | false | `true`, `false`, `yes`, `no` |
| CMA_RADARR_URL | URL of radarr | http://radarr:7878 | Radarr hostname/url. Always use HTTPS link if using public hostname/url |
| CMA_RADARR_ENCODED_API_KEY | Base64 encoded Radarr API key | | Any base64 encoded key |
| CMA_RADARR_ENABLE_NOTIFICATION | Send notification related to Radarr maintenance ? | false | `true`, `false`, `yes`, `no` |
| CMA_MONITOR_SONARR | Monitor sonarr for maintenance | false | `true`, `false`, `yes`, `no` |
| CMA_SONARR_URL | URL of sonarr | http://sonarr:8989 | Sonarr hostname/url. Always use HTTPS link if using public hostname/url |
| CMA_SONARR_ENCODED_API_KEY | Base64 encoded Sonarr API key | | Any base64 encoded key |
| CMA_SONARR_ENABLE_NOTIFICATION | Send notification related to Sonarr maintenance ? | false | `true`, `false`, `yes`, `no` |
| | | |


All the credential and sensitive values must be added after base64 encoding it. You can base64 encode any string using below command in Linux terminal.

```
echo 'str1ng-to_enc0de' | base64 -w 0
```

## Setup telegram bot for telegram notification

Follow below article on medium to learn about setting up telegram bot for notifications and alerts.

[Setup Telegram bot to get alert notifications](https://medium.com/linux-shots/setup-telegram-bot-to-get-alert-notifications-90be7da4444)

