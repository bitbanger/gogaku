#!/bin/sh

go build trainer.go dfe.go
go build recog.go dfe.go
go build contour.go dfe.go
