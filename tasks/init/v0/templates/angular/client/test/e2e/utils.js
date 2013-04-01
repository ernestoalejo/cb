'use strict';


// Return a references to the tested app frame.
function app(obj) {
  if (angular.isUndefined(obj))
    throw 'null application object';

  return obj.application.getWindow_();
}


// Use this DSL to delay some actions to the correct point
// in the chain of futures.
angular.scenario.dsl('defer', function() {
  return function(f) {
    return this.addFuture('deferred function', function(done) {
      f.call(this);
      done();
    });
  };
});
