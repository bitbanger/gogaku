# -*- coding: utf-8 -*-

import Image, ImageFont, ImageDraw, sys

joyo = raw_input().decode("utf-8")
kanji_limit = len(joyo)

kanji_buf = [] # buffer of kanji to output to stdout after length

num = 0
for kanji in joyo:
	img = Image.new("RGB", (64, 64), "white")

	draw = ImageDraw.Draw(img)

	font = ImageFont.truetype("../misc/ARIALUNI.TTF", 64)

	kanji_buf.append(kanji)

	draw.text((0, -12), kanji, (0, 0, 0), font=font)
	# draw.text((0, 0), u"\u250c", (0, 0, 0), font=font)

	img.save("../img/training/%s.png" % (kanji.encode("utf-8")))
	
	num += 1

