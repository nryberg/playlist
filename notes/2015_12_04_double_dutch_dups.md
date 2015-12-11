Dropping the dupes 

For some reason, the feed is giving me junk repeats.   So the trick is to flag dups by time/stationid/songid, and then nail them all _except_ one.  They all have a playid - so flag them all, then go back and tweak the min playid as a non-drop.


## Stacking the queries


### Dups
vw_dups - identifies the songs that are doubled up within the timeframe for the station
```
SELECT play."time",
    play.stationid,
    play.songid
   FROM play
  GROUP BY play."time", play.stationid, play.songid
 HAVING count(*) > 1;
```

### Lead to duped play ids
vw_duped_playids - Now that you've got the time, station and songs identified, go back and fetch
the play id's:


```
 SELECT pl.playid
   FROM play pl,
    vw_dups dup
  WHERE pl."time" = dup."time" AND pl.stationid = dup.stationid AND pl.songid = dup.songid;
```

Go through and update all dupes to drop => TRUE

UPDATE play SET drop = TRUE WHERE playid in (SELECT playid from vw_duped_playids);

### Go back and undo one of each set 
Now go back and set the very first duped playid to FALSE.  For each timeblock,
you've got 2 or more dups.  You need to update just one of those and the rest
can lie.  

Let's take a look at what we've got:

```
SELECT playid, time, stationid, songid, drop FROM play WHERE playid IN (SELECT playid FROM vw_duped_playids) ORDER BY
stationid, time, playid;
```

| playid |        time         | stationid |  songid  | drop 
| -------+---------------------+-----------+----------+------
| 17561 | 2015-11-16 22:37:02 |        61 | 35245009 | t
| 17566 | 2015-11-16 22:37:02 |        61 | 35245009 | t
| 21950 | 2015-11-17 02:07:01 |        61 | 34715937 | t
| 21959 | 2015-11-17 02:07:01 |        61 | 34715937 | t
| 22563 | 2015-11-17 02:37:02 |        61 | 34715937 | t

So go fetch the playid's for the blocks - second query, and then for each
block, grab the minimum 

``` 
SELECT MIN(playid) FROM vw_duped_blocks GROUP BY time, stationid, songid ORDER BY stationid, time;
```

  min   
--------
  17561
  21950
  22563
  23185
  34418
  99943

Now set those drops to FALSE

``` 
UPDATE play SET drop = FALSE WHERE playid in (SELECT MIN(playid) FROM vw_duped_blocks GROUP BY time, stationid, songid ORDER BY stationid, time);
```

and flip it to drop => FALSE
