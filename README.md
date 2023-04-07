# cookunity-exercise
Take home exercise required for appylying to Cook Unity as Senior Backend Engineer.
## Problem statement
As a Senior Backend Candidate I need to Design and Implement a service that can expose a RESTfull API that will trace the location information associated with a determined IP address received by the service. Also it is required that the API will return statistics about several metrics associated with the resoluted location.
### Requirements 
- [x] Upload the solution to a versioning tool (Github, Bitbucket, GitLab).
- [ ] Document how your API is to be used. ( Upload swagger file to swagger UI ? ).
### Optional
- [ ]Deploy the solution to a free cloud platform ( Google App Engine, AWS , etc ).
### Recommendations
- [ ] Leverage using the following free APIs to complete the exercise.
	a. IP Geolocation: https://ip-api.complete
	b. Currency conversions: https://fixer.io.
- [ ] Bear in mind that expected througput may vary from 1k to 5M requests per minute,
   The solution has to be deployed in a high-concurrency environment.
- [X] Assume reasonable answers to all the questions you may have, and document them clearly in the readme file.
## Resolution process summary
### Understanding the problem - Duh !
Yes, Understanding the problem might sound trivial but even going for the first intuition which means to do a simple handover between the request service and [ip-api](https://ip-api.com/#pricing) service requires a quick research. Just looking at the root of the page ) one can see that the free tier of the service only allows 45 rpm then requests will be throttled. Moreover looking at the [documentation](https://ip-api.com/docs/api:json) there is some warnings on going beyond limits too often, IP addresses could be banned.
In the bright side, we know that the limits are IP based and nothing more. 
Last but not least, we found that the API has in most countries an average latency of 50ms ( we would love to have percentile metrics in order to analyze latency behavior, a mail was sent for asking further details and response is expected ).

#### So then, What is the solution that will be persue?
For the sake of being reasonable and trying to not get caught in any complexity I'm inclined to begin with the following analysis:

1. We want a cost effective solution meaning that although is not ideal to increase the cost of the solution, IP-API PRO plan does not seem to be very expensive [this](https://members.ip-api.com/#pricing) most if we compare to the potential cost that would take to recreate the entire database scraping the free tier API. Same rationale is be applied for the usage of [fixer](https://fixer.io) but first we want to focus on the first third party solution (and according to the nature of the exercise a key integration).
2. Before moving forward a deep analysis will be done trying to follow an hipotetical "in house" solution and see if there is any benefit on that approach. My No-brain impulse tells me that the infraestructure for making something from scratch will be more expensive than paying to IP-API. We will see ...

#### In house from scratch deep analysis
What extra cost we would incure if we choose doing everything from scratch? Lets summarize the main points
First of all, lets try to describe how such "from scratch" solution would looks like. Lets start with the simplest approach:

If we don't want to pay any extra fee to IP-API we could laverage the usage of the batch endpoint that allows to handle 100 IPs per call, as was menthion IP-API has some rate limiting 
and throttling after reaching certain limit, we need to take that into account in this solution.

So, how this proposal will work then? Given the API contract of the exercise, we find that the entity associated with the IP that we need to trace is ===country===. If we can find
how the IPv4 ( I'm simplifying the problem disregarding IPv6 ) subnetting asigment was made by IANA , we can then use those indexes for scrapping the api, if we know a subnet associated with a country we can go and iterate over the subnet and use the [batch](https://ip-api.com/docs/api:batch) to get all the information related for each IP that belongs to a certain range ( in this case a country). We can say that we have a k


Next table compares the infraestructure cost that will take for handling 1k to 5m RPS with a cost free API usage.


2) We 


Last, but not least.
This is far for being ideal because looking at the 



### 
