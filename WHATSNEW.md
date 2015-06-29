Version 0.6.0 - Jun 28, 2015
============================
Update
- Rewrite gomemcached client with Consistent Hashing lib.

Fix
- Rename TestMain() to init() in ring_test.go to avoid the error of two TestMain() when running "go test".

Version 0.5.5 - Jun 24, 2015
============================
Update
- Implement Consistent Hashing lib but not used in gomemcached client yet.

Version 0.5.0 - Jun 12, 2015
============================
First workable go client of memcached.
