#!/bin/sh
md5sum valk* | awk '{print "md5list[\""$2"\"] = \""$1"\""}'
