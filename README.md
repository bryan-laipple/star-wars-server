# Star Wars DB Server

The Star Wars DB Server is a poject to help its author experiment and continue to gain working knowledge in the [Go](https://golang.org/) programming language.

The project has two main areas of development:
- Logic for ETL of Star Wars data
- A web server providing a RESTful API to that data (read only at the time of this writing)

For an example of a small browser based client using the API, see the Star Wars DB Client [repo](https://github.com/bryan-laipple/star-wars-client).

## Deployment

The **`deploy-to-aws-ecs.sh`** script builds a Docker image to run the server process.  The image is tagged and pushed to AWS ECR.  An ECS task definition is updated and corresponding service is modifed to use this new task definition.

Initial setup of the ECR repository, ECS cluster, task definition and service are necessary before running the deploy script.

## The Data

The data and images that make up the content served by the API are from the [Star Wars API](http://swapi.co/) as well as [Wookieepedia](http://starwars.wikia.com/wiki/Main_Page).

The **`etl`** package contains logic to extract data from the SWAPI and scrape images from the public Wiki.  A transform and load functions are provided to build a DynamoDB table, modify the JSON data to be loaded into [DynamoDB](https://aws.amazon.com/dynamodb/), then performs batch updates to the table.

Although the extraction step does most of the data gathering, it was necessary to manually modify the output of this script (**`extracted.json`**) to massage some of the images and links.

## The Server

At the time of this writing, [Iris](https://github.com/kataras/iris), is used as the underlying web framework.  Admittedly, the current demands of the types of requests supported (only GETs) is quite small so any package, or the built-in tooling, would be sufficient.  There are plenty of other flavors of web frameworks available and as this project evolves other implementation decisions may be made.

### DynamoDB client and local cache

The web server uses a client object found in the **`storage`** package which caches the reads from DynamoDB to reduce latency (as well as minimize the requests to AWS).

## Disclosure

Although repeating what was already mentioned in the [data ETL section](#user-content-the-data), more disclosure is better than none. I do not take any credit for the actual content served through the API.  The data and images have been provided/stolen from the [Star Wars API](http://swapi.co/) as well as [Wookieepedia](http://starwars.wikia.com/wiki/Main_Page).
