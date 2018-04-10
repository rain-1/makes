# makes

This tool is used as a replacement for make.

> I wanted to remove stress from my life by not using makefiles anymore.

Most make implementations are around 20-30 thousand lines of code. This one is less than 1 thousand.

# Usage

`makes <output> <input-1> <input-2> ...`

`makes -v <output> <input-1> <input-2> ...`

It exits with status 0 if either
* the output file does not exist
* one of the inputs is newer than result

Otherwise it exists with code 1.

You can use this in a shell script like this:

```
makes data.o \
        data.c data.h &&
        $CC $CFLAGS -o data.o data.c
```

More complete examples can be found:
* [build tarot](https://github.com/rain-1/tarot-vm/blob/master/makesfile).
* [build jq](https://gist.github.com/rain-1/bd0d745c3bd0c4a643a49b74a8c5eb4a)
