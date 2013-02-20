'use strict';


var m = angular.module('errorHandler', ['ng']);


var insideErr = false;
var limitErr = 0;
m.factory('$exceptionHandler', ['$injector', '$log', function($injector, $log) {
  return function(ex, cause) {
    // Log errors to the console too
    $log.error.apply($log, arguments);

    // Protect agains recursive errors
    if (insideErr)
      return;

    // Development mode should trigger errors ot the server
    if (location.hostname && location.hostname == 'localhost')
      return;

    limitErr++;
    if (limitErr <= 3) {
      insideErr = true;

      // Retrieve the info of the exception
      var message = (ex && ex.message) ? ex.message : '~message~';
      var name = (ex && ex.name) ? ex.name : '~name~';
      var stack = (ex && ex.stack) ? ex.stack : '~stack~';

      // Send the info to the server
      var http = $injector.get('$http');
      http.post('/_/reporter', {
        error: ex,
        message: message,
        name: name,
        stack: stack
      });
    }

    // Show some feedback to the user
    $('#http-error').modal();

    insideErr = false;
  }
}]);
