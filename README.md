# image-tools

## Installation

Download it from

- [image-tools  for linux](https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-dev-local/istag_and_image_management/image-tools-linux/image-tools/image-tools))
- [image-tools .exe for windows](https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-dev-local/istag_and_image_management/image-tools-windows/image-tools.exe/image-tools.exe))

and store it somewhere in your PATH. It is a statically linked go program and no installation is neccessary.

## DESCRIPTION

image-tools  reports image or istag details for a application family (eg. pkp, fpc, aps, ssp)
             or
             generate a bash script as output to delete istags

- For __existing Is, IsTags and Images__ it operates cluster and family specific.
    That means it works for __one cluster__ like
    'cid, int, ppr, vpt or pro' and for __one family__ like 'pkp, sps, fpc, aps, ...'
    The cluster must be defined by the mandatory parametter '-cluster=[cid|int|ppr|vpt|pro]'
    The family must be defined by the mandatory parameter '-family=[aps|fpc|pkp|ssp]

- For __used IsTags__ it looks in __all clusters__ and reports the istags used by any

    - deploymentconfig,
    - job,
    - cronjob or
    - pod

    of __all namespaces that belong to the application family__.

- __Generate reports__ about imagestreamtags, imagestreams and images or all together. The content, what
    should be collected from the cluster for the content of the report can be defined by __one or more__ of
    the mandatory parameter '-is', '-istag', '-image', '-used' or '-all'.

- __Variable report output format__. The output format of the report can be chosen ba on of:
  - -json,
  - -yaml,
  - -csv,
  - -table or - tabgroup (table with grouped rows for identical content).
    Output as table or tabgroup is best used when piped into less and is am

- __filter data for reports__. Define parameter -isname=..., -istagname=..., tagname=... or -shaname=...

- _(not yet implemented) delete istags based on filters like 'older than n days' and/or 'istag name like pattern' (not yet implemented)_

For this reports the data is collected from the oc cluster defined by parameter '-cluster=...' and
the parameter 'family=...'. For type 'used' (also included in type 'all') from all clusters.


- __generate delete script for istags__.  Define parameter __-delete__ to switch in delete mode (generate a script as output for deleting istags)
  - further define the details with the additional parameters:
    - -snapshot
    - -nonbuild (pkp specific)
    - -delfilter=regexp-string
    - -isname=string
    - -tagname=string
    - -istagname=string
    - -namespace=string

## Usage

    execute image-tools with parameter -h to get help and examples
### Command

    ./image-tools  -family=... -cluster=... -all|-image|-is|-istag|-used [output format (default json)] [filter (default none)]

the parameters can be specified in any order

#### Parameters

##### Define source, both parameters are mandatory

    -family=string    family name, eg. pkp, aps, ssp or fpc

    -cluster=string   shortname of cluster, eg.
      cid-scp0,  ppr-scp0 or pro-scp0
      cid-scp0, ppr-scp0, vpt-scp0, pro-scp0

##### Type of objects to collect and report. One of them is mandatory

    -all              collect and output all
                      _imageStreams_, _imageStreamTags_, _image's__
                      of the given cluster and the
                      _used-istags_ from all clusters.

    -image            collect and output of Image's

    -istag            collect and output of imageStreamTags

    -is               collect and output of imageStreams

    -used             collect and output used imageStreams imageStreamTags
                      and Image's from all clusters

#### Options

    -delete           activate the delete mode. The same data as for reports is collected from the clusters, but
                      it is used to generate a script as output for istag deleting instead a report.

    -delpattern=regexp-string    optional, can be combined with the following filter parameters
##### Filter (usable for reports and delete mode)

    -namespace=string namespace to look for istags

    -isname=string    filter output of one imageStream as json, eg. -is=wvv-service

    -istagname=string filter output of one imageStreamTag

    -shaname=string   filter output of a Image with this SHA

    -tagname=string   filter output all istags with this Tag

##### Output Format (for reports)

    -json             defines JSON as the output format for the reported data.
                      This is the DEFAULT

    -yaml             defines YAML as the output format for the reported data

    -csv              defines CSV as the output format for the reported data

    -csvfile=string   define the common filename-part for the output files in csv format.
                      For every type of openshift objects a separate file is generated
                      with the following names schema: '<common-filename>-<type>.csv'.

    -table            defines formatted ASCI TABLE as the output format for the
                      reported data

    -tabgroup         defines formatted ASCII TABLE WITH GROUPED ROWS as the
                      output format for the reported data.

##### Other Options

    -occlient         use oc client instead of api call for cluster communication

    -statcfg          us statically defined config for families instead of dynamic generated config based on files of config-tool

    -noproxy          disable the usage of a proxy for OpenShift API requests

    -socks5=string    enable socks5 usage. E.g. -socks5=localhost:65022, this is the default. To disable socks5 set -socks5=no

## EXAMPLES

Report all information for family pkp in cluster cid as json
(which is the default output format)

    ./image-tools  -cluster=cid-scp0 -family=pkp -all

or as table

        ./image-tools  -cluster=cid-scp0 -family=pkp -all -table

or csv in different files for each type of information

    ./image-tools  -cluster=cid-scp0 -family=pkp -all -csvfile=prefix

writes the output to different files 'prefix-type' in current directory

Report only __used__ istags for family pkp as pretty printed table
(the output is paginated to fit your screen size and piped to
    the pager define in the environment variable $PAGER/%PAGER%.
If $PAGER is not set, it try to use 'more')

        ./image-tools  -family=pkp -used -table

or json

        ./image-tools  -family=pkp -used

or yaml

    ./image-tools  -family=pkp -used -yaml

or csv

    ./image-tools  -family=pkp -used -csv

Report istags with tag=latest for family pkp in cluster cid as yaml report

    ./image-tools  -cluster=cid-scp0 -family=aps -istag -yaml -tagname=latest

Report ImageStreams for family aps in cluster int as yaml report

    ./image-tools  -cluster=cid-scp0 -family=aps -is -yaml

Report ImageStreams with name=webcode-service for family pkp in cluster cid as table report

    ./image-tools  -cluster=cid-scp0 -family=pkp -is -isname=webcode-service -table

Delete: Generate a shell script as output to delete old istags(60 days, the default) for family pkp in cluster cid
    and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service and zahlungsstoerung-service

        image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'

Delete: To use the script output to really delete the istags, you can use the following line (you must be an openshift admin):

        image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'|xargs -n 1 -I{} bash -c "{}"

Delete: To only generate a script to delete old snapshot istags:

        image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot

Delete: To delete all not used images of family 'aps' in cluster cid

        image-tools -family=aps -cluster=cid-scp0 -delete  -minage=0 -delpattern='.'

Delete: To delete all hybris istags of family pkp older than 45 days

        image-tools -family=pkp -cluster=cid-scp0 -delete -isname=hybris -minage=45
