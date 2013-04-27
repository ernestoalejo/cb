'use strict';


var m = angular.module('errorHandler', ['ng']);


var insideErr = false;
var limitErr = 0;
m.factory('$exceptionHandler', function($injector, $log, ErrorRegister) {
  return function(ex, cause) {
    // Log errors to the console too
    $log.error.apply($log, arguments);

    // Protect against recursive errors
    if (insideErr)
      return;

    limitErr++;
    if (limitErr <= 3) {
      insideErr = true;

      var message = (ex && ex.message) ? ex.message : '~message~';
      var name = (ex && ex.name) ? ex.name : '~name~';
      var stack = (ex && ex.stack) ? ex.stack : '~stack~';

      // Filter some common errors and ignore others
      if (message.indexOf('missing hash prefix') != -1) {
        insideErr = false;
        return;
      }
      if (message.indexOf('Circular dependency') != -1) {
        insideErr = false;
        return;
      }
      if (message.indexOf('$digest') != -1) {
        message += 'browser:||' + $injector.get('$browser').url() + '||  ';
        message += 'location:||' + $injector.get('$location').absUrl() + '||';
      }

      $injector.get('$http').post('/_/reporter', {
        error: ex,
        message: message,
        name: name,
        stack: stack
      });
    }

    ErrorRegister.set('exception');

    insideErr = false;
  }
});
