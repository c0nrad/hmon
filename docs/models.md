
# Models
Host
  Name
  Url
  Description

Request
  Host_id
  ts
  Headers
    [string]
  Values
    [string]

# Queries

### I need all count vs ts for all header types

For each header in headers:
  scans = Scan.find(header in Headers)
  buckets['header'] = bucket(scans, day)

merge buckets
assuming forwarding till next point.

Same x: ["mondya, tuesday, wednesday, thursday"]
[ {header: 'XSS', y: [120, 125, 128, 131]},
  {header: "CSP", y: [2, 2, 2, 3, 4] },
  ...
]

### I need the id for a specific host

Text search of Hosts

### I need all the scans for a spcific host

Scans.find({ hostid: host.id })

### I need all the hosts that returns true for a specific header

Scans.find({header in "header"})
