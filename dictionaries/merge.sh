#!/bin/bash

cat old_dict.txt <(sed '/[A-Z]/d' websters.txt | grep -v '^.$' | grep -v '^..$') | sort -u > merged.txt 
