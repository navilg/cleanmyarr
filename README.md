# cleanmyarr
A lightweight utility to delete movies and shows from Radarr and Sonarr after specified time

# Plan

- A docker container will run.
- Docker container will run a process which will periodically check the movies and series in Radarr and Sonarr.
- Movies/Series will be tagged with format `cleanmyarr-3` (Which means clean automatically after 3 days)
- Movies/Series can be tagged with `donotclean` (Which means this movie/series will be ignored by cleanarr)
- If there none of the tag is added in a movie/series, It will be cleaned after default cleanup time set.
- Movie/Series will first be first tagged as `markedfordeletion` on T-3 days. A notification would be sent (Checking technical possibility) (Probably using smtp server)
- For movies/series tagged to be deleted in less than 3 days, No `markedfordeletion` tag will be added.


# Configuration options

- Period: number type Every month, Every 7 days, Every 3 days, Every day,  Default: Every day
- Default cleanup time - number type, 0, any number of days, Default: 0 (0 means do not clean until movie/series is tagged)
- SMTP configuration for notification.
- Radarr and Sonarr URL and API key

# Application working plan

- When process starts, It checkss for configuration file under $HOME/config/
- Config directory can be changed by setting environment variable CMA_CONFIG_DIR
- If no config.yaml file found in config dir, Creates sample config.yaml file with everything disabled.
- It runs first check soon after it starts.
- Repeats as per period defined in config.yaml
