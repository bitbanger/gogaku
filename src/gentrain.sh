#!/bin/sh

python drawimg.py < ../txt/joyo.txt 
./trainer ../txt/db.txt ../img/training/
