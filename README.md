# fluxy
A gRPC rate-limiting service which defines an easy interface to define and use any existing rate limiting algorithm with redis. Fluxy also supports easy kubernetes deployment with a helm chart and can support multiple clients at once.
- **RegisterClient:** Registers a new client with custom limits into fluxy database
- **GetClient:** Retrieves existing client limits from the fluxy database
- **UpdateClient:** Updates an existing client rate limits
- **DeleteClient:** Deletes the client from fluxy
- **VerifyLimit:** This is where the rate limiting magic happens. This handler verifies if the application is conforming to the rate limits that it was registered with or not. It can use any rate limiting algorithm that implements the fluxy Strategy interface hence making it extremely easy to implement and try out new algorithms. Adding new algorithms also does not need knowledge of how fluxy works, its as simple as adding a single go file.

Database IO calls are minimized using redis for rate limiting checks. Whenever a new client is registered or an existing one is updated the cache is updated with the client limits hence VerifyLimit need not look into the database for client limits. This also means that rate limits for a client can be updated on the fly since it is being maintained as a database entry in fluxy.

## Implemented Rate Limiting Algorithms
### Fixed Window Algorithm

The fixed window algorithm is a rate-limiting strategy where once the request limit has been reached, a client will be blocked from making further requests until the expiration time elapses. For instance, if a client has a limit of 50 requests per minute and makes all 50 requests within the first 5 seconds of the minute, it will encounter a block, unable to make additional requests until 55 seconds have passed.

However, a significant downside of this approach is its susceptibility to bursty traffic patterns. For example, allowing a client to exhaust its entire limit rapidly can potentially overload the service. In scenarios where traffic is concentrated within a short timeframe, it could lead to unexpected spikes in load, disrupting the anticipated distribution of requests across the limiting period.

While the fixed window algorithm is straightforward and easy to implement, its inherent limitation in handling bursty traffic requires careful consideration in managing and balancing the load on the service. Alternative rate-limiting algorithms might better accommodate fluctuating or uneven request patterns.

### Rolling Window Algorithm

The rolling window strategy maintains a historical record of traffic within a defined time window, unlike the fixed window strategy. It continuously tracks the traffic generated over a specified period (e.g., the last 5 minutes) to determine if a client should be allowed further requests, offering better protection against sudden bursts of traffic.

This approach requires higher CPU and memory resources due to storing more request information with timestamps. However, it significantly improves resilience against traffic spikes by evenly distributing requests over time.

Utilizing a sorted set with timestamps allows for efficient removal of expired requests. Using UUIDs minimizes conflicts within the set, ensuring smooth operation even under heavy traffic loads.

Implementing the rolling window algorithm enhances the stability and reliability of rate limiting by adapting better to fluctuating traffic patterns.

## High level overview of the kubernetes layout for fluxy:
![CleanShot 2023-12-18 at 13 41 51](https://github.com/thegeekywanderer/fluxy/assets/30985448/92ad1a49-48fd-46a0-8b40-da4c2750db4f)
