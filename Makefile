.PHONY: load-epics load-issues clean

all: load-epics load-issues

load-epics:
	make -C load-epics

load-issues:
	make -C load-issues

clean:
	make -C load-epics clean
	make -C load-issues clean