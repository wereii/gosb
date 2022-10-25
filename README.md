## Go SponsorBlock

This is **unofficial** SponsorBlock server implementation that does only segments lookup
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

Instance of this server:

https://sb.doubleuu.win


### Notes

The underlying database is modified and the schema has a lot of (small) changes.    
The code was done in a quick manner (e.g. it sucks).  
Currently, it also ignores all the extra options of `/api/skipSegments/:shaHash` 
endpoint like `categories`, `service` and `actionType` (TBD?).  
Service is (currently) forced to `YouTube`. 

### Running

As of now the final DB size is about 7GB (as reported by postgres).

- Clone the repo
- Download `sponsorTimes.csv` from https://sponsor.ajay.app/database and put it into `db-dumps/`
- `docker compose up postgres` - This will create schema and import data from the csv.
    - It might take some time - mainly seems to depend on storage speed (?)
    - Took about 5 minutes on Ryzen 5 3600, 6*2 Threads with 64 GB ram and decent ~500GB NVMe SSD
- `docker compose up -d` - This will also start `gosb` server
- You can test if `gosb` is running and reachable by opening the index (http://127.0.0.1:8000/)

**To use this server with the extension:**

- Go to your browser, click SponsorBlock extension and open Settings,
  replace the `SponsorBlock Server Address:` address at the end of Miscellaneous tab 
- Running the server locally, without fully qualified domain (localhost, 127.0.0.1)
  might not work due to browser security reasons/CORS (?).


### Environment options
- The `POSTGRES_DSN` is required
- `HTTP_PORT` - listening port, 8000 by default
- `DEBUG` - a bit of extra logging (the value is ignored, if it is set, it's enabled)
