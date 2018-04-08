#include <stdlib.h>
#include <stdio.h>

#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

// makes output.o input-1.txt input-2.txt ...
// return success (0) when
//   output does not exist
//   or any input is newer than the output
// otherwise fails, signalling the need for a rebuild

int main(int argc, char **argv) {
        struct stat s1, s2;

        if(stat(argv[1], &s1))
                return 0;

        for(int i = 2; i < argc; i++) {
                if(!stat(argv[i], &s2))
                        if(s1.st_mtime > s2.st_mtime)
                                return 0;
        }

        return 1;
}
