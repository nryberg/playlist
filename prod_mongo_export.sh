#!/bin/bash
mongoexport --host linus.mongohq.com:10031 --db shorten -u $MONGUSER -p $MONGPWD --collection tracks --type=csv --fieldFile fields.txt --out ./tracks.csv
