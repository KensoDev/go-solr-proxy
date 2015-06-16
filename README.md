Go SOLR Proxy

This is a WIP, not ready for prouction or even testing at this point.

This proxy is meant to "sit" in the way, between your application and your SOLR
cluster and save all the document passed into it to S3 as well.

The purpose of this is to make SOLR partial updates with local copies and not
worry about what you have stored or not.

This should speed up the change of a single field across all documents (For
example, realculating score for objects)