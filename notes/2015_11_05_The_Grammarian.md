The Grammarian of Data

There should be a clean break between pipelines of data and the mechanics of each individual pipe - going through process artists is making me realize how messy all of this is.  While Golang can be very expressive, the overall syntax is cluttered and hard to read. 

The main function should essentially be a couple of chunks of pipe gluing the rest of it together.  Right now it's a combo platter from the local all you can eat buffet, and the quality just isn't there.

A lot of the data cleanup can and should be done in SQL - it's the right thing to do rather than building iterating loops around huge arrays in Go - let's start migrating this stuff to a platform that gets it.

First steps first - build a pg_fetcher.