# Go SOLR Proxy

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/KensoDev/go-solr-proxy?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Build Status](https://travis-ci.org/KensoDev/go-solr-proxy.svg?branch=master)](https://travis-ci.org/KensoDev/go-solr-proxy)

## Why?

We at [Gogobot](http://www.gogobot.com) are working heavily with SOLR and we had some really painful use cases with it that we wanted to solve.

### 1. Load Balancing

![Solr Load Balancing](http://aviioblog.s3.amazonaws.com/solr-load-balancing.png)

We use multiple SOLR slaves and we want a way to load balance between them.

Also, while backing up an instance we take down Jetty for that instance, meaning SOLR is basically down as far as the application is concerned, we want to put this in Standby mode until the backup is finished.

### 2. Single configurable endpoint

Since we are using SOLR heavily, we want to be able to switch out servers without deploying production or changing configuration files. We want a single access point to all servers that will be load balanced and the same for the master and all the slaves.

This way, we can reindex a full cluster from scratch and just switch out production with no configuration change in the app.

### 3. Partial update document cache

Implementing Partial updates in SOLR is crucial for indexing speeds.

Say you want to update a single field across all documents, you have to reindex your entire cluster, which can take days.

SOLR partially supports this allowing you to update a field in a document, however there's a limitation [Read Here](https://cwiki.apache.org/confluence/display/solr/Updating+Parts+of+Documents)

This linked documentation mentions this limitation pretty clearly

![SOLR Partial updates](http://aviioblog.s3.amazonaws.com/screen-shot-2015-06-16-gh6de.png)


### Getting around that limitation

ElasticSearch bypasses this limitation smartly

> In Updating a Whole Document, we said that the way to update a document is to retrieve it, change it, and then reindex the whole document. This is true. However, using the update API, we can make partial updates like incrementing a counter in a single request.

> We also said that documents are immutable: they cannot be changed, only replaced. The update API must obey the same rules. Externally, it appears as though we are partially updating a document in place. Internally, however, the update API simply manages the same retrieve-change-reindex process that we have already described. The difference is that this process happens within a shard, thus avoiding the network overhead of multiple requests. By reducing the time between the retrieve and reindex steps, we also reduce the likelihood of there being conflicting changes from other processes.


Here's the plan here:

Each of our documents has this: 

```xml
<add><doc boost="0.5"><field name="id">Restaurant 4000000690324</field>
```

As you can see, there's a unique id for that document.

![Proxy save document](http://aviioblog.s3.amazonaws.com/proxy-save-document.png)

Now that the proxy is in the way, it will save the document to s3 under `Restaurant/4000000690324` and also send the request to SOLR for the actual update.

Now, you have a document cache that you can grab from and update a single field. Since you send that document through the same pipeline, you will also have the new document on S3 and in SOLR.

Since there's no DB connection when reindexing the document, and there's no rebuilding of that XML, the update process is super fast.

## Usage

```

usage: solr_proxy [<flags>]

Flags:
  --help           Show help (also see --help-long and --help-man).
  --listen-port=LISTEN-PORT
                   Which port should the proxy listen on
  --master=MASTER  Location to your master server
  --slaves=SLAVES  Comma separated list of servers that act as slaves
  --aws-region="us-west-2"
                   Which AWS region should it use for the cache
  --aws-endpoint="https://s3-us-west-2.amazonaws.com"
                   AWS Endpoint for s3
  --bucket-name=BUCKET-NAME
                   What's the bucket name you want to save the documents in
  --log-location="stdout"
                   Where do you want to keep logs
  --bucket-prefix=BUCKET-PREFIX
                   Prefix after the bucket name before the filename

```

### View configuration when running

When the proxy is running it has a web-accesible configuration JSON.  
You can access it by going to `/proxy/configuration`.

It looks like this:

![Proxy Configuration](http://aviioblog.s3.amazonaws.com/screen-shot-2015-06-30-i4isg.png)

## CHANGELOG

[View Here](CHANGELOG.md)


## Status

This is under active development at [Gogobot](http://gogobot.com).

We are currently testing it on staging and a small percentage of production traffic.  
Once this is running all production traffic, I will release 1.0 with some benchmarks and data.

## Contibutors

* [@kensodev](http://twitter.com/kensodev) [Github](http://github.com/KensoDev)  
* [@kenegozi](http://twitter.com/kenegozi) [Github](http://github.com/kenegozi)
  Even though Ken is not in the commit logs (yet) he contributed a lot. We paired and he contributed a lot of feedback and insights. 
