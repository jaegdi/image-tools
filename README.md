# report-istags

## DESCRIPTION

istag-mgmt reports image date for a application family (eg. pkp, fpc, aps, ssp)

- For existing Images it operates cluster and family specific. That means it works for one cluster like
    'cid, int, ppr, vpt or pro' and for families like 'pkp, sps, fpc, aps, ...'
    The cluster must be defined by the mandatory parametter '-cluster=[cid|int|ppr|vpt|pro]'
    The family must be defined by the mandatory parameter '-family=[aps|fpc|pkp|ssp]

- For used images it looks in all clusters and reports the istags used by any deploymentconfig, job, 
    cronjob or pod of all namespaces that belong to the application family.

- Generate JSON reports about imagestreamtags, imagestreams and images. The content of the JSON 
    report can be defined by the mandatory parameter '-output=[is|istag|image|used|all]'.

- Variable output format: json, yaml, csv, table and tabgroup (table with grouped rows for identical content).
    Output as table or tabgroup is best used when piped into less

- filter data for reports. Define parameter -isname=..., -istagname=..., tagname=... or -shaname=...

- delete istags based on filters like 'older than n days' and/or 'istag name like pattern' (not yet implemented)

For this reports the data is collected from the oc cluster defined by parameter '-cluster=...' and
the parameter 'family=...'. For type 'used' (also included in type 'all') from all clusters.

## Usage

### command

 ./report-istags [OPTIONS] 

#### options

##### define source

  -family=string
      family name, eg. pkp, aps, ssp or fpc 

  -cluster=string
      shortname of cluster, eg. cid,int, ppr or pro

##### What to report, which types

  -all
      collect and output all imageStreams imageStreamTags and Image's

  -image
      collect and output of Image's

  -istag
      collect and output of imageStreamTags

  -is
      collect and output of imageStreams

  -used
      collect and output used imageStreams imageStreamTags and Image's from all clusters

##### filter

  -namespace=string
      namespace to look for istags

  -isname=string
      filter output of one imageStream as json, eg. -is=wvv-service

  -istagname=string
      filter output of one imageStreamTag

  -shaname=string
      filter output of a Image with this SHA

  -tagname=string
      filter output all istags with this Tag

##### output format

  -json
      defines JSON as the output format for the reported data. This is the DEFAULT

  -yaml
      defines YAML as the output format for the reported data

  -csv
      defines CSV as the output format for the reported data

  -table
      defines formated ASCI TABLE as the output format for the reported data

  -tabgroup
      defines formated ASCII TABLE WITH GROUPED ROWS as the output format for the reported data.

##### other soptions

  -occlient
      use oc client instead of api call for cluster communication

## EXAMPLES

Report all information for family pkp in cluster cid as json(which is the default output format)

```sh
./report-istags -cluster=cid -family=pkp -all
```

Report only used istags for family pkp as pretty printed table (the output is paginated to fit your screen size
so it is best use with less. Then you can go up or down with the page key)

```sh
./report-istags -cluster=cid -family=pkp -used -table | less
```

Report istags for family aps in cluster int as yaml report

```sh
./report-istags -cluster=int -family=aps -istag -yaml
```
