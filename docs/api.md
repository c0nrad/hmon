
GET /api/hosts
  // all scans in the last 24 hours
  Returns all hosts, with their latest scan

GET /api/hosts/google.com
  Get all scans related to the lost

GET /api/headers
  // counts over today
  Returns all headers,

GET /api/header/xss


Returns { header: "X-XSS-Protection",
          x: { ts, ts, ts, ts },
          y: { 30, 10, 50, 20 }
}

GET /api/
