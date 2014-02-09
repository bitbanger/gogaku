# -*- coding: utf-8 -*-

import Image, ImageFont, ImageDraw

joyo = raw_input()
kanji_limit = 100

uni_len = 3 # number of bytes in unicode; awful hack

for i in range(0, kanji_limit * uni_len, uni_len):
	img = Image.new("RGB", (64, 64), "white")

	draw = ImageDraw.Draw(img)

	font = ImageFont.truetype("ARIALUNI.TTF", 64)

	kanji = joyo[i:i+uni_len]
	print kanji

	draw.text((0, -12), kanji.decode("utf-8"), (0, 0, 0), font=font)
	# draw.text((0, 0), u"\u250c", (0, 0, 0), font=font)

	img.save("img/training/" + str(i) + ".png")
