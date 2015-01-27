# Monarch

Monarch is a simple app for mocking a large chain of dependent services. We use it to test
scenarios on Vamp.  

## Running Monarch

1. Clone to your local machine and build

        $ git clone https://github.com/magneticio/monarch.git
        $ cd monarch
        $ go build
        $ ./monarch

2. Use Docker

        $ docker run -p 9012:8080 magneticio/monarch:latest

## Startup flag and environment variables

      --port | MONARCH_PORT                             # sets the API port
      --depends_on_host | MONARCH_DEPENDS_ON_HOST       # sets the depedency's host
      --depends_on_port | MONARCH_DEPENDS_ON_PORT       # sets the depedency's port

## API


        GET   /ping                     # responds with a 'pong' message
        GET   /randomwait               # responds after a random wait time
        GET   /host                     # responds with the hostname
        POST  /work                     # responds with a monarch and some extra data. See the rest of the documentation.

## Get some monarchs!

Posting the supplied monarchs.json file to the `/work` endpoint will trigger some "work".
As a response a random monarch will be selected and given in the response. The full response also details some timing details and the hops this requests has seen, i.e. all the hosts in the chain. See the next paragraph for details.

        $ http POST localhost:8080/work < monarchs.json 
        HTTP/1.1 200 OK
        Content-Length: 223
        Content-Type: application/json
        Date: Tue, 27 Jan 2015 12:38:48 GMT

        {
            "backendTime": 1422362328423, 
            "cty": "GB", 
            "endTime": 1422362328423, 
            "hops": [
                {
                    "host": "Tims-MacBook-Pro-2.local", 
                    "id": 0, 
                    "timeStamp": 1422362328423
                }
            ], 
            "hse": "Commonwealth", 
            "nm": "Oliver Cromwell", 
            "roundTripTime": 0, 
            "yrs": "1653-1658"
        }



## Set up a chain of Monarchs

Monarch's raison d'être is using it in a chain. You provide a dependent host/port combination as startup flags or environment variables and all "work" posted to one monarch will be relayed to the other dependent monarch. In the eventual response all "in between hops" will be recorded.

    $ monarch --port=9011 --depends_on_host=localhost --depends_on_port=9012
    $ monarch --port=9012
    $ http POST localhost:9011/work < monarchs.json 
      HTTP/1.1 200 OK
      Content-Length: 300
      Content-Type: application/json
      Date: Tue, 27 Jan 2015 12:50:36 GMT

      {
          "backendTime": 1422363036944, 
          "cty": "GB", 
          "endTime": 1422363036944, 
          "hops": [
              {
                  "host": "Tims-MacBook-Pro-2.local", 
                  "id": 0, 
                  "timeStamp": 1422363036944
              }, 
              {
                  "host": "Tims-MacBook-Pro-2.local", 
                  "id": 1, 
                  "timeStamp": 1422363036944
              }
          ], 
          "hse": "House of Wessex", 
          "nm": "Edward the Confessor", 
          "roundTripTime": 0, 
          "yrs": "1042-1066"
      }

