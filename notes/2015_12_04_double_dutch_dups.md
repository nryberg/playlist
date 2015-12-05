Dropping the dupes 

For some reason, the feed is giving me junk repeats.   So the trick is to flag dups by time/stationid/songid, and then nail them all _except_ one.  They all have a playid - so flag them all, then go back and tweak the min playid as a non-drop.


## Stacking the queries


vw_dups - identifies the songs that are doubled up within the timeframe for the station
```
SELECT play."time",
    play.stationid,
    play.songid
   FROM play
  GROUP BY play."time", play.stationid, play.songid
 HAVING count(*) > 1;
```
