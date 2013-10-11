'use strict';


var m = angular.module('app', [
  'controllers.global',
  'controllers.home',
  'directives.match',
  'directives.placeholder',
  'errorHandler',
  'httpInterceptor',
  'ngSanitize',
  'services.global',
  'services.modals',
  'services.scroll',
  'ui.bootstrap.modal',
  'ui.bootstrap.tpls'
]);


m.config(function($routeProvider, $locationProvider) {
  $routeProvider
      .when('/', {
        templateUrl: '/views/home/home.html',
        controller: 'HomeCtrl'
      })

      /*.when('/accounts/login', {
        templateUrl: '/views/accounts/login.html',
        controller: LoginCtrl,
        resolve: {r: require('notlogged')}
      })*/

      .otherwise({
        templateUrl: '/views/404.html',
        controller: 'NotFoundCtrl'
      });
});
