# Build rule

Each build rule is:

* output string
* inputs []string
* command string

Terminology:
* "final build product" - An output which is not the input of any build rule. (the parentmost).
* "source" - An input which is not the output of any build rule. (the childmost).

declaratively a rule means `to create <output> from <inputs>, execute <command>`.

collectively the rules form a directed acyclic graph. (*1) We add one final 'done' node at the very top, it's children being those final build products.

operationally we have 3 goals:

* incremental: Don't perform a build command unless we need to. This allows a build to be stopped and restarted. (*2)
* dependency: An output may need to be rebuilt if one of its transitive children(*3) has been modified. The build process will wait until all the children have been completed to do that.
* parallelism: Run up to N simultaneous build commands that don't interfere when possible. This gives enormous speed-ups when building something large like GCC on a multicore. We also need to be careful about partially created files: Many build tools will create a file then fill its content in over time. We cannot make use of a file until the process is complete.

* (*1) We should reject cycles. The system may deadlock if we don't.
* (*2) If killed during creation of a file, delete that partial build product. Otherwise we leave behind corrupt data which will cause the next build to fail.
* (*3) childrens children, childrens childrens children, etc.


# Implementation

The build starts from the 'done' node. We create N>=1 build slots for parallelism.
There will be a global hashtable to mark filenames as built. Initially the only things that are considered built are sources.

To build a node:

- check if this output is already considered built using the global hashtable. If so finish. (This will happen for sources, this will also happen if another worker has already built this object before we got here).
- open a lock on building this item.
- build all it's children in a pool of parallel goroutines
- We are going to run a build command if: the output does not exist OR any input is newer than the output. Otherwise finish.
- wait for that pool to complete.
- wait for a build slot to open.
- use that build slot run the build command. 
- wait for the command to exit.
- if the command exited with failure the entire build failed. Otherwise finish.

To finish: mark this object as built. close any locks/mutexes.

* build slots ensure no more than N commands run at once.
* branching out many goroutines ensures we run as many commands at once as possible.
* waiting on the pool of children ensures files are completely finished before used by the next command.
* the build lock ensures each command will only be run once.

