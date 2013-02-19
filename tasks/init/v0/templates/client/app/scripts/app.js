'use strict';


var m = angular.module('testcb', ['errorHandler', 'ngSanitize',
  'directives.match', 'httpInterceptor', 'services.global']);


m.config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
    $locationProvider.hashPrefix('!');

    $routeProvider
        .when('/', {
          templateUrl: '/views/home/home.html',
          controller: HomeCtrl
        })

        .when('/accounts/login', {
          templateUrl: '/views/accounts/login.html',
          controller: LoginCtrl,
          resolve: {r: require('notlogged')}
        })

        .otherwise({
          templateUrl: '/home/global/404.html',
          controller: NotFoundCtrl
        });
  }]);


var requirements = {
  notlogged: function(User) { return !User.isLogged(); },
  logged: function(User) { return User.isLogged(); },
  admin: function(User) { return User.isAdmin(); }
};


function require(arr) {
  if (!angular.isArray(arr))
    arr = [arr];

  return ['$q', '$timeout', 'User', function($q, $timeout, User) {
    var defer = $q.defer();
    for (var i = 0; i < arr.length; i++) {
      var r = requirements[arr[i]];
      if (!r)
        throw new Error('unknown requirement: ' + arr[i]);

      if (!r(User))
        defer.reject(arr[i]);
      else
        defer.resolve(true);
    }
    return defer.promise;
  }];
}

