# cookunity-exercise
Take home exercise required for appylying to Cook Unity as Senior Backend Engineer.
## Problem statement
As a Senior Backend Candidate I need to Design and Implement a service that can expose a RESTfull API that will trace the location information associated with a determine IP address received by the service. Also it is required that the API will return statistics about several metrics associated with the resoluted location.
### Requirements 
- [x] Upload the solution to a versioning tool (Github, Bitbucket, GitLab).
- [ ] Document how your API is to be used. ( Upload swagger file to swagger UI ? ).
### Optional
- [ ] Deploy the solution to a free cloud platform ( Google App Engine, AWS , etc ).
### Recommendations
- Leverage using the following free APIs to complete the exercise,
	1. [IP Geolocation](https://ip-api.complete),
	2. [Currency conversions](https://fixer.io.).
- Bear in mind that expected througput may vary from 1k to 5M requests per minute, the solution has to be deployed in a high-concurrency environment,
- Assume reasonable answers to all the questions you may have, and document them clearly in the readme file.
## Resolution process summary
### Understanding the problem - Duh !
Yes, Understanding the problem might sound trivial but even going for the first intuition which means to do a simple handover between the request service and [IP-API](https://ip-api.com/) service requires a quick research. Just looking at the root of the page one can see that the free tier of the service only allows 45 rpm then requests will be throttled. Moreover looking at the [IP-API docs](https://ip-api.com/docs/api:json) there is some warnings on going beyond limits too often, IP addresses could be banned.

In the bright side, we know that the limits are IP based and nothing more. 
Last but not least, we found that the API has in most countries an average latency of 50ms ( we would love to have percentile metrics in order to analyze latency behavior, a mail was sent for asking further details and response is expected ).

#### So then, What is the solution that will persue?
For the sake of being reasonable and trying to not get caught in any complexity I'm inclined to begin with the following analysis:

1. We want a cost effective solution meaning that although is not ideal to increase the cost of the solution, IP-API PRO plan does not seem to be very expensive [this](https://members.ip-api.com/#pricing) most if we compare to the potential cost that would take to recreate the entire database scraping the free tier API. Same rationale is be applied for the usage of [FIXER](https://fixer.io) but first we want to focus on the first third party solution (and according to the nature of the exercise a key integration),
2. Before moving forward a deep analysis will be done trying to follow an hipotetical "in house" solution and see if there is any benefit on that approach. My No-brain impulse tells me that the infraestructure for making something from scratch will be more expensive than paying to IP-API. 

We will see ...

#### In house from scratch analysisP
What extra cost we would incure if we choose doing everything from scratch? Lets summarize the main points.

First of all, lets try to describe how such "from scratch" solution would look like. Lets start with the simplest approach; If we don't want to pay any extra fee to IP-API we could laverage the usage of the batch endpoint that allows to handle 100IPs per call, as was menthion IP-API has some rate limiting and throttling after reaching certain limit, we need to take that into account in the code and handle it accordingly but it won´t be explain in deep given that our focus is to budget the solution. So ... Why is this relevate Man? just because given the API contract of ===the take home exercise===, was found that the entity associated with the IP that we need to trace is ===country===. If we can find how IPv4 ( I'm simplifying the problem disregarding IPv6 ) subneting asignment was made by IANA we can find which country took that range and then we can go and iterate over that range and use the [batch](https://ip-api.com/docs/api:batch) to get all the information related for that IP that belongs to a country. Why we need to iterate over all the IP-RANGE if find a way to know the country? Because *latitude* and *longitude* is an exact location and is related to the IP rather than the range or the country. What do we mean by this? Basically given an IP, `/traces` returns country wide information but *latitude* and *longitude* which is specific to each IP and of course is more granular. Nevertheless it can happen that for a specific IP the relative information to location could be the same for the range or adyadcents IPs, but in those cases holding the same information to the ip-key can do the job.

Let's focus on the cost that would take to find ( or buy ) an IP2Country database, the options that we find are the following:

1. [IP2Location IP-Country Database](https://www.ip2location.com/database/db1-ip-country). This is a solid solution, not this one specifically but one similar that has more info was used in a high-concurrecy and low latency environment and worked as expected. ( We can have discussion about this on live ).
2. [RFC 8805 Format for Self-Published IP Geolocation Feeds](https://www.rfc-editor.org/rfc/rfc8805). There is a standard that already address the problem of Self-Published IP Geolocation information. There are some file around with the Feeds, for example [LACNIC Feeds](https://www.lacnic.net/3106/2/lacnic/ip-geolocation),latam organization which is in charge for IP range assignation for the region. 

Just for the sake of being extremelly optimal in terms of solution cost, and disregarding any human cost associated with the scrapping or searching of *GEOFeeds*, we choose the second option.

I'm not even will try to do the whole financial cost, but let's assume that [Fixer](https://fixer.com) can be use it for free if we use some type of cache ( More on this later ... ) given that you have 100 API's calls per month and with just one API call you can retrieve all currency rates symbols to USD. Let's assume that we already have all the IP-Range's per country needed for scrapping IP-API. Let's first iterate over all countries without taking into account *latitude* nor *longitude*, and just get country related information. We will use the IP-API `/batch/` endpoint free tier which has 15 request per minute of allowance and you can request for 100 ips at once. If we make the assumption that on average each country could have 100 IP-RANGE's we can do a full iteration of the all 195 countries in the world and will take approximate 15 minutes. We would probably need a mechanism to run this periodically so at least we can guarantee that we are up-to-date at most on the periodicity that we decide to configure. So if we say that we need to run this daily we will have to calculate 15 minutes every day every month of computer power. Let's also add the cost that it will take to store in a NO-SQL database such REDIS to mantein the low latency high-throughput scenario. We will use [Amazon Web Services](https://aws.amazon.com/) for do the pricing estimation in a Cloud environment and we will use the cheapest hardware available to calculate the price. 

Using [AWS Redis Calculator](https://calculator.aws/#/addService/ElastiCache) we found that for having a NO-SQL storage for this solution will take 35.04 USD per month using the smallest REDIS instance type. 
Using [AWS EC2 Calculator](https://calculator.aws/#/addService/ec2-enhancement) we found that for having 15 minutes/day nano instance ( which is the smallest one ) for doing the IP-API free tier scrapping would cost around 3 USD per month.

What's is the positive thing about this solution? The main thing that I can think about is about execution speed and overall latency. What do I mean by this? IP-API according to the documentation has a response time in the order of 50 ms using dedicated servers and accelerators througout the globe, that's means if we only use IP-API as a handover service the overhead per request will be in the order of 50ms ( like a very slow database ), but the in-house solution given using REDIS will have one order of magnitude lower of *DATA* latency. For answer this question ( which solution is better ) we would need to do an extra calculation and understand the scale of the solution, basically how many request per minute the in house solution vs the IP-API Proxy solution can manage.

### Scale, the tie-breaker

Let's take one concept as granted, we need to have 24/7 availability in the API so at least 1 server running all the time will be neccesary. Serverless discussion will be disregared given that the pattern of potential traffic that the service will handle is unknown, and if the service reaches 1M rpm often serverless probably will be very expensive compared to having a small server running all the time.

So, at the end, given that the in-house solution vs IP-API can scale horizontally ( Ok .. Elastic Cache not but lets keep with simplicy ) the comparation will be 1:1 between the most cheapest server that we can compare. The following described test will be run in a 




In appearence, doing a very simple cost calculation 


2) We 


Last, but not least.
This is far for being ideal because looking at the 



### 
