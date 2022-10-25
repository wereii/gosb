## Go SponsorBlock

This is **unofficial**, **read-only** SponsorBlock server implementation that does only segments lookup
(`/api/skipSegments/:shaHash`).

It's trying to be a very quick and small (temporary) replacement of the main server.
There is no other functionality like adding segments, voting
or even looking up the user id of the submitted segments.
Basically it will load the segments but nothing else will work.

**Though** you can create segments as the extension will store them locally,
and it should be able to upload them if you go back to official server.

Official links:

https://github.com/ajayyy/SponsorBlock  
https://sponsor.ajay.app/

### How To Use With SponsorBlock Extension

- Click on the SponsorBlock extension button then open Settings
- Open Miscellaneous tab on the left and replace the 
`SponsorBlock Server Address:` with an address of a mirror

Instance of this server:

https://sb.doubleuu.win

### Running

As of now the final DB size is about 7GB (as reported by postgres).  
You will also need to host this behind full (fully qualified) domain, valid ssl might also be required.  
Localhost will mostly not work due to browser access limitations.

- Clone the repo
- Download `sponsorTimes.csv` from https://sponsor.ajay.app/database into `db-dumps/`
- `docker compose up postgres` - This will create schema and import data from the csv.
    - It might take some time - mainly seems to depend on storage speed (?)
    - Took about 5 minutes on Ryzen 5 3600, 6*2 Threads with 64 GB ram and decent ~500GB NVMe SSD
- `docker compose up -d` - This will also start `gosb` server
- You can test if `gosb` is running and reachable by opening the index (http://127.0.0.1:8000/)

```shell
git clone https://github.com/wereii/gosb.git
cd gosb
wget -U 'wereii/gosb' https://sponsor.ajay.app/database/sponsorTimes.csv -O db-dumps/sponsorTimes.csv
docker compose up postgres
# docker compose up -d # The gosb server is not exposed by default, edit docker-compose before running
```

### Environment options

- The `POSTGRES_DSN` is required
- `HTTP_PORT` - listening port, 8000 by default
- `DEBUG` - a bit of extra logging (the value is ignored, if it is set, it's enabled)

### Notes

The underlying database is modified and the schema has a lot of (small) changes.    
The code was done in a quick manner (e.g. it sucks).  
Currently, it also ignores all the extra options of `/api/skipSegments/:shaHash`
endpoint like `categories`, `service` and `actionType` (TBD?).  
Service is (currently) forced to `YouTube`. 