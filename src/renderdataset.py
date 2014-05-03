# -*- coding: utf-8 -*-

import imp, sys

use_pillow = False
try:
	imp.find_module("Image")
except:
	use_pillow = True

if use_pillow:
	from PIL import Image, ImageFont, ImageDraw
else:
	import Image, ImageFont, ImageDraw, sys


joyo = raw_input().decode("utf-8")

for kanji in joyo:
	img = Image.new("RGB", (64, 64), "white")

	draw = ImageDraw.Draw(img)

	font = ImageFont.truetype("../misc/ARIALUNI.TTF", 64)

	draw.text((0, -12), kanji, (0, 0, 0), font=font)

	img.save("../img/training/%s.png" % (kanji.encode("utf-8")))

