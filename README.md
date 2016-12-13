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
```json
"post-processors": [
  {
    "type": "shell-with-imageid",
    "script": "post-processor"
  }
]
```

post-processor
```bash
#!/bin/bash

# $1: builder [amazon-ebs, googlecompute], $2: image_id

if [ "$1" == "amazon-ebs" ]
then
        python scripts/add_launching_configuration.py $2
fi

if [ "$1" == "googlecompute" ]
then
        python scripts/add_instance_template.py $2
fi
```
