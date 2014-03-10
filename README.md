gogetmeta
=========

Background
----------

Our internal git server matches all **git** to the git http backend. By convention all repositories end with **.git**, so our URLs will look like this:

	http://internal.net/project.git

To use `go get` one has to append **.git** to the import path, otherwise it won't know what VCS to use:

	go get internal.net/project				// won't work	

`go get` will try to fetch meta data about the kind of repository in use that doesn't exist.
For more information see `go help importpath`.

	go get internal.net/project.git 		// works

It's only a minor nuisance, as the import path won't change. Only the directory where the project is stored has `.git` appended:

	src/internal.net/project.git

instead of

	src/internal.net/project

This tool creates the required meta data on the fly, based on the URL of the incoming request.

Configuration
-------------

gogetmeta reads the following environment variables on startup:

* **GOGETMETA_ADDRESS**: the address it will listen on, defaults to **:8080**
* **GOGETMETA_DOMAIN**: the domain that will be used to create the meta data, set this to the name or alias your server can be reached (e.g. internal.net)

Apache Configuration
--------------------

I'm by no means an expert, so take this with a grain of salt. `mod_rewrite` is used to proxy all `go get` metadata requests to **gogetmeta**:

	RewriteEngine on
	RewriteCond %{QUERY_STRING}     ^go-get=1$    [NC]
	RewriteRule ^/(.*)$ http://localhost:8080/$1?go-get=1 [P]

Here is the Git configuration we use:

	SetEnv GIT_PROJECT_ROOT d:/gitrepos
	SetEnv GIT_HTTP_EXPORT_ALL
	SetEnv REMOTE_USER=$REDIRECT_REMOTE_USER

	#make sure apache can access git stuff
	<Directory "c:/Git/libexec/git-core*">
		Order allow,deny
		Allow from all
	</Directory>

	# alias all /git/ to the git http backend
	ScriptAliasMatch \
	    "(?x)^/(.*/(HEAD | \
	                    info/refs | \
	                    objects/(info/[^/]+ | \
	                             [0-9a-f]{2}/[0-9a-f]{38} | \
	                             pack/pack-[0-9a-f]{40}\.(pack|idx)) | \
	                    git-(upload|receive)-pack))$" \
	                    "c:/Git/libexec/git-core/git-http-backend.exe/$1"

	<LocationMatch "^.*git-upload-pack$">
		AuthType Basic
		AuthName git
		AuthUserFile c:/Apache2.2/conf/authuser
		Require valid-user
	</LocationMatch>

	<LocationMatch "^.*git-receive-pack$">
		AuthType Basic
		AuthName git
		AuthUserFile c:/Apache2.2/conf/authuser
		Require valid-user
	</LocationMatch>