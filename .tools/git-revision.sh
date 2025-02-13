#!/bin/bash
revisioncount=`git log $(git describe --tags --abbrev=0)..HEAD --oneline | wc -l`
projectversion=`git describe --tags`
cleanversion=${projectversion%%-*}

echo "$cleanversion.$revisioncount"
