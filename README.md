# Usage

 NOTE: cache_clone ONLY worjks with git HTTPS connections. It does not work with SSH connections

Every time an agent needs to clone a git repo from a remote service it takes time. If the agent had a local mirror of the repo, it would take much less time.  cache_clone addresses this problem with two subcommands: clone and push 

As an added convenience (and to avoid passing credentials around in shells/shell scripts), cache_clone also accesses AWS secret manager to get the http credentials used to access the remote git server

Running the clone command uses an existing local mirror or creates one if necessary clone a local working repo.

Running the push command assumes the lcoal working repo is clean and pushes it to the local mirror, then pushes the local mirror to the remote

NOTE: cache_clone expectes to create the local directory and will fail if it exists
NOTE: git must be installed. cache_clone just runs git commands

## Accessing the remote
The program will access AWS secret manager to get the username and token for the git remote before running commands. it requires:
 - a Secret Manager secretId path (which returns a JSON  document in a map structure)
 - the username key that contains the git remote username
 - the token key that contains the git remote token

##  The Clone command
This is an example of how it works

```bash

export AWS_ACCESS_KEY_ID=****
export AWS_SECRET_ACCESS_KEY=****
export AWS_DEFAULT_REGION=us-east-1
#Assumption: A secret exists at secretId path /my/secretId/path
#Assumption: the stash username is stored with the key 'stash_user_name'
#Assumption: the stash token is stored with the key 'stash_token'

# Create a temporary directory for our experiment. In  practice this should NOT be a /tmp directory becuase they  can
# get cleaned up. Use the agent home directory or some other safe location the agent can use
ROOT="$(mktemp -d)"
mkdir -p "${ROOT}/mirror" "${ROOT}/local"


# here we use the time command to measure the performance improvement

# first clone has no mirror, so it has to be created
time cache_clone clone --verbose \
--remote=https://my.git.com/my/project.git\
--mirror="${ROOT}/mirror" \
--local="${ROOT}/local/project" \
--secretID="/aws/secretmanager/secret/path" \
--userKey="AWS_SM_username_key" \
--tokenKey="AWS_SM_token_key"


# ...
# This took ~18 seconds
# build/v0.0.7-28-g457b1d8/darwin/amd64/cache_clone clone --verbose     6.67s user 12.03s system 32% cpu 56.994 total

# the second clone detects the existing local mirror 
time cache_clone clone --verbose \
--remote=https://my.git.com/my/project.git\
--mirror="${ROOT}/mirror" \
--local="${ROOT}/local/project2" \
--secretID="/aws/secretmanager/secret/path" \
--userKey="AWS_SM_username_key" \
--tokenKey="AWS_SM_token_key"

# ...
# the clone is local and takes much less time
# build/v0.0.7-28-g457b1d8/darwin/amd64/cache_clone clone --verbose     0.78s user 1.04s system 69% cpu 2.605 total

```


# Running the Go Tests
AWS_PROFILE or AWS_SECRET_ACCESS_KEY and AWS_SECRET_ACCESS_KEY for an account that has the git remote repo credentials stored in an secret manager secret document

AWSSMSECRETID : This is the AWS Secret Manager secret id (path)
USERNAMEKEY : This is the key for the remote repo user name in the AWS M Secret JSON map
TOKENKEY : This is the key for the remote repo token in the AWS M Secret JSON map 


