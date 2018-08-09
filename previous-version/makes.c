#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

#include <limits.h>
#include <libgen.h>
#include <sys/inotify.h>

#include <sys/stat.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <errno.h>

// makes output.o input-1.txt input-2.txt ...
// return success (0) when
//   output does not exist
//   or any input is newer than the output
// otherwise fails with status (1), signalling that no build needs performed
// if an error occurs status (2) will be produced

#define NEXT do { argc--; argv++; } while(0)
int main(int argc, char **argv) {
	int verbose = 0;
	int rebuild = 0;
	struct stat s1, s2;
	
	NEXT;
	if(!strcmp("-v", argv[0])) {
		verbose = 1;
		NEXT;
	}
	if(getenv("MAKES_VERBOSE")) {
		verbose = 1;
	}
	
	if(stat(argv[0], &s1)) {
		if(verbose) fprintf(stderr, "[makes] output file \"%s\" did not exist\n", argv[0]);
		rebuild = 1;
	}
	
	for(int i = 1; i < argc; i++) {
		if(!stat(argv[i], &s2)) {
			if(s2.st_mtime > s1.st_mtime) {
				if(verbose) fprintf(stderr, "[makes] input file \"%s\" was newer\n", argv[i]);
				rebuild = 1;

			}
		}
	}
	
	return !rebuild;
}
