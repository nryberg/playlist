Simple is good

The current architecture is cutting edge and at the heart of it, is pretty cool that it's portable and fast (at least in theory).  But at the end of the day, I'm carrying a lot of mental baggage and using up new technology credits that could be better used in the actual development of the analytics. 

I'm not using LevelDb - I can keep going down the path of simplifying the whole BoltDB algorithm, or I can just switch to a simpler SQL set up.  While that feels a little like cheating, at least it would be comprehensible.  The worst of it is that I'm neck deep in a current structure that's going to have to take some work to get out of the pokey.

1. Switch out to a new fetch paradigm.
2. Don't loose the current data
3. Figure out a transition plan that cuts over between one and the other
4. Start building analytics on new data - develop summary queries
5. Stop worrying about byte conversions.  While it's cool, it's way too low level to get anything actually done. 