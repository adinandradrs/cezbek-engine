[![cezbek-engine-build](https://github.com/adinandradrs/cezbek-engine/actions/workflows/ci.yml/badge.svg?branch=master) ](https://github.com/adinandradrs/cezbek-engine/actions/workflows/ci.yml)[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=adinandradrs_cezbek-engine&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=adinandradrs_cezbek-engine) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=adinandradrs_cezbek-engine&metric=coverage)](https://sonarcloud.io/summary/new_code?id=adinandradrs_cezbek-engine)

# Cezbek Engine

Cezbek engine is a new prototype of core back-end system that helps Kezbek to manage their business, at scale. Currently we only develop for proof of concept on development environment. As per this document is created, Kezbek is officially partnered with LinkSaja, GoPaid, and Josvo as wallet service providers. Then it only works for B2B partner's customer who have MSISDN as identifier. To integrate with wallet service provider Kezbek could connect through a direct H2H or 3rd party payment provider API such as Middletrans and Xenit, based on the lowest service charge. And with the new Cezbek engine hopefully one day Kezbek could have 

- Business to Business (B2B) onboarding for enterprise partner.
- B2B cashback workflow for B2B Partner transaction by their customer.
- B2B loyalty reward workflow for "tiering" gamification.
- Notification to boost customer engagement.
- **Budget friendly payment providers** to give Kezbek finance team a better experience.
- **Dashboard for reporting & analytics** for invoice and reconciliation purpose.

With some notes from the author that

- In development environment we will use third party Mock API server that will be pretended as H2H and 3rd party payment provider. 
- Kezbek does not manage customer whitelists, product, and campaign that partner have and considered as out of scope.
- Cezbek in current phase does not provide any parameter and configuration management so it will be manually handle on prototype version.
- As per the first statement, Kezbek will take MSISDN as the unique identifier.

### Contents

------

- [Technology Background](#technology-background)
	- [Infrastructure Architecture](#infrastructure-architecture)
	- [Data Structure](#data-structure)
	- [System Migration](#system-migration)
 - [Documentation](#documentation)
   - [Installation](#installation)
   - [Unit Test](#installation)
   - [Surrounding](#surrounding)


### Technology Background

------

Cezbek engine is using **monorepo**, a version-controlled code repository that holds many projects. While these projects may be related, they are often logically independent and run by different teams. And looks perfect for a small development team that will begin a development pattern. The new Cezbek engine relies on cloud environment. As this apps that born in the cloud with its limited development resource, we proudly choose Go as the national language to talk. Why Go is selected as the core backend? 

1. Go is gaining its popularity and is **[getting increased in 2022](https://www.tiobe.com/tiobe-index/)** and have a lot of men-resource market in Indonesia.
2. Support multi thread and single thread concurrency called Goroutine to develop an optimize distributed system.
3. Go have many characteristics such as lightweight, cloud agnostic, build as native, and small memory footprint on startup are the advantages.
4. For the HTTP router to serve Cezbek APIs we choose Fiber that have **shown an impressive benchmark** from **[Tech Empower](https://www.techempower.com/benchmarks/#section=data-r21)**.  
5. Like an Alfamidi or Indomaret, it has lot of standard library and many popular plugins that ready to use without need to lot of deep research the compatibility e.g database driver, ORM, memory cache, cloud library, message queue, logging, distributed systems, etc.

For the naming convention in our code, we are following the nature of Go that has been written on **[Effective Go](https://go.dev/doc/effective_go)** and emitted by community. As the vision is growing up on each feature and the business need to minimize the impact of any changes on any features at the earliest stage, development team decide to split into several services instead of a single binary. Technically Cezbek engine itself is a multipurpose service to provide :

- Cezbek APIs to serve APIs that related to Kezbek business operation and B2B partner.
- Cezbek Job to serve backend job and scheduled task to support asynchronous business process.  
- Cezbek Analytics to serve reporting as summary for Kezbek organization and the B2B partner.

Services in Cezbek Engine are not divided into small things (atomic service) due to its commitment to consistent with the **[domain boundary context](https://learn.microsoft.com/en-us/azure/architecture/microservices/model/domain-analysis)** and to prevent **[over-layer service](https://betterprogramming.pub/is-your-microservices-architecture-more-pinball-or-mcdonalds-9d2a79224da7)**. Over-layered service could be the cause of many negative side effect such as higher latency due to a communication flood over protocol, waste of infrastructure resource due to too many deploy, scattered and will be hard to maintain.

Postgres has been choose as Object Relational Database over MySQL and Oracle due to its open-source and widely used on the market (Tokopedia, Gojek, Tiket, BRI, Mandiri, and Traveloka). Redis also has been choose as it is free on cloud and mandatory to have a distributed caching between services, the nice part is Redis has a pooling mechanism. We could image if we are relying on disk I/O to centralize our operation for store, fetch, and put an event in a large scale of process. Tremendously it can be a chaos due to limited queue of I/O operation and getting worse if the operation come in concurrent. 

Below is the table of standard open-source libraries that we used in Cezbek Engine. Each of them have been selected based on "*not so many manual activity*" and we are trying not too have so many variant as much as possible due to vendor library management restriction. 

| Library                             | Version            | Category                                              | Description                                                  |
| ----------------------------------- | ------------------ | ----------------------------------------------------- | ------------------------------------------------------------ |
| Pgxpool                             | v4.17.2            | Storage                                               | Postgres library that support transaction, pool mechanism, cache, and pgxscan to optimize scanning into struct with a lesser ops allocation. Docs for [reference](https://github.com/efectn/go-orm-benchmarks/blob/master/results.md). We avoid to use ORM if we are not smart enough to use a raw SQL. |
| UberZap                             | v1.22.0            | Logger                                                | Log library that used by Uber, the output data statically by default is a JSON format. So it becomes developer friendly if the log is stored into Elastic tool. Docs for [reference](http://hackemist.com/logbench/). |
| Go-Redis                            | v6.15.9            | Storage                                               | Redis library that support pool mechanism. Docs for [reference](https://levelup.gitconnected.com/fastest-redis-client-library-for-go-7993f618f5ab). |
| Consul                              | v1.18.0            | Service                                               | Consul library to support service discovery and config management system as we are not using K8s on development. |
| Viper                               | v1.14.0            | Config                                                | Viper is a library that could do a multiform of configuration. It could fetch the config from a config server, OS environment variable, and YAML format. |
| Fiber                               | v23.6.0            | Router                                                | Built-in library in Fiber framework for HTTP router, middleware (interceptor), and limiter. |
| Fiber                               | v23.6.0            | Monitor                                               | Built in library in Fiber framework for resource monitoring. |
| Fiber w/ stdlib                     | v23.6.0 w/ 1.17.12 | Message Parser                                        | Built-in library in Fiber framework and combined with standard library to support message encoding decoding such as JSON, Struct, and interface. |
| Fiber w/ stdlib                     | v23.6.0 w/ 1.17.12 | File Operation                                        | Built-in library in Fiber framework for disk operation and stdlib for file operation. |
| AWS SDK                             | v1.44.81           | Object Storage, Message Queue, Notification, Security | S3 object storage library to store and fetch file, legally supported by Amazon Web Service.<br /><br />SQS library to do message queue ops, legally supported by Amazon Web Service.<br /><br />SES Notification library to send message and notification such as email<br /><br />Cognito library to connect with Customer Identity Access Management and handle customer authentication logic or data |
| Gojek Heimdall                      | v7.0.2             | External Adaptor                                      | Extended of standard Go HTTP client library for circuit breaker. Actually this library is using Hystrix. We use Heimdall because the sophisticated of casual typing by using its interface, we do not have to define of every context on each HTTP client function. |
| Go Cron w/ redislock                | v1.17.0 w/ v0.8.2  | Job                                                   | Job library that support UNIX cron expression, quartz, and resource lock for distributed job. One of the top usage by community and mostly have many custom time options, we can not benchmark this one but can be trusted due to its active maintainer.<br /><br />Redislock is a library that enable to prevent multi instance execute a same process (distributed scheduler) by using a standard redis lock operation |
| Gosec                               | v2.14.0            | CLI Toolkit                                           | Go CLI SAST checker, standard Go AST by go.dev and has been a basic standard plugin in many CI cloud provider (GitLab, GitHub, CircleCI, and Codacy) |
| Swaggo                              | v1.8.4             | CLI Toolkit and OpenAPI                               | Go CLI Swagger Open API code generator                       |
| Testify, Mockgen, and Counterfeiter | v1.8.1 and v1.6.0  | Unit test and CLI Toolkit                             | Go CLI and library unit test, mock, and stub code generator. Counterfeiter is an extended version of gomock and mockery to generate stub code |

#### Infrastructure Architecture

------

Amazon Web Service (AWS) will be the provider of Platform as a Service and Infrastructure as a Service because it provides many product and configuration to ease Cezbek Engine setup. AWS does not force to develop and do coding on the cloud, thus the developer still can run on their own machine. AWS free features that should be mandatory and included are : 

1. **Amazon S3** low performant on development environment, DevOps team does not have to deploy separate EC2 instance and storage to store file as SFTP
2. **Amazon Cognito** Customer Identity Access Management (CIAM), developer and DevOps team does not have to deploy or develop separate security layer or service 
3. **Amazon SQS** DevOps team does not have to deploy separate MQ product and its manual setup
4. **Amazon SES** DevOps team could connect to SMTP that been used by the organization infrastructure
5. **Amazon Lambda** event driven function to be called on every AWS product such as Cognito so it becomes serverless
6. **Amazon RDS** Postgres Server free-tier with t3.micro
7. **Amazon Elastic Cache** VM server free-tier with t3.micro
8. **Amazon EC2** VM server free-tier with t3.micro

![high-level-arch](https://github.com/adinandradrs/cezbek-engine/blob/master/docs/high-level-arch.jpg?raw=true)

Below is the detail of development budget estimation per month, seems AWS provides a budget friendly environment if compared to another cloud provider.

| Infrastructure      | Provider   | Description                                | Est Price/Month                                              |
| ------------------- | ---------- | ------------------------------------------ | ------------------------------------------------------------ |
| Database            | AWS        | AWS RDS **t3.micro**                       | $0 / month                                                   |
| In-Memory Database  | AWS        | AWS Elastic Cache non Cluster **t3.micro** | $0 / month                                                   |
| Code Repository     | GitHub     | Cloud Git versioning                       | $0 / month with 2000 minutes of CI / month                   |
| Package Artifactory | GitHub     | Cloud Git artifactory                      | $0 / month with limited up to 2GB of storage                 |
| App Server          | AWS        | AWS EC2 **t3.micro**                       | $0 / month                                                   |
| Config Server       | AWS        | AWS EC2 **t3.micro**                       | $0 / month                                                   |
| MQ                  | AWS        | AWS SQS                                    | $0 - $2 / month                                              |
| CIAM                | AWS        | AWS Cognito                                | $0 - $5 / month                                              |
| SonarQube           | SonarCloud | Sonar scanner on cloud by SonarCloud       | $0 for a public repository or $12 / month for a private repository |

Our code delivery to operation is using GitHub action as the CI/CD, and should passed the 4 steps :

1. Merge Request (MR) : Legitimate to push from feature branch into release/sprint-{x} branch.
2. Unit Test : Unit test should pass and have at least 60% of code coverage, Kezbek does not accept for a non testable code or code coverage less than 60% to embrace refactoring in the future. 
3. SAST : secure and no harm as the result and it more than a library audit, Kezbek wants a bug free software with a secure source code.
4. Code Scan : SonarQube Cloud will scan the code quality, men and code should be improvise together along with Kezbek commitment to develop their employee skill and their product quality. A rating for maintain level, A rating for security, 0 of code smell, and 0 of bugs.
5. Package : Compilation, packaging, and send to the machine. For a reason, GitHub Package will be used as artifactory and release versioning.

In contrast for the production environment, we recommend to use docker and K8s as the application backbone and heavily backed up by AWS. Because by using K8s we could scale our apps with a more easy way. So even we are not using it on development environment due to current economy condition, it still have so many similarities because the docker itself is just for a container that non mainly used by the developer. Still it does not break 12 factor apps for dev/prod parity. The packaging will be handled by CI/CD based on the selected branch.

We still recommend Consul and Vault as the configmaps and secrets source, due to its ability as a distributed config server and can be secured using authentication.

![](https://github.com/adinandradrs/cezbek-engine/blob/master/docs/archie-1.0-HLA-PROD-HA.jpg?raw=true)

#### Technology Architecture



| Library         | Version            | Category                  | Description                                                  |
| --------------- | ------------------ | ------------------------- | ------------------------------------------------------------ |
| Pgxpool         | v4.17.2            | Storage                   | Postgres library that support transaction, pool mechanism, cache, and pgxscan to optimize scanning into struct with a lesser ops allocation. Docs for [reference](https://github.com/efectn/go-orm-benchmarks/blob/master/results.md). We avoid to use ORM if we are not smart enough to use a raw SQL. |
| UberZap         | v1.22.0            | Logger                    | Log library that used by Uber, the output data statically by default is a JSON format. So it becomes developer friendly if one day we will store our log into Elastic tool. Docs for [reference](http://hackemist.com/logbench/). |
| Go-Redis        | v6.15.9            | Storage                   | Redis library that support pool mechanism. Docs for [reference](https://levelup.gitconnected.com/fastest-redis-client-library-for-go-7993f618f5ab). |
| Consul          |                    | Service                   | Consul library to support service discovery and config management system as we are not using K8s on development. |
| Viper           |                    | Config                    | Viper is a library that could do a multiform of configuration. It could fetch the config from a config server, OS environment variable, and YAML format. |
| Fiber           | v23.6.0            | Router                    | Built-in library in Fiber framework for HTTP router, middleware (interceptor), and limiter. |
| Fiber           | v23.6.0            | Monitor                   | Built in library in Fiber framework for resource monitoring. |
| Fiber w/ stdlib | v23.6.0 w/ 1.17.12 | Message Parser            | Built-in library in Fiber framework and combined with standard library to support message encoding decoding such as JSON, Struct, and interface. |
| Fiber w/ stdlib | v23.6.0 w/ 1.17.12 | File Operation            | Built-in library in Fiber framework for disk operation and stdlib for file operation. |
| S3 AWS SDK      | v1.44.81           | Object Storage            | S3 object storage library to store and fetch file, legally supported by Amazon Web Service. |
| SQS AWS SDK     | v1.44.81           | Message Queue             | SQS library to do message queue ops, legally supported by Amazon Web Service. |
| Gojek Heimdall  | v7.0.2             | External Adaptor          | Extended of standard Go HTTP client library for circuit breaker. Actually this library is using Hystrix. We use Heimdall because the sophisticated of casual typing by using its interface, we do not have to define of every context on each HTTP client function. |
| Go Cron         | v1.17.0            | Job                       | Job library that support UNIX cron expression, quartz, and resource lock for distributed job. One of the top usage by community and mostly have many custom time options, we can not benchmark this one but can be trusted due to its active maintainer. |
| Gosec           | v2.14.0            | CLI Toolkit               | Go CLI SAST checker, standard Go AST by go.dev and has been a basic standard plugin in many CI cloud provider (GitLab, GitHub, CircleCI, and Codacy). |
| Swaggo          | v1.8.4             | CLI Toolkit and OpenAPI   | Go CLI Swagger Open API code generator.                      |
| Mockgen         | v1.6.0             | CLI Toolkit and Unit Test | Go CLI and library unit test, mock, and stub code generator. |

#### Data Structure

------

![high-level-arch](https://github.com/adinandradrs/cezbek-engine/blob/master/docs/erd.png?raw=true)

#### System Migration

------

For configmap and secrets we will use a repository called cezbek-sre-ops. For data to generate table on a schema and data manipulation we will use a repository called cezbek-sre-data. Data migration will be automated by CI/CD that run Flyway.

### Documentation

------

To begin with Cezbek Engine product, all developer can refer to this section to start. Our product requirement(s) are : 

1. [Consul](https://www.consul.io/) as configuration server and service discovery
2. [WireMock](https://wiremock.org/) with at least JRE11 as mock tool for external API
3. Postgres 13 as the database
4. Redis 6 as the memory cache
5. AWS S3 as object storage
6. AWS SQS as queue service
7. AWS Cognito as Customer Identity Access Management (CIAM)

#### Installation

------

At least have installed Go 1.17 or above. Some toolkit that need also to be installed are :

1. [Swaggo](https://github.com/swaggo/swag) to generate Open API specs
2. Mockgen to generate code mock

For those who just checkout Cezbek Engine should run command below to **download libraries dependencies**

```
go mod tidy
```

**To run API router** on local could run the command below, just make sure the port to run the HTTP server not yet used

```
go run cmd/api/router.go .
```

**To run on scheduler** on local could run the command below

```
go run cmd/job/cron.go .
```

**To generate OpenAPI specification on router** could run the command below 

```
//swag init
```

