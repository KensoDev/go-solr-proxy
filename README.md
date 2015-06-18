# Go SOLR Proxy

This project is using [Readme Driven Development](http://tom.preston-werner.com/2010/08/23/readme-driven-development.html).

None of the features described here is fully/partialy working, once it is, the README will reflect it. If you have contributions or feedbacks please feel welcome to do so.

## Why?

### Load Balancing

![Solr Load Balancing](http://aviioblog.s3.amazonaws.com/solr-load-balancing.png)

We use multiple SOLR slaves and we want a way to load balance between them.

Also, while backing up an instance we take down Jetty for that instance, meaning SOLR is basically down as far as the application is concerned, we want to put this in Standby mode until the backup is finished.

### Partial update document cache

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

## Development Status

This is under active development at [Gogobot](http://gogobot.com), I definitely want to get this done in the upcoming weeks.

## Contibutors

[@kensodev](http://twitter.com/kensodev) [Github](http://github.com/KensoDev)  
[@kenegozi](http://twitter.com/kenegozi) [Github](http://github.com/kenegozi)
