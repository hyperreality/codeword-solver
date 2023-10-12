#!/bin/bash

cat merged.txt <(sed 's/-//g' lexicon.txt | sed '/[^[:alnum:]_@]/d' | grep -v '^.$' | grep -v '^..$') | tr '[:upper:]' '[:lower:]' | sort -u > merged2.txt 
