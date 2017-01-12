## 1.0.0

* [@KensoDev] - 2016/01/11 - Multiple big changes
* Added support to Skip the S3 write with the `--disable-upload` flag.
* All documents are now written with the collection prefix after the env prefix. This is now ready to accept a mult-core SOLR installation
* Some code refactoring

At this point, this has been running in production for over a year with zero issues

## 0.5.0

* [@KensoDev] - 2015/06/30 - Added a complete docker workflow
  * `/proxy/configuration` will now respond with the configuration JSON so you
  can view the config and make sure it's right.
  * Refactoring 


## 0.2.0

* [@KensoDev] - 2015/06/29 - Added the `BucketPrefix` to `AWSConfig`.
  Added release.sh script

## v0.1.0

* [@KensoDev] - 2015/06/29 - First production ready release. Pre stress
  testing.


