# packer-post-processor-shell-with-imageid
This is a plugin of packer post-processor to run script after provioners on googlecompute or amazon.

Usage
=====
* go get https://github.com/odg0318/packer-post-processor-shell-with-imageid
* cd $GOPATH/github.com/odg0318/packer-post-processor-shell-with-imageid
* make install

Example
=======
template.json
```
"post-processors": [
  {
    "type": "shell-with-imageid",
    "script": "post-processor"
  }
]
```

post-processor
```
#!/bin/bash

# $1: builder [amazon-ebs, googlecompute], $2: image_id

if [ "$1" == "amazon-ebs" ]; then
fi

if [ "$1" == "googlecompute" ]; then
fi
```
