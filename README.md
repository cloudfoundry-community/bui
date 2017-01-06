# BOSH Admin UI

BOSH Admin UI to help operators get information quickly. Still a huge WIP

## Getting Started
#### Requirements
* BOSH Director with 1 or more deployments

#### Get required libraries and build
First, use glide to get the dependencies

`
glide install -v -s
`

Next, make the binary

`
make
`

#### Make a config file and run!
Create a simple config to point your BOSH similar to below

`bosh-lite-config.yml`
```yml
listen_addr: :9304
web_root: ui
skip_ssl_validation: true
bosh_addr: https://192.168.50.4:25555
```

Lastly, run the binary and point to the config file

`./bui -c bosh-lite-config.yml`

#### Access Bui on browser 
Go on to your favorite browser and put in 

`127.0.0.1:9304`

The default credentials is admin/admin.



## Features

* List stemcells
* List releases
* List deployments
* SSH Support! ![](https://github.com/cloudfoundry-community/bui/raw/master/images/BOSH_Admin_UI.png)

## Update dependencies

```
glide install -v -s
```
