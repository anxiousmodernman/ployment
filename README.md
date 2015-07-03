# ployment

Webhook-based deployment tool.

The idea:

Webhook -> *ployment* -> download .zip -> unzip -> copy to directory -> run commands


It works because GitHub makes downloading a repo as a zip file easy:

```
wget https://github.com/{username}/{repo}/zipball/master
```

Thanks, GitHub!

### Usage

Build the `/ployment` server binary and run:

```
ployment -config yourConfig.json
```

Right now the config file is this simple json:

```
{
	"repositoryUrl": "https://github.com/.../.../zipball/master",
	"targetDirectory": "~/temp"
}
```

_TODO_

* post-webhook commands
* token-based auth
