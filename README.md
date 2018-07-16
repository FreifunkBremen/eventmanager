# FreifunkManager [![Build Status](https://travis-ci.org/FreifunkBremen/freifunkmanager.svg?branch=master)](https://travis-ci.org/FreifunkBremen/freifunkmanager) [![Coverage Status](https://coveralls.io/repos/github/FreifunkBremen/freifunkmanager/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/freifunkmanager?branch=master)
is a little software to manage gluon nodes with the help of ssh and yanic
(used on the Breminale since 2017)

## Getting Started
- install Go (version 1.10 or higher)
  - note that system packages are usually too old; see https://golang.org/doc/install for install instructions
  - set $GOPATH (`export GOPATH=$HOME/go`)
  - full install:
    - mkdir -p ~/inst
    - cd ~/inst
    - wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
    - tar xf go1.10.3.linux-amd64.tar.gz
    - export PATH=~/inst/go/bin/:$PATH
    - export GOROOT=~/inst/go/

- install nodejs >= 4.8
  - mkdir -p ~/inst
  - cd ~/inst
  - wget https://nodejs.org/dist/v8.11.3/node-v8.11.3-linux-x64.tar.xz
  - tar xf node-v8.11.3-linux-x64.tar.xz
  - export PATH=~/inst/node-v8.11.3-linux-x64/bin/:$PATH

- install yarn (https://yarnpkg.com/en/docs/install)
  - mkdir -p ~/inst
  - cd ~/inst
  - wget https://github.com/yarnpkg/yarn/releases/download/v1.7.0/yarn-v1.7.0.tar.gz
  - tar xf yarn-v1.7.0.tar.gz
  - export PATH=~/inst/yarn-v1.7.0/bin/:$PATH

- download and build freifunkmanager:
  - go get -t github.com/FreifunkBremen/freifunkmanager/...
  - go get github.com/mattn/goveralls
  - go get "golang.org/x/tools/cmd/cover"
  - cd ~/go/src/github.com/FreifunkBremen/freifunkmanager/
  - go build
  - cd webroot
  - yarn install
  - yarn gulp build
- run:
  - ./freifunkmanager -config config_example.conf


## Usage
Visit http://localhost:8080/

Navigation bar at top of page:
- marker icon: ???
- List: show list of all known nodes
  - use Edit link in last column of a node to edit its details; changes made on the Edit page are saved immediately
  - to change just the hostname, double-click on hostname field in list and make your change
- Map: show map of nodes
  - use Layers icon in upper right corner to enable geojson overlay and view clients
- Statistics: show statistics about nodes, clients, used channels...
- Login with text field: enter password (value of "secret" in config file) and click "Login" to log in
  - this is necessary to make any changes
  - there is no user management; anybody who has the password has full access
- blue rectangle on the far right: ??? (connection status?)
  - click to reconnect?


## Technical Details

List of known APs will be retrieved from Yanic (ie. representing live data).
Additionally, APs can be added manually by visiting a page like /#/n/apname (where apname is the name of the new AP), and then 

Each browser tab has a websocket connection to the server, so changes made in one tab will appear immediately in other tabs as well

All changes are saved to state file (eg. /tmp/freifunkmanager.json - can be changed in config file).
