# fluxy
A gRPC rate-limiting service which defines an easy interface to define and use any existing rate limiting algorithm with redis. Fluxy also supports easy kubernetes deployment with a helm chart and can support multiple clients at once.
- **RegisterClient:** Registers a new client with custom limits into fluxy database
- **GetClient:** Retrieves existing client limits from the fluxy database
- **UpdateClient:** Updates an existing client rate limits
- **DeleteClient:** Deletes the client from fluxy
- **VerifyLimit:** This is where the rate limiting magic happens. This handler verifies if the application is conforming to the rate limits that it was registered with or not. It can use any rate limiting algorithm that implements the fluxy Strategy interface hence making it extremely easy to implement and try out new algorithms. Adding new algorithms also does not need knowledge of how fluxy works, its as simple as adding a single go file.

Database IO calls are minimized using redis for rate limiting checks. Whenever a new client is registered or an existing one is updated the cache is updated with the client limits hence VerifyLimit need not look into the database for client limits. This also means that rate limits for a client can be updated on the fly since it is being maintained as a database entry in fluxy.

Following is a high level overview of the kubernetes layout for fluxy:
![CleanShot 2023-12-18 at 13 41 51](https://github.com/thegeekywanderer/fluxy/assets/30985448/92ad1a49-48fd-46a0-8b40-da4c2750db4f)
