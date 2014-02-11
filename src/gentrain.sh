#!/bin/sh

if [ ! -d "../img" ]; then mkdir ../img; fi
if [ ! -d "../img/training" ]; then mkdir ../img/training; fi

python renderdataset.py < ../txt/joyo.txt 
./trainer ../txt/db.txt ../img/training/
