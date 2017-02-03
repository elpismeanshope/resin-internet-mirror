# resin-internet-mirror

This is a resin.io application. Beside installing the required software and
configuration for the access point, the `Dockerfile.template` also downloads
a list of websites recursively. The websites then get served on a webserver.


## Content
The main webserver is running on port 80, serving the following:

- /www contains mirrored websites
- /zim serves raw zim data for download
- /kiwix serves raw kiwix index and library
- / redirects to /www/elpissite.weebly.com/


Other services:

- kiwix webserver on port 8080, serving zim files like wikipedia
- ka-lite on port 8008

Zim files can be added by editing `ZIMS`, sites by editing `SITES` in
[Dockerfile.template](Dockerfile.template).

For a list of available ZIM files, see [the kiwix
wiki](http://wiki.kiwix.org/wiki/Content_in_all_languages). Use the filename
from the "Non-indexed ZIM" URLs.


## Configuration
To set the ka-lite admin account (username = admin), the following environment
variables need to be set:

- KA\_EMAIL
- KA\_PASSWORD

Without these, ka-lite will ask to create a admin account when first started.

## Development

Run `make` to build a local docker image. It can run locally by executing:

```
docker run -p 8080:80 resin-internet-mirror
```

And open http://localhost:8080 in your browser.
