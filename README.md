Header Monitor
-------------

Monitor the security headers for a list of hosts.

Continually scans a list of hosts for changes in their security headers. Reports and security recommendations are then generated.

## Models

Host
  Name
  Url
  Description

Scan
  Host_id
  ts
  Headers
    [string]
  Values
    [string]

Report
  Header
  Description
  Possible Mitigations

## API

### Hosts
| Method | Route | Description |
| GET  | /api/hosts/:id         | Returns a specific host |
| GET  | /api/hosts/:id/latest  | Returns lastest scan for a  host |
| POST | /api/hosts/:id/scan    | Runs a scan on specific host |
| POST | /api/hosts/            | Adds a host to the scan list |
| GET  | /api/hosts/            | Returns all hosts            |

### Headers
| Method | Route | Description|
| GET | /api/header/xss | Returns all host using X-XSS-Protection |
| GET | /api/header/csp | Returns all host using CSP              |
...


### Scans
