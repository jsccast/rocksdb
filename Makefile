include $(GOROOT)/src/Make.inc

TARG=levigo
CGOFILES=\
	batch.go\
	comparator.go\
	cache.go\
	db.go\
	env.go\
	iterator.go\
	options.go\
	conv.go

HFILES=\
	levigo.h

include $(GOROOT)/src/Make.pkg