"use strict";

var headers = ['xss', 'csp', 'hsts', 'cors', 'xfo'];

var longHeaders = ['access-control-allow-origin', 'content-security-policy',
  'x-permitted-cross-domain-policies', 'x-download-options',
  'x-content-type-options',
  'public-key-pins', 'server', 'strict-transport-security', 'content-type',
  'x-frame-options', 'x-powered-by', 'x-xss-protection'
];

var headerDescriptions = {
  'access-control-allow-origin': "Allows web request to bypass the same origin policy.",
  'content-security-policy': 'Restricts the list of accept domains to load and execute resources from.',
  'x-permitted-cross-domain-policies': 'Specifies how Adobe products should handle cross domain policies.',
  'x-content-type-options': 'IE sometimes tries to guess the content of a file, called MIME-snigging. Sometimes IE can be tricked into making incorrect decisions on how to handle a file.',
  'x-download-options': 'Used in IE to specify how a file should handled. Usually used to remove the "open with" option'
}

var app = angular.module('hmon', ['ui.router', 'ngResource', 'nvd3']);

app.config(function($stateProvider, $urlRouterProvider, $resourceProvider) {
  $resourceProvider.defaults.stripTrailingSlashes = false;

  $urlRouterProvider.otherwise('/');

  $stateProvider
    .state('home', {
      url: '/',
      templateUrl: 'views/home.html',
      controller: 'HomeController'
    })

  .state('header', {
    url: '/header/:header',
    templateUrl: 'views/header.html',
    controller: 'HeaderController'
  })

  .state('host', {
    url: '/host/:domain',
    templateUrl: 'views/host.html',
    controller: 'HostController'
  })

  .state('hosts', {
    url: '/hosts/',
    templateUrl: 'views/hosts.html',
    controller: 'HostsController'
  });

});

app.controller('SidePanel', function($scope) {
  $scope.headers = longHeaders;
});

app.factory('Scan', function($resource) {
  return $resource('/api/scans/');
});

app.factory('Search', function($resource) {
  return $resource('/api/hosts/search/:search');
})

app.controller('HostsController', function($scope, Search) {
  $scope.search = "";
  $scope.$watch('search', function(newSearch) {
    if (newSearch != "") {
      Search.query({search: newSearch}, function(scans) {
        $scope.scans = scans;
      })
    } else {
      $scope.scans = [];
    }
  });
});

app.controller('HostController', function($scope, Scan, $stateParams) {
  $scope.headers = longHeaders
  $scope.domain = $stateParams.domain
  $scope.scans = Scan.query({ domain: $stateParams.domain, duration: "total"}, function(scans) {
    $scope.recentScan = _.max(scans, function(s) { return s.TS; });
    $scope.headerScanGroups = groupByHeaderChanges(scans);
  });

  $scope.headerDescriptions = headerDescriptions;
});

app.controller('HomeController', function(Scan, $scope) {
  Scan.query({
    duration: 'week'
  }, function(scans) {

    var data = [];
    for (var i = 0; i < longHeaders.length; ++i) {
      var header = longHeaders[i];
      var out = {
        key: header,
        values: bucketCount(scans, header, 'hour'),
      };

      data.push(out);
    }
    $scope.data = normalizeSeries(data, 'hour');
  });

  $scope.options = {
    chart: {
      type: 'stackedAreaChart',
      height: 450,
      margin: {
        top: 20,
        right: 20,
        bottom: 60,
        left: 40
      },
      x: function(d) {
        return d.ts;
      },
      y: function(d) {
        return d.count;
      },
      useVoronoi: false,
      clipEdge: true,
      transitionDuration: 500,
      useInteractiveGuideline: true,
      xAxis: {
        showMaxMin: false,
        tickFormat: function(d) {
          return moment.unix(d).format('L');
        }
      },
      yAxis: {
        tickFormat: function(d) {
          return d3.format(',.2f')(d);
        }
      }
    }
  };

});

