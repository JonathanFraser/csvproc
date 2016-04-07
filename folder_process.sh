#!/bin/bash
for f in $(ls -1 $1); do
	cat $f | ./convertcsv $2
done
