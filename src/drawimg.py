# -*- coding: utf-8 -*-

import Image, ImageFont, ImageDraw, sys

joyo = raw_input().decode("utf-8")
kanji_limit = len(joyo)

kanji_buf = [] # buffer of kanji to output to stdout after length

num = 0
for kanji in joyo:
	img = Image.new("RGB", (64, 64), "white")

	draw = ImageDraw.Draw(img)

	font = ImageFont.truetype("../ARIALUNI.TTF", 64)

	kanji_buf.append(kanji)

	draw.text((0, -12), kanji, (0, 0, 0), font=font)
	# draw.text((0, 0), u"\u250c", (0, 0, 0), font=font)

	img.save("../img/training/" + str(num) + ".png")
	
	num += 1

print len(kanji_buf)

for k in kanji_buf:
	print k.encode("utf-8")
