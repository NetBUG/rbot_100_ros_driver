
# Readme for the Feature registry

## Documentation
* Requirments/Design document: Cisco box repository under FeaturePipeline
  https://cisco.app.box.com/notes/133543562267

##Build:
export GOPATH=.
mkdir -p src/github.com/CiscoZeus
cd src/github.com/CiscoZeus
git clone github.com/CiscoZeus/zeus0-analytics OR ln -s ../../../zeus-analytics zeus-analytics  (Create a link to your zeus-analytics repo)
cd zeus-analytics/registry
rm -rf vendor
glide install
go build

##Testing
Unit tests code/scripts in the tests directory

## TBD
1. API:
 1.1 Delete operations:
 	Delete cross references:
		Delete of feature should delete all its jobs/data/schedule
		Delete of job should remove entry from feature
		Delete of data should remove entry from feature and job
		Delete of schedule should remove entry from feature

2. Structures:
	Converted id reference strings to arrays/hashes

3. Clone feature
	- For user FE (see design document)

4. Provisioning library
