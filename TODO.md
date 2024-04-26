separate path management from git operations
slone:

clone or pull the mirror repo
- if the mirror repo does not exist, clone the remote repo to the mirror repo
- if the mirror repo exists, pull the remote repo to the mirror repo
clone the local remo from the mirror repo


clone the remote repo to a local directory:
- ensure mirror host path exists, which is the mirror pat + remote host and path (without protocol)
ex. /home/nmarks/tmp/deleteme.j65Rr2/mirror + stash.imprivata.com/scm/cor_ng

now run the clone

ensure mirror parent path exists:
ex. /home/nmarks/tmp/deleteme.j65Rr2/mirror/stash.imprivata.com/scm/cor_ng
mirror path + remote host and path (without protocol) + repo name
[/home/nmarks/tmp/deleteme.j65Rr2/mirror] + [stash.imprivata.com/scm/cor_ng] + [ng.git]

/home/nmarks/tmp/deleteme.j65Rr2/mirror/stash.imprivata.com/scm/cor_ng/ng.git