app.controller("HeaderController", function(Scan, $scope, $stateParams) {
  $scope.header = $stateParams.header;
  $scope.scans = Scan.query({
    header: $stateParams.header,
    duration: 'month'
  }, function(scans) {
    var out = {
      key: $scope.header,
      values: bucketCount(scans, $scope.header, 'hour'),
    };

    $scope.data = normalizeSeries([out], 'hour');
  });

  $scope.options = {
    chart: {
      type: 'stackedAreaChart',
      height: 450,
      margin: {
        top: 20,
        right: 20,
        bottom: 60,
        left: 40
      },
      x: function(d) {
        return d.ts;
      },
      y: function(d) {
        return d.count;
      },
      useVoronoi: false,
      clipEdge: true,
      transitionDuration: 500,
      useInteractiveGuideline: true,
      xAxis: {
        showMaxMin: false,
        tickFormat: function(d) {
          return moment.unix(d).format('L');
        }
      },
      yAxis: {
        tickFormat: function(d) {
          return d3.format(',.2f')(d);
        }
      }
    }
  };


});

function bucketCount(scans, header, period) {
  var buckets = bucket(scans, header, period);
  var out = []
  for (var ts in buckets) {

    out.push({
      ts: ts,
      count: buckets[ts].length
    });
  }
  return out
}

function bucket(scans, header, period) {
  scans = _.filter(scans, function(s) {
    return _.contains(s.Headers, header)
  })

  scans = groupByTimeUniqDomains(scans, period)
  return scans
}

function groupByTimeUniqDomains(scans, period) {
  scans = _.groupBy(scans, function(s) {
    return moment.unix(s.TS).startOf(period).unix();
  });

  for (var ts in scans) {
    scans[ts] = _.uniq(scans[ts], function(s) {
      return s.Domain
    })
  }

  return scans
}

function normalizeSeries(series, period) {
  var minTS = moment().unix() + 10
  var maxTS = 0
  for (var i = 0; i < series.length; i++) {
    var values = series[i].values;

    for (var v = 0; v < values.length; v++) {
      if (values[v].ts > maxTS) {
        maxTS = values[v].ts
      }

      if (values[v].ts < minTS) {
        minTS = values[v].ts
      }
    }
  }

  // fill values
  for (var i = 0; i < series.length; i++) {
    var values = series[i].values;
    for (var ts = moment.unix(minTS); ts <= moment.unix(maxTS); ts = ts.add(1,
        period)) {
      if (!_.some(values, function(v) {
          return v.ts == ts.unix()
        })) {
        values.push({
          ts: ts.unix(),
          count: 0
        })
      }
    }
    values = _.sortBy(values, function(v) {
      return v.ts
    });
    series[i].values = values
  }
  return series
}

// Given [scans]
// Returns { x-xss-protection: [{value: script-src abc.com, ts: 3}, {value: style-src def.com, ts: 6}]}
function groupByHeaderChanges(scans) {
  var out = {};
  for (var i = 0; i < scans.length; ++i) {
    for (var h = 0; h < scans[i].Headers.length; ++h) {
      var header = scans[i].Headers[h];
      var value = scans[i].Values[h].join(' ');

      var scanResult = {value: value, ts: moment.unix(scans[i].TS).format('LL') }
      if (header in out) {

        // Check if value already exist. If it does, and the date is earlier, then replace
        var prevExist = false
        for (var prevScanIndex in out[header]) {
          var prevScan = out[header][prevScanIndex]
          if (prevScan.value === scanResult.value) {
            prevExist = true
            if (prevScan.ts > scanResult.ts) {
              prevScan.ts = scanResult.ts
            }
            break
          }
        }

        if (!prevExist) {
          out[header].push(scanResult)
        }
      } else {
        out[header] = [scanResult]
      }
    }
  }


  for (var header in out) {
    out[header] = _.sortBy(out[header], function(s) { return -s.ts} );
  }

  return out
}

app.filter('capitalize', function() {
  return function(input, all) {
    return (!!input) ? input.replace(/([^\W_]+[^\s-]*) */g, function(txt) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    }) : '';
  }
});
